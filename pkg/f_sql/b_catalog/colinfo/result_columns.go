package colinfo

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_catalog/descpb"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// ResultColumn contains the name and type of a SQL "cell".
type ResultColumn struct {
	Name string
	Typ  *types.T

	// If set, this is an implicit column; used internally.
	Hidden bool

	// TableID/PGAttributeNum identify the source of the column, if it is a simple
	// reference to a column of a base table (or view). If it is not a simple
	// reference, these fields are zeroes.
	TableID        descpb.ID // OID of column's source table (pg_attribute.attrelid).
	PGAttributeNum uint32    // Column's number in source table (pg_attribute.attnum).
}

// ResultColumns is the type used throughout the sql module to
// describe the column types of a table.
type ResultColumns []ResultColumn
