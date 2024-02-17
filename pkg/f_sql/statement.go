package sql

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser/statements"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
)

// Statement contains a statement with optional expected result columns and metadata.
type Statement struct {
	statements.Statement[tree.Statement]
}

func makeStatement(parserStmt statements.Statement[tree.Statement],
) Statement {
	return Statement{
		Statement: parserStmt,
	}
}
