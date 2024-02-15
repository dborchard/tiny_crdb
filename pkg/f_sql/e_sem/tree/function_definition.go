package tree

// FunctionProperties defines the properties of the built-in
// functions that are common across all overloads.
type FunctionProperties struct {
	Category string
}

// ResolvedFunctionDefinition is similar to FunctionDefinition but with all the
// overloads qualified with schema name.
type ResolvedFunctionDefinition struct {
	Name      string
	Overloads []QualifiedOverload
}

func (d ResolvedFunctionDefinition) MergeWith(routine *ResolvedFunctionDefinition) (*ResolvedFunctionDefinition, error) {
	return nil, nil
}

// QualifiedOverload is a wrapper of Overload prefixed with a schema name.
// It indicates that the overload is defined with the specified schema.
type QualifiedOverload struct {
	Schema string
	*Overload
}

var ResolvedBuiltinFuncDefs map[string]*ResolvedFunctionDefinition

func GetBuiltinFuncDefinition(
	fName RoutineName, searchPath SearchPath,
) (*ResolvedFunctionDefinition, error) {
	if fName.ExplicitSchema {
		return ResolvedBuiltinFuncDefs[fName.Schema()+"."+fName.Object()], nil
	}
	return nil, nil
}
