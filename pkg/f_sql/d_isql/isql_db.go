// Copyright 2016 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package isql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
)

// DB enables clients to create and execute sql transactions from code inside
// the database. Multi-statement transactions should leverage the Txn method.
type DB interface {

	// KV returns the underlying *w_kv.DB.
	KV() *kv.DB

	// Txn enables callers to run transactions with a *Collection such that all
	// retrieved immutable descriptors are properly leased and all mutable
	// descriptors are handled. The function deals with verifying the two version
	// invariant and retrying when it is violated. Callers need not worry that they
	// write mutable descriptors multiple times. The call will explicitly wait for
	// the leases to drain on old versions of descriptors modified or deleted in the
	// transaction; callers do not need to call lease.WaitForOneVersion.
	// It also enables using internal executor to run sql queries in a txn manner.
	Txn(context.Context, func(context.Context, Txn) error, ...TxnOption) error

	// Executor constructs an internal executor not bound to a transaction.
	Executor(...ExecutorOption) Executor
}

// Txn is an internal sql transaction.
type Txn interface {

	// KV returns the underlying w_kv.Txn.
	KV() *kv.Txn

	// SessionData returns the transaction's SessionData.
	SessionData() *sessiondata.SessionData

	// Executor allows the user to execute transactional SQL statements.
	Executor
}

// Executor is meant to be used by layers below SQL in the system that
// nevertheless want to execute SQL queries (presumably against system tables).
// It is extracted in this "isql" package to avoid circular references and
// is implemented by *v_sql.InternalExecutor.
//
// TODO(ajwerner): Remove the txn argument from all the functions. They are
// now implicit -- if you have your hands on an c_isql.Txn, you know it's
// transactional. If you just have an Executor, you don't know, but you
// cannot assume one way or the other.
type Executor interface {

	// ExecEx is like Exec, but allows the caller to override some session data
	// fields.
	ExecEx(
		ctx context.Context,
		opName string,
		txn *kv.Txn,
		o sessiondata.InternalExecutorOverride,
		stmt string,
		qargs ...interface{},
	) (int, error)

	// QueryBufferedExWithCols is like QueryBufferedEx, additionally returning the computed
	// ResultColumns of the input query.
	QueryBufferedExWithCols(
		ctx context.Context,
		opName string,
		txn *kv.Txn,
		session sessiondata.InternalExecutorOverride,
		stmt string,
		qargs ...interface{},
	) ([]tree.Datums, colinfo.ResultColumns, error)
}

// Rows is an iterator interface that's exposed by the internal
// executor. It provides access to the rows from a query.
type Rows interface {
	// Next advances the iterator by one row, returning false if there are no
	// more rows in this iterator or if an error is encountered (the latter is
	// then returned).
	//
	// The iterator is automatically closed when false is returned, consequent
	// calls to Next will return the same values as when the iterator was
	// closed.
	Next(context.Context) (bool, error)

	// Cur returns the row at the current position of the iterator. The row is
	// safe to hold onto (meaning that calling Next() or Close() will not
	// invalidate it).
	Cur() tree.Datums

	// RowsAffected returns the count of rows affected by the statement.
	// This is only guaranteed to be accurate after Next() has returned
	// false (no more rows).
	RowsAffected() int

	// Close closes this iterator, releasing any resources it held open. Close
	// is idempotent and *must* be called once the caller is done with the
	// iterator.
	Close() error

	// Types returns the types of the columns returned by this iterator. The
	// returned array is guaranteed to correspond 1:1 with the tree.Datums rows
	// returned by Cur().
	//
	// WARNING: this method is safe to call anytime *after* the first call to
	// Next() (including after Close() was called).
	Types() colinfo.ResultColumns

	// HasResults returns true if there are results to the query, false otherwise.
	HasResults() bool
}
