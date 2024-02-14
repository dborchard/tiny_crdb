package server

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/h_storage"
)

// Config holds the parameters needed to set up a combined KV and SQL server.
type Config struct {
}

// Engines is a container of engines, allowing convenient closing.
type Engines []storage.Engine

// CreateEngines creates Engines based on the specs in cfg.Stores.
func (cfg *Config) CreateEngines(ctx context.Context) (Engines, error) {
	var location storage.Location
	var engines Engines
	eng, _ := storage.Open(ctx, location)
	engines = append(engines, eng)
	return engines, nil
}
