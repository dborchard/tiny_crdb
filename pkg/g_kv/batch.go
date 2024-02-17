package kv

import (
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvpb"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
)

// Batch provides for the parallel execution of a number of database
// operations. Operations are added to the Batch and then the Batch is executed
// via either DB.Run, Txn.Run or Txn.Commit.
//
// TODO(pmattis): Allow a timestamp to be specified which is applied to all
// operations within the batch.
type Batch struct {
	// The Txn the batch is associated with. This field may be nil if the batch
	// was not created via Txn.NewBatch.
	txn *Txn
	// Results contains an entry for each operation added to the batch. The order
	// of the results matches the order the operations were added to the
	// batch. For example:
	//
	//   b := db.NewBatch()
	//   b.Put("a", "1")
	//   b.Put("b", "2")
	//   _ = db.Run(b)
	//   // string(b.Results[0].Rows[0].Key) == "a"
	//   // string(b.Results[1].Rows[0].Key) == "b"
	Results []Result

	// approxMutationReqBytes tracks the approximate size of keys and values in
	// mutations added to this batch via Put, CPut, InitPut, Del, etc.
	approxMutationReqBytes int
	// Set when AddRawRequest is used, in which case using the "other"
	// operations renders the batch unusable.
	raw bool
	// Once received, the response from a successful batch.
	response *kvpb.BatchResponse
	// Once received, any error encountered sending the batch.
	pErr *kvpb.Error

	// We use pre-allocated buffers to avoid dynamic allocations for small batches.
	resultsBuf    [8]Result
	rowsBuf       []KeyValue
	rowsStaticBuf [8]KeyValue
	rowsStaticIdx int
	reqs          []kvpb.RequestUnion
}

// Put sets the value for a key.
//
// A new result will be appended to the batch which will contain a single row
// and Result.Err will indicate success or failure.
//
// key can be either a byte slice or a string. value can be any key type, a
// protoutil.Message or any Go primitive type (bool, int, etc).
func (b *Batch) Put(key, value interface{}) {
	if value == nil {
		// Empty values are used as deletion tombstones, so one can't write an empty
		// value. If the intention was indeed to delete the key, use Del() instead.
		panic("can't Put an empty Value; did you mean to Del() instead?")
	}
	//b.put(key, value, false)
}

// Result holds the result for a single DB or Txn operation (e.g. Get, Put,
// etc).
type Result struct {
	calls int
	// Err contains any error encountered when performing the operation.
	Err error
	// Rows contains the key/value pairs for the operation. The number of rows
	// returned varies by operation. For Get, Put, CPut, and Inc the number
	// of rows returned is the number of keys operated on. For Scan the number of
	// rows returned is the number or rows matching the scan capped by the
	// maxRows parameter and other options. For Del and DelRange Rows is nil.
	Rows []KeyValue

	// Keys is set by Del and DelRange instead of returning the rows themselves.
	Keys []roachpb.Key
}

// KeyValue represents a single key/value pair. This is similar to
// roachpb.KeyValue except that the value may be nil. The timestamp
// in the value will be populated with the MVCC timestamp at which this
// value was read if this struct was produced by a GetRequest or
// ScanRequest which uses the KEY_VALUES ScanFormat. Values created from
// a ScanRequest which uses the BATCH_RESPONSE ScanFormat will contain a
// zero Timestamp.
type KeyValue struct {
	Key   roachpb.Key
	Value *roachpb.Value
}
