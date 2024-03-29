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
	schemaResolver
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

	// Corresponding Statement for this query.
	stmt Statement

	// curPlan collects the properties of the current plan being prepared. This state
	// is undefined at the beginning of the planning of each new statement, and cannot
	// be reused for an old prepared statement after a new statement has been prepared.
	curPlan planTop
}

// internalPlannerParams encapsulates configurable planner fields. The defaults
// are set in newInternalPlanner.
type internalPlannerParams struct {
	collection *descs.Collection
}

// InternalPlannerParamsOption is an option that can be passed to
// NewInternalPlanner.
type InternalPlannerParamsOption func(*internalPlannerParams)

// NewInternalPlanner is an exported version of newInternalPlanner. It
// returns an interface{} so it can be used outside of the sql package.
func NewInternalPlanner(
	opName string,
	txn *kv.Txn,
	user username.SQLUsername,
	execCfg *ExecutorConfig,
	sessionData *sessiondata.SessionData,
	opts ...InternalPlannerParamsOption,
) (interface{}, func()) {
	return newInternalPlanner(opName, txn, user, execCfg, sessionData, opts...)
}

// newInternalPlanner creates a new planner instance for internal usage. This
// planner is not associated with a sql session.
//
// Since it can't be reset, the planner can be used only for planning a single
// statement.
//
// Returns a cleanup function that must be called once the caller is done with
// the planner.
func newInternalPlanner(
	// TODO(yuzefovich): make this redact.RedactableString.
	opName string,
	txn *kv.Txn,
	user username.SQLUsername,
	execCfg *ExecutorConfig,
	sd *sessiondata.SessionData,
	opts ...InternalPlannerParamsOption,
) (*planner, func()) {
	return nil, nil
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

// makeOptimizerPlan generates a plan using the cost-based optimizer.
// On success, it populates p.curPlan.
func (p *planner) makeOptimizerPlan(ctx context.Context) error {
	return nil
	//opc := &p.optPlanningCtx

	// If we got here, we did not create a plan above.
	//return opc.runExecBuilder(
	//	ctx,
	//	&p.curPlan,
	//	&p.stmt,
	//	newExecFactory(ctx, p),
	//	execMemo,
	//	p.SemaCtx(),
	//	p.EvalContext(),
	//	p.autoCommit,
	//)
}

func (p *planner) ExtendedEvalContext() *extendedEvalContext {
	return &p.extendedEvalCtx
}
