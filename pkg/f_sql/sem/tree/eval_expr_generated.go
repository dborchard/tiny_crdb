package tree

import "context"

// ExprEvaluator is used to evaluate TypedExpr expressions.
type ExprEvaluator interface {
	EvalOrExpr(context.Context, *OrExpr) (Datum, error)

	//EvalAllColumnsSelector(context.Context, *AllColumnsSelector) (Datum, error)
	//EvalAndExpr(context.Context, *AndExpr) (Datum, error)
	//EvalArray(context.Context, *Array) (Datum, error)
	//EvalArrayFlatten(context.Context, *ArrayFlatten) (Datum, error)
	//EvalBinaryExpr(context.Context, *BinaryExpr) (Datum, error)
	//EvalCaseExpr(context.Context, *CaseExpr) (Datum, error)
	//EvalCastExpr(context.Context, *CastExpr) (Datum, error)
	//EvalCoalesceExpr(context.Context, *CoalesceExpr) (Datum, error)
	//EvalCollateExpr(context.Context, *CollateExpr) (Datum, error)
	//EvalColumnAccessExpr(context.Context, *ColumnAccessExpr) (Datum, error)
	//EvalColumnItem(context.Context, *ColumnItem) (Datum, error)
	//EvalComparisonExpr(context.Context, *ComparisonExpr) (Datum, error)
	//EvalDefaultVal(context.Context, *DefaultVal) (Datum, error)
	//EvalFuncExpr(context.Context, *FuncExpr) (Datum, error)
	//EvalIfErrExpr(context.Context, *IfErrExpr) (Datum, error)
	//EvalIfExpr(context.Context, *IfExpr) (Datum, error)
	//EvalIndexedVar(context.Context, *IndexedVar) (Datum, error)
	//EvalIndirectionExpr(context.Context, *IndirectionExpr) (Datum, error)
	//EvalIsNotNullExpr(context.Context, *IsNotNullExpr) (Datum, error)
	//EvalIsNullExpr(context.Context, *IsNullExpr) (Datum, error)
	//EvalIsOfTypeExpr(context.Context, *IsOfTypeExpr) (Datum, error)
	//EvalNotExpr(context.Context, *NotExpr) (Datum, error)
	//EvalNullIfExpr(context.Context, *NullIfExpr) (Datum, error)
	//EvalParenExpr(context.Context, *ParenExpr) (Datum, error)
	//EvalPlaceholder(context.Context, *Placeholder) (Datum, error)
	//EvalRangeCond(context.Context, *RangeCond) (Datum, error)
	//EvalRoutineExpr(context.Context, *RoutineExpr) (Datum, error)
	//EvalSubquery(context.Context, *Subquery) (Datum, error)
	//EvalTuple(context.Context, *Tuple) (Datum, error)
	//EvalTupleStar(context.Context, *TupleStar) (Datum, error)
	//EvalTypedDummy(context.Context, *TypedDummy) (Datum, error)
	//EvalUnaryExpr(context.Context, *UnaryExpr) (Datum, error)
	//EvalUnqualifiedStar(context.Context, UnqualifiedStar) (Datum, error)
	//EvalUnresolvedName(context.Context, *UnresolvedName) (Datum, error)
}
