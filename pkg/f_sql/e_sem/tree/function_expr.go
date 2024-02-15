package tree

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// FuncExpr represents a function call.
type FuncExpr struct {
	Func  ResolvableFunctionReference
	Type  funcType
	Exprs Exprs
	// Filter is used for filters on aggregates: SUM(k) FILTER (WHERE k > 0)
	Filter Expr
	// AggType is used to specify the type of aggregation.
	AggType AggType
	// OrderBy is used for aggregations which specify an order. This same field
	// is used for any type of aggregation.
	OrderBy OrderBy

	// InCall is true when the FuncExpr is part of a CALL statement.
	InCall bool

	typeAnnotation
	fnProps *FunctionProperties
	fn      *Overload
}

var _ Expr = new(FuncExpr)

func (node *FuncExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (node *FuncExpr) Format(ctx *FmtCtx) {
	//TODO implement me
	panic("implement me")
}

func (node *FuncExpr) Walk(visitor Visitor) Expr {
	//TODO implement me
	panic("implement me")
}

func (node *FuncExpr) TypeCheck(ctx context.Context, semaCtx *SemaContext, desired *types.T) (TypedExpr, error) {
	//TODO implement me
	panic("implement me")
}
