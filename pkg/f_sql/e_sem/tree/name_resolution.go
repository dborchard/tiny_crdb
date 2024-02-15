package tree

// SearchPath encapsulates the ordered list of schemas in the current database
// to search during name resolution.
type SearchPath interface {
	// NumElements returns the number of elements in the SearchPath.
	NumElements() int

	// GetSchema returns the schema at the ord offset in the SearchPath.
	// Note that it will return the empty string if the ordinal is out of range.
	GetSchema(ord int) string
}
