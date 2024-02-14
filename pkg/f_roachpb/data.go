package roachpb

import (
	"github.com/dborchard/tiny_crdb/pkg/f_kv/kvserver/concurrency/isolation"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// Key is a custom type for a byte string in proto
// messages which refer to Cockroach keys.
type Key []byte

type Value struct {
	RawBytes  []byte
	Timestamp hlc.Timestamp
}

type Span struct {
}

type KeyValue struct {
	Key   []byte
	Value Value
}

// MakeTransaction creates a new transaction. The transaction key is
// composed using the specified baseKey (for locality with data
// affected by the transaction) and a random ID to guarantee
// uniqueness. The specified user-level priority is combined with a
// randomly chosen value to yield a final priority, used to settle
// write conflicts in a way that avoids starvation of long-running
// transactions (see Replica.PushTxn).
//
// coordinatorNodeID is provided to track the SQL (or possibly KV) node
// that created this transaction, in order to be used (as
// of this writing) to enable observability on contention events
// between different transactions.
//
// baseKey can be nil, in which case it will be set when sending the first
// write.
//
// omitInRangefeeds controls whether the transaction's writes are exposed via
// rangefeeds. When set to true, all the transaction's writes will be
// filtered out by rangefeeds, and will not be available in changefeeds.
func MakeTransaction(
	name string,
	baseKey Key,
	isoLevel isolation.Level,
	userPriority UserPriority,
	now hlc.Timestamp,
) Transaction {
	return Transaction{}
}
