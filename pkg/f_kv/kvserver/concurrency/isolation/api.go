package isolation

type Level int32

const (
	Serializable Level = 0
	Snapshot     Level = 1
)
