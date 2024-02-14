package storage

import "context"

// Engine is the interface that wraps the core operations of a key/value store.
type Engine interface {
	Reader
	Writer

	// Compact forces compaction over the entire database.
	Compact() error
	// Flush causes the engine to write all in-memory data to disk
	// immediately.
	Flush() error
}

// Reader is the read interface to an engine's data. Certain implementations
// of Reader guarantee consistency of the underlying engine state across the
// different iterators created by NewMVCCIterator, NewEngineIterator:
//   - pebbleSnapshot, because it uses an engine snapshot.
//   - pebbleReadOnly, pebbleBatch: when the IterOptions do not specify a
//     timestamp hint (see IterOptions). Note that currently the engine state
//     visible here is not as of the time of the Reader creation. It is the time
//     when the first iterator is created, or earlier if
//     PinEngineStateForIterators is called.
//
// The ConsistentIterators method returns true when this consistency is
// guaranteed by the Reader.
// TODO(sumeer): this partial consistency can be a source of bugs if future
// code starts relying on it, but rarely uses a Reader that does not guarantee
// it. Can we enumerate the current cases where KV uses Engine as a Reader?
type Reader interface {
	// Close closes the reader, freeing up any outstanding resources. Note that
	// various implementations have slightly different behaviors. In particular,
	// Distinct() batches release their parent batch for future use while
	// Engines, Snapshots and Batches free the associated C++ resources.
	Close()
	// Closed returns true if the reader has been closed or is not usable.
	// Objects backed by this reader (e.g. Iterators) can check this to ensure
	// that they are not using a closed engine. Intended for use within package
	// engine; exported to enable wrappers to exist in other packages.
	Closed() bool
	// MVCCIterate scans from the start key to the end key (exclusive), invoking
	// the function f on each key value pair. The inputs are copies, and safe to
	// retain beyond the function call. It supports interleaved iteration over
	// point and/or range keys, providing any overlapping range keys for each
	// point key if requested. If f returns an error or if the scan itself
	// encounters an error, the iteration will stop and return the error.
	//
	// Note that this method is not expected take into account the timestamp of
	// the end key; all MVCCKeys at end.Key are considered excluded in the
	// iteration.
	MVCCIterate(
		ctx context.Context, start, end roachpb.Key, iterKind MVCCIterKind, keyTypes IterKeyType,
		readCategory ReadCategory, f func(MVCCKeyValue, MVCCRangeKeyStack) error,
	) error
	// NewMVCCIterator returns a new instance of an MVCCIterator over this engine.
	// The caller must invoke Close() on it when done to free resources.
	//
	// Write visibility semantics:
	//
	// 1. An iterator has a consistent view of the reader as of the time of its
	//    creation. Subsequent writes are never visible to it.
	//
	// 2. All iterators on readers with ConsistentIterators=true have a consistent
	//    view of the _engine_ (not reader) as of the time of the first iterator
	//    creation or PinEngineStateForIterators call: newer engine writes are
	//    never visible. The opposite holds for ConsistentIterators=false: new
	//    iterators see the most recent engine state at the time of their creation.
	//
	// 3. Iterators on unindexed batches never see batch writes, but satisfy
	//    ConsistentIterators for engine write visibility.
	//
	// 4. Iterators on indexed batches see all batch writes as of their creation
	//    time, but they satisfy ConsistentIterators for engine writes.
	NewMVCCIterator(
		ctx context.Context, iterKind MVCCIterKind, opts IterOptions) (MVCCIterator, error)
	// NewEngineIterator returns a new instance of an EngineIterator over this
	// engine. The caller must invoke EngineIterator.Close() when finished
	// with the iterator to free resources. The caller can change IterOptions
	// after this function returns.
	NewEngineIterator(ctx context.Context, opts IterOptions) (EngineIterator, error)
	// ScanInternal allows a caller to inspect the underlying engine's InternalKeys
	// using a visitor pattern, while also allowing for keys in shared files to be
	// skipped if a visitor is provided for visitSharedFiles. Useful for
	// fast-replicating state from one Reader to another. Point keys are collapsed
	// such that only one internal key per user key is exposed, and rangedels and
	// range keys are collapsed and defragmented with each span being surfaced
	// exactly once, alongside the highest seqnum for a rangedel on that span
	// (for rangedels) or all coalesced rangekey.Keys in that span (for range
	// keys). A point key deleted by a rangedel will not be exposed, but the
	// rangedel would be exposed.
	//
	// Note that ScanInternal does not obey the guarantees indicated by
	// ConsistentIterators.
	ScanInternal(
		ctx context.Context, lower, upper roachpb.Key,
		visitPointKey func(key *pebble.InternalKey, value pebble.LazyValue, info pebble.IteratorLevel) error,
		visitRangeDel func(start, end []byte, seqNum uint64) error,
		visitRangeKey func(start, end []byte, keys []rangekey.Key) error,
		visitSharedFile func(sst *pebble.SharedSSTMeta) error,
	) error
	// ConsistentIterators returns true if the Reader implementation guarantees
	// that the different iterators constructed by this Reader will see the same
	// underlying Engine state. This is not true about Batch writes: new iterators
	// will see new writes made to the batch, existing iterators won't.
	ConsistentIterators() bool

	// PinEngineStateForIterators ensures that the state seen by iterators
	// without timestamp hints (see IterOptions) is pinned and will not see
	// future mutations. It can be called multiple times on a Reader in which
	// case the state seen will be either:
	// - As of the first call.
	// - For a Reader returned by Engine.NewSnapshot, the pinned state is as of
	//   the time the snapshot was taken.
	// So the semantics that are true for all Readers is that the pinned state
	// is somewhere in the time interval between the creation of the Reader and
	// the first call to PinEngineStateForIterators.
	// REQUIRES: ConsistentIterators returns true.
	PinEngineStateForIterators(readCategory ReadCategory) error
}

