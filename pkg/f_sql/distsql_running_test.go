package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/e_security/username"
	parser "github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvclient/kvcoord"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_testutils/serverutils"
	"github.com/stretchr/testify/require"
	"testing"
)

// Test that we don't attempt to create flows in an aborted transaction.
// Instead, a retryable error is created on the gateway. The point is to
// simulate a race where the heartbeat loop finds out that the txn is aborted
// just before a plan starts execution and check that we don't create flows in
// an aborted txn (which isn't allowed). Note that, once running, each flow can
// discover on its own that its txn is aborted - that's handled separately. But
// flows can't start in a txn that's already known to be aborted.
//
// We test this by manually aborting a txn and then attempting to execute a plan
// in it. We're careful to not use the transaction for anything but running the
// plan; planning will be performed outside of the transaction.
func TestDistSQLRunningInAbortedTxn(t *testing.T) {

	ctx := context.Background()
	s, sqlDB, db := serverutils.StartServer()
	defer s.Stopper().Stop(ctx)

	if _, err := sqlDB.ExecContext(
		ctx, "create database test; create table test.t(a int)"); err != nil {
		t.Fatal(err)
	}
	key := roachpb.Key("a")

	// Plan a statement.
	execCfg := s.ExecutorConfig().(ExecutorConfig)
	sd := NewInternalSessionData(ctx, "test")
	internalPlanner, cleanup := NewInternalPlanner(
		"test",
		kv.NewTxn(ctx, db),
		username.RootUserName(),
		&execCfg,
		sd,
	)
	defer cleanup()
	p := internalPlanner.(*planner)
	query := "select * from test.t"
	stmt, err := parser.ParseOne(query)
	if err != nil {
		t.Fatal(err)
	}

	push := func(ctx context.Context, key roachpb.Key) error {
		// Conflicting transaction that pushes another transaction.
		conflictTxn := kv.NewTxn(ctx, db)
		// Push through a Put, as opposed to a Get, so that the pushee gets aborted.
		if err := conflictTxn.Put(ctx, key, "pusher was here"); err != nil {
			return err
		}
		err = conflictTxn.Commit(ctx)
		require.NoError(t, err)
		t.Log(conflictTxn.Rollback(ctx))
		return err
	}

	// Make a db with a short heartbeat interval, so that the aborted txn finds
	// out quickly.
	ambient := context.Background()
	tsf := kvcoord.NewTxnCoordSenderFactory(s.DistSenderI().(*kvcoord.DistSender))
	shortDB := kv.NewDB(ambient, tsf, s.Clock(), s.Stopper())

	iter := 0
	err = shortDB.Txn(context.Background(), func(ctx context.Context, txn *kv.Txn) error {
		iter++
		if iter == 1 {
			// On the first iteration, abort the txn.
			if err := txn.Put(ctx, key, "val"); err != nil {
				t.Fatal(err)
			}

			if err := push(ctx, key); err != nil {
				t.Fatal(err)
			}

		}

		// Create and run a DistSQL plan.
		rw := NewCallbackResultWriter(func(ctx context.Context, row tree.Datums) error {
			return nil
		})
		recv := MakeDistSQLReceiver(
			ctx,
			rw,
			stmt.AST.StatementReturnType(),
			txn,
		)

		// We need to re-plan every time, since the plan is closed automatically
		// by PlanAndRun() below making it unusable across retries.
		p.stmt = makeStatement(stmt)
		if err := p.makeOptimizerPlan(ctx); err != nil {
			t.Fatal(err)
		}
		defer p.curPlan.close(ctx)

		evalCtx := p.ExtendedEvalContext()
		// We need distribute = true so that executing the plan involves marshaling
		// the root txn meta to leaf txns. Local flows can start in aborted txns
		// because they just use the root txn.
		planCtx := execCfg.DistSQLPlanner.NewPlanningCtx(ctx, evalCtx, p, nil,
			DistributionTypeSystemTenantOnly)

		execCfg.DistSQLPlanner.PlanAndRun(ctx, evalCtx, planCtx, txn, p.curPlan.main, recv, nil)
		return rw.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
	if iter != 2 {
		t.Fatalf("expected two iterations, but txn took %d to succeed", iter)
	}
}
