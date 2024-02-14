package kvtenant

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvpb"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"sync"
)

type connectorFactory struct {
}

// NewConnector creates a new loopback only tenant connector with the given
// configuration. A loopback address is required, and an error is returned if
// one is not provided. If any remote addresses are provided, an error is
// returned.
func (connectorFactory) NewConnector(
	cfg ConnectorConfig, addressConfig KVAddressConfig,
) (Connector, error) {
	return NewConnector(cfg, []string{addressConfig.LoopbackAddress}), nil
}

// Connector mediates the communication of cluster-wide state to sandboxed
// SQL-only tenant servers through a restricted interface. A Connector is
// seeded with a set of one or more network addresses that reference existing
// KV nodes in the cluster (or a load-balancer which fans out to some/all KV
// nodes). On startup, it establishes contact with one of these nodes to learn
// about the topology of the cluster and bootstrap the rest of SQL <-> KV
// network communication.
type Connector interface {
	// Start starts the connector.
	Start(context.Context) error
}

// connector mediates the communication of cluster-wide state to sandboxed
// SQL-only tenant servers through a restricted interface.
//
// A connector is instantiated inside a tenant's SQL server and is seeded with a
// set of one or more network addresses that reference existing KV nodes in the
// storage cluster (or a load-balancer which fans out to some/all KV nodes). On
// startup, it establishes contact with one of these nodes to learn about the
// topology of the cluster and bootstrap the rest of SQL <-> KV network
// communication.
//
// The connector communicates with the storage cluster through the
// roachpb.Internal API.
//
// See below for the connector's roles.
type connector struct {
	addrs                              []string
	earlyShutdownIfMissingTenantRecord bool

	mu struct {
		sync.RWMutex
		client               *client
		nodeDescs            map[roachpb.NodeID]*roachpb.NodeDescriptor
		storeDescs           map[roachpb.StoreID]*roachpb.StoreDescriptor
		systemConfigChannels map[chan<- struct{}]struct{}
	}
}

// NewConnector creates a new connector.
// NOTE: Calling Start will set cfg.RPCContext.ClusterID.
func NewConnector(cfg ConnectorConfig, addrs []string) Connector {
	c := &connector{
		addrs:                              addrs,
		earlyShutdownIfMissingTenantRecord: cfg.ShutdownTenantConnectorEarlyIfNoRecordPresent,
	}

	c.mu.nodeDescs = make(map[roachpb.NodeID]*roachpb.NodeDescriptor)
	c.mu.storeDescs = make(map[roachpb.StoreID]*roachpb.StoreDescriptor)
	c.mu.systemConfigChannels = make(map[chan<- struct{}]struct{})
	return c
}

// Start launches the connector's worker thread and waits for it to successfully
// connect to a KV node. Start returns once the connector has determined the
// cluster's ID and set connector.rpcContext.ClusterID.
func (c *connector) Start(ctx context.Context) error {
	return nil
}

// client represents an RPC client that proxies to a KV instance.
type client struct {
	kvpb.InternalClient
	//serverpb.StatusClient
	//serverpb.AdminClient
	//tspb.TimeSeriesClient
}
