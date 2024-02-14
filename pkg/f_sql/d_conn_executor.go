package sql

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"io"
)

type closeType int

const (
	normalClose closeType = iota
	panicClose
	// externalTxnClose means that the connExecutor has been used within a
	// higher-level txn (through the InternalExecutor).
	externalTxnClose
)

type executorType int

const (
	executorTypeExec executorType = iota
	executorTypeInternal
)

type connExecutor struct {
	// executorType is set to whether this executor is an ordinary executor which
	// responds to user queries or an internal one.
	executorType executorType
	server       *Server
	stmtBuf      *StmtBuf
	clientComm   ClientComm
	planner      planner

	// activated determines whether activate() was called already.
	// When this is set, close() must be called to release resources.
	activated      bool
	queryCancelKey string
	curStmtAST     tree.Statement

	// Finity "the machine" Automaton is the state machine controlling the state
	// below.
	machine fsm.Machine

	// The metrics to which the statement metrics should be accounted.
	// This is different whether the executor is for regular client
	// queries or for "internal" queries.
	metrics       *Metrics
	transitionCtx transitionCtx

	// state encapsulates fields related to the ongoing SQL txn. It is mutated as
	// the machine's ExtendedState.
	state txnState
}

// initConnEx creates a connExecutor and runs it on a separate goroutine. It
// takes in a StmtBuf into which commands can be pushed and a WaitGroup that
// will be signaled when connEx.run() returns.
//
// If txn is not nil, the statement will be executed in the respective txn.
//
// The ieResultWriter coordinates communicating results to the client. It may
// block execution when rows are being sent in order to prevent hazardous
// concurrency.
//
// sd will constitute the executor's session state.
func (ie *InternalExecutor) initConnEx(
	ctx context.Context,
	txn *kv.Txn,
	w ieResultWriter,
	mode ieExecutionMode,
	stmtBuf *StmtBuf,
	syncCallback func([]*streamingCommandResult),
) (*connExecutor, error) {
	clientComm := &internalClientComm{
		w:    w,
		mode: mode,
		sync: syncCallback,
	}
	clientComm.rowsAffectedState.rewind = func() {
		var zero int
		_ = w.addResult(ctx, ieIteratorResult{rowsAffected: &zero})
	}

	var ex *connExecutor
	var err error
	if txn == nil {
		postSetupFn := func(ex *connExecutor) {
		}
		ie.s = &Server{
			cfg: &ExecutorConfig{}}
		ex = ie.s.newConnExecutor(
			ctx,
			stmtBuf,
			clientComm,
			ie.s.cfg.GenerateID(),
			postSetupFn,
		)
	} else {
		ex, err = ie.newConnExecutorWithTxn(
			ctx,
			txn,
			stmtBuf,
			clientComm,
		)
		if err != nil {
			return nil, err
		}
	}

	ex.executorType = executorTypeInternal
	return ex, nil

}

func (ex *connExecutor) run(
	ctx context.Context,
	onCancel context.CancelFunc,
) (err error) {
	if !ex.activated {
		ex.activated = true
	}

	//sessionID := ex.planner.extendedEvalCtx.SessionID
	//ex.server.cfg.SessionRegistry.register(sessionID, ex.queryCancelKey, ex)

	defer func() {
	}()

	for {
		ex.curStmtAST = nil
		if err := ctx.Err(); err != nil {
			return err
		}

		var err error
		if err = ex.execCmd(); err != nil {
			if errors.Is(err, errDrainingComplete) || errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func (ex *connExecutor) close(ctx context.Context, closeType closeType) {
	//  NONE.
}
