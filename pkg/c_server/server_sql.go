package server

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/d_upgrade/upgrademanager"
	"github.com/dborchard/tiny_crdb/pkg/sql"
	"github.com/dborchard/tiny_crdb/pkg/sql/pgwire"
)

// SQLServer encapsulates the part of a CRDB server that is dedicated to SQL
// processing. All SQL commands are reduced to primitive operations on the
// lower-level KV layer. Multi-tenant installations of CRDB run zero or more
// standalone SQLServer instances per tenant (the KV layer is shared across all
// tenants).
type SQLServer struct {

	// upgradeManager deals with cluster version upgrades on bootstrap and on
	// `set cluster setting version = <v>`.
	upgradeManager *upgrademanager.Manager
	execCfg        *sql.ExecutorConfig
	internalDB     *sql.InternalDB
}

type sqlServerArgs struct {
	internalDB *sql.InternalDB
}

// newSQLServer constructs a new SQLServer. The caller is responsible for
// listening to the server's serverctl.ShutdownRequested() channel (which is the same as
// cfg.stopTrigger.C()) and stopping cfg.stopper when signaled.
func newSQLServer(ctx context.Context, cfg sqlServerArgs) (*SQLServer, error) {
	execCfg := &sql.ExecutorConfig{}

	// Initialize the pgwire server which handles connections
	// established via the pgPreServer.
	pgServer := pgwire.MakeServer(
		//cfg.AmbientCtx,
		//cfg.Config,
		//cfg.Settings,
		//sqlMemMetrics,
		//rootSQLMemoryMonitor,
		//cfg.HistogramWindowInterval(),
		execCfg,
	)

	internalDB := sql.NewInternalDB(
		pgServer.SQLServer,
	)

	*cfg.internalDB = *internalDB
	execCfg.InternalDB = internalDB

	return &SQLServer{
		execCfg:    execCfg,
		internalDB: cfg.internalDB,
	}, nil
}