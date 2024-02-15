package tree

// TableName corresponds to the name of a table in a FROM clause,
// INSERT or UPDATE statement, etc.
//
// This is constructed for incoming SQL queries from an UnresolvedObjectName,
//
// Internal uses of this struct should not construct instances of
// TableName directly, and instead use the NewTableName /
// MakeTableName functions underneath.
//
// TableName is a public type for objName. It exposes the fields
// and can be default-constructed but cannot be instantiated with a
// non-default value; this encourages the use of the constructors below.
type TableName struct {
	objName
}

// objName is the internal type for a qualified object.
type objName struct {
	// ObjectName is the unqualified name for the object
	// (table/view/sequence/function/type).
	ObjectName Name

	// ObjectNamePrefix is the path to the object.  This can be modified
	// further by name resolution, see name_resolution.go.
	ObjectNamePrefix
}

// A Name is an SQL identifier.
//
// In general, a Name is the result of parsing a name nonterminal, which is used
// in the grammar where reserved keywords cannot be distinguished from
// identifiers. A Name that matches a reserved keyword must thus be quoted when
// formatted. (Names also need quoting for a variety of other reasons; see
// isBareIdentifier.)
//
// For historical reasons, some Names are instead the result of parsing
// `unrestricted_name` nonterminals. See UnrestrictedName for details.
type Name string

// ObjectNamePrefix corresponds to the path prefix of an object name.
type ObjectNamePrefix struct {
	CatalogName Name
	SchemaName  Name

	// ExplicitCatalog is true iff the catalog was explicitly specified
	// or it needs to be rendered during pretty-printing.
	ExplicitCatalog bool
	// ExplicitSchema is true iff the schema was explicitly specified
	// or it needs to be rendered during pretty-printing.
	ExplicitSchema bool
}

// An UnrestrictedName is a Name that does not need to be escaped when it
// matches a reserved keyword.
//
// In general, an UnrestrictedName is the result of parsing an unrestricted_name
// nonterminal, which is used in the grammar where reserved keywords can be
// unambiguously interpreted as identifiers. When formatted, an UnrestrictedName
// that matches a reserved keyword thus does not need to be quoted.
//
// For historical reasons, some unrestricted_name nonterminals are instead
// parsed as Names. The only user-visible impact of this is that we are too
// aggressive about quoting names in certain positions. New grammar rules should
// prefer to parse unrestricted_name nonterminals into UnrestrictedNames.
type UnrestrictedName string
