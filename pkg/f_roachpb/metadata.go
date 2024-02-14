package roachpb

// NodeID is a custom type for a cockroach node ID. (not a raft node ID)
// 0 is not a valid NodeID.
type NodeID int32

// StoreID is a custom type for a cockroach store ID.
type StoreID int32

type NodeDescriptor struct {
}

type StoreDescriptor struct {
}
