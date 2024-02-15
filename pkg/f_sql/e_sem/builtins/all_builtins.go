package builtins

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/builtins/builtinsregistry"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/lib/pq/oid"
	"strings"
)

// registerBuiltin adds the given builtin to the builtins registry. All
// overloads of the Generator class are updated to have Fn and FnWithExprs
// fields to be functions that return assertion errors upon execution (to
// prevent misuse).
//
// If enforceClass is true, then it panics if at least one overload is not of
// the expected class.
//
// Note that additional sanity checks are also performed in eval/overload.go.
func registerBuiltin(
	name string, def builtinDefinition, expectedClass tree.FunctionClass, enforceClass bool,
) {
	for i := range def.overloads {
		overload := &def.overloads[i]
		if enforceClass {
			if overload.Class != expectedClass {
				panic("unexpected class")
			}
		}
	}
	builtinsregistry.Register(name, &def.props, def.overloads)
}

func init() {
	tree.ResolvedBuiltinFuncDefs = make(map[string]*tree.ResolvedFunctionDefinition)

	builtinsregistry.AddSubscription(func(name string, props *tree.FunctionProperties, overloads []tree.Overload) {
		fDef := tree.NewFunctionDefinition(name, props, overloads)
		addResolvedFuncDef(tree.ResolvedBuiltinFuncDefs, tree.OidToQualifiedBuiltinOverload, fDef)
	})
}

func addResolvedFuncDef(
	resolved map[string]*tree.ResolvedFunctionDefinition,
	oidToOl map[oid.Oid]tree.QualifiedOverload,
	def *tree.FunctionDefinition,
) {
	parts := strings.Split(def.Name, ".")
	if len(parts) > 2 || len(parts) == 0 {
		panic("malformed function name")
	}

	var fd *tree.ResolvedFunctionDefinition
	if len(parts) == 2 {
		fd = tree.QualifyBuiltinFunctionDefinition(def, parts[0])
		resolved[def.Name] = fd
		return
	}
	for _, o := range fd.Overloads {
		oidToOl[o.Oid] = o
	}
}
