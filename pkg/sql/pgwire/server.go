package pgwire

import (
	"github.com/dborchard/tiny_crdb/pkg/sql"
)

// Server implements the server side of the PostgreSQL wire protocol for one
// specific tenant (i.e. its configuration is specific to one tenant).
type Server struct {
	SQLServer *sql.Server
}

// MakeServer creates a Server.
//
// Start() needs to be called on the Server so it begins processing.
func MakeServer(
	//ambientCtx log.AmbientContext,
	//cfg *base.Config,
	//st *cluster.Settings,
	//sqlMemMetrics sql.MemoryMetrics,
	//parentMemoryMonitor *mon.BytesMonitor,
	//histogramWindow time.Duration,
	executorConfig *sql.ExecutorConfig,
) *Server {
	server := &Server{
		//AmbientCtx: ambientCtx,
		//cfg:        cfg,
		//execCfg:    executorConfig,

		//tenantMetrics: newTenantSpecificMetrics(sqlMemMetrics, histogramWindow),
	}
	return server
}
