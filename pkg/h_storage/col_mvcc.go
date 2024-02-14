package storage

import (
	"context"
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
)

// NextKVer can fetch a new KV from somewhere. If MVCCDecodingStrategy is set
// to MVCCDecodingRequired, the returned KV will include a timestamp.
type NextKVer interface {
	// Init initializes the NextKVer. It returns a boolean indicating whether
	// the KVs returned by NextKV are stable (i.e. whether they will not be
	// invalidated by calling NextKV again).
	Init(getter FirstKeyOfRowGetter) (stableKVs bool)
	// NextKV returns the next kv from this NextKVer.
	// - ok=false indicates that there are no more kvs to fetch,
	// - partialRow indicates whether the fetch stopped in the middle of a SQL
	// row (in this case ok will be set to false),
	// - the kv that was fetched,
	// - any errors that may have occurred.
	//
	// When (ok=false,partialRow=true) is returned, the caller is expected to
	// discard all KVs that were part of the last SQL row that was incomplete.
	// The scan will be resumed from the key provided by the FirstKeyOfRowGetter
	// (provided in Init by the caller) obtained during this NextKV call.
	NextKV(context.Context, MVCCDecodingStrategy) (ok bool, partialRow bool, kv roachpb.KeyValue, err error)
}

// This file defines several interfaces as well as introduces a couple of
// components that power the direct columnar scans. The main idea of this
// feature is to use the injected decoding logic from SQL in order to process
// each KV and keep only the needed parts (i.e. necessary SQL columns). Those
// needed parts are then propagated back to the KV client as coldata.Batch'es
// (serialized in the Apache Arrow format).
//
// Here is an example outline of all components involved:
//
//      ┌────────────────────────────────────────────────┐
//      │                       SQL                      │
//      │________________________________________________│
//      │          colfetcher.ColBatchDirectScan         │
//      │                        │                       │
//      │                        ▼                       │
//      │                 row.txnKVFetcher               │
//      │    (behind the row.KVBatchFetcher interface)   │
//      └────────────────────────────────────────────────┘
//                               │
//                               ▼
//      ┌────────────────────────────────────────────────┐
//      │                    KV Client                   │
//      └────────────────────────────────────────────────┘
//                               │
//                               ▼
//      ┌────────────────────────────────────────────────┐
//      │                    KV Server                   │
//      │________________________________________________│
//      │           colfetcher.cFetcherWrapper           │
//      │ (behind the storage.CFetcherWrapper interface) │
//      │                        │                       │
//      │                        ▼                       │
//      │              colfetcher.cFetcher               │
//      │                        │                       │
//      │                        ▼                       │
//      │          storage.mvccScanFetchAdapter ────────┐│
//      │    (behind the storage.NextKVer interface)    ││
//      │                        │                      ││
//      │                        ▼                      ││
//      │           storage.pebbleMVCCScanner           ││
//      │ (which put's KVs into storage.singleResults) <┘│
//      └────────────────────────────────────────────────┘
//
// On the KV client side, row.txnKVFetcher issues Scans and ReverseScans with
// the COL_BATCH_RESPONSE format and returns the response (which contains the
// columnar data) to the colfetcher.ColBatchDirectScan.
//
// On the KV server side, we create a storage.CFetcherWrapper that asks the
// colfetcher.cFetcher for the next coldata.Batch. The cFetcher, in turn,
// fetches the next KV, decodes it, and keeps only values for the needed SQL
// columns, discarding the rest of the KV. The KV is emitted by the
// mvccScanFetchAdapter which - via the singleResults struct - exposes access to
// the current KV that the pebbleMVCCScanner is pointing at.
//
// Note that there is an additional "implicit synchronization" between
// components that is not shown on this diagram. In particular,
// storage.singleResults.maybeTrimPartialLastRow must be in sync with the
// colfetcher.cFetcher which is achieved by
// - the cFetcher exposing access to the first key of the last incomplete SQL
//   row via the FirstKeyOfRowGetter,
// - the singleResults using that key as the resume key for the response,
// - and the cFetcher removing that last partial SQL row when NextKV() returns
//   partialRow=true.
// This "upstream" link (although breaking the layering a bit) allows us to
// avoid a performance penalty for handling the case with multiple column
// families. (This case is handled by the storage.pebbleResults via tracking
// offsets into the pebbleResults.repr.)
//
// This code structure deserves some elaboration. First, there is a mismatch
// between the "push" mode in which the pebbleMVCCScanner operates and the
// "pull" mode that the NextKVer exposes. The adaption between two different
// modes is achieved via the mvccScanFetcherAdapter grabbing (when the control
// returns to it) the current unstable KV pair from the singleResults struct
// which serves as a one KV pair buffer that the pebbleMVCCScanner `put`s into.
// Second, in order be able to use the unstable KV pair without performing a
// copy, the pebbleMVCCScanner stops at the current KV pair and returns the
// control flow (which is exactly what pebbleMVCCScanner.getOne does) back to
// the mvccScanFetcherAdapter, with the adapter advancing the scanner only when
// the next KV pair is needed.

// FirstKeyOfRowGetter returns the first key included into the last incomplete
// SQL row by the user of NextKVer. If the last row is complete, then nil is
// returned.
type FirstKeyOfRowGetter func() roachpb.Key

// MVCCDecodingStrategy controls if and how the fetcher should decode MVCC
// timestamps from returned KV's.
type MVCCDecodingStrategy int

const (
	// MVCCDecodingNotRequired is used when timestamps aren't needed.
	MVCCDecodingNotRequired MVCCDecodingStrategy = iota
	// MVCCDecodingRequired is used when timestamps are needed.
	MVCCDecodingRequired
)
