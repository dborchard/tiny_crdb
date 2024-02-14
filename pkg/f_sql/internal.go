package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/isql"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
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
	if sd := cfg.GetSessionData(); sd != nil {
		ie.SetSessionData(sd)
	}
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

func MakeInternalExecutor(
	s *Server, //memMetrics MemoryMetrics, monitor *mon.BytesMonitor,
) InternalExecutor {
	return InternalExecutor{
		s: s,
		//mon:        monitor,
		//memMetrics: memMetrics,
	}
}
func (ie *InternalExecutor) Exec(ctx context.Context, opName string, txn *kv.Txn, statement string, params ...interface{}) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) ExecEx(ctx context.Context, opName string, txn *kv.Txn, o sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryRow(ctx context.Context, opName string, txn *kv.Txn, statement string, qargs ...interface{}) (tree.Datums, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryRowEx(ctx context.Context, opName string, txn *kv.Txn, session sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) (tree.Datums, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryRowExWithCols(ctx context.Context, opName string, txn *kv.Txn, session sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) (tree.Datums, colinfo.ResultColumns, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryBuffered(ctx context.Context, opName string, txn *kv.Txn, stmt string, qargs ...interface{}) ([]tree.Datums, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryBufferedEx(ctx context.Context, opName string, txn *kv.Txn, session sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) ([]tree.Datums, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryIterator(ctx context.Context, opName string, txn *kv.Txn, stmt string, qargs ...interface{}) (isql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryIteratorEx(ctx context.Context, opName string, txn *kv.Txn, session sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) (isql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (ie *InternalExecutor) QueryBufferedExWithCols(ctx context.Context, opName string, txn *kv.Txn, session sessiondata.InternalExecutorOverride, stmt string, qargs ...interface{}) ([]tree.Datums, colinfo.ResultColumns, error) {
	//TODO implement me
	panic("implement me")
}

// SetSessionData binds the session variables that will be used by queries
// performed through this executor from now on. This creates a new session stack.
// It is recommended to use SetSessionDataStack.
//
// SetSessionData cannot be called concurrently with query execution.
func (ie *InternalExecutor) SetSessionData(sessionData *sessiondata.SessionData) {
	if sessionData != nil {
	}
}

type rowsIterator struct {
	r            ieResultReader
	rowsAffected int
	resultCols   colinfo.ResultColumns
	lastRow      tree.Datums
	lastErr      error
	done         bool
	// wg can be used to wait for the connExecutor's goroutine to exit.
	wg *sync.WaitGroup
	// depth tracks the current depth of recursion in Next(). Once it exceeds
	// iteratorDepthLimit, an error is returned to prevent stack overflow.
	depth int64
}

func (r *rowsIterator) Next(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) Cur() tree.Datums {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) RowsAffected() int {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) Close() error {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) Types() colinfo.ResultColumns {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) HasResults() bool {
	//TODO implement me
	panic("implement me")
}

type ieIteratorResult struct {
	// Exactly one of these 4 fields will be set.
	row          tree.Datums
	rowsAffected *int
	cols         colinfo.ResultColumns
	err          error
}
