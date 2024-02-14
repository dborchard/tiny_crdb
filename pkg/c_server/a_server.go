package server

import (
	"context"
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/c_server/authserver"
	"github.com/dborchard/tiny_crdb/pkg/c_server/serverctl"
	kv "github.com/dborchard/tiny_crdb/pkg/f_kv"
	"github.com/dborchard/tiny_crdb/pkg/f_kv/kvclient/kvcoord"
	sql "github.com/dborchard/tiny_crdb/pkg/i_sql"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
	"github.com/dborchard/tiny_crdb/pkg/z_util/netutil"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
	"net"
	"time"
)

// topLevelServer is the cockroach server node.
type topLevelServer struct {
	stopper        *stop.Stopper
	cfg            Config
	clock          *hlc.Clock
	db             *kv.DB
	http           *httpServer
	sqlServer      *SQLServer
	startTime      time.Time
	engines        Engines
	authentication authserver.Server

	pgL         net.Listener
	loopbackPgL *netutil.LoopbackListener
}

func (s *topLevelServer) PreStart(ctx context.Context) error {
	//// Start the RPC server. This opens the RPC/SQL listen socket,
	//// and dispatches the server worker for the RPC.
	//// The SQL listener is returned, to start the SQL server later
	//// below when the server has initialized.
	//pgL, loopbackPgL, rpcLoopbackDialFn, startRPCServer, err := startListenRPCAndSQL()
	//if err != nil {
	//	return err
	//}
	//s.pgL = pgL
	//s.loopbackPgL = loopbackPgL

	// Connect the HTTP endpoints. This also wraps the privileged HTTP
	// endpoints served by gwMux by the HTTP cookie authentication
	// check.
	if err := s.http.setupRoutes(ctx, s.authentication); err != nil {
		return err
	}

	return nil
}

func (s *topLevelServer) AcceptInternalClients(ctx context.Context) error {
	return nil
}

func (s *topLevelServer) RunInitialSQL(ctx context.Context, startSingleNode bool, adminUser, adminPassword string) error {
	if adminUser != "" {
		if err := s.createAdminUser(ctx, adminUser, adminPassword); err != nil {
			return err
		}
	}
	return nil
}

// createAdminUser creates an admin user with the given name.
func (s *topLevelServer) createAdminUser(
	ctx context.Context, adminUser, adminPassword string,
) error {
	ie := s.sqlServer.internalExecutor
	_, err := ie.Exec(
		ctx, "admin-user", nil,
		fmt.Sprintf("CREATE USER %s WITH PASSWORD $1", adminUser),
		adminPassword,
	)
	if err != nil {
		return err
	}
	_, err = ie.Exec(ctx, "admin-user", nil, fmt.Sprintf("GRANT admin TO %s", adminUser))
	return err
}

func (s *topLevelServer) AcceptClients(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *topLevelServer) ShutdownRequested() <-chan serverctl.ShutdownRequest {
	//TODO implement me
	panic("implement me")
}

// NewServer creates a Server from a server.Config.
//
// The caller is responsible for listening on the server's ShutdownRequested()
// channel and calling stopper.Stop().
func NewServer(cfg Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error) {
	ctx := context.Background()

	clock, err := newClockFromConfig()

	engines, err := cfg.CreateEngines(ctx)

	// The Executor will be further initialized later, as we create more
	// of the server's components. There's a circular dependency - many things
	// need an Executor, but the Executor needs an executorConfig,
	// which in turn needs many things. That's why everybody that needs an
	// Executor uses this one instance.
	internalExecutor := &sql.InternalExecutor{}

	_dbCtx := kv.DefaultDBContext(stopper)
	_distSender := kvcoord.NewDistSender()
	_tcsFactory := kvcoord.NewTxnCoordSenderFactory(_distSender)
	db := kv.NewDBWithContext(_tcsFactory, clock, _dbCtx)
	insqlDB := sql.NewShimInternalDB(db)
	sqlServer, err := newSQLServer(ctx, sqlServerArgs{
		db:                       db,
		internalDB:               insqlDB,
		circularInternalExecutor: internalExecutor,
	})
	if err != nil {
		return nil, err
	}

	sAuth := authserver.NewServer(sqlServer)
	sHTTP := newHTTPServer()

	lateBoundServer := &topLevelServer{
		cfg:            cfg,
		clock:          clock,
		stopper:        stopper,
		db:             db,
		sqlServer:      sqlServer,
		http:           sHTTP,
		startTime:      time.Now(),
		engines:        engines,
		authentication: sAuth,
	}
	return lateBoundServer, nil
}

func newClockFromConfig() (*hlc.Clock, error) {
	var clock *hlc.Clock
	return clock, nil
}
