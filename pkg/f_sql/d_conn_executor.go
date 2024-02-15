package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/a_parser/statements"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"github.com/dborchard/tiny_crdb/pkg/z_util/fsm"
)

// ResultBase is the common interface implemented by all the different command
// results.
type ResultBase interface {
	CommandResultErrBase
	CommandResultClose
}

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
	executorType   executorType
	server         *Server
	stmtBuf        *StmtBuf
	clientComm     ClientComm
	planner        planner
	activated      bool
	queryCancelKey string
	curStmtAST     tree.Statement
	machine        fsm.Machine
}

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
	var ex *connExecutor
	if txn == nil {
		postSetupFn := func(ex *connExecutor) {
		}
		ie.s = &Server{cfg: &ExecutorConfig{}}
		ex = ie.s.newConnExecutor(
			ctx,
			stmtBuf,
			clientComm,
			postSetupFn,
		)
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

	for {
		ex.curStmtAST = nil
		if err := ctx.Err(); err != nil {
			return err
		}

		var err error
		if err = ex.execCmd(); err != nil {
			return err
		}
	}
}

func (ex *connExecutor) close(ctx context.Context, closeType closeType) {
	//  NONE.
}

func (ex *connExecutor) initPlanner(ctx context.Context, p *planner) {
	//  NONE.
}

func (ex *connExecutor) execCmd() (retErr error) {
	ctx := context.Background()
	cmd, pos, err := ex.stmtBuf.CurCmd()
	if err != nil {
		return err // err could be io.EOF
	}

	var ev fsm.Event
	var payload fsm.EventPayload
	var res ResultBase

	switch tcmd := cmd.(type) {
	case ExecStmt:
		err := func() error {
			stmtRes := ex.clientComm.CreateStatementResult(tcmd.AST, pos, 0, "", true)
			res = stmtRes
			ev, payload, err = ex.execStmt(ctx, tcmd.Statement, nil, nil, stmtRes, true)
			return err
		}()
		if err != nil {
			return err
		}
	default:
		panic(errors.New("unknown command type"))
	}

	// If an event was generated, feed it to the state machine.
	var advInfo advanceInfo
	if ev != nil {
		var err error
		advInfo, err = ex.txnStateTransitionsApplyWrapper(ev, payload, res, pos)
		if err != nil {
			return err
		}
	} else {
		advInfo = advanceInfo{code: advanceOne}
	}

	// Move the cursor according to what the state transition told us to do.
	switch advInfo.code {
	case advanceOne:
		ex.stmtBuf.AdvanceOne()
	default:
		panic(errors.New("unexpected advance code: %s"))
	}
	return nil
}

func (ex *connExecutor) txnStateTransitionsApplyWrapper(
	ev fsm.Event, payload fsm.EventPayload, res ResultBase, pos CmdPos,
) (advanceInfo, error) {
	return advanceInfo{code: advanceOne}, nil
}

func (ex *connExecutor) execStmt(
	ctx context.Context,
	parserStmt statements.Statement[tree.Statement],
	portal any,
	pinfo any,
	res RestrictedCommandResult,
	canAutoCommit bool,
) (fsm.Event, fsm.EventPayload, error) {
	var ev fsm.Event
	var payload fsm.EventPayload
	var err error

	switch ex.machine.CurState().(type) {
	case stateNoTxn:
		ev, payload = ex.execStmtInNoTxnState(parserStmt)
	default:
		panic(errors.New(fmt.Sprintf("unexpected state: %s", ex.machine.CurState())))
	}

	return ev, payload, err
}
