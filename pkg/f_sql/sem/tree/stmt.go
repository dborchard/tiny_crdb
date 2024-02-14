package tree

import "fmt"

// StatementReturnType is the enumerated type for Statement return styles on
// the wire.
type StatementReturnType int

// StatementType is the enumerated type for Statement type return
type StatementType int

//go:generate stringer -type=StatementReturnType
const (
	// Ack indicates that the statement does not have a meaningful
	// return. Examples include SET, BEGIN, COMMIT.
	Ack StatementReturnType = iota
	// DDL indicates that the statement mutates the database schema.
	//
	// Note: this is the type indicated back to the client; it is not a
	// sufficient test for schema mutation for planning purposes. There
	// are schema-modifying statements (e.g. CREATE TABLE AS) which
	// report RowsAffected to the client, not DDL.
	// Use CanModifySchema() below instead.
	DDL
	// RowsAffected indicates that the statement returns the count of
	// affected rows.
	RowsAffected
	// Rows indicates that the statement returns the affected rows after
	// the statement was applied.
	Rows
	// CopyIn indicates a COPY FROM statement.
	CopyIn
	// CopyOut indicates a COPY TO statement.
	CopyOut
	// Replication indicates a replication protocol statement.
	Replication
	// Unknown indicates that the statement does not have a known
	// return style at the time of parsing. This is not first in the
	// enumeration because it is more convenient to have Ack as a zero
	// value, and because the use of Unknown should be an explicit choice.
	// The primary example of this statement type is EXECUTE, where the
	// statement type depends on the statement type of the prepared statement
	// being executed.
	Unknown
)

// Statement represents a statement.
type Statement interface {
	fmt.Stringer

	// StatementReturnType is the return styles on the wire
	// (Ack, DDL, RowsAffected, Rows, CopyIn or Unknown)
	StatementReturnType() StatementReturnType
	// StatementType identifies whether the statement is a DDL, DML, DCL, or TCL.
	StatementType() StatementType
	// StatementTag is a short string identifying the type of statement
	// (usually a single verb). This is different than the Stringer output,
	// which is the actual statement (including args).
	// TODO(dt): Currently tags are always pg-compatible in the future it
	// might make sense to pass a tag format specifier.
	StatementTag() string
}
