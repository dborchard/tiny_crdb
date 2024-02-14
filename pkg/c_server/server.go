package server

import (
	"github.com/dborchard/tiny_crdb/pkg/c_server/serverctl"
	"github.com/dborchard/tiny_crdb/pkg/y_util/stop"
)

// topLevelServer is the cockroach server node.
type topLevelServer struct {
	// The following fields are populated in NewServer.

}

// NewServer creates a Server from a server.Config.
//
// The caller is responsible for listening on the server's ShutdownRequested()
// channel and calling stopper.Stop().
func NewServer(cfg Config, stopper *stop.Stopper) (serverctl.ServerStartupInterface, error) {
	lateBoundServer := &topLevelServer{}
	return lateBoundServer, nil
}
