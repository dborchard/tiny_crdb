package eval

import "github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"

// Context defines the context in which to evaluate an expression, allowing
// the retrieval of state such as the node ID or statement start time.
//
// ATTENTION: Some fields from this struct (particularly, but not exclusively,
// from SessionData) are also represented in execinfrapb.EvalContext. Whenever
// something that affects DistSQL execution is added, it needs to be marshaled
// through that proto too.
// TODO(andrei): remove or limit the duplication.
//
// NOTE(andrei): Context is dusty; it started as a collection of fields
// needed by expression evaluation, but it has grown quite large; some of the
// things in it don't seem to belong in this low-level package (e.g. Planner).
// In the sql package it is embedded by extendedEvalContext, which adds some
// more fields from the sql package. Through that extendedEvalContext, this
// struct now generally used by planNodes.
type Context struct {
	Planner          Planner
	SessionDataStack *sessiondata.Stack
}

// SessionData returns the SessionData the current EvalCtx should use to eval.
func (ec *Context) SessionData() *sessiondata.SessionData {
	if ec.SessionDataStack == nil {
		return nil
	}
	return ec.SessionDataStack.Top()
}
