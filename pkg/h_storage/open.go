package storage

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/h_storage/vfs"
)

// A Location describes where the storage engine's data will be written. A
// Location may be in-memory or on the filesystem.
type Location struct {
	dir string
	fs  vfs.FS
}

type engineConfig struct {
	Dir string
	FS  vfs.FS
}

// Open opens a new Pebble storage engine, reading and writing data to the
// provided Location, configured with the provided options.
func Open(ctx context.Context, loc Location) (*Pebble, error) {
	var cfg engineConfig
	cfg.Dir = loc.dir
	cfg.FS = loc.fs
	p, err := NewPebble(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}
