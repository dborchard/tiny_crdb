package storage

import (
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// SimpleMVCCIterator is an interface for iterating over key/value pairs in an
// engine. SimpleMVCCIterator implementations are thread safe unless otherwise
// noted. SimpleMVCCIterator is a subset of the functionality offered by
// MVCCIterator.
//
// API invariants are asserted via assertSimpleMVCCIteratorInvariants().
//
// The iterator exposes both point keys and range keys. Range keys are only
// emitted when enabled via IterOptions.KeyTypes. Currently, all range keys are
// MVCC range tombstones, and this is enforced during writes.
//
// Range keys and point keys exist separately in Pebble. A specific key position
// can have both a point key and multiple range keys overlapping it. Their
// properties are accessed via:
//
// HasPointAndRange(): Key types present at the current position.
// UnsafeKey():        Current position (point key if any).
// UnsafeValue():      Current point key value (if any).
// RangeBounds():      Start,end bounds of range keys at current position.
// RangeKeys():        All range keys/values overlapping current position.
//
// Consider the following point keys and range keys:
//
//	4: a4  b4
//	3: [-------)
//	2: [-------)
//	1:     b1  c1
//	   a   b   c
//
// Range keys cover a span between two roachpb.Key bounds (start inclusive, end
// exclusive) and contain timestamp/value pairs. They overlap *all* point key
// versions within their key bounds regardless of timestamp. For example, when
// the iterator is positioned on b@4, it will also expose [a-c)@3 and [a-c)@2.
//
// During iteration with IterKeyTypePointsAndRanges, range keys are emitted at
// their start key and at every overlapping point key. For example, iterating
// across the above span would emit this sequence:
//
// UnsafeKey HasPointAndRange UnsafeValue RangeKeys
// a         false,true       -           [a-c)@3 [a-c)@2
// a@4       true,true        a4          [a-c)@3 [a-c)@2
// b@4       true,true        b4          [a-c)@3 [a-c)@2
// b@1       true,true        b1          [a-c)@3 [a-c)@2
// c@1       true,false       c1          -
//
// MVCCIterator reverse iteration yields the above sequence in reverse.
// Notably, bare range keys are still emitted at their start key (not end key),
// so they will be emitted last in this example.
//
// When using SeekGE within range key bounds, the iterator may land on the bare
// range key first, unless seeking exactly to an existing point key. E.g.:
//
// SeekGE UnsafeKey HasPointAndRange UnsafeValue RangeKeys
// b      b         false,true       -           [a-c)@3 [a-c)@2
// b@5    b@5       false,true       -           [a-c)@3 [a-c)@2
// b@4    b@4       true,true        b@4         [a-c)@3 [a-c)@2
// b@3    b@3       false,true       -           [a-c)@3 [a-c)@2
//
// Note that intents (with timestamp 0) encode to a bare roachpb.Key, so they
// will be colocated with a range key start bound. For example, if there was an
// intent on a in the above example, then both SeekGE(a) and forward iteration
// would land on a@0 and [a-c)@3,[a-c)@2 simultaneously, instead of the bare
// range keys first.
//
// Range keys do not have a stable, discrete identity, and should be
// considered a continuum: they may be merged or fragmented by other range key
// writes, split and merged along with CRDB ranges, partially removed by GC,
// and truncated by iterator bounds.
//
// Range keys are fragmented by Pebble such that all overlapping range keys
// between two keys form a stack of range key fragments at different timestamps.
// For example, writing [a-e)@1 and [c-g)@2 will yield this fragment structure:
//
//	2:     |---|---|
//	1: |---|---|
//	   a   c   e   g
//
// Fragmentation makes all range key properties local, which avoids incurring
// unnecessary access costs across SSTs and CRDB ranges. It is deterministic
// on the current range key state, and does not depend on write history.
// Stacking allows easy access to all range keys overlapping a point key.
//
// For more information on MVCC range keys, see this tech note:
// https://github.com/cockroachdb/cockroach/blob/master/docs/tech-notes/mvcc-range-tombstones.md
type SimpleMVCCIterator interface {
	// Close frees up resources held by the iterator.
	Close()
	// SeekGE advances the iterator to the first key in the engine which is >= the
	// provided key. This may be in the middle of a bare range key straddling the
	// seek key.
	SeekGE(key MVCCKey)
	// Valid must be called after any call to Seek(), Next(), Prev(), or
	// similar methods. It returns (true, nil) if the iterator points to
	// a valid key (it is undefined to call Key(), Value(), or similar
	// methods unless Valid() has returned (true, nil)). It returns
	// (false, nil) if the iterator has moved past the end of the valid
	// range, or (false, err) if an error has occurred. Valid() will
	// never return true with a non-nil error.
	Valid() (bool, error)
	// Next advances the iterator to the next key in the iteration. After this
	// call, Valid() will be true if the iterator was not positioned at the last
	// key.
	Next()
	// NextKey advances the iterator to the next MVCC key. This operation is
	// distinct from Next which advances to the next version of the current key
	// or the next key if the iterator is currently located at the last version
	// for a key. NextKey must not be used to switch iteration direction from
	// reverse iteration to forward iteration.
	//
	// If NextKey() lands on a bare range key, it is possible that there exists a
	// versioned point key at the start key too. Calling NextKey() again would
	// skip over this point key, since the start key was already emitted. If the
	// caller wants to see it, it must call Next() to check for it. Note that
	// this is not the case with intents: they don't have a timestamp, so the
	// encoded key is identical to the range key's start bound, and they will
	// be emitted together at that position.
	NextKey()
	// MVCCValueLenAndIsTombstone should be called only for MVCC (i.e.,
	// UnsafeKey().IsValue()) point values, when the actual point value is not
	// needed, for example when updating stats and making GC decisions, and it
	// is sufficient for the caller to know the length (len(UnsafeValue()), and
	// whether the underlying MVCCValue is a tombstone
	// (MVCCValue.IsTombstone()). This is an optimization that can allow the
	// underlying storage layer to avoid retrieving the value.
	// REQUIRES: HasPointAndRange() has returned (true, *).
	MVCCValueLenAndIsTombstone() (int, bool, error)
	// ValueLen can be called for MVCC or non-MVCC values, when only the value
	// length is needed. This is an optimization that can allow the underlying
	// storage layer to avoid retrieving the value.
	// REQUIRES: HasPointAndRange() has returned (true, *).
	ValueLen() int
	// HasPointAndRange returns whether the current iterator position has a point
	// key and/or a range key. Must check Valid() first. At least one of these
	// will always be true for a valid iterator. For details on range keys, see
	// comment on SimpleMVCCIterator.
	HasPointAndRange() (bool, bool)
	// RangeBounds returns the range bounds for the current range key, or an
	// empty span if there are none. The returned keys are valid until the
	// range key changes, see RangeKeyChanged().
	RangeBounds() roachpb.Span
	// RangeKeys returns a stack of all range keys (with different timestamps) at
	// the current key position. When at a point key, it will return all range
	// keys overlapping that point key. The stack is valid until the range key
	// changes, see RangeKeyChanged().
	//
	// For details on range keys, see SimpleMVCCIterator comment, or tech note:
	// https://github.com/cockroachdb/cockroach/blob/master/docs/tech-notes/mvcc-range-tombstones.md
	RangeKeys() MVCCRangeKeyStack
	// RangeKeyChanged returns true if the previous seek or step moved to a
	// different range key (or none at all). Requires a valid iterator, but an
	// exhausted iterator is considered to have had no range keys when calling
	// this after repositioning.
	RangeKeyChanged() bool
}

// MVCCIterator is an interface for iterating over key/value pairs in an
// engine. It is used for iterating over the key space that can have multiple
// versions, and if often also used (due to historical reasons) for iterating
// over the key space that never has multiple versions (i.e.,
// MVCCKey.Timestamp.IsEmpty()).
//
// MVCCIterator implementations are thread safe unless otherwise noted. API
// invariants are asserted via assertMVCCIteratorInvariants().
//
// For details on range keys and iteration, see comment on SimpleMVCCIterator.
type MVCCIterator interface {
	SimpleMVCCIterator
}

// IterOptions contains options used to create an {MVCC,Engine}Iterator.
//
// For performance, every {MVCC,Engine}Iterator must specify either Prefix or
// UpperBound.
type IterOptions struct {
	// If Prefix is true, Seek will use the user-key prefix of the supplied
	// {MVCC,Engine}Key (the Key field) to restrict which sstables are searched,
	// but iteration (using Next) over keys without the same user-key prefix
	// will not work correctly (keys may be skipped).
	Prefix bool
	// LowerBound gives this iterator an inclusive lower bound. Attempts to
	// SeekReverse or Prev to a key that is strictly less than the bound will
	// invalidate the iterator.
	LowerBound roachpb.Key
	// UpperBound gives this iterator an exclusive upper bound. Attempts to Seek
	// or Next to a key that is greater than or equal to the bound will invalidate
	// the iterator. UpperBound must be provided unless Prefix is true, in which
	// case the end of the prefix will be used as the upper bound.
	UpperBound roachpb.Key
	// MinTimestamp and MaxTimestamp, if set, indicate that only keys
	// within the time range formed by [MinTimestamp, MaxTimestamp] should be
	// returned. The underlying iterator may be able to efficiently skip over
	// keys outside of the hinted time range, e.g., when a block handle
	// indicates that the block contains no keys within the time range. Intents
	// will not be visible to such iterators at all. This is only relevant for
	// MVCCIterators.
	//
	// Note that time-bound iterators previously were only a performance
	// optimization but now guarantee that no keys outside of the [start, end]
	// time range will be returned.
	//
	// NB: Range keys are not currently subject to timestamp filtering due to
	// complications with MVCCIncrementalIterator. See:
	// https://github.com/cockroachdb/cockroach/issues/86260
	MinTimestamp, MaxTimestamp hlc.Timestamp
	// KeyTypes specifies the types of keys to surface: point and/or range keys.
	// Use HasPointAndRange() to determine which key type is present at a given
	// iterator position, and RangeBounds() and RangeKeys() to access range keys.
	// Defaults to IterKeyTypePointsOnly. For more details on range keys, see
	// comment on SimpleMVCCIterator.
	KeyTypes IterKeyType
	// RangeKeyMaskingBelow enables masking (hiding) of point keys by range keys.
	// Any range key with a timestamp at or below RangeKeyMaskingBelow
	// will mask point keys below it, preventing them from being surfaced.
	// Consider the following example:
	//
	// 4          o---------------o    RangeKeyMaskingBelow=4 emits b3
	// 3      b3      d3               RangeKeyMaskingBelow=3 emits b3,d3,f2
	// 2  o---------------o   f2       RangeKeyMaskingBelow=2 emits b3,d3,f2
	// 1  a1  b1          o-------o    RangeKeyMaskingBelow=1 emits a1,b3,b1,d3,f2
	//    a   b   c   d   e   f   g
	//
	// Range keys themselves are not affected by the masking, and will be
	// emitted as normal.
	RangeKeyMaskingBelow hlc.Timestamp
	// ReadCategory is used to map to a user-understandable category string, for
	// stats aggregation and metrics, and a Pebble-understandable QoS.
	ReadCategory ReadCategory
	// useL6Filters allows the caller to opt into reading filter blocks for
	// L6 sstables. Only for use with Prefix = true. Helpful if a lot of prefix
	// Seeks are expected in quick succession, that are also likely to not
	// yield a single key. Filter blocks in L6 can be relatively large, often
	// larger than data blocks, so the benefit of loading them in the cache
	// is minimized if the probability of the key existing is not low or if
	// this is a one-time Seek (where loading the data block directly is better).
	useL6Filters bool
}
