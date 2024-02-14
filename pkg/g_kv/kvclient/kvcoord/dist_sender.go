package kvcoord

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvpb"
)

// A DistSender provides methods to access Cockroach's monolithic,
// distributed key value store. Each method invocation triggers a
// lookup or lookups to find replica metadata for implicated key
// ranges. RPCs are sent to one or more of the replicas to satisfy
// the method invocation.
type DistSender struct {
}

// NewDistSender returns a batch.Sender instance which connects to the
// Cockroach cluster via the supplied gossip instance. Supplying a
// DistSenderContext or the fields within is optional. For omitted values, sane
// defaults will be used.
func NewDistSender() *DistSender {
	ds := &DistSender{}
	return ds
}

// Send implements the batch.Sender interface. It subdivides the Batch
// into batches admissible for sending (preventing certain illegal
// mixtures of requests), executes each individual part (which may
// span multiple ranges), and recombines the response.
//
// When the request spans ranges, it is split by range and a partial
// subset of the batch request is sent to affected ranges in parallel.
func (ds *DistSender) Send(
	ctx context.Context, ba *kvpb.BatchRequest,
) (*kvpb.BatchResponse, *kvpb.Error) {
	//TODO implement me
	panic("implement me")
}
