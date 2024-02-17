package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
)

// planTop is the struct that collects the properties
// of an entire plan.
// Note: some additional per-statement state is also stored in
// semaCtx (placeholders).
// TODO(jordan): investigate whether/how per-plan state like
// placeholder data can be concentrated in a single struct.
type planTop struct {
	// stmt is a reference to the current statement (AST and other metadata).
	stmt *Statement
	main planMaybePhysical
}

func (t *planTop) close(ctx context.Context) {

}

// planMaybePhysical is a utility struct representing a plan. It can currently
// use either planNode or DistSQL spec representation, but eventually will be
// replaced by the latter representation directly.
type planMaybePhysical struct {
	planNode planNode
	// physPlan (when non-nil) contains the physical plan that has not yet
	// been finalized.
	physPlan *physicalPlanTop
}

// physicalPlanTop is a utility wrapper around PhysicalPlan that allows for
// storing planNodes that "power" the processors in the physical plan.
type physicalPlanTop struct {
	// PhysicalPlan contains the physical plan that has not yet been finalized.
	*PhysicalPlan
	// planNodesToClose contains the planNodes that are a part of the physical
	// plan (via planNodeToRowSource wrapping). These planNodes need to be
	// closed explicitly since we don't have a planNode tree that performs the
	// closure.
	planNodesToClose []planNode
	// onClose, if non-nil, will be called when closing this object.
	onClose func()
}

// planNode defines the interface for executing a query or portion of a query.
//
// The following methods apply to planNodes and contain special cases
// for each type; they thus need to be extended when adding/removing
// planNode instances:
// - planVisitor.visit()           (walk.go)
// - planNodeNames                 (walk.go)
// - setLimitHint()                (limit_hint.go)
// - planColumns()                 (plan_columns.go)
type planNode interface {
	startExec(params runParams) error

	// Next performs one unit of work, returning false if an error is
	// encountered or if there is no more work to do. For statements
	// that return a result set, the Values() method will return one row
	// of results each time that Next() returns true.
	//
	// Available after startPlan(). It is illegal to call Next() after it returns
	// false. It is legal to call Next() even if the node implements
	// planNodeFastPath and the FastPathResults() method returns true.
	Next(params runParams) (bool, error)

	// Values returns the values at the current row. The result is only valid
	// until the next call to Next().
	//
	// Available after Next().
	Values() tree.Datums

	// Close terminates the planNode execution and releases its resources.
	// This method should be called if the node has been used in any way (any
	// methods on it have been called) after it was constructed. Note that this
	// doesn't imply that startExec() has been necessarily called.
	//
	// This method must not be called during execution - the planNode
	// tree must remain "live" and readable via walk() even after
	// execution completes.
	//
	// The node must not be used again after this method is called. Some nodes put
	// themselves back into memory pools on Close.
	Close(ctx context.Context)
}

// runParams is a struct containing all parameters passed to planNode.Next() and
// startPlan.
type runParams struct {
	// context.Context for this method call.
	ctx context.Context

	// extendedEvalCtx groups fields useful for this execution.
	// Used during local execution and distsql physical planning.
	extendedEvalCtx *extendedEvalContext

	// planner associated with this execution. Only used during local
	// execution.
	p *planner
}
