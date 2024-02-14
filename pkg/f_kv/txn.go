package kv

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/f_kv/kvpb"
	"github.com/dborchard/tiny_crdb/pkg/f_kv/kvserver/concurrency/isolation"
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
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
