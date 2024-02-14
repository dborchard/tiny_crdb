package tree

import (
	"context"
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// Expr represents an expression.
type Expr interface {
	fmt.Stringer
	NodeFormatter
	// Walk recursively walks all children using WalkExpr. If any children are changed, it returns a
	// copy of this node updated to point to the new children. Otherwise the receiver is returned.
	// For childless (leaf) Exprs, its implementation is empty.
	Walk(Visitor) Expr
	// TypeCheck transforms the Expr into a well-typed TypedExpr, which further permits
	// evaluation and type introspection, or an error if the expression cannot be well-typed.
	// When type checking is complete, if no error was reported, the expression and all
	// sub-expressions will be guaranteed to be well-typed, meaning that the method effectively
	// maps the Expr tree into a TypedExpr tree.
	//
	// The semaCtx parameter defines the context in which to perform type checking.
	// The desired parameter hints the desired type that the method's caller wants from
	// the resulting TypedExpr. It is not valid to call TypeCheck with a nil desired
	// type. Instead, call it with wildcard type types.Any if no specific type is
	// desired. This restriction is also true of most methods and functions related
	// to type checking.
	TypeCheck(ctx context.Context, semaCtx *SemaContext, desired *types.T) (TypedExpr, error)
}

// TypedExpr represents a well-typed expression.
type TypedExpr interface {
	Expr

	// ResolvedType provides the type of the TypedExpr, which is the type of Datum
	// that the TypedExpr will return when evaluated.
	ResolvedType() *types.T

	// Eval evaluates an SQL expression. Expression evaluation is a
	// mostly straightforward walk over the parse tree. The only
	// significant complexity is the handling of types and implicit
	// conversions. See binOps and cmpOps for more details. Note that
	// expression evaluation returns an error if certain node types are
	// encountered: Placeholder, VarName (and related UnqualifiedStar,
	// UnresolvedName and AllColumnsSelector) or Subquery. These nodes
	// should be replaced prior to expression evaluation by an
	// appropriate WalkExpr. For example, Placeholder should be replaced
	// by the argument passed from the client.
	Eval(context.Context, ExprEvaluator) (Datum, error)
}
