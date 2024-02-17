package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/e_security/username"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/descs"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/eval"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/privilege"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"github.com/lib/pq/oid"
)

type planner struct {
	extendedEvalCtx extendedEvalContext
	// autoCommit indicates whether the plan is allowed (but not required) to
	// commit the transaction along with other KV operations. Committing the txn
	// might be beneficial because it may enable the 1PC optimization. Note that
	// autocommit may be false for implicit transactions; for example, an implicit
	// transaction is used for all the statements sent in a batch at the same
	// time.
	//
	// NOTE: plan node must be configured appropriately to actually perform an
	// auto-commit. This is dependent on information from the optimizer.
	autoCommit bool
	// optPlanningCtx stores the optimizer planning context, which contains
	// data structures that can be reused between queries (for efficiency).
	optPlanningCtx optPlanningCtx
	txn            *kv.Txn
	// internalSQLTxn corresponds to the object returned from InternalSQLTxn.
	// It is here to avoid the need to allocate another structure. The value
	// is initialized lazily. The assumption is that that method is called
	// during statement execution when the planner is in a valid state.
	// The internalSQLTxn may hold on to a stale txn reference and should
	// never be accessed directly. Nothing explicitly resets this field.
	internalSQLTxn internalTxn
}

// extendedEvalContext extends eval.Context with fields that are needed for
// distsql planning.
type extendedEvalContext struct {
	eval.Context
	ExecCfg *ExecutorConfig
}

func (p *planner) ResolveTableName(ctx context.Context, tn *tree.TableName) (tree.ID, error) {
	panic("implement me")
}

func (p *planner) SchemaExists(ctx context.Context, dbName, scName string) (found bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) HasAnyPrivilegeForSpecifier(ctx context.Context, specifier eval.HasPrivilegeSpecifier, user username.SQLUsername, privs []privilege.Privilege) (eval.HasAnyPrivilegeResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveType(ctx context.Context, name *tree.UnresolvedObjectName) (*types.T, error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveTypeByOID(ctx context.Context, oid oid.Oid) (*types.T, error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveOIDFromString(ctx context.Context, resultType *types.T, toResolve *tree.DString) (_ *tree.DOid, errSafeToIgnore bool, _ error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveOIDFromOID(ctx context.Context, resultType *types.T, toResolve *tree.DOid) (_ *tree.DOid, errSafeToIgnore bool, _ error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveFunction(ctx context.Context, name tree.UnresolvedRoutineName, path tree.SearchPath) (*tree.ResolvedFunctionDefinition, error) {
	//TODO implement me
	panic("implement me")
}

func (p *planner) ResolveFunctionByOID(ctx context.Context, oid oid.Oid) (*tree.RoutineName, *tree.Overload, error) {
	//TODO implement me
	panic("implement me")
}

// QueryRowEx executes the supplied SQL statement and returns a single row, or
// nil if no row is found, or an error if more that one row is returned.
//
// The fields set in session that are set override the respective fields if
// they have previously been set through SetSessionData().
func (p *planner) QueryRowEx(
	ctx context.Context,
	opName string,
	override sessiondata.InternalExecutorOverride,
	stmt string,
	qargs ...interface{},
) (tree.Datums, error) {
	return p.InternalSQLTxn().QueryRowEx(ctx, opName, p.Txn(), override, stmt, qargs...)
}

func (p *planner) InternalSQLTxn() descs.Txn {
	if p.txn == nil {
		return nil
	}

	// We lazily initialize the internalSQLTxn structure so that we don't have
	// to pay to initialize this structure if the statement being executed does
	// not execute internal sql statements.
	if p.internalSQLTxn.txn != p.txn {
		ief := p.ExecCfg().InternalDB
		ie := MakeInternalExecutor(ief.server)
		p.internalSQLTxn.init(p.txn, ie)
	}
	return &p.internalSQLTxn
}

// ExecCfg implements the PlanHookState interface.
func (p *planner) ExecCfg() *ExecutorConfig {
	return p.extendedEvalCtx.ExecCfg
}

// Optimizer is part of the eval.Planner interface.
func (p *planner) Optimizer() interface{} {
	return p.optPlanningCtx.Optimizer()
}

// AutoCommit is part of the eval.Planner interface.
func (p *planner) AutoCommit() bool {
	return p.autoCommit
}

func (p *planner) Txn() *kv.Txn {
	return p.txn
}
