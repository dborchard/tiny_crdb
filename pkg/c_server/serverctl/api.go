package serverctl

import "context"

type ServerStartupInterface interface {
	ServerShutdownInterface
	// PreStart starts the server on the specified port(s) and
	// initializes subsystems.
	// It does not activate the pgwire listener over the network / unix
	// socket, which is done by the AcceptClients() method. The separation
	// between the two exists so that SQL initialization can take place
	// before the first client is accepted.
	PreStart(ctx context.Context) error

	// AcceptInternalClients starts listening for incoming internal SQL clients over the
	// loopback interface.
	AcceptInternalClients(ctx context.Context) error

	// RunInitialSQL runs the SQL initialization for brand new clusters,
	// if the cluster is being started for the first time.
	// The arguments are:
	// - startSingleNode is used by 'demo' and 'start-single-node'.
	// - adminUser/adminPassword is used for 'demo'.
	RunInitialSQL(ctx context.Context, startSingleNode bool, adminUser, adminPassword string) error

	// AcceptClients starts listening for incoming SQL clients over the network.
	AcceptClients(ctx context.Context) error
}

// ServerShutdownInterface is the subset of the APIs on a server
// object that's sufficient to run a server shutdown.
type ServerShutdownInterface interface {
	AnnotateCtx(context.Context) context.Context
	ShutdownRequested() <-chan ShutdownRequest
}

// ShutdownRequest is used to signal a request to shutdown the server through
// server.stopTrigger. It carries the reason for the shutdown.
type ShutdownRequest struct {
	// Err is populated for reason ServerStartupError and FatalError.
	Err error
}
