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

// Stack represents a stack of SessionData objects.
// This is used to support transaction-scoped variables, where SET LOCAL only
// affects the top of the stack.
// There is always guaranteed to be one element in the stack.
type Stack struct {
	// Use an internal variable to prevent abstraction leakage.
	stack []*SessionData
	// base is a pointer to the first element of the stack.
	// This avoids a race with stack being reassigned, as the first element
	// is *always* set.
	base *SessionData
}

// Top returns the top element of the stack.
func (s *Stack) Top() *SessionData {
	if len(s.stack) == 0 {
		return nil
	}
	return s.stack[len(s.stack)-1]
}
