package storage

import roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"

// MVCCValue is a versioned value, stored at an associated MVCCKey with a
// non-zero version timestamp.
//
// MVCCValue wraps a roachpb.Value and extends it with MVCC-level metadata which
// is stored in an enginepb.MVCCValueHeader struct.
//
// The mvcc value has a "simple" and an "extended" encoding scheme, depending on
// whether the value's header is empty or not. If the value's header is empty,
// it is omitted in the encoding and the mvcc value's encoding is identical to
// that of roachpb.Value. This provided backwards compatibility and ensures that
// the MVCCValue optimizes away in the common case. If the value's header is not
// empty, it is prepended to the roachpb.Value encoding. The encoding scheme's
// variants are:
//
// Simple (identical to the roachpb.Value encoding):
//
//	<4-byte-checksum><1-byte-tag><encoded-data>
//
// Extended (header prepended to roachpb.Value encoding):
//
//	<4-byte-header-len><1-byte-sentinel><mvcc-header><4-byte-checksum><1-byte-tag><encoded-data>
//
// The two encoding scheme variants are distinguished using the 5th byte, which
// is either the roachpb.Value tag (which has many values) or a sentinel tag not
// used by the roachpb.Value encoding which indicates the extended encoding
// scheme.
//
// For a deletion tombstone, the encoding of roachpb.Value is special cased to
// be empty, i.e., no checksum, tag, or encoded-data. In that case the extended
// encoding above is simply:
//
//	<4-byte-header-len><1-byte-sentinel><mvcc-header>
//
// To identify a deletion tombstone from an encoded MVCCValue, callers should
// decode the value using DecodeMVCCValue and then use the IsTombstone method.
// For example:
//
//	valRaw := iter.UnsafeValue()
//	val, err := DecodeMVCCValue(valRaw)
//	if err != nil { ... }
//	isTombstone := val.IsTombstone()
type MVCCValue struct {
	//enginepb.MVCCValueHeader
	Value roachpb.Value
}

// IsTombstone returns whether the MVCCValue represents a deletion tombstone.
func (v MVCCValue) IsTombstone() bool {
	return len(v.Value.RawBytes) == 0
}
