package pgwire

import (
	"context"
	sql "github.com/dborchard/tiny_crdb/pkg/f_sql"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
	"net"
)

// Server implements the server side of the PostgreSQL wire protocol for one
// specific tenant (i.e. its configuration is specific to one tenant).
type Server struct {
	SQLServer *sql.Server
	execCfg   *sql.ExecutorConfig
}

// MakeServer creates a Server.
//
// Start() needs to be called on the Server so it begins processing.
func MakeServer(
	executorConfig *sql.ExecutorConfig,
) *Server {
	server := &Server{
		execCfg: executorConfig,
	}
	return server
}

// Start makes the Server ready for serving connections.
func (s *Server) Start(ctx context.Context, stopper *stop.Stopper) {
	s.SQLServer.Start(ctx, stopper)
}

func (s *Server) ServeConn(ctx context.Context, conn net.Conn) (err error) {
	return nil
}
