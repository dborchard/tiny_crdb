package coldata

// Nulls represents a list of potentially nullable values using a bitmap. It is
// intended to be used alongside a slice (e.g. in the Vec interface) -- if the
// ith bit is off, then the ith element in that slice should be treated as NULL.
type Nulls struct {
	nulls []byte
	// maybeHasNulls is a best-effort representation of whether or not the
	// vector has any null values set. If it is false, there definitely will be
	// no null values. If it is true, there may or may not be null values.
	maybeHasNulls bool
}
