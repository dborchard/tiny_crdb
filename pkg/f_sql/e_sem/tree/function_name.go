package tree

import (
	"context"
	"github.com/lib/pq/oid"
	"strings"
)

type FunctionReferenceResolver interface {
	ResolveFunction(ctx context.Context, name UnresolvedRoutineName, path SearchPath) (*ResolvedFunctionDefinition, error)
	ResolveFunctionByOID(ctx context.Context, oid oid.Oid) (*RoutineName, *Overload, error)
}

// ResolvableFunctionReference implements the editable reference call of a
// FuncExpr.
type ResolvableFunctionReference struct {
	// ReferenceByName keeps track of the name that was used to resolve the
	// function, if one was used. This is used for metadata dependency tracking.
	ReferenceByName *UnresolvedObjectName
	FunctionReference
}

// Resolve converts a ResolvableFunctionReference into a *FunctionDefinition. If
// the reference has already been resolved, it simply returns the definition. If
// a FunctionReferenceResolver is provided, it will be used to resolve the
// function definition. Otherwise, the default resolution of
// UnresolvedName.ResolveFunction is used.
func (ref *ResolvableFunctionReference) Resolve(
	ctx context.Context, path SearchPath, resolver FunctionReferenceResolver,
) (*ResolvedFunctionDefinition, error) {
	switch t := ref.FunctionReference.(type) {
	case *ResolvedFunctionDefinition:
		return t, nil
	case *FunctionDefinition:
		// TODO(Chengxiong): get rid of FunctionDefinition entirely.
		parts := strings.Split(t.Name, ".")
		if len(parts) > 2 {
			// In theory, this should not happen since all builtin functions are
			// defined within virtual schema and don't belong to any database catalog.
			return nil, errors.AssertionFailedf("invalid builtin function name: %q", t.Name)
		}
		fullName := t.Name
		if len(parts) == 1 {
			fullName = catconstants.PgCatalogName + "." + t.Name
		}
		fd := ResolvedBuiltinFuncDefs[fullName]
		ref.FunctionReference = fd
		return fd, nil
	case *UnresolvedName:
		if resolver == nil {
			// If a resolver is not provided, just try to fetch a builtin function.
			fn, err := t.ToRoutineName()
			if err != nil {
				return nil, err
			}
			fd, err := GetBuiltinFuncDefinitionOrFail(fn, path)
			if err != nil {
				return nil, err
			}
			ref.FunctionReference = fd
			return fd, nil
		}
		// Use the resolver if it is provided.
		fd, err := resolver.ResolveFunction(ctx, MakeUnresolvedFunctionName(t), path)
		if err != nil {
			return nil, err
		}
		referenceByName, _ := t.ToUnresolvedObjectName(NoAnnotation)
		ref.ReferenceByName = &referenceByName
		ref.FunctionReference = fd
		return fd, nil
	case *FunctionOID:
		if resolver == nil {
			return GetBuiltinFunctionByOIDOrFail(t.OID)
		}
		fnName, o, err := resolver.ResolveFunctionByOID(ctx, t.OID)
		if err != nil {
			return nil, err
		}
		fd := &ResolvedFunctionDefinition{
			Name:      fnName.Object(),
			Overloads: []QualifiedOverload{{Schema: fnName.Schema(), Overload: o}},
		}
		ref.FunctionReference = fd
		return fd, nil
	default:
		return nil, errors.AssertionFailedf("unknown resolvable function reference type %s", t)
	}
}
