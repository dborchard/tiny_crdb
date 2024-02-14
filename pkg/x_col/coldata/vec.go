package coldata

import "github.com/dborchard/tiny_crdb/pkg/i_sql/types"

// Column is an interface that represents a raw array of a Go native type.
type Column interface {
	// Len returns the number of elements in the Column.
	Len() int
}

// SliceArgs represents the arguments passed in to Vec.Append and Nulls.set.
type SliceArgs struct {
	// Src is the data being appended.
	Src Vec
	// Sel is an optional slice specifying indices to append to the destination
	// slice. Note that Src{Start,End}Idx apply to Sel.
	Sel []int
	// DestIdx is the first index that Append will append to.
	DestIdx int
	// SrcStartIdx is the index of the first element in Src that Append will
	// append.
	SrcStartIdx int
	// SrcEndIdx is the exclusive end index of Src. i.e. the element in the index
	// before SrcEndIdx is the last element appended to the destination slice,
	// similar to Src[SrcStartIdx:SrcEndIdx].
	SrcEndIdx int
}

// Vec is an interface that represents a column vector that's accessible by
// Go native types.
type Vec interface {
	// Type returns the type of data stored in this Vec. Consider whether
	// CanonicalTypeFamily() should be used instead.
	Type() *types.T
	// CanonicalTypeFamily returns the canonical type family of data stored in
	// this Vec.
	CanonicalTypeFamily() types.Family

	// Bool returns a bool list.
	Bool() Bools
	// Int32 returns an int32 slice.
	Int32() Int32s
	// Col returns the raw, typeless backing storage for this Vec.
	Col() Column

	// SetCol sets the member column (in the case of mutable columns).
	SetCol(Column)

	// TemplateType returns an []interface{} and is used for operator templates.
	// Do not call this from normal code - it'll always panic.
	TemplateType() []interface{}

	// Append uses SliceArgs to append elements of a source Vec into this Vec.
	// It is logically equivalent to:
	// destVec = append(destVec[:args.DestIdx], args.Src[args.SrcStartIdx:args.SrcEndIdx])
	// An optional Sel slice can also be provided to apply a filter on the source
	// Vec.
	// Refer to the SliceArgs comment for specifics and TestAppend for examples.
	//
	// Note: Append()'ing from a Vector into itself is not supported.
	Append(SliceArgs)

	// Copy uses SliceArgs to copy elements of a source Vec into this Vec. It is
	// logically equivalent to:
	// copy(destVec[args.DestIdx:], args.Src[args.SrcStartIdx:args.SrcEndIdx])
	// An optional Sel slice can also be provided to apply a filter on the source
	// Vec.
	// Refer to the SliceArgs comment for specifics and TestCopy for examples.
	Copy(SliceArgs)

	// Window returns a "window" into the Vec. A "window" is similar to Golang's
	// slice of the current Vec from [start, end), but the returned object is NOT
	// allowed to be modified (the modification might result in an undefined
	// behavior).
	Window(start int, end int) Vec

	// MaybeHasNulls returns true if the column possibly has any null values, and
	// returns false if the column definitely has no null values.
	MaybeHasNulls() bool

	// Nulls returns the nulls vector for the column.
	Nulls() *Nulls

	// SetNulls sets the nulls vector for this column.
	SetNulls(Nulls)

	// Length returns the length of the slice that is underlying this Vec.
	Length() int

	// Capacity returns the capacity of the Golang's slice that is underlying
	// this Vec. Note that if there is no "slice" (like in case of flat bytes),
	// then "capacity" of such object is equal to the number of elements.
	Capacity() int
}

// memColumn is a simple pass-through implementation of Vec that just casts
// a generic interface{} to the proper type when requested.
type memColumn struct {
	t                   *types.T
	canonicalTypeFamily types.Family
	col                 Column
	nulls               Nulls
}

var _ Vec = &memColumn{}

func (m *memColumn) Type() *types.T {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) CanonicalTypeFamily() types.Family {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Bool() Bools {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Int32() Int32s {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Col() Column {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) SetCol(column Column) {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) TemplateType() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Append(args SliceArgs) {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Copy(args SliceArgs) {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Window(start int, end int) Vec {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) MaybeHasNulls() bool {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Nulls() *Nulls {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) SetNulls(nulls Nulls) {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Length() int {
	//TODO implement me
	panic("implement me")
}

func (m *memColumn) Capacity() int {
	//TODO implement me
	panic("implement me")
}
