package sql

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

// SetError is part of the RestrictedCommandResult interface.
func (r *streamingCommandResult) SetError(err error) {
	r.err = err
	// Note that we intentionally do not send the error on the channel (when it
	// is present) since we might replace the error with another one later which
	// is allowed by the interface. An example of this is queryDone() closure
	// in execStmtInOpenState().
}

// Err is part of the RestrictedCommandResult interface.
func (r *streamingCommandResult) Err() error {
	return r.err
}
