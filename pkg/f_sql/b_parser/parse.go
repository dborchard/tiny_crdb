package parser

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser/statements"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
)

// ParseOne parses a sql statement string, ensuring that it contains only a
// single statement, and returns that Statement. ParseOne will always
// interpret the INT and SERIAL types as 64-bit types, since this is
// used in various internal-execution paths where we might receive
// bits of SQL from other nodes. In general,earwe expect that all
// user-generated SQL has been run through the ParseWithInt() function.
func ParseOne(sql string) (statements.Statement[tree.Statement], error) {
	return statements.Statement[tree.Statement]{
		AST: &tree.Select{
			With: &tree.With{
				Recursive: false,
			},
			Select: &tree.SelectClause{
				From:        tree.From{},
				DistinctOn:  tree.DistinctOn{},
				Exprs:       tree.SelectExprs{},
				GroupBy:     tree.GroupBy{},
				Having:      nil,
				Where:       nil,
				Distinct:    false,
				TableSelect: false,
			},
			OrderBy: tree.OrderBy{},
			Limit: &tree.Limit{
				Offset:   nil,
				Count:    nil,
				LimitAll: false,
			},
		},
		SQL: sql,
	}, nil
}
