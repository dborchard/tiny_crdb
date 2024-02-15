package builtins

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/builtins/builtinsregistry"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
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
