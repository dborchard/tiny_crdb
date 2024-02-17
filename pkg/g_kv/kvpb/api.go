package kvpb

type BatchRequest struct {
	Requests interface{}
}

type BatchResponse struct {
}

type Error struct {
	Error error
}

type InternalClient struct {
}

type RequestUnion struct {
}
