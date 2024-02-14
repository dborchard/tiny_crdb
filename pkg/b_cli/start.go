package cli

import (
	"context"
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/c_server"
	"github.com/dborchard/tiny_crdb/pkg/c_server/serverctl"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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
	return runStart(cmd, args, true)
}

func runStart(cmd *cobra.Command, args []string, startSingleNode bool) error {
	newServerFn := func(_ context.Context, serverCfg server.Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error) {
		s, err := server.NewServer(serverCfg, stopper)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return runStartInternal(cmd, newServerFn, startSingleNode)
}

func runStartInternal(
	cmd *cobra.Command,
	newServerFn newServerFn,
	startSingleNode bool,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up the signal handlers. This also ensures that any of these
	// signals received beyond this point do not interrupt the startup
	// sequence until the point signals are checked below.
	// We want to set up signal handling before starting logging, because
	// logging uses buffering, and we want to be able to sync
	// the buffers in the signal handler below. If we started capturing
	// signals later, some startup logging might be lost.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, DrainSignals...)
	if exitAbruptlySignal != nil {
		signal.Notify(signalCh, exitAbruptlySignal)
	}

	// Set up the logging and profiling output.
	//
	// We want to do this as early as possible, because most of the code
	// in CockroachDB may use logging, and until logging has been
	// initialized log files will be created in $TMPDIR instead of their
	// expected location.
	//
	// This initialization uses the various configuration parameters
	// initialized by flag handling (before runStart was called). Any
	// additional server configuration tweaks for the startup process
	// must be necessarily non-logging-related, as logging parameters
	// cannot be picked up beyond this point.
	stopper, err := setupAndInitializeLoggingAndProfiling(ctx, cmd, true /* isServerCmd */)
	if err != nil {
		return err
	}

	var serverCfg = func() server.Config {
		return server.Config{}
	}()
	// Beyond this point, the configuration is set and the server is
	// ready to start.

	srvStatus, serverShutdownReqC := createAndStartServerAsync(ctx, &serverCfg, stopper, newServerFn, startSingleNode)

	return waitForShutdown(stopper, serverShutdownReqC, signalCh, srvStatus)
}

func setupAndInitializeLoggingAndProfiling(ctx context.Context, cmd *cobra.Command, b bool) (stopper *stop.Stopper, err error) {
	stopper = stop.NewStopper()
	return stopper, nil
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
//   - serverCfg: the server configuration.
//   - stopper: the stopper used to start all the async tasks. This is the stopper
//     used by the shutdown logic.
//   - newServerFn: a constructor function for the server object.
func createAndStartServerAsync(
	ctx context.Context,
	serverCfg *server.Config,
	stopper *stop.Stopper,
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
				return err
			}

			// Attempt to start the server.
			if err := s.PreStart(ctx); err != nil {
				return err
			}

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
			return reportServerInfo(ctx, serverCfg)
		}(); err != nil {
			shutdownReqC <- serverctl.MakeShutdownRequest(err)
		} else {
			// Start a goroutine that watches for shutdown requests and notifies
			// errChan.
			go func() {
				select {
				case req := <-s.ShutdownRequested():
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

// reportServerInfo prints out the server version and network details
// in a standardized format.
func reportServerInfo(
	ctx context.Context,
	serverCfg *server.Config,
) error {
	fmt.Println("CockroachDB node started")
	return nil
}

func waitForShutdown(
	stopper *stop.Stopper,
	shutdownC <-chan serverctl.ShutdownRequest,
	signalCh <-chan os.Signal,
	serverStatusMu *serverStatus,
) (returnErr error) {
	select {
	case shutdownRequest := <-shutdownC:
		panic(shutdownRequest)
	case sig := <-signalCh:
		panic(sig)
	}
}
