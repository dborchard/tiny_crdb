package tree

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/volatility"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// UnaryOp is a unary operator.
type UnaryOp struct {
	Typ        *types.T
	ReturnType *types.T
	EvalOp     UnaryEvalOp
	Volatility volatility.V

	types   TypeList
	retType ReturnTyper

	// OnTypeCheck is called when the op is type checked.
	OnTypeCheck func()
}
