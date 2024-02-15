package tree

// TransactionModes holds the transaction modes for a transaction.
type TransactionModes struct {
	Isolation     IsolationLevel
	UserPriority  UserPriority
	ReadWriteMode ReadWriteMode
	//AsOf          AsOfClause
	//Deferrable    DeferrableMode
}

// IsolationLevel holds the isolation level for a transaction.
type IsolationLevel int

// IsolationLevel values
const (
	UnspecifiedIsolation IsolationLevel = iota
	ReadUncommittedIsolation
	ReadCommittedIsolation
	RepeatableReadIsolation
	SnapshotIsolation
	SerializableIsolation
)

// ReadWriteMode holds the read write mode for a transaction.
type ReadWriteMode int

// ReadWriteMode values
const (
	UnspecifiedReadWriteMode ReadWriteMode = iota
	ReadOnly
	ReadWrite
)

// UserPriority holds the user priority for a transaction.
type UserPriority int

// UserPriority values
const (
	UnspecifiedUserPriority UserPriority = iota
	Low
	Normal
	High
)
