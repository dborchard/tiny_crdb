package descpb

import "github.com/dborchard/tiny_crdb/pkg/i_sql/sem/catid"

// ID, ColumnID, FamilyID, and IndexID are all uint32, but are each given a
// type alias to prevent accidental use of one of the types where
// another is expected.

// ID is a custom type for {Database,Table}Descriptor IDs.
type ID = catid.DescID
