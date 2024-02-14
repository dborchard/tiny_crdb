package storage

import (
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
)

// MVCCKey is a versioned key, distinguished from roachpb.Key with the addition
// of a "version" timestamp.
//
// The version timestamp dictates the key's visibility to readers. Readers with
// read timestamps equal to or greater than the version timestamp observe the
// key. Readers with read timestamps below the version timestamp ignore the key.
// Keys are stored in decreasing version order, with the exception of version
// zero (timestamp 0), which is referred to as a "meta" version and is stored
// before all other versions of the same key.
type MVCCKey struct {
	Key       roachpb.Key
	Timestamp hlc.Timestamp
}

// MVCCRangeKeyStack represents a stack of range key fragments as returned
// by SimpleMVCCIterator.RangeKeys(). All fragments have the same key bounds,
// and are ordered from newest to oldest.
type MVCCRangeKeyStack struct {
	//Bounds   roachpb.Span
	Versions MVCCRangeKeyVersions
}

// MVCCRangeKeyVersions represents a stack of range key fragment versions.
type MVCCRangeKeyVersions []MVCCRangeKeyVersion

// MVCCRangeKeyVersion represents a single range key fragment version.
type MVCCRangeKeyVersion struct {
	Timestamp hlc.Timestamp
	Value     []byte
}
