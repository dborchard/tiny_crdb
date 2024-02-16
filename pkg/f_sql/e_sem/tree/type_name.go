package tree

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	"github.com/lib/pq/oid"
)

// TypeReferenceResolver is the interface that will provide the ability
// to actually look up type metadata and transform references into
// *types.T's. Implementers of TypeReferenceResolver should also implement
// descpb.TypeDescriptorResolver is sqlbase.TypeDescriptorInterface is the
// underlying representation of a user defined type.
type TypeReferenceResolver interface {
	ResolveType(ctx context.Context, name *UnresolvedObjectName) (*types.T, error)
	ResolveTypeByOID(ctx context.Context, oid oid.Oid) (*types.T, error)
}
