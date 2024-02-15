package tree

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// ColumnItem corresponds to the name of a column in an expression.
type ColumnItem struct {
	// TableName holds the table prefix, if the name refers to a column. It is
	// optional.
	//
	// This uses UnresolvedObjectName because we need to preserve the
	// information about which parts were initially specified in the SQL
	// text. ColumnItems are intermediate data structures anyway, that
	// still need to undergo name resolution.
	TableName *UnresolvedObjectName
	// ColumnName names the designated column.
	ColumnName Name
}

var _ Expr = new(ColumnItem)

func (c *ColumnItem) String() string {
	//TODO implement me
	panic("implement me")
}

func (c *ColumnItem) Format(ctx *FmtCtx) {
	//TODO implement me
	panic("implement me")
}

func (c *ColumnItem) Walk(visitor Visitor) Expr {
	//TODO implement me
	panic("implement me")
}

func (c *ColumnItem) TypeCheck(ctx context.Context, semaCtx *SemaContext, desired *types.T) (TypedExpr, error) {
	//TODO implement me
	panic("implement me")
}

// Column retrieves the unqualified column name.
func (c *ColumnItem) Column() string {
	return string(c.ColumnName)
}
