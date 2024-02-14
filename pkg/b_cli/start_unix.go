//go:build !windows

package cli

import (
	"golang.org/x/sys/unix"
	"os"
)

// DrainSignals is the list of signals that trigger the start of a shutdown
// sequence ("server drain").
//
// The first time they're received, both signals initiate a drain just the same.
// The behavior between the two differs if they're received a second time (or,
// more generally, after the drain had started):
// - a second SIGTERM is ignored.
// - a second SIGINT terminates the process abruptly.
var DrainSignals = []os.Signal{unix.SIGINT, unix.SIGTERM}

// exitAbruptlySignal is the signal to make the process exit immediately. It is
// preferable to SIGKILL when running with coverage instrumentation because the
// coverage profile gets dumped on exit.
var exitAbruptlySignal os.Signal = unix.SIGUSR1
