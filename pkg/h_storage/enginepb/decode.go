package enginepb

import (
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// ScanDecodeKeyValue decodes a key/value pair from a binary stream, such as in
// an MVCCScan "batch" (this is not the RocksDB batch repr format), returning
// the key/value, the timestamp, and the suffix of data remaining in the batch.
func ScanDecodeKeyValue(
	repr []byte,
) (key []byte, ts hlc.Timestamp, value []byte, orepr []byte, err error) {
	return nil, hlc.Timestamp{}, nil, nil, nil
}

// ScanDecodeKeyValueNoTS decodes a key/value pair from a binary stream, such as
// in an MVCCScan "batch" (this is not the RocksDB batch repr format), returning
// the key/value and the suffix of data remaining in the batch.
func ScanDecodeKeyValueNoTS(repr []byte) (key []byte, value []byte, orepr []byte, err error) {
	return nil, nil, nil, nil
}
