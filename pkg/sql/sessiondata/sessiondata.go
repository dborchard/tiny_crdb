package sessiondata

import "time"

// SessionData contains session parameters. They are all user-configurable.
// A SQL Session changes fields in SessionData through sql.sessionDataMutator.
type SessionData struct {
	// All session parameters below must be propagated to the remote nodes but
	// are not easily serializable. They require custom serialization
	// (MarshalNonLocal) and deserialization (UnmarshalNonLocal).
	//
	// Location indicates the current time zone.
	Location       *time.Location
	DefaultIntSize int
}
