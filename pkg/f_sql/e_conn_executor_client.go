package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
)

// ClientComm is the interface used by the connExecutor for creating results to
// be communicated to client and for exerting some control over this
// communication.
//
// ClientComm is implemented by the pgwire connection.
type ClientComm interface {
	// CreateStatementResult creates a StatementResult for stmt.
	CreateStatementResult(
		stmt tree.Statement,
		pos CmdPos,
		limit int,
		portalName string,
		implicitTxn bool,
	) CommandResult

	// Flush delivers all the previous results to the client. The results might
	// have been buffered, in which case this flushes the buffer.
	Flush(pos CmdPos) error
}

// CommandResult represents the result of a statement. It which needs to be
// ultimately delivered to the client. pgwire.conn implements this.
type CommandResult interface {
	RestrictedCommandResult
	CommandResultClose
}

// RestrictedCommandResult is a subset of CommandResult meant to make it clear
// that its clients don't close the CommandResult.
type RestrictedCommandResult interface {
	CommandResultErrBase
	SetColumns(context.Context, colinfo.ResultColumns)
	AddRow(ctx context.Context, row tree.Datums) error
	AddBatch(ctx context.Context, batch coldata.Batch) error
	SupportsAddBatch() bool
	SetRowsAffected(ctx context.Context, n int)
	RowsAffected() int
}
type CommandResultClose interface {
	Discard()
}

type internalClientComm struct {
	results []*streamingCommandResult
	w       ieResultWriter
	mode    ieExecutionMode
	sync    func([]*streamingCommandResult)
}

func (icc *internalClientComm) CreateStatementResult(stmt tree.Statement, pos CmdPos, limit int, portalName string, implicitTxn bool) CommandResult {
	return icc.createRes(pos)
}

func (icc *internalClientComm) createRes(pos CmdPos) *streamingCommandResult {
	res := &streamingCommandResult{
		pos:             pos,
		w:               icc.w,
		discardCallback: func() {},
	}
	icc.results = append(icc.results, res)
	return res
}

func (icc *internalClientComm) Flush(pos CmdPos) error {
	//TODO implement me
	panic("implement me")
}
