package kv

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvpb"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvserver/concurrency/isolation"
	"github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// TxnType specifies whether a transaction is the root (parent)
// transaction, or a leaf (child) in a tree of kv.Txns, as is
// used in a DistSQL flow.
type TxnType int

const (
	_ TxnType = iota
	// RootTxn specifies this sender is the root transaction, and is
	// responsible for aggregating all transactional state and
	// finalizing the transaction. The root txn is responsible for
	// heartbeating the transaction record.
	RootTxn
	// LeafTxn specifies this sender is for one of potentially many
	// distributed client transactions. The state from this transaction
	// must be propagated back to the root transaction and used to
	// augment its state before the transaction can be finalized. Leaf
	// transactions do not heartbeat the transaction record.
	//
	// Note: As leaves don't perform heartbeats, the transaction might
	// be cleaned up while this leaf is executing an operation. We rely
	// on the cleanup process poisoning the AbortSpans for all intents
	// so that reads performed through a leaf txn don't miss writes
	// previously performed by the transaction (at least not until the
	// expiration of the GC period / abort span entry timeout).
	LeafTxn
)

// TxnSenderFactory is the interface used to create new instances
// of TxnSender.
type TxnSenderFactory interface {
	// RootTransactionalSender returns a root sender to be used for
	// transactional requests. txn contains the transaction whose
	// requests this sender will carry.
	RootTransactionalSender(txn *roachpb.Transaction, pri roachpb.UserPriority) TxnSender

	// LeafTransactionalSender returns a leaf sender to be used for
	// transactional requests on behalf of a root.
	LeafTransactionalSender(tis *roachpb.LeafTxnInputState) TxnSender

	// NonTransactionalSender returns a sender to be used for
	// non-transactional requests. Generally this is a sender that
	// TransactionalSender() wraps.
	NonTransactionalSender() Sender
}

// Sender is implemented by modules throughout the crdb stack, on both
// the "client" and the "server", involved in passing along and
// ultimately evaluating requests (batches). The interface is now
// considered regrettable because it's too narrow and at times leaky.
// Notable implementors: w_kv.Txn, kvcoord.TxnCoordSender, server.Node,
// kvserver.Store, kvserver.Replica.
type Sender interface {
	// Send sends a batch for evaluation. Either a response or an error is
	// returned.
	//
	// The caller retains ownership of all the memory referenced by the
	// BatchRequest; the callee is not allowed to hold on to any parts
	// of it past after it returns from the call (this is so that the
	// client module can allocate requests from a pool and reuse
	// them). For example, the DistSender makes sure that, if there are
	// concurrent requests, it waits for all of them before returning,
	// even in error cases.
	//
	// Once the request reaches the `transport` module, another
	// restriction applies (particularly relevant for the case when the
	// node that the transport is talking to is local, and so there's
	// not gRPC marshaling/unmarshaling):
	// - the callee has to treat everything inside the BatchRequest as
	// read-only. This is so that the client module retains the right to
	// pass pointers into its internals, like for example the
	// Transaction. This wouldn't work if the server would be allowed to
	// change the Transaction willy-nilly.
	//
	// TODO(andrei): The client does not currently use this last
	// guarantee; it clones the txn for every request. Given that a
	// w_kv.Txn can be used concurrently, in order for the client to
	// take advantage of this, it would need to switch to a
	// copy-on-write scheme so that its updates to the txn do not race
	// with the server reading it. We should do this to avoid the
	// cloning allocations. And to be frank, it'd be a good idea for the
	// BatchRequest/Response to generally stop carrying transactions;
	// the requests usually only need a txn id and some timestamp. The
	// responses would ideally contain a list of targeted instructions
	// about what the client should update, as opposed to a full txn
	// that the client is expected to diff with its copy and apply all
	// the updates.
	Send(context.Context, *kvpb.BatchRequest) (*kvpb.BatchResponse, *kvpb.Error)
}

// TxnSender is the interface used to call into a CockroachDB instance
// when sending transactional requests. In addition to the usual
// Sender interface, TxnSender facilitates marshaling of transaction
// metadata between the "root" w_kv.Txn and "leaf" instances.
type TxnSender interface {
	Sender

	// SetIsoLevel sets the txn's isolation level.
	SetIsoLevel(isolation.Level) error

	// IsoLevel returns the txn's isolation level.
	IsoLevel() isolation.Level

	// TxnStatus exports the txn's status.
	TxnStatus() roachpb.TransactionStatus

	// ReadTimestamp returns the transaction's current read timestamp.
	// Note a transaction can be internally pushed forward in time
	// before committing so this is not guaranteed to be the commit
	// timestamp. Use CommitTimestamp() when needed.
	ReadTimestamp() hlc.Timestamp

	// ReadTimestampFixed returns true if the read timestamp has been fixed
	// and cannot be pushed forward.
	ReadTimestampFixed() bool

	// CommitTimestamp returns the transaction's commit timestamp.
	//
	// If the transaction is committed, the method returns the timestamp at
	// which the transaction performed all of its writes.
	//
	// If the transaction is aborted, the method returns an error.
	//
	// If the transaction is pending and running under serializable isolation,
	// the method returns the transaction's current provisional commit
	// timestamp. It also fixes the transaction's read timestamp to ensure
	// that the transaction cannot be pushed to a later timestamp and still
	// commit. It does so by disabling read refreshes. As a result, using this
	// method just once increases the likelihood that a retry error will
	// bubble up to a client.
	//
	// If the transaction is pending and running under a weak isolation level,
	// the method returns an error. Fixing the commit timestamp early is not
	// supported for transactions running under weak isolation levels.
	CommitTimestamp() (hlc.Timestamp, error)
}
