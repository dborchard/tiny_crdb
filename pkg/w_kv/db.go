package kv

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/y_util/hlc"
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
}

func (db *DB) Txn(ctx context.Context, retryable func(context.Context, *Txn) error) error {
	return nil
}
