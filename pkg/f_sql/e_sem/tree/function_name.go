package tree

import (
	"context"
	"errors"
	"fmt"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/catconstants"
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

// FunctionReference is the common interface to UnresolvedName and QualifiedFunctionName.
type FunctionReference interface {
	fmt.Stringer
	NodeFormatter
	functionReference()
}

var _ FunctionReference = &UnresolvedName{}
var _ FunctionReference = &FunctionDefinition{}
var _ FunctionReference = &ResolvedFunctionDefinition{}

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
		parts := strings.Split(t.Name, ".")
		if len(parts) > 2 {
			return nil, errors.New("function name has too many parts")
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
			fd, err := GetBuiltinFuncDefinition(fn, path)
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
		ref.FunctionReference = fd
		return fd, nil
	default:
		return nil, errors.New("unknown function reference type")
	}
}
