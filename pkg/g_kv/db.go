package kv

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
)

// DB is a database handle to a single cockroach cluster. A DB is safe for
// concurrent use by multiple goroutines.
type DB struct {
	clock *hlc.Clock
	ctx   DBContext

	factory TxnSenderFactory
	crs     CrossRangeTxnWrapperSender
}

// NewDB returns a new DB.
func NewDB(
	actx context.Context, factory TxnSenderFactory, clock *hlc.Clock, stopper *stop.Stopper,
) *DB {
	return NewDBWithContext(factory, clock, DefaultDBContext(stopper))
}

// NewDBWithContext returns a new DB with the given parameters.
func NewDBWithContext(factory TxnSenderFactory, clock *hlc.Clock, ctx DBContext) *DB {
	return &DB{
		clock:   clock,
		ctx:     ctx,
		factory: factory,
	}
}

func (db *DB) Txn(ctx context.Context, retryable func(context.Context, *Txn) error) error {
	return db.TxnWithAdmissionControl(ctx, retryable)
}

// TxnWithAdmissionControl is like Txn, but uses a configurable admission
// control source and priority.
func (db *DB) TxnWithAdmissionControl(
	ctx context.Context,
	retryable func(context.Context, *Txn) error,
) error {
	txn := NewTxnWithAdmissionControl(ctx, db)
	return runTxn(ctx, txn, retryable)
}

type DBContext struct {
	// Stopper is used for async tasks.
	Stopper *stop.Stopper
}

func DefaultDBContext(stopper *stop.Stopper) DBContext {
	return DBContext{
		Stopper: stopper,
	}
}

type CrossRangeTxnWrapperSender struct {
	db      *DB
	wrapped Sender
}
