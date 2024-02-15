package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
)

// CommandResultErrBase is the subset of CommandResult dealing with setting a
// query execution error.
type CommandResultErrBase interface {
	// SetError accumulates an execution error that needs to be reported to the
	// client. No further calls other than SetError(), Close() and Discard() are
	// allowed.
	//
	// Calling SetError() a second time overwrites the previously set error.
	SetError(error)

	// Err returns the error previously set with SetError(), if any.
	Err() error
}

// streamingCommandResult is a CommandResult that streams rows on the channel
// and can call a provided callback when closed.
type streamingCommandResult struct {
	pos CmdPos

	//// All the data (the rows and the metadata) are written into w. The
	//// goroutine writing into this streamingCommandResult might block depending
	//// on the synchronization strategy.
	w ieResultWriter

	// cannotRewind indicates whether this result has communicated some data
	// (rows or metadata) such that the corresponding command cannot be rewound.
	cannotRewind bool

	err          error
	rowsAffected int

	// closeCallback, if set, is called when Close() is called.
	closeCallback func()

	// discardCallback, if set, is called when Discard() is called.
	discardCallback func()
}

// SetColumns is part of the RestrictedCommandResult interface.
func (r *streamingCommandResult) SetColumns(ctx context.Context, cols colinfo.ResultColumns) {
	if cols == nil {
		cols = colinfo.ResultColumns{}
	}
	_ = r.w.addResult(ctx, ieIteratorResult{cols: cols})
}

// AddRow is part of the RestrictedCommandResult interface.
func (r *streamingCommandResult) AddRow(ctx context.Context, row tree.Datums) error {
	r.rowsAffected++
	rowCopy := make(tree.Datums, len(row))
	copy(rowCopy, row)
	r.cannotRewind = true
	return r.w.addResult(ctx, ieIteratorResult{row: rowCopy})
}

func (r *streamingCommandResult) AddBatch(ctx context.Context, batch coldata.Batch) error {
	//TODO implement me
	panic("implement me")
}

func (r *streamingCommandResult) SupportsAddBatch() bool {
	//TODO implement me
	panic("implement me")
}

func (r *streamingCommandResult) SetRowsAffected(ctx context.Context, n int) {
	r.rowsAffected = n
	if r.w != nil {
		_ = r.w.addResult(ctx, ieIteratorResult{rowsAffected: &n})
	}
}

func (r *streamingCommandResult) RowsAffected() int {
	//TODO implement me
	panic("implement me")
}

func (r *streamingCommandResult) Discard() {
	//TODO implement me
	panic("implement me")
}

func (r *streamingCommandResult) SetError(err error) {
	r.err = err
}

// Err is part of the RestrictedCommandResult interface.
func (r *streamingCommandResult) Err() error {
	return r.err
}
