package tree

import "context"

// UnaryEvalOp is a unary operation which can be evaluated.
type UnaryEvalOp interface {
	Eval(context.Context, OpEvaluator, Datum) (Datum, error)
}

// OpEvaluator is an evaluator for UnaryEvalOp and BinaryEvalOp operations.
type OpEvaluator interface {
	UnaryOpEvaluator
	BinaryOpEvaluator
}

type UnaryOpEvaluator interface {
	EvalCbrtFloatOp(context.Context, *CbrtFloatOp, Datum) (Datum, error)
}

type BinaryOpEvaluator interface {
	EvalMinusIntOp(context.Context, *MinusIntOp, Datum, Datum) (Datum, error)
}
