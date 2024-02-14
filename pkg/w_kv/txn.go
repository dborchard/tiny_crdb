package kv

// Txn is an in-progress distributed database transaction. A Txn is safe for
// concurrent use by multiple goroutines.
type Txn struct {
	db *DB

	// typ indicates the type of transaction.
	typ TxnType
}

// TxnType specifies whether a transaction is the root (parent)
// transaction, or a leaf (child) in a tree of kv.Txns, as is
// used in a DistSQL flow.
type TxnType int
