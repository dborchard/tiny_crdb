package serverutils

import (
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
	"github.com/dborchard/tiny_crdb/pkg/z_util/stop"
)

type TestServerInterface interface {
	// Stopper returns the stopper used by the server.
	// TODO(knz): replace uses by Stop().
	Stopper() *stop.Stopper

	// ExecutorConfig returns a copy of the server's ExecutorConfig.
	// The real return type is sql.ExecutorConfig.
	ExecutorConfig() interface{}
	DistSenderI() interface{}
	Clock() *hlc.Clock
}
