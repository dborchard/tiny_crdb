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

func (d *ResolvedFunctionDefinition) String() string {
	//TODO implement me
	panic("implement me")
}

func (d *ResolvedFunctionDefinition) Format(ctx *FmtCtx) {
	//TODO implement me
	panic("implement me")
}

func (d *ResolvedFunctionDefinition) functionReference() {
	//TODO implement me
	panic("implement me")
}

func (d *ResolvedFunctionDefinition) MergeWith(routine *ResolvedFunctionDefinition) (*ResolvedFunctionDefinition, error) {
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

// FunctionDefinition implements a reference to the (possibly several)
// overloads for a built-in function.
// TODO(Chengxiong): Remove this struct entirely. Instead, use overloads from
// function resolution or use "GetBuiltinProperties" if the need is to only look
// at builtin functions(there are such existing use cases). Also change "Name"
// of ResolvedFunctionDefinition to Name type.
type FunctionDefinition struct {
	// Name is the short name of the function.
	Name string

	// Definition is the set of overloads for this function name.
	Definition []*Overload

	// FunctionProperties are the properties common to all overloads.
	FunctionProperties
}

func (f *FunctionDefinition) String() string {
	//TODO implement me
	panic("implement me")
}

func (f *FunctionDefinition) Format(ctx *FmtCtx) {
	//TODO implement me
	panic("implement me")
}

func (f *FunctionDefinition) functionReference() {
	//TODO implement me
	panic("implement me")
}
