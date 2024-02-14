package coldata

// Batch is the type that columnar operators receive and produce. It
// represents a set of column vectors (partial data columns) as well as
// metadata about a batch, like the selection vector (which rows in the column
// batch are selected).
type Batch interface {
	// Length returns the number of values in the columns in the batch.
	Length() int
	// SetLength sets the number of values in the columns in the batch. Note
	// that if the selection vector will be set or updated on the batch, it must
	// be set **before** setting the length.
	SetLength(int)
	// Capacity returns the maximum number of values that can be stored in the
	// columns in the batch. Note that it could be a lower bound meaning some
	// of the Vecs could actually have larger underlying capacity (for example,
	// if they have been appended to).
	Capacity() int
	// Width returns the number of columns in the batch.
	Width() int
	// ColVec returns the ith Vec in this batch.
	ColVec(i int) Vec
	// ColVecs returns all of the underlying Vecs in this batch.
	ColVecs() []Vec

	// AppendCol appends the given Vec to this batch.
	AppendCol(Vec)
	// ReplaceCol replaces the current Vec at the provided index with the
	// provided Vec. The original and the replacement vectors *must* be of the
	// same type.
	ReplaceCol(Vec, int)
	// String returns a pretty representation of this batch.
	String() string
}

// MemBatch is an in-memory implementation of Batch.
type MemBatch struct {
	// length is the length of batch or sel in tuples.
	length int
	// capacity is the maximum number of tuples that can be stored in this
	// MemBatch.
	capacity int
	// b is the slice of columns in this batch.
	b      []Vec
	useSel bool
	// sel is - if useSel is true - a selection vector from upstream. A
	// selection vector is a list of selected tuple indices in this memBatch's
	// columns (tuples for which indices are not in sel are considered to be
	// "not present").
	sel []int
}

var _ Batch = &MemBatch{}

func (m *MemBatch) Length() int {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) SetLength(i int) {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) Capacity() int {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) Width() int {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) ColVec(i int) Vec {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) ColVecs() []Vec {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) AppendCol(vec Vec) {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) ReplaceCol(vec Vec, i int) {
	//TODO implement me
	panic("implement me")
}

func (m *MemBatch) String() string {
	//TODO implement me
	panic("implement me")
}
