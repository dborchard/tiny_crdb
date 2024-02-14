package server

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/c_server/authserver"
	"github.com/dborchard/tiny_crdb/pkg/c_server/serverctl"
	sql "github.com/dborchard/tiny_crdb/pkg/v_sql"
	kv "github.com/dborchard/tiny_crdb/pkg/w_kv"
	"github.com/dborchard/tiny_crdb/pkg/w_kv/kvclient/kvcoord"
	"github.com/dborchard/tiny_crdb/pkg/y_util/hlc"
	"github.com/dborchard/tiny_crdb/pkg/y_util/stop"
	"net"
	"time"
)

// topLevelServer is the cockroach server node.
type topLevelServer struct {
	stopper   *stop.Stopper
	cfg       Config
	clock     *hlc.Clock
	db        *kv.DB
	http      *httpServer
	sqlServer *SQLServer
	startTime time.Time

	engines        Engines
	pgL            net.Listener
	authentication authserver.Server
}

func (s *topLevelServer) ShutdownRequested() <-chan serverctl.ShutdownRequest {
	//TODO implement me
	panic("implement me")
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
	if err := s.http.setupRoutes(ctx, s.authentication /* authnServer */); err != nil {
		return err
	}

	return nil
}

func (s *topLevelServer) AcceptInternalClients(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *topLevelServer) RunInitialSQL(ctx context.Context, startSingleNode bool, adminUser, adminPassword string) error {
	//TODO implement me
	panic("implement me")
}

func (s *topLevelServer) AcceptClients(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// NewServer creates a Server from a server.Config.
//
// The caller is responsible for listening on the server's ShutdownRequested()
// channel and calling stopper.Stop().
func NewServer(cfg Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error) {
	ctx := context.Background()

	dbCtx := kv.DefaultDBContext(stopper)
	clock, err := newClockFromConfig()

	distSender := kvcoord.NewDistSender()
	tcsFactory := kvcoord.NewTxnCoordSenderFactory(distSender)

	db := kv.NewDBWithContext(tcsFactory, clock, dbCtx)
	insqlDB := sql.NewShimInternalDB(db)

	sHTTP := newHTTPServer()

	// Instantiate the SQL server proper.
	sqlServer, err := newSQLServer(ctx, sqlServerArgs{
		db:         db,
		internalDB: insqlDB,
	})
	if err != nil {
		return nil, err
	}

	lateBoundServer := &topLevelServer{
		cfg:       cfg,
		clock:     clock,
		stopper:   stopper,
		db:        db,
		sqlServer: sqlServer,
		http:      sHTTP,
		startTime: time.Now(),
	}
	return lateBoundServer, nil
}

func newClockFromConfig() (*hlc.Clock, error) {
	var clock *hlc.Clock
	return clock, nil
}
