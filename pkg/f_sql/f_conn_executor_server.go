package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
)

type Server struct {
	cfg *ExecutorConfig
}

func (s *Server) newConnExecutor(
	ctx context.Context,
	stmtBuf *StmtBuf,
	clientComm ClientComm,
	postSetupFn func(ex *connExecutor),
) *connExecutor {

	ex := &connExecutor{
		server:       s,
		stmtBuf:      stmtBuf,
		clientComm:   clientComm,
		executorType: executorTypeExec,
	}

	if postSetupFn != nil {
		postSetupFn(ex)
	}

	ex.initPlanner(ctx, &ex.planner)
	return ex
}

func (s *Server) Start(ctx context.Context, stopper *stop.Stopper) {

}
