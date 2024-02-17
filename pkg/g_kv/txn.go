package kv

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvpb"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvserver/concurrency/isolation"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// Txn is an in-progress distributed database transaction. A Txn is safe for
// concurrent use by multiple goroutines.
type Txn struct {
	db  *DB
	typ TxnType
	mu  struct {
		debugName    string
		userPriority roachpb.UserPriority
		sender       TxnSender
		deadline     hlc.Timestamp
	}
}

// NewTxn returns a new RootTxn.
// Note: for SQL usage, prefer NewTxnWithSteppingEnabled() below.
// Note: for KV usage that should be subject to admission control, prefer
// NewTxnRootKV() below.
//
// If the transaction is used to send any operations, Commit() or Rollback()
// should eventually be called to commit/rollback the transaction (including
// stopping the heartbeat loop).
//
// gatewayNodeID: If != 0, this is the ID of the node on whose behalf this
//
//	transaction is running. Normally this is the current node, but in the case
//	of Txns created on remote nodes by DistSQL this will be the gateway.
//	If 0 is passed, then no value is going to be filled in the batches sent
//	through this txn. This will have the effect that the DistSender will fill
//	in the batch with the current node's ID.
//	If the gatewayNodeID is set and this is a root transaction, we optimize
//	away any clock uncertainty for our own node, as our clock is accessible.
//
// See also db.NewTxn().
func NewTxn(ctx context.Context, db *DB) *Txn {
	return NewTxnWithAdmissionControl(ctx, db)
}

// NewTxnWithAdmissionControl creates a new transaction with the specified
// admission control source and priority. See NewTxn() for details.
func NewTxnWithAdmissionControl(
	ctx context.Context,
	db *DB,
) *Txn {

	now := db.clock.NowAsClockTimestamp()
	kvTxn := roachpb.MakeTransaction(
		"unnamed",
		nil, // baseKey
		isolation.Serializable,
		roachpb.NormalUserPriority,
		now.ToTimestamp(),
	)
	txn := NewTxnFromProto(ctx, db, now, RootTxn, &kvTxn)
	return txn
}

// NewTxnFromProto is like NewTxn but assumes the Transaction object is already initialized.
// Do not use this directly; use NewTxn() instead.
// This function exists for testing only.
func NewTxnFromProto(
	ctx context.Context,
	db *DB,
	now hlc.ClockTimestamp,
	typ TxnType,
	proto *roachpb.Transaction,
) *Txn {
	txn := &Txn{db: db, typ: typ}
	txn.mu.userPriority = roachpb.NormalUserPriority
	txn.mu.sender = db.factory.RootTransactionalSender(proto, txn.mu.userPriority)
	return txn
}

// runTxn runs the given retryable transaction function using the given *Txn.
func runTxn(ctx context.Context, txn *Txn, retryable func(context.Context, *Txn) error) error {
	err := txn.exec(ctx, retryable)
	if err != nil {
		if rollbackErr := txn.Rollback(ctx); rollbackErr != nil {
			println("rollbackErr: ", rollbackErr.Error())
		}
	}
	return err
}

// exec executes fn in the context of a distributed transaction. The closure is
// retried on retriable errors.
// If no error is returned by the closure, an attempt to commit the txn is made.
//
// When this method returns, txn might be in any state; exec does not attempt
// to clean up the transaction before returning an error. In case of
// TransactionAbortedError, txn is reset to a fresh transaction, ready to be
// used.
func (txn *Txn) exec(ctx context.Context, fn func(context.Context, *Txn) error) (err error) {
	// Run fn in a retry loop until we encounter a success or
	// error condition this loop isn't capable of handling.
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = fn(ctx, txn)

		if err != nil {
			return err
		}
		err = txn.Commit(ctx)
		break
	}

	return nil
}

// Commit sends an EndTxnRequest with Commit=true.
func (txn *Txn) Commit(ctx context.Context) error {
	if txn.typ != RootTxn {
		return errors.New("commit() called on leaf txn")
	}

	return txn.commit(ctx)
}

func (txn *Txn) commit(ctx context.Context) error {
	ba := &kvpb.BatchRequest{}
	_, pErr := txn.Send(ctx, ba)
	return pErr.Error
}

// Send runs the specified calls synchronously in a single batch and
// returns any errors. If the transaction is read-only or has already
// been successfully committed or aborted, a potential trailing
// EndTxn call is silently dropped, allowing the caller to always
// commit or clean-up explicitly even when that may not be required
// (or even erroneous). Returns (nil, nil) for an empty batch.
func (txn *Txn) Send(
	ctx context.Context, ba *kvpb.BatchRequest,
) (*kvpb.BatchResponse, *kvpb.Error) {
	return txn.mu.sender.Send(ctx, ba)
}

// Rollback sends an EndTxnRequest with Commit=false.
// txn is considered finalized and cannot be used to send any more commands.
func (txn *Txn) Rollback(ctx context.Context) error {
	return txn.rollback(ctx).Error
}

func (txn *Txn) rollback(ctx context.Context) *kvpb.Error {
	ba := &kvpb.BatchRequest{}
	_, pErr := txn.Send(ctx, ba)
	return pErr
}

// Put sets the value for a key
//
// key can be either a byte slice or a string. value can be any key type, a
// protoutil.Message or any Go primitive type (bool, int, etc).
func (txn *Txn) Put(ctx context.Context, key, value interface{}) error {
	b := txn.NewBatch()
	b.Put(key, value)
	return txn.Run(ctx, b)
}

// NewBatch creates and returns a new empty batch object for use with the Txn.
func (txn *Txn) NewBatch() *Batch {
	return &Batch{txn: txn}
}

// Run executes the operations queued up within a batch. Before executing any
// of the operations the batch is first checked to see if there were any errors
// during its construction (e.g. failure to marshal a proto message).
//
// The operations within a batch are run in parallel and the order is
// non-deterministic. It is an unspecified behavior to modify and retrieve the
// same key within a batch.
//
// Upon completion, Batch.Results will contain the results for each
// operation. The order of the results matches the order the operations were
// added to the batch.
func (txn *Txn) Run(ctx context.Context, b *Batch) error {
	return sendAndFill(ctx, txn.Send, b)
}

func (txn *Txn) Sender() TxnSender {
	return txn.mu.sender
}

// sendAndFill is a helper which sends the given batch and fills its results,
// returning the appropriate error which is either from the first failing call,
// or an "internal" error.
func sendAndFill(ctx context.Context, send SenderFunc, b *Batch) error {
	// Errors here will be attached to the results, so we will get them from
	// the call to fillResults in the regular case in which an individual call
	// fails. But send() also returns its own errors, so there's some dancing
	// here to do because we want to run fillResults() so that the individual
	// result gets initialized with an error from the corresponding call.
	ba := &kvpb.BatchRequest{}
	ba.Requests = b.reqs
	b.response, b.pErr = send(ctx, ba)
	return b.pErr.Error
}
