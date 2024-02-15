package sql

import (
	"context"
	parser "github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/colinfo"
	isql "github.com/dborchard/tiny_crdb/pkg/f_sql/d_isql"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"sync"
)

// InternalDB stored information needed to construct a new
// internal executor.
type InternalDB struct {
	server *Server
	db     *kv.DB
}

var _ isql.DB = &InternalDB{}

// NewInternalDB returns a new InternalDB.
func NewInternalDB(s *Server) *InternalDB {
	return &InternalDB{
		server: s,
		db:     s.cfg.DB,
	}
}

// NewShimInternalDB is used to bootstrap the server which needs access to
// components which will ultimately have a handle to an InternalDB. Some of
// those components may attempt to access the *kv.DB before the InternalDB
// has been fully initialized. To get around this, we initially construct
// an InternalDB with just a handle to a *kv.DB and then we'll fill in the
// object during sql server construction.
func NewShimInternalDB(db *kv.DB) *InternalDB {
	return &InternalDB{db: db}
}

func (ief *InternalDB) KV() *kv.DB {
	return ief.db
}

// Txn is used to run queries with internal executor in a transactional
// manner.
func (ief *InternalDB) Txn(
	ctx context.Context, f func(context.Context, isql.Txn) error, opts ...isql.TxnOption,
) error {
	wrapped := func(ctx context.Context, txn *internalTxn) error { return f(ctx, txn) }
	return ief.txn(ctx, wrapped, opts...)
}

func (ief *InternalDB) txn(
	ctx context.Context, f func(context.Context, *internalTxn) error, opts ...isql.TxnOption,
) error {
	var cfg isql.TxnConfig
	cfg.Init(opts...)

	db := ief.server.cfg.DB
	run := db.Txn
	for {
		if err := run(ctx, func(ctx context.Context, kvTxn *kv.Txn) (err error) {
			ie, commitTxnFn := ief.newInternalExecutorWithTxn(
				ctx,
				cfg.GetSessionData(),
				//cf.GetClusterSettings(),
				kvTxn,
				//descsCol,
			)
			txn := internalTxn{txn: kvTxn}
			txn.InternalExecutor = ie
			if err := f(ctx, &txn); err != nil {
				return err
			}
			return commitTxnFn(ctx)
		}); err != nil {
			return err
		}
	}
}

// newInternalExecutorWithTxn creates an internal executor with txn-related info,
// such as descriptor collection and schema change job records, etc.
// This function should only be used under
// InternalDB.DescsTxnWithExecutor().
// TODO (janexing): This function will be soon refactored after we change
// the internal executor infrastructure with a single conn executor for all
// sql statement executions within a txn.
func (ief *InternalDB) newInternalExecutorWithTxn(
	ctx context.Context,
	sd *sessiondata.SessionData,
	//settings *cluster.Settings,
	txn *kv.Txn,
	// descCol *descs.Collection,
) (InternalExecutor, internalExecutorCommitTxnFunc) {
	return InternalExecutor{}, nil
}

// internalExecutorCommitTxnFunc is to commit the txn associated with an
// internal executor.
type internalExecutorCommitTxnFunc func(ctx context.Context) error

// Executor returns an Executor not bound with any txn.
func (ief *InternalDB) Executor(opts ...isql.ExecutorOption) isql.Executor {
	var cfg isql.ExecutorConfig
	cfg.Init(opts...)
	ie := MakeInternalExecutor(ief.server) //, ief.memMetrics, ief.monitor)
	return &ie
}

type internalTxn struct {
	internalExecutor
	txn *kv.Txn
}

func (txn *internalTxn) SessionData() *sessiondata.SessionData {
	return txn.sessionDataStack.Top()
}

func (txn *internalTxn) KV() *kv.Txn { return txn.txn }

type internalExecutor struct {
	InternalExecutor
}

// InternalExecutor can be used internally by code modules to execute SQL
// statements without needing to open a SQL connection.
//
// InternalExecutor can execute one statement at a time. As of 03/2018, it
// doesn't offer a session interface for maintaining session state or for
// running explicit SQL transactions. However, it supports running SQL
// statements inside a higher-lever (KV) txn and inheriting session variables
// from another session.
//
// Methods not otherwise specified are safe for concurrent execution.
type InternalExecutor struct {
	s *Server
	// sessionDataStack, if not nil, represents the session variable stack used by
	// statements executed on this internalExecutor. Note that queries executed
	// by the executor will run on copies of the top element of this data.
	sessionDataStack *sessiondata.Stack
}

func MakeInternalExecutor(s *Server) InternalExecutor {
	return InternalExecutor{s: s}
}

func (ie *InternalExecutor) ExecEx(ctx context.Context, opName string, txn *kv.Txn, o sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) (int, error) {
	// We will run the query to completion, so we can use an async result
	// channel.
	rw := newAsyncIEResultChannel()
	// Since we only return the number of rows affected as given by the
	// rowsIterator, we execute this stmt in "rows affected" mode allowing the
	// internal executor to transparently retry.
	const mode = rowsAffectedIEExecutionMode
	it, err := ie.execInternal(ctx, opName, rw, mode, txn, stmt, qargs...)
	if err != nil {
		return 0, err
	}
	// We need to exhaust the iterator so that it can count the number of rows
	// affected.
	var ok bool
	for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
	}
	if err != nil {
		return 0, err
	}
	return it.rowsAffected, nil
}

