package kvstorage

import (
	"context"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	storage "github.com/dborchard/tiny_crdb/pkg/h_storage"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// InitEngine writes a new store ident to the underlying engine. To
// ensure that no crufty data already exists in the engine, it scans
// the engine contents before writing the new store ident. The engine
// should be completely empty save for a cluster version, which must
// already have been persisted to it. Returns an error if this is not
// the case.
func InitEngine(ctx context.Context, eng storage.Engine, ident roachpb.StoreIdent) error {
	batch := eng.NewBatch()
	if err := storage.MVCCPutProto(
		ctx,
		batch,
		roachpb.Key{},
		hlc.Timestamp{},
	); err != nil {
		batch.Close()
		return err
	}
	if err := batch.Commit(true /* sync */); err != nil {
		return err
	}

	return nil
}
