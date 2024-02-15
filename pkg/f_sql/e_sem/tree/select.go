package tree

import "github.com/dborchard/tiny_crdb/pkg/y_col/coldata"

// Select represents a SelectStatement with an ORDER and/or LIMIT.
type Select struct {
	With    *With
	Select  SelectStatement
	OrderBy OrderBy
	Limit   *Limit
}

// With represents a WITH statement.
type With struct {
	Recursive bool
}

// SelectStatement represents any SELECT statement.
type SelectStatement interface {
	Statement
	selectStatement()
}

// func (*ParenSelect) selectStatement()         {}
func (*SelectClause) selectStatement() {}

// func (*UnionClause) selectStatement()         {}
func (*ValuesClause) selectStatement()        {}
func (*LiteralValuesClause) selectStatement() {}

// SelectClause represents a SELECT statement.
type SelectClause struct {
	From        From
	DistinctOn  DistinctOn
	Exprs       SelectExprs
	GroupBy     GroupBy
	Having      *Where
	Where       *Where
	Distinct    bool
	TableSelect bool
}

func (c *SelectClause) String() string {
	//TODO implement me
	panic("implement me")
}

func (c *SelectClause) StatementReturnType() StatementReturnType {
	//TODO implement me
	panic("implement me")
}

func (c *SelectClause) StatementType() StatementType {
	//TODO implement me
	panic("implement me")
}

func (c *SelectClause) StatementTag() string {
	//TODO implement me
	panic("implement me")
}

// GroupBy represents a GROUP BY clause.
type GroupBy []Expr

// Where represents a WHERE or HAVING clause.
type Where struct {
	Type string
	Expr Expr
}

// SelectExprs represents SELECT expressions.
type SelectExprs []SelectExpr

// SelectExpr represents a SELECT expression.
type SelectExpr struct {
	Expr Expr
	As   UnrestrictedName
}

// DistinctOn represents a DISTINCT ON clause.
type DistinctOn []Expr

// From represents a FROM clause.
type From struct {
	Tables TableExprs
	AsOf   AsOfClause
}

// TableExprs represents a list of table expressions.
type TableExprs []TableExpr

// TableExpr represents a table expression.
type TableExpr interface {
	NodeFormatter
	tableExpr()
	WalkTableExpr(Visitor) TableExpr
}

// AsOfClause represents an as of time.
type AsOfClause struct {
	Expr Expr
}

// LiteralValuesClause is like ValuesClause but values have been typed checked
// and evaluated and are assumed to be ready to use Datums.
type LiteralValuesClause struct {
	Rows ExprContainer
}

// ExprContainer represents an abstract container of Exprs
type ExprContainer interface {
	// NumRows returns number of rows.
	NumRows() int
	// NumCols returns number of columns.
	NumCols() int
	// Get returns the Expr at row i column j.
	Get(i, j int) Expr
}

// VectorRows lets us store a Batch in a tree.LiteralValuesClause.
type VectorRows struct {
	Batch coldata.Batch
}

// NumRows implements the ExprContainer interface.
func (r VectorRows) NumRows() int {
	return r.Batch.Length()
}

// NumCols implements the ExprContainer interface.
func (r VectorRows) NumCols() int {
	return r.Batch.Width()
}

// Get implements the ExprContainer interface.
func (r VectorRows) Get(i, j int) Expr {
	//return DNull
	return nil
}

var _ ExprContainer = VectorRows{}

// ValuesClause represents a VALUES clause.
type ValuesClause struct {
	Rows []Exprs
}

// Exprs represents a list of value expressions. It's not a valid expression
// because it's not parenthesized.
type Exprs []Expr

// OrderBy represents an ORDER BY clause.
type OrderBy []*Order

// Order represents an ordering expression.
type Order struct {
	OrderType OrderType
	Expr      Expr
	//Direction  Direction
	//NullsOrder NullsOrder
	// Table/Index replaces Expr when OrderType = OrderByIndex.
	Table TableName
	// If Index is empty, then the order should use the primary key.
	Index UnrestrictedName
}

// OrderType indicates which type of expression is used in ORDER BY.
type OrderType int

const (
	// OrderByColumn is the regular "by expression/column" ORDER BY specification.
	OrderByColumn OrderType = iota
	// OrderByIndex enables the user to specify a given index' columns implicitly.
	OrderByIndex
)

// Limit represents a LIMIT clause.
type Limit struct {
	Offset, Count Expr
	LimitAll      bool
}

var _ Statement = &Select{}

func (s *Select) String() string {
	//TODO implement me
	panic("implement me")
}

// StatementReturnType implements the Statement interface.
func (*Select) StatementReturnType() StatementReturnType { return Rows }

func (s *Select) StatementType() StatementType {
	//TODO implement me
	panic("implement me")
}

func (s *Select) StatementTag() string {
	//TODO implement me
	panic("implement me")
}
