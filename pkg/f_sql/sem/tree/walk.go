package tree

// Visitor defines methods that are called for nodes during an expression or statement walk.
type Visitor interface {
	// VisitPre is called for each node before recursing into that subtree. Upon return, if recurse
	// is false, the visit will not recurse into the subtree (and VisitPost will not be called for
	// this node).
	//
	// The returned Expr replaces the visited expression and can be used for rewriting expressions.
	// The function should NOT modify nodes in-place; it should make copies of nodes. The Walk
	// infrastructure will automatically make copies of parents as needed.
	VisitPre(expr Expr) (recurse bool, newExpr Expr)

	// VisitPost is called for each node after recursing into the subtree. The returned Expr
	// replaces the visited expression and can be used for rewriting expressions.
	//
	// The returned Expr replaces the visited expression and can be used for rewriting expressions.
	// The function should NOT modify nodes in-place; it should make and return copies of nodes. The
	// Walk infrastructure will automatically make copies of parents as needed.
	VisitPost(expr Expr) (newNode Expr)
}
