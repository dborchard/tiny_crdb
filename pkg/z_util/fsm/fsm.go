package fsm

// Event is something that happens to a Machine which may or may not trigger a
// state transition.
type Event interface {
	Event()
}

type EventPayload interface{}

type State interface {
	State()
}

type Machine struct {
	cur State
}

func (m *Machine) CurState() State {
	return m.cur
}

// Bool represents a boolean pattern.
type Bool interface {
	bool()
	// Get returns the value of a Bool.
	Get() bool
}

type b bool

func (b b) bool() {
	//TODO implement me
	panic("implement me")
}

func (b b) Get() bool {
	//TODO implement me
	panic("implement me")
}

var (
	// True is a pattern that matches true booleans.
	True Bool = b(true)
	// False is a pattern that matches false booleans.
	False Bool = b(false)
)
