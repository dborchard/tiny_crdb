package server

import "github.com/dborchard/tiny_crdb/pkg/storage"

// Config holds the parameters needed to set up a combined KV and SQL server.
type Config struct {
}

// Engines is a container of engines, allowing convenient closing.
type Engines []storage.Engine
