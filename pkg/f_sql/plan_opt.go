package sql

import "github.com/dborchard/tiny_crdb/pkg/f_sql/opt/xform"

type optPlanningCtx struct {
	p *planner

	optimizer xform.Optimizer
}

// Optimizer returns the Optimizer associated with this planning context.
func (opc *optPlanningCtx) Optimizer() interface{} {
	return &opc.optimizer
}
