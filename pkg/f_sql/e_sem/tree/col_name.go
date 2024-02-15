package tree

import "context"

// ComputeColNameInternal is the workhorse for GetRenderColName.
// The return value indicates the strength of the confidence in the result:
// 0 - no information
// 1 - second-best name choice
// 2 - good name choice
//
// The algorithm is borrowed from FigureColnameInternal in PostgreSQL 10,
// to be found in src/backend/parser/parse_target.c.
func ComputeColNameInternal(
	ctx context.Context, sp SearchPath, target Expr, funcResolver FunctionReferenceResolver,
) (int, string, error) {
	switch e := target.(type) {
	case *UnresolvedName:
		if e.Star {
			return 0, "", nil
		}
		return 2, e.Parts[0], nil

	case *ColumnItem:
		return 2, e.Column(), nil
	case *FuncExpr:
		fd, err := e.Func.Resolve(ctx, sp, funcResolver)
		if err != nil {
			return 0, "", err
		}
		return 2, fd.Name, nil

	default:
		return 0, "", nil
	}
}
