package tree

import "github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/catid"

// ID is a custom type for {Database,Table}Descriptor IDs.
type ID = catid.ColumnID
