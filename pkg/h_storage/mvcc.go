package storage

import (
	"context"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// MVCCKeyValue contains the raw bytes of the value for a key.
type MVCCKeyValue struct {
	Key MVCCKey
	// if Key.IsValue(), Value is an encoded MVCCValue.
	//             else, Value is an encoded MVCCMetadata.
	Value []byte
}

// ReadWriter is the read/write interface to an engine's data.
type ReadWriter interface {
	Reader
	Writer
}

// MVCCPutProto sets the given key to the protobuf-serialized byte
// string of msg and the provided timestamp.
func MVCCPutProto(
	ctx context.Context,
	rw ReadWriter,
	key roachpb.Key,
	timestamp hlc.Timestamp,
) error {
	value := roachpb.Value{}
	err := MVCCPut(ctx, rw, key, timestamp, value)
	return err
}

// MVCCPut sets the value for a specified key. It will save the value
// with different versions according to its timestamp and update the
// key metadata. The timestamp must be passed as a parameter; using
// the Timestamp field on the value results in an error.
//
// Note that, when writing transactionally, the txn's timestamps
// dictate the timestamp of the operation, and the timestamp parameter is
// confusing and redundant. See the comment on mvccPutInternal for details.
//
// If the timestamp is specified as hlc.Timestamp{}, the value is
// inlined instead of being written as a timestamp-versioned value. A
// zero timestamp write to a key precludes a subsequent write using a
// non-zero timestamp and vice versa. Inlined values require only a
// single row and never accumulate more than a single value. Successive
// zero timestamp writes to a key replace the value and deletes clear
// the value. In addition, zero timestamp values may be merged.
func MVCCPut(
	ctx context.Context,
	rw ReadWriter,
	key roachpb.Key,
	timestamp hlc.Timestamp,
	value roachpb.Value,
) error {
	return nil
}
