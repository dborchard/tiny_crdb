package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
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

func (r *streamingCommandResult) SetColumns(ctx context.Context, columns colinfo.ResultColumns) {
	//TODO implement me
	panic("implement me")
}

func (r *streamingCommandResult) AddRow(ctx context.Context, row tree.Datums) error {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
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
