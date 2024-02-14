package storage

import (
	"context"
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
)

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
	NewMVCCIterator(ctx context.Context, iterKind MVCCIterKind, opts IterOptions) (MVCCIterator, error)
}

// Writer is the write interface to an engine's data.
type Writer interface {
	// PutMVCC sets the given key to the value provided. It requires that the
	// timestamp is non-empty (see {PutUnversioned,PutIntent} if the timestamp
	// is empty).
	//
	// It is safe to modify the contents of the arguments after PutMVCC returns.
	PutMVCC(key MVCCKey, value MVCCValue) error
	// BufferedSize returns the size of the underlying buffered writes if the
	// Writer implementation is buffered, and 0 if the Writer implementation is
	// not buffered. Buffered writers are expected to always give a monotonically
	// increasing size.
	BufferedSize() int
}

// MVCCIterKind is used to inform Reader about the kind of iteration desired
// by the caller.
type MVCCIterKind int

// "Intent" refers to non-inline meta, that can be interleaved or separated.
const (
	// MVCCKeyAndIntentsIterKind specifies that intents must be seen, and appear
	// interleaved with keys, even if they are in a separated lock table.
	// Iterators of this kind are not allowed to span from local to global keys,
	// since the physical layout has the separated lock table in-between the
	// local and global keys. These iterators do strict error checking and panic
	// if the caller seems that to be trying to violate this constraint.
	// Specifically:
	// - If both bounds are set they must not span from local to global.
	// - Any bound (lower or upper), constrains the iterator for its lifetime to
	//   one of local or global keys. The iterator will not tolerate a seek that
	//   violates this constraint.
	// We could, with significant code complexity, not constrain an iterator for
	// its lifetime, and allow a seek that specifies a global (local) key to
	// change the constraint to global (local). This would allow reuse of the
	// same iterator with a large global upper-bound. But a Next call on the
	// highest local key (Prev on the lowest global key) would still not be able
	// to transparently skip over the intermediate lock table. We deem that
	// behavior to be more surprising and bug-prone (for the caller), than being
	// strict.
	MVCCKeyAndIntentsIterKind MVCCIterKind = iota
	// MVCCKeyIterKind specifies that the caller does not need to see intents.
	// Any interleaved intents may be seen, but no correctness properties are
	// derivable from such partial knowledge of intents. NB: this is a performance
	// optimization when iterating over (a) MVCC keys where the caller does
	// not need to see intents, (b) a key space that is known to not have multiple
	// versions (and therefore will never have intents), like the raft log.
	MVCCKeyIterKind
)

// IterKeyType configures which types of keys an iterator should surface.
type IterKeyType int8

const (
	// IterKeyTypePointsOnly configures an iterator to iterate over point keys
	// only.
	IterKeyTypePointsOnly IterKeyType = iota
	// IterKeyTypeRangesOnly configures an iterator to iterate over range keys
	// only.
	IterKeyTypeRangesOnly
	// IterKeyTypePointsAndRanges configures an iterator iterate over both point
	// keys and range keys simultaneously.
	IterKeyTypePointsAndRanges
)

// ReadCategory is used to export metrics and maps to a QoS understood by
// Pebble. Categories are being introduced lazily, since more categories
// result in more metrics.
type ReadCategory int8

const (
	// UnknownReadCategory are requests that are not categorized. If the metric
	// for this category becomes a high fraction of reads, we will need to
	// investigate and break out more categories.
	UnknownReadCategory ReadCategory = iota
	// BatchEvalReadCategory includes evaluation of most BatchRequests. It
	// excludes scans and reverse scans. If scans and reverse scans are mixed
	// with other requests in a batch, we may currently assign the category
	// based on the first request.
	BatchEvalReadCategory
	// ScanRegularBatchEvalReadCategory are BatchRequest (reverse) scans that
	// have admission priority NormalPri or higher.
	ScanRegularBatchEvalReadCategory
	// ScanBackgroundBatchEvalReadCategory are BatchRequest (reverse) scans that
	// have admission priority lower than NormalPri. This includes backfill
	// scans for changefeeds (see changefeedccl/kvfeed/scanner.go, which sends
	// ScanRequests).
	ScanBackgroundBatchEvalReadCategory
	// MVCCGCReadCategory are reads for MVCC GC.
	MVCCGCReadCategory
	// RangeSnapshotReadCategory are reads for sending range snapshots.
	RangeSnapshotReadCategory
	// RangefeedReadCategory are reads for rangefeeds, including catchup scans.
	RangefeedReadCategory
	// ReplicationReadCategory are reads related to Raft replication.
	ReplicationReadCategory
	// IntentResolutionReadCategory are reads for intent resolution.
	IntentResolutionReadCategory
	// BackupReadCategory are reads for backups.
	BackupReadCategory
)
