package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
)

// schemaResolver implements the resolver.SchemaResolver interface.
// Currently, this is only being embedded in the planner but also a convenience
// for inejcting it into the declarative schema changer.
// It holds sessionDataStack and a transaction handle which are reset when
// planner is reset.
// TODO (Chengxiong) refactor this out into a separate package.
type schemaResolver struct {
}

func (sr *schemaResolver) ResolveFunction(
	ctx context.Context, name tree.UnresolvedRoutineName, path tree.SearchPath,
) (*tree.ResolvedFunctionDefinition, error) {
	uname := name.UnresolvedName()
	fn, err := uname.ToRoutineName()
	if err != nil {
		return nil, err
	}
	// Get builtin and udf functions if there is any match.
	builtinDef, err := tree.GetBuiltinFuncDefinition(fn, path)
	if err != nil {
		return nil, err
	}
	routine, err := maybeLookupRoutine(ctx, sr, path, fn)
	if err != nil {
		return nil, err
	}

	switch {
	case builtinDef != nil && routine != nil:
		return builtinDef.MergeWith(routine)
	case builtinDef != nil:
		return builtinDef, nil
	case routine != nil:
		return routine, nil
	default:
		return nil, err
	}
}

func maybeLookupRoutine(
	ctx context.Context, sr *schemaResolver, path tree.SearchPath, fn tree.RoutineName,
) (*tree.ResolvedFunctionDefinition, error) {
	return nil, nil
}
