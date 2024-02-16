package eval

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	"github.com/lib/pq/oid"
	"strings"
)

// ParseDOid parses and returns an Oid family datum.
func ParseDOid(ctx context.Context, evalCtx *Context, s string, t *types.T) (*tree.DOid, error) {

	switch t.Oid() {
	case oid.T_regproc:
		// To be compatible with postgres, we always treat the trimmed input string
		// as a function name.
		substrs := strings.Split(strings.TrimSpace(s), ".")
		name := tree.UnresolvedName{NumParts: len(substrs)}
		for i := 0; i < len(substrs); i++ {
			name.Parts[i] = substrs[len(substrs)-1-i]
		}
		funcDef, err := evalCtx.Planner.ResolveFunction(
			ctx, tree.MakeUnresolvedFunctionName(&name), &evalCtx.SessionData().SearchPath,
		)
		if err != nil {
			return nil, err
		}
		overload := funcDef.Overloads[0]
		return tree.NewDOidWithTypeAndName(overload.Oid, t, funcDef.Name), nil
	default:
		d, _ /* errSafeToIgnore */, err := evalCtx.Planner.ResolveOIDFromString(ctx, t, tree.NewDString(s))
		return d, err
	}
}
