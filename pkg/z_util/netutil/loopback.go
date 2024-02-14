package netutil

import (
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
	"net"
	"sync"
)

// LoopbackListener implements a local listener
// that delivers net.Conns via its Connect() method
// based on the other side calls to its Accept() method.
type LoopbackListener struct {
	stopper *stop.Stopper

	closeOnce sync.Once
	active    chan struct{}

	// requests are tokens from the Connect() method to the
	// Accept() method.
	requests chan struct{}
	// conns are responses from the Accept() method
	// to the Connect() method.
	conns chan net.Conn
}

var _ net.Listener = (*LoopbackListener)(nil)

// ErrLocalListenerClosed is returned when the listener
// is shutting down.
// note that we need to use cmux.ErrListenerClosed as base (leaf)
// error so that it is recognized as special case in
// netutil.IsClosedConnection.
var ErrLocalListenerClosed = errors.New("loopback listener")

// Accept waits for and returns the next connection to the listener.
func (l *LoopbackListener) Accept() (conn net.Conn, err error) {
	select {
	case <-l.stopper.ShouldQuiesce():
		return nil, ErrLocalListenerClosed
	case <-l.active:
		return nil, ErrLocalListenerClosed
	case <-l.requests:
	}
	c1, c2 := net.Pipe()
	select {
	case l.conns <- c1:
		return c2, nil
	case <-l.stopper.ShouldQuiesce():
	case <-l.active:
	}
	err = ErrLocalListenerClosed
	return nil, err
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *LoopbackListener) Close() error {
	l.closeOnce.Do(func() {
		close(l.active)
	})
	return nil
}

// Addr returns the listener's network address.
func (l *LoopbackListener) Addr() net.Addr {
	return loopbackAddr{}
}

type loopbackAddr struct{}

var _ net.Addr = loopbackAddr{}

func (loopbackAddr) Network() string { return "pipe" }
func (loopbackAddr) String() string  { return "loopback" }
