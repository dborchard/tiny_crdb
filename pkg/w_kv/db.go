package kv

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/y_util/hlc"
	"github.com/dborchard/tiny_crdb/pkg/y_util/stop"
)

// DB is a database handle to a single cockroach cluster. A DB is safe for
// concurrent use by multiple goroutines.
type DB struct {
	clock *hlc.Clock
	ctx   DBContext

	factory TxnSenderFactory
	crs     CrossRangeTxnWrapperSender
}

type DBContext struct {
	// Stopper is used for async tasks.
	Stopper *stop.Stopper
}

// CrossRangeTxnWrapperSender is a Sender whose purpose is to wrap
// non-transactional requests that span ranges into a transaction so they can
// execute atomically.
//
// TODO(andrei, bdarnell): This is a wart. Our semantics are that batches are
// atomic, but there's only historical reason for that. We should disallow
// non-transactional batches and scans, forcing people to use transactions
// instead. And then this Sender can go away.
type CrossRangeTxnWrapperSender struct {
	db      *DB
	wrapped Sender
}

func (db *DB) Txn(ctx context.Context, retryable func(context.Context, *Txn) error) error {
	return nil
}

// DefaultDBContext returns (a copy of) the default options for
// NewDBWithContext.
func DefaultDBContext(stopper *stop.Stopper) DBContext {
	return DBContext{
		Stopper: stopper,
	}
}

// NewDBWithContext returns a new DB with the given parameters.
func NewDBWithContext(factory TxnSenderFactory, clock *hlc.Clock, ctx DBContext) *DB {
	return &DB{
		clock:   clock,
		ctx:     ctx,
		factory: factory,
	}
}
