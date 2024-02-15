package builtins

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/builtins/builtinconstants"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/eval"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/volatility"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

func init() {
	for k, v := range regularBuiltins {
		const enforceClass = true
		registerBuiltin(k, v, tree.NormalClass, enforceClass)
	}
}

// builtinDefinition represents a built-in function before it becomes
// a tree.FunctionDefinition.
type builtinDefinition struct {
	props     tree.FunctionProperties
	overloads []tree.Overload
}

// builtins contains the built-in functions indexed by name.
//
// For use in other packages, see AllBuiltinNames and GetBuiltinProperties().
var regularBuiltins = map[string]builtinDefinition{
	"bit_length": makeBuiltin(tree.FunctionProperties{Category: builtinconstants.CategoryString},
		stringOverload1(
			func(_ context.Context, _ *eval.Context, s string) (tree.Datum, error) {
				return tree.NewDInt(tree.DInt(len(s) * 8)), nil
			},
			types.Int,
			"Calculates the number of bits used to represent `val`.",
			volatility.Immutable,
		),
		bytesOverload1(
			func(_ context.Context, _ *eval.Context, s string) (tree.Datum, error) {
				return tree.NewDInt(tree.DInt(len(s) * 8)), nil
			},
			types.Int,
			"Calculates the number of bits used to represent `val`.",
			volatility.Immutable,
		),
	),
}

func makeBuiltin(props tree.FunctionProperties, overloads ...tree.Overload) builtinDefinition {
	return builtinDefinition{
		props:     props,
		overloads: overloads,
	}
}

func stringOverload1(
	f func(context.Context, *eval.Context, string) (tree.Datum, error),
	returnType *types.T,
	info string,
	volatility volatility.V,
) tree.Overload {
	return tree.Overload{
		Types:      tree.ParamTypes{{Name: "val", Typ: types.String}},
		ReturnType: tree.FixedReturnType(returnType),
		Fn: func(ctx context.Context, evalCtx *eval.Context, args tree.Datums) (tree.Datum, error) {
			return f(ctx, evalCtx, string(*(args[0].(*tree.DString))))
		},
		Info:       info,
		Volatility: volatility,
	}
}

func bytesOverload1(
	f func(context.Context, *eval.Context, string) (tree.Datum, error),
	returnType *types.T,
	info string,
	volatility volatility.V,
) tree.Overload {
	return tree.Overload{
		Types:      tree.ParamTypes{{Name: "val", Typ: types.Bytes}},
		ReturnType: tree.FixedReturnType(returnType),
		Fn: func(ctx context.Context, evalCtx *eval.Context, args tree.Datums) (tree.Datum, error) {
			return f(ctx, evalCtx, string(*args[0].(*tree.DBytes)))
		},
		Info:       info,
		Volatility: volatility,
	}
}
