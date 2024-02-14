package storage

// MVCCKeyValue contains the raw bytes of the value for a key.
type MVCCKeyValue struct {
	Key MVCCKey
	// if Key.IsValue(), Value is an encoded MVCCValue.
	//             else, Value is an encoded MVCCMetadata.
	Value []byte
}
