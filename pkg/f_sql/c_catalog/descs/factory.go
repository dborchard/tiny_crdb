package descs

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/c_server/serverpb"
	isql "github.com/dborchard/tiny_crdb/pkg/f_sql/d_isql"
)

type Txn interface {
	isql.Txn
	Descriptors() *Collection
	Regions() RegionProvider
}

// RegionProvider abstracts the lookup of regions. It is used to implement
// crdb_internal.regions, which ultimately drives `SHOW REGIONS` and the
// logic in the commands to manipulate multi-region features.
type RegionProvider interface {
	// GetRegions provides access to the set of regions available to the
	// current tenant.
	GetRegions(ctx context.Context) (*serverpb.RegionsResponse, error)
}
