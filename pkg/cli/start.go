package cli

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/server"
	"github.com/spf13/cobra"
)

type newServerFn func(ctx context.Context, serverCfg server.Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error)

// startSingleNodeCmd starts a node by initializing the stores.
var startSingleNodeCmd = &cobra.Command{
	Use:     "start-single-node",
	Short:   "start a single-node cluster",
	Example: `  cockroach start-single-node --insecure --store=attrs=ssd,path=/mnt/ssd1`,
	Args:    cobra.NoArgs,
	RunE:    runStartSingleNode,
}

func runStartSingleNode(cmd *cobra.Command, args []string) error {
	return runStart(cmd, args, true /*startSingleNode*/)
}

// runStart starts the cockroach node using --store as the list of
// storage devices ("stores") on this machine and --join as the list
// of other active nodes used to join this node to the cockroach
// cluster, if this is its first time connecting.
//
// The argument startSingleNode is morally equivalent to `cmd ==
// startSingleNodeCmd`, and triggers special initialization specific
// to one-node clusters. See server/initial_sql.go for details.
//
// We need a separate argument instead of solely relying on cmd
// because we cannot refer to startSingleNodeCmd under
// runStartInternal: there would be a cyclic dependency between
// runStart, runStartSingleNode and runStartSingleNodeCmd.
func runStart(cmd *cobra.Command, args []string, startSingleNode bool) error {

	newServerFn := func(_ context.Context, serverCfg server.Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error) {
		// Beware of not writing simply 'return server.NewServer()'. This is
		// because it would cause the serverctl.ServerStartupInterface reference to
		// always be non-nil, even if NewServer returns a nil pointer (and
		// an error). The code below is dependent on the interface
		// reference remaining nil in case of error.
		s, err := server.NewServer(serverCfg, stopper)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return runStartInternal(cmd, serverType, serverCfg.InitNode, newServerFn, startSingleNode)
}

// runStartInternal contains the code common to start a regular server
// or a SQL-only server.
func runStartInternal(
	cmd *cobra.Command,
	initConfigFn func(context.Context) error,
	newServerFn newServerFn,
	startSingleNode bool,
) error {
	// Beyond this point, the configuration is set and the server is
	// ready to start.

	// Run the rest of the startup process in a goroutine separate from
	// the main goroutine to avoid preventing proper handling of signals
	// if we get stuck on something during initialization (#10138).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srvStatus, serverShutdownReqC := createAndStartServerAsync(ctx, newServerFn, startSingleNode)

	return waitForShutdown(
		// NB: we delay the access to s, as it is assigned
		// asynchronously in a goroutine above.
		stopper, serverShutdownReqC, signalCh,
		srvStatus)
}

// createAndStartServerAsync starts an async goroutine which instantiates
// the server and starts it.
// We run it in a separate goroutine because the instantiation&start
// could block, and we want to retain the option to start shutting down
// the process (e.g. via Ctrl+C on the terminal) even in that case.
// The shutdown logic thus starts running asynchronously, via waitForShutdown,
// concurrently with createAndStartServerAsync.
//
// The arguments are as follows:
//   - tBegin: time when startup began; used to report statistics at the end of startup.
//   - serverCfg: the server configuration.
//   - stopper: the stopper used to start all the async tasks. This is the stopper
//     used by the shutdown logic.
//   - startupSpan: the tracing span for the context that was started earlier
//     during startup. It needs to be finalized when the async goroutine completes.
//   - newServerFn: a constructor function for the server object.
//   - serverType: a title used for the type of server. This is used
//     when reporting the startup messages on the terminal & logs.
func createAndStartServerAsync(
	ctx context.Context,
	newServerFn newServerFn,
	startSingleNode bool,
) (srvStatus *serverStatus, serverShutdownReqC <-chan serverctl.ShutdownRequest) {
	var serverStatusMu serverStatus
	var s serverctl.ServerStartupInterface
	shutdownReqC := make(chan serverctl.ShutdownRequest, 1)

	go func() {
		// Any error beyond this point is reported through shutdownReqC.
		if err := func() error {
			// Instantiate the server.
			var err error
			s, err = newServerFn(ctx, *serverCfg, stopper)
			if err != nil {
				return errors.Wrap(err, "failed to start server")
			}

			// Have we already received a signal to terminate? If so, just
			// stop here.
			if serverStatusMu.shutdownInProgress() {
				return nil
			}

			// Attempt to start the server.
			if err := s.PreStart(ctx); err != nil {
				if le := (*server.ListenError)(nil); errors.As(err, &le) {
					const errorPrefix = "consider changing the port via --%s"
					if le.Addr == serverCfg.Addr {
						err = errors.Wrapf(err, errorPrefix, cliflags.ListenAddr.Name)
					} else if le.Addr == serverCfg.HTTPAddr {
						err = errors.Wrapf(err, errorPrefix, cliflags.ListenHTTPAddr.Name)
					}
				}

				return errors.Wrap(err, "cockroach server exited with error")
			}
			// Server started, notify the shutdown monitor running concurrently.
			if shutdownInProgress := serverStatusMu.setStarted(s, stopper); shutdownInProgress {
				// A shutdown was requested already, e.g. by sending SIGTERM to the process:
				// maybeWaitForShutdown (which runs concurrently with this goroutine) has
				// called serverStatusMu.startShutdown() already.
				// However, because setStarted() had not been called before,
				// maybeWaitForShutdown did not call Stop on the stopper.
				// So we do it here.
				stopper.Stop(ctx)
				return nil
			}
			// After this point, if a shutdown is requested concurrently
			// with the startup steps below, the stopper.Stop() method will
			// be called by the shutdown goroutine, which in turn will cause
			// all these startup steps to fail. So we do not need to look at
			// the "shutdown status" in serverStatusMu any more.

			// Accept internal clients early, as RunInitialSQL might need it.
			if err := s.AcceptInternalClients(ctx); err != nil {
				return err
			}

			// Run one-off cluster initialization.
			if err := s.RunInitialSQL(ctx, startSingleNode, "" /* adminUser */, "" /* adminPassword */); err != nil {
				return err
			}

			// Now let SQL clients in.
			if err := s.AcceptClients(ctx); err != nil {
				return err
			}

			// Now inform the user that the server is running and tell the
			// user about its run-time derived parameters.
			return reportServerInfo(ctx, tBegin, serverCfg, s.ClusterSettings(),
				serverType, s.InitialStart(), s.LogicalClusterID())
		}(); err != nil {
			shutdownReqC <- serverctl.MakeShutdownRequest(serverctl.ShutdownReasonServerStartupError, errors.Wrapf(err, "server startup failed"))
		} else {
			// Start a goroutine that watches for shutdown requests and notifies
			// errChan.
			go func() {
				select {
				case req := <-s.ShutdownRequested():
					shutdownCtx := s.AnnotateCtx(context.Background())
					log.Infof(shutdownCtx, "server requesting spontaneous shutdown: %v", req.ShutdownCause())
					shutdownReqC <- req
				case <-stopper.ShouldQuiesce():
				}
			}()
		}

	}()
	serverShutdownReqC = shutdownReqC
	srvStatus = &serverStatusMu
	return srvStatus, serverShutdownReqC
}

// serverStatus coordinates the async goroutine that starts the server
// up (e.g. in runStart) and the async goroutine that stops the server
// (in waitForShutdown).
//
// We need this intermediate coordination because it isn't safe to try
// to drain a server that doesn't exist or is in the middle of
// starting up, or to start a server after shutdown has begun.
type serverStatus struct {
}