// Writer is the write interface to an engine's data.
type Writer interface {
	// ApplyBatchRepr atomically applies a set of batched updates. Created by
	// calling Repr() on a batch. Using this method is equivalent to constructing
	// and committing a batch whose Repr() equals repr. If sync is true, the batch
	// is synchronously flushed to the OS and written to disk. It is an error to
	// specify sync=true if the Writer is a Batch.
	//
	// It is safe to modify the contents of the arguments after ApplyBatchRepr
	// returns.
	ApplyBatchRepr(repr []byte, sync bool) error

	// ClearMVCC removes the point key with the given MVCCKey from the db. It does
	// not affect range keys. It requires that the timestamp is non-empty (see
	// ClearUnversioned or ClearIntent if the timestamp is empty). Note that clear
	// actually removes entries from the storage engine, rather than inserting
	// MVCC tombstones.
	//
	// If the caller knows the size of the value that is being cleared, they
	// should set ClearOptions.{ValueSizeKnown, ValueSize} accordingly to
	// improve the storage engine's ability to prioritize compactions.
	//
	// It is safe to modify the contents of the arguments after it returns.
	ClearMVCC(key MVCCKey, opts ClearOptions) error
	// ClearUnversioned removes an unversioned item from the db. It is for use
	// with inline metadata (not intents) and other unversioned keys (like
	// Range-ID local keys). It does not affect range keys.
	//
	// If the caller knows the size of the value that is being cleared, they
	// should set ClearOptions.{ValueSizeKnown, ValueSize} accordingly to
	// improve the storage engine's ability to prioritize compactions.
	//
	// It is safe to modify the contents of the arguments after it returns.
	ClearUnversioned(key roachpb.Key, opts ClearOptions) error
	// ClearEngineKey removes the given point key from the engine. It does not
	// affect range keys.  Note that clear actually removes entries from the
	// storage engine. This is a general-purpose and low-level method that should
	// be used sparingly, only when the other Clear* methods are not applicable.
	//
	// If the caller knows the size of the value that is being cleared, they
	// should set ClearOptions.{ValueSizeKnown, ValueSize} accordingly to
	// improve the storage engine's ability to prioritize compactions.
	//
	// It is safe to modify the contents of the arguments after it returns.
	ClearEngineKey(key EngineKey, opts ClearOptions) error

	// ClearRawRange removes point and/or range keys from start (inclusive) to end
	// (exclusive) using Pebble range tombstones. It can be applied to a range
	// consisting of MVCCKeys or the more general EngineKeys -- it simply uses the
	// roachpb.Key parameters as the Key field of an EngineKey. This implies that
	// it does not clear intents unless the intent lock table is targeted
	// explicitly.
	//
	// Similar to the other Clear* methods, this method actually removes entries
	// from the storage engine. It is safe to modify the contents of the arguments
	// after it returns.
	ClearRawRange(start, end roachpb.Key, pointKeys, rangeKeys bool) error
	// ClearMVCCRange removes MVCC point and/or range keys (including intents)
	// from start (inclusive) to end (exclusive) using Pebble range tombstones.
	//
	// Similar to the other Clear* methods, this method actually removes entries
	// from the storage engine. It is safe to modify the contents of the arguments
	// after it returns.
	ClearMVCCRange(start, end roachpb.Key, pointKeys, rangeKeys bool) error
	// ClearMVCCVersions removes MVCC point key versions from start (inclusive) to
	// end (exclusive) using a Pebble range tombstone. It is meant for efficiently
	// clearing a subset of versions of a key, since the parameters are MVCCKeys
	// and not roachpb.Keys, but it can also be used across multiple keys. It will
	// ignore intents and range keys, leaving them in place.
	//
	// Similar to the other Clear* methods, this method actually removes entries
	// from the storage engine. It is safe to modify the contents of the arguments
	// after it returns.
	ClearMVCCVersions(start, end MVCCKey) error
	// ClearMVCCIteratorRange removes all point and/or range keys in the given
	// span using an MVCC iterator, by clearing individual keys (including
	// intents).
	//
	// Similar to the other Clear* methods, this method actually removes entries
	// from the storage engine. It is safe to modify the contents of the arguments
	// after it returns.
	//
	// TODO(erikgrinaker): This should be a separate function rather than an
	// interface method, but we keep it for now to make use of UnsafeRawKey() when
	// clearing keys.
	ClearMVCCIteratorRange(start, end roachpb.Key, pointKeys, rangeKeys bool) error

	// ClearMVCCRangeKey deletes an MVCC range key from start (inclusive) to end
	// (exclusive) at the given timestamp. For any range key that straddles the
	// start and end boundaries, only the segments within the boundaries will be
	// cleared. Range keys at other timestamps are unaffected.  Clears are
	// idempotent.
	//
	// This method is primarily intended for MVCC garbage collection and similar
	// internal use.
	ClearMVCCRangeKey(rangeKey MVCCRangeKey) error

	// PutMVCCRangeKey writes an MVCC range key. It will replace any overlapping
	// range keys at the given timestamp (even partial overlap). Only MVCC range
	// tombstones, i.e. an empty value, are currently allowed (other kinds will
	// need additional handling in MVCC APIs and elsewhere, e.g. stats and GC).
	//
	// Range keys must be accessed using special iterator options and methods,
	// see SimpleMVCCIterator.RangeKeys() for details.
	//
	// For more information on MVCC range keys, see this tech note:
	// https://github.com/cockroachdb/cockroach/blob/master/docs/tech-notes/mvcc-range-tombstones.md
	PutMVCCRangeKey(MVCCRangeKey, MVCCValue) error

	// PutRawMVCCRangeKey is like PutMVCCRangeKey, but accepts an encoded
	// MVCCValue. It can be used to avoid decoding and immediately re-encoding an
	// MVCCValue, but should generally be avoided due to the lack of type safety.
	//
	// It is safe to modify the contents of the arguments after PutRawMVCCRangeKey
	// returns.
	PutRawMVCCRangeKey(MVCCRangeKey, []byte) error

	// PutEngineRangeKey sets the given range key to the values provided. This is
	// a general-purpose and low-level method that should be used sparingly, only
	// when the other Put* methods are not applicable.
	//
	// It is safe to modify the contents of the arguments after it returns.
	PutEngineRangeKey(start, end roachpb.Key, suffix, value []byte) error

	// ClearEngineRangeKey clears the given range key. This is a general-purpose
	// and low-level method that should be used sparingly, only when the other
	// Clear* methods are not applicable.
	//
	// It is safe to modify the contents of the arguments after it returns.
	ClearEngineRangeKey(start, end roachpb.Key, suffix []byte) error

	// Merge is a high-performance write operation used for values which are
	// accumulated over several writes. Multiple values can be merged
	// sequentially into a single key; a subsequent read will return a "merged"
	// value which is computed from the original merged values. We only
	// support Merge for keys with no version.
	//
	// Merge currently provides specialized behavior for three data types:
	// integers, byte slices, and time series observations. Merged integers are
	// summed, acting as a high-performance accumulator.  Byte slices are simply
	// concatenated in the order they are merged. Time series observations
	// (stored as byte slices with a special tag on the roachpb.Value) are
	// combined with specialized logic beyond that of simple byte slices.
	//
	//
	// It is safe to modify the contents of the arguments after Merge returns.
	Merge(key MVCCKey, value []byte) error

	// PutMVCC sets the given key to the value provided. It requires that the
	// timestamp is non-empty (see {PutUnversioned,PutIntent} if the timestamp
	// is empty).
	//
	// It is safe to modify the contents of the arguments after PutMVCC returns.
	PutMVCC(key MVCCKey, value MVCCValue) error
	// PutRawMVCC is like PutMVCC, but it accepts an encoded MVCCValue. It
	// can be used to avoid decoding and immediately re-encoding an MVCCValue,
	// but should generally be avoided due to the lack of type safety.
	//
	// It is safe to modify the contents of the arguments after PutRawMVCC
	// returns.
	PutRawMVCC(key MVCCKey, value []byte) error
	// PutUnversioned sets the given key to the value provided. It is for use
	// with inline metadata (not intents) and other unversioned keys (like
	// Range-ID local keys).
	//
	// It is safe to modify the contents of the arguments after Put returns.
	PutUnversioned(key roachpb.Key, value []byte) error
	// PutEngineKey sets the given key to the value provided. This is a
	// general-purpose and low-level method that should be used sparingly,
	// only when the other Put* methods are not applicable.
	//
	// It is safe to modify the contents of the arguments after Put returns.
	PutEngineKey(key EngineKey, value []byte) error

	// LogData adds the specified data to the RocksDB WAL. The data is
	// uninterpreted by RocksDB (i.e. not added to the memtable or sstables).
	//
	// It is safe to modify the contents of the arguments after LogData returns.
	LogData(data []byte) error
	// LogLogicalOp logs the specified logical mvcc operation with the provided
	// details to the writer, if it has logical op logging enabled. For most
	// Writer implementations, this is a no-op.
	LogLogicalOp(op MVCCLogicalOpType, details MVCCLogicalOpDetails)

	// SingleClearEngineKey removes the most recent write to the item from the db
	// with the given key. Whether older writes of the item will come back
	// to life if not also removed with SingleClear is undefined. See the
	// following:
	//   https://github.com/facebook/rocksdb/wiki/Single-Delete
	// for details on the SingleDelete operation that this method invokes. Note
	// that clear actually removes entries from the storage engine, rather than
	// inserting MVCC tombstones. This is a low-level interface that must not be
	// called from outside the storage package. It is part of the interface
	// because there are structs that wrap Writer and implement the Writer
	// interface, that are not part of the storage package.
	//
	// It is safe to modify the contents of the arguments after it returns.
	SingleClearEngineKey(key EngineKey) error

	// ShouldWriteLocalTimestamps is only for internal use in the storage package.
	// This method is temporary, to handle the transition from clusters where not
	// all nodes understand local timestamps.
	ShouldWriteLocalTimestamps(ctx context.Context) bool

	// BufferedSize returns the size of the underlying buffered writes if the
	// Writer implementation is buffered, and 0 if the Writer implementation is
	// not buffered. Buffered writers are expected to always give a monotonically
	// increasing size.
	BufferedSize() int
}
