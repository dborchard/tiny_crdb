package sessiondata

// SearchPath represents a list of namespaces to search builtins in.
// The names must be normalized (as per Name.Normalize) already.
type SearchPath struct {
	paths                []string
	containsPgCatalog    bool
	containsPgExtension  bool
	containsPgTempSchema bool
	tempSchemaName       string
	userSchemaName       string
}

// NumElements returns the number of elements in the search path.
func (s *SearchPath) NumElements() int {
	return 0
}

// GetSchema returns the ith schema element if it is in range.
func (s *SearchPath) GetSchema(ord int) string {
	return ""
}
