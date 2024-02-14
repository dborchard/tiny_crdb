package tree

import "time"

// Datum represents a SQL value.
type Datum interface {
	TypedExpr

	// AmbiguousFormat indicates whether the result of formatting this Datum can
	// be interpreted into more than one type. Used with
	// fmtFlags.disambiguateDatumTypes.
	AmbiguousFormat() bool

	// Compare returns -1 if the receiver is less than other, 0 if receiver is
	// equal to other and +1 if receiver is greater than other.
	// TODO(rafi): Migrate all usages of this to CompareError, then delete this.
	Compare(ctx CompareContext, other Datum) int

	// CompareError is the same as Compare, but it returns an error instead of
	// panicking.
	CompareError(ctx CompareContext, other Datum) (int, error)

	// Prev returns the previous datum and true, if one exists, or nil and false.
	// The previous datum satisfies the following definition: if the receiver is
	// "b" and the returned datum is "a", then for every compatible datum "x", it
	// holds that "x < b" is true if and only if "x <= a" is true.
	//
	// The return value is undefined if IsMin(_ *EvalContext) returns true.
	//
	// TODO(#12022): for DTuple, the contract is actually that "x < b" (SQL order,
	// where NULL < x is unknown for all x) is true only if "x <= a"
	// (.Compare/encoding order, where NULL <= x is true for all x) is true. This
	// is okay for now: the returned datum is used only to construct a span, which
	// uses .Compare/encoding order and is guaranteed to be large enough by this
	// weaker contract. The original filter expression is left in place to catch
	// false positives.
	Prev(ctx CompareContext) (Datum, bool)

	// IsMin returns true if the datum is equal to the minimum value the datum
	// type can hold.
	IsMin(ctx CompareContext) bool

	// Next returns the next datum and true, if one exists, or nil and false
	// otherwise. The next datum satisfies the following definition: if the
	// receiver is "a" and the returned datum is "b", then for every compatible
	// datum "x", it holds that "x > a" is true if and only if "x >= b" is true.
	//
	// The return value is undefined if IsMax(_ *EvalContext) returns true.
	//
	// TODO(#12022): for DTuple, the contract is actually that "x > a" (SQL order,
	// where x > NULL is unknown for all x) is true only if "x >= b"
	// (.Compare/encoding order, where x >= NULL is true for all x) is true. This
	// is okay for now: the returned datum is used only to construct a span, which
	// uses .Compare/encoding order and is guaranteed to be large enough by this
	// weaker contract. The original filter expression is left in place to catch
	// false positives.
	Next(ctx CompareContext) (Datum, bool)

	// IsMax returns true if the datum is equal to the maximum value the datum
	// type can hold.
	IsMax(ctx CompareContext) bool

	// Max returns the upper value and true, if one exists, otherwise
	// nil and false. Used By Prev().
	Max(ctx CompareContext) (Datum, bool)

	// Min returns the lower value, if one exists, otherwise nil and
	// false. Used by Next().
	Min(ctx CompareContext) (Datum, bool)

	// Size returns a lower bound on the total size of the receiver in bytes,
	// including memory that is pointed at (even if shared between Datum
	// instances) but excluding allocation overhead.
	//
	// It holds for every Datum d that d.Size().
	Size() uintptr
}

// Datums is a slice of Datum values.
type Datums []Datum

// CompareContext represents the dependencies used to evaluate comparisons
// between datums.
type CompareContext interface {

	// UnwrapDatum will unwrap the OIDs and potentially the placeholders.
	UnwrapDatum(d Datum) Datum
	GetLocation() *time.Location
	GetRelativeParseTime() time.Time

	// MustGetPlaceholderValue is used to compare Datum
	//MustGetPlaceholderValue(p *Placeholder) Datum
}
