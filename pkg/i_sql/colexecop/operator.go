package colexecop

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/i_sql/execinfra/execopnode"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
	"time"
)

// KVReader is an operator that performs KV reads.
// TODO(yuzefovich): consider changing the contract to remove the mention of
// concurrency safety once stats are only retrieved from Next goroutines.
type KVReader interface {
	// GetBytesRead returns the number of bytes read from KV by this operator.
	// It must be safe for concurrent use.
	GetBytesRead() int64
	// GetKVPairsRead returns the number of key-values pairs read from KV by
	// this operator. It must be safe for concurrent use.
	GetKVPairsRead() int64
	// GetRowsRead returns the number of rows read from KV by this operator.
	// It must be safe for concurrent use.
	GetRowsRead() int64
	// GetBatchRequestsIssued returns the number of BatchRequests issued to KV
	// by this operator. It must be safe for concurrent use.
	GetBatchRequestsIssued() int64
	// GetContentionTime returns the amount of time KV reads spent
	// contending. It must be safe for concurrent use.
	GetContentionTime() time.Duration
	// GetConsumedRU returns the number of RUs that were consumed during the
	// KV reads.
	GetConsumedRU() uint64
	// GetKVCPUTime returns the CPU time consumed *on the current goroutine* by
	// KV requests. It must be safe for concurrent use. It is used to calculate
	// the SQL CPU time.
	GetKVCPUTime() time.Duration
	// UsedStreamer returns whether the Streamer API was used by the KVReader.
	UsedStreamer() bool
}

// ClosableOperator is an Operator that needs to be Close()'d.
// NOTE: even if the Operator wasn't Init()'ed properly, it must still be safe
// to Close().
type ClosableOperator interface {
	Operator
	Closer
}

// Operator is a column vector operator that produces a Batch as output.
type Operator interface {
	// Init initializes this operator. It will be called once at operator setup
	// time. Second, third, etc calls should be noops. If an operator has any
	// input operators, it's responsible for calling Init on all of those input
	// operators as well.
	//
	// Canceling the provided context results in forceful termination of
	// execution. The operators are expected to hold onto the provided context
	// (and derive a new one if needed) that is then used for Next(),
	// DrainMeta(), and Close() calls (when applicable).
	//
	// It might panic with an expected error, so there must be a "root"
	// component that will catch that panic.
	Init(ctx context.Context)

	// Next returns the next Batch from this operator. Once the operator is
	// finished, it will return a Batch with length 0. Subsequent calls to
	// Next at that point will always return a Batch with length 0.
	//
	// Calling Next may invalidate the contents of the last Batch returned by
	// Next.
	//
	// It might panic with an expected error, so there must be a "root"
	// component that will catch that panic.
	Next() coldata.Batch

	execopnode.OpNode
}

// Closer is an object that releases resources when Close is called. Note that
// this interface must be implemented by all operators that could be planned on
// top of other operators that do actually need to release the resources (e.g.
// if we have a simple project on top of a disk-backed operator, that simple
// project needs to implement this interface so that Close() call could be
// propagated correctly).
type Closer interface {
	// Close releases the resources associated with this Closer. If this Closer
	// is an Operator, the implementation of Close must be safe to execute even
	// if Operator.Init wasn't called. Multiple calls to Close() are allowed,
	// and most of the implementations should make all calls except for the
	// first one no-ops.
	//
	// Unless the Closer derives its own context with a separate tracing span,
	// the argument context rather than the one from Init() must be used
	// (wherever necessary) by the implementation. This is so since the span in
	// the context from Init() might be already finished when Close() is called
	// whereas the argument context will contain an unfinished span.
	Close(context.Context) error
}
