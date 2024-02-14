package tree

// DatumAlloc provides batch allocation of datum pointers, amortizing the cost
// of the allocations.
// NOTE: it *must* be passed in by a pointer.
type DatumAlloc struct {
}
