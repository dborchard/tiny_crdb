package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/flowinfra"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
)

// rowResultWriter is a subset of CommandResult to be used with the
// DistSQLReceiver. It's implemented by RowResultWriter.
type rowResultWriter interface {
	// AddRow writes a result row.
	// Note that the caller owns the row slice and might reuse it.
	AddRow(ctx context.Context, row tree.Datums) error
	SetRowsAffected(ctx context.Context, n int)
	SetError(error)
	Err() error
}

// CallbackResultWriter is a rowResultWriter that runs a callback function
// on AddRow.
type CallbackResultWriter struct {
	fn           func(ctx context.Context, row tree.Datums) error
	rowsAffected int
	err          error
}

func (c *CallbackResultWriter) AddRow(ctx context.Context, row tree.Datums) error {
	//TODO implement me
	panic("implement me")
}

func (c *CallbackResultWriter) SetRowsAffected(ctx context.Context, n int) {
	//TODO implement me
	panic("implement me")
}

func (c *CallbackResultWriter) SetError(err error) {
	//TODO implement me
	panic("implement me")
}

func (c *CallbackResultWriter) Err() error {
	//TODO implement me
	panic("implement me")
}

// NewCallbackResultWriter creates a new CallbackResultWriter.
func NewCallbackResultWriter(
	fn func(ctx context.Context, row tree.Datums) error,
) *CallbackResultWriter {
	return &CallbackResultWriter{fn: fn}
}

// batchResultWriter is a subset of CommandResult to be used with the
// DistSQLReceiver when the consumer can operate on columnar batches directly.
type batchResultWriter interface {
	AddBatch(context.Context, coldata.Batch) error
}

// MakeDistSQLReceiver creates a DistSQLReceiver.
//
// ctx is the Context that the receiver will use throughout its
// lifetime. resultWriter is the container where the results will be
// stored. If only the row count is needed, this can be nil.
//
// txn is the transaction in which the producer flow runs; it will be updated
// on errors. Nil if the flow overall doesn't run in a transaction.
func MakeDistSQLReceiver(
	ctx context.Context,
	resultWriter rowResultWriter,
	stmtType tree.StatementReturnType,
	txn *kv.Txn,
) *DistSQLReceiver {
	var batchWriter batchResultWriter
	if commandResult, ok := resultWriter.(RestrictedCommandResult); ok {
		if commandResult.SupportsAddBatch() {
			batchWriter = commandResult
		}
	}
	r := &DistSQLReceiver{
		ctx:          context.Background(),
		resultWriter: resultWriter,
		batchWriter:  batchWriter,
		// At the time of writing, there is only one concurrent goroutine that
		// might send at most one error.
		concurrentErrorCh: make(chan error, 1),
		txn:               txn,
		stmtType:          stmtType,
	}
	return r
}

// DistSQLReceiver is an execinfra.RowReceiver and execinfra.BatchReceiver that
// writes results to a rowResultWriter and batchResultWriter, respectively. This
// is where the DistSQL execution meets the SQL Session - the result writer
// comes from a client Session.
//
// DistSQLReceiver also update the RangeDescriptorCache in response to DistSQL
// metadata about misplanned ranges.
type DistSQLReceiver struct {
	ctx context.Context

	// These two interfaces refer to the same object, but batchWriter might be
	// unset (resultWriter is always set). These are used to send the results
	// to.
	resultWriter rowResultWriter
	batchWriter  batchResultWriter
	// concurrentErrorCh is a buffered channel that allows for concurrent
	// goroutines to tell the main execution goroutine that there is an error.
	// The main goroutine will read from this channel on the next call to Push
	// or PushBatch.
	//
	// This channel is needed since rowResultWriter is not thread-safe and will
	// only be sent on during distributed plan execution.
	concurrentErrorCh chan error
	stmtType          tree.StatementReturnType
	// The transaction in which the flow producing data for this
	// receiver runs. The DistSQLReceiver updates the transaction in
	// response to RetryableTxnError's and when distributed processors
	// pass back LeafTxnFinalState objects via ProducerMetas. Nil if no
	// transaction should be updated on errors (i.e. if the flow overall
	// doesn't run in a transaction).
	txn *kv.Txn
}

// PlanAndRun generates a physical plan from a planNode tree and executes it. It
// assumes that the tree is supported (see checkSupportForPlanNode).
//
// All errors encountered are reported to the DistSQLReceiver's resultWriter.
// Additionally, if the error is a "communication error" (an error encountered
// while using that resultWriter), the error is also stored in
// DistSQLReceiver.commErr. That can be tested to see if a client session needs
// to be closed.
//
// An allow-list of errors that are encountered during the distributed query
// execution are transparently retried by re-planning and re-running the query
// as local (as long as no data has been communicated to the result writer).
//
// - finishedSetupFn, if non-nil, is called synchronously after all the local
// processors have been created but haven't started running yet. If the query is
// re-planned as local after having encountered an error during distributed
// execution, then finishedSetupFn will be called twice.
func (dsp *DistSQLPlanner) PlanAndRun(
	ctx context.Context,
	evalCtx *extendedEvalContext,
	planCtx *PlanningCtx,
	txn *kv.Txn,
	plan planMaybePhysical,
	recv *DistSQLReceiver,
	finishedSetupFn func(localFlow flowinfra.Flow),
) {

}
