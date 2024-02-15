package netutil

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
	"net"
	"sync"
)

// TCPServer is wrapper around a map of active connections.
type TCPServer struct {
	mu          sync.Mutex
	stopper     *stop.Stopper
	activeConns map[net.Conn]struct{}
}

func MakeTCPServer(ctx context.Context, stopper *stop.Stopper) *TCPServer {
	server := &TCPServer{
		stopper:     stopper,
		activeConns: make(map[net.Conn]struct{}),
	}

	return server
}