// QueryBufferedExWithCols is like QueryBufferedEx, additionally returning the computed
// ResultColumns of the input query.
func (ie *InternalExecutor) QueryBufferedExWithCols(
	ctx context.Context,
	opName string,
	txn *kv.Txn,
	session sessiondata.InternalExecutorOverride,
	stmt string,
	qargs ...interface{},
) ([]tree.Datums, colinfo.ResultColumns, error) {
	datums, cols, err := ie.queryInternalBuffered(ctx, opName, txn, session, stmt, 0, qargs...)
	return datums, cols, err
}

func (ie *InternalExecutor) queryInternalBuffered(
	ctx context.Context,
	opName string,
	txn *kv.Txn,
	sessionDataOverride sessiondata.InternalExecutorOverride,
	stmt string,
	limit int,
	qargs ...interface{},
) ([]tree.Datums, colinfo.ResultColumns, error) {
	rw := newAsyncIEResultChannel()
	it, err := ie.execInternal(ctx, opName, rw, defaultIEExecutionMode, txn, stmt, qargs...)
	if err != nil {
		return nil, nil, err
	}
	var rows []tree.Datums
	var ok bool
	for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
		rows = append(rows, it.Cur())
		if limit != 0 && len(rows) == limit {
			err = it.Close()
			break
		}
	}
	if err != nil {
		return nil, nil, err
	}
	return rows, it.Types(), nil
}

func (ie *InternalExecutor) execInternal(
	ctx context.Context,
	opName string,
	rw *ieResultChannel,
	mode ieExecutionMode,
	txn *kv.Txn,
	stmt string,
	qargs ...interface{}) (r *rowsIterator, retErr error) {
	parsed, err := parser.ParseOne(stmt)
	if err != nil {
		return nil, err
	}

	stmtBuf := NewStmtBuf()
	var wg sync.WaitGroup

	syncCallback := func(results []*streamingCommandResult) {
		// Close the stmtBuf so that the connExecutor exits its run() loop.
		stmtBuf.Close()
		for _, res := range results {
			if res.Err() != nil {
				// If we encounter an error, there's no point in looking
				// further; the rest of the commands in the batch have been
				// skipped.
				_ = rw.addResult(ctx, ieIteratorResult{err: res.Err()})
				return
			}
		}
	}
	// errCallback is called if an error is returned from the connExecutor's
	// run() loop.
	errCallback := func(err error) {
		_ = rw.addResult(ctx, ieIteratorResult{err: err})
	}

	err = ie.runWithEx(ctx, txn, rw, mode, stmtBuf, &wg, syncCallback, errCallback)
	if err != nil {
		return nil, err
	}

	if len(qargs) == 0 {
		if err := stmtBuf.Push(
			ctx,
			ExecStmt{
				Statement:   parsed,
				LastInBatch: true,
			}); err != nil {
			return nil, err
		}

		if err := stmtBuf.Push(ctx, Sync{ExplicitFromClient: false}); err != nil {
			return nil, err
		}
	}
	r = &rowsIterator{
		r:       rw,
		mode:    mode,
		stmtBuf: stmtBuf,
		wg:      &wg,
	}

	if parsed.AST.StatementReturnType() != tree.Rows {
		r.resultCols = colinfo.ResultColumns{
			colinfo.ResultColumn{
				Name: "rows_affected",
				Typ:  types.Int,
			},
		}
	}

	// Now we need to block the reader goroutine until the query planning has
	// been performed by the connExecutor goroutine. We do so by waiting until
	// the first object is sent on the data channel.
	{
		var first ieIteratorResult
		if first, r.done, r.lastErr = rw.firstResult(ctx); !r.done {
			r.first = &first
		}
	}
	for !r.done && r.first.cols != nil {
		// If the query is of ROWS statement type, the very first thing sent on
		// the channel will be the column schema. This will occur before the
		// query is given to the execution engine, so we actually need to get
		// the next piece from the data channel.
		//
		// We also need to keep on looping until we get the first actual result
		// with rows. In theory, it is possible for a stmt of ROWS type to
		// encounter a retry error after sending the column schema but before
		// going into the execution engine. In such a scenario we want to keep
		// the latest column schema (in case there was a schema change
		// in-between retries).
		//
		// Note that only statements of ROWS type should send the cols, but we
		// choose to be defensive and don't assert that.
		if parsed.AST.StatementReturnType() == tree.Rows {
			r.resultCols = r.first.cols
		}
		var first ieIteratorResult
		first, r.done, r.lastErr = rw.nextResult(ctx)
		if !r.done {
			r.first = &first
		}
	}

	return r, nil
}

func (ie *InternalExecutor) runWithEx(
	ctx context.Context,
	txn *kv.Txn,
	w ieResultWriter,
	mode ieExecutionMode,
	stmtBuf *StmtBuf,
	wg *sync.WaitGroup,
	syncCallback func([]*streamingCommandResult),
	errCallback func(error),
) error {
	ex, err := ie.initConnEx(ctx, txn, w, mode, stmtBuf, syncCallback)
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		if err := ex.run(ctx, nil); err != nil {
			errCallback(err)
		}
		w.finish()
		closeMode := normalClose
		if txn != nil {
			closeMode = externalTxnClose
		}
		ex.close(ctx, closeMode)
		wg.Done()
	}()
	return nil
}
