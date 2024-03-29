package sql

import (
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser/statements"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"time"
)

type Command interface {
	fmt.Stringer
	command() string
	isExtendedProtocolCmd() bool
}

var _ Command = &ExecStmt{}
var _ Command = &Sync{}
var _ Command = &PrepareStmt{}
var _ Command = &BindStmt{}
var _ Command = &ExecPortal{}

// ExecStmt is the command for running a query sent through the "simple" pgwire
// protocol.
type ExecStmt struct {
	// Information returned from parsing: AST, SQL, NumPlaceholders.
	// Note that AST can be nil, in which case executing it should produce an
	// "empty query response" message.
	statements.Statement[tree.Statement]

	// TimeReceived is the time at which the exec message was received
	// from the client. Used to compute the service latency.
	TimeReceived time.Time
	// ParseStart/ParseEnd are the timing info for parsing of the query. Used for
	// stats reporting.
	ParseStart time.Time
	ParseEnd   time.Time

	// LastInBatch indicates if this command contains the last query in a
	// simple protocol Query message that contains a batch of 1 or more queries.
	LastInBatch bool
	// LastInBatchBeforeShowCommitTimestamp indicates that this command contains
	// the second-to-last query in a simple protocol Query message that contains
	// a batch of 2 or more queries and the last query is SHOW COMMIT TIMESTAMP.
	// Detecting this case allows us to treat this command as the LastInBatch
	// such that the SHOW COMMIT TIMESTAMP statement can return the timestamp of
	// the transaction which applied to all the other statements in the batch.
	// Note that SHOW COMMIT TIMESTAMP is not permitted in any other position in
	// such a multi-statement implicit transaction.
	LastInBatchBeforeShowCommitTimestamp bool
}

func (e ExecStmt) isExtendedProtocolCmd() bool {
	//TODO implement me
	panic("implement me")
}

// Sync is a command that serves two purposes:
// 1) It marks the end of one batch of commands and the beginning of the next.
// stmtBuf.seekToNextBatch will seek to this marker.
// 2) It generates a ReadyForQuery protocol message.
//
// A Sync command is generated for both the simple and the extended pgwire
// protocol variants. So, it doesn't strictly correspond to a pgwire sync
// message - those are not sent in the simple protocol. We synthesize Sync
// commands though because their handling matches the simple protocol too.
type Sync struct {
	// ExplicitFromClient specifies whether this Sync command was generated by the
	// client during the extended protocol, or implicitly created to handle a
	// simple protocol query.
	ExplicitFromClient bool
}

func (s Sync) isExtendedProtocolCmd() bool {
	//TODO implement me
	panic("implement me")
}

// PrepareStmt is the command for creating a prepared statement.
type PrepareStmt struct {
	Name string
	statements.Statement[tree.Statement]
	ParseStart time.Time
	ParseEnd   time.Time
}

func (p PrepareStmt) isExtendedProtocolCmd() bool {
	//TODO implement me
	panic("implement me")
}

// BindStmt is the Command for creating a portal from a prepared statement.
type BindStmt struct {
	PreparedStatementName string
	PortalName            string
	// OutFormats contains the requested formats for the output columns.
	// It either contains a bunch of format codes, in which case the number will
	// need to match the number of output columns of the portal, or contains a single
	// code, in which case that code will be applied to all columns.
	//OutFormats []pgwirebase.FormatCode
	// Args are the arguments for the prepared statement.
	// They are passed in without decoding because decoding requires type
	// inference to have been performed.
	//
	// A nil element means a tree.DNull argument.
	Args [][]byte
	// ArgFormatCodes are the codes to be used to deserialize the Args.
	// It either contains a bunch of format codes, in which case the number will
	// need to match the number of arguments for the portal, or contains a single
	// code, in which case that code will be applied to all arguments.
	//ArgFormatCodes []pgwirebase.FormatCode

	// internalArgs, if not nil, represents the arguments for the prepared
	// statements as produced by the internal clients. These don't need to go
	// through encoding/decoding of the args. However, the types of the datums
	// must correspond exactly to the inferred types (but note that the types of
	// the datums are passes as type hints to the PrepareStmt command, so the
	// inferred types should reflect that).
	// If internalArgs is specified, Args and ArgFormatCodes are ignored.
	internalArgs []tree.Datum
}

func (b BindStmt) isExtendedProtocolCmd() bool {
	//TODO implement me
	panic("implement me")
}

// ExecPortal is the Command for executing a portal.
type ExecPortal struct {
	Name string
	// limit is a feature of pgwire that we don't really support. We accept it and
	// don't complain as long as the statement produces fewer results than this.
	Limit int
	// TimeReceived is the time at which the exec message was received
	// from the client. Used to compute the service latency.
	TimeReceived time.Time
	// FollowedBySync is true if the next command after this is a Sync. This is
	// used to enable the 1PC txn fast path in the extended protocol.
	FollowedBySync bool
}

func (e ExecPortal) isExtendedProtocolCmd() bool {
	//TODO implement me
	panic("implement me")
}

// ================================ Functions ================================
func (e ExecPortal) String() string {
	//TODO implement me
	panic("implement me")
}

func (e ExecPortal) command() string {
	//TODO implement me
	panic("implement me")
}

func (b BindStmt) String() string {
	//TODO implement me
	panic("implement me")
}

func (b BindStmt) command() string {
	//TODO implement me
	panic("implement me")
}

func (p PrepareStmt) String() string {
	//TODO implement me
	panic("implement me")
}

func (p PrepareStmt) command() string {
	//TODO implement me
	panic("implement me")
}

func (e ExecStmt) String() string {
	//TODO implement me
	panic("implement me")
}

func (e ExecStmt) command() string {
	//TODO implement me
	panic("implement me")
}

func (s Sync) String() string {
	//TODO implement me
	panic("implement me")
}

func (s Sync) command() string {
	//TODO implement me
	panic("implement me")
}
