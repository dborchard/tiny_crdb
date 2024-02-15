package tree

import "github.com/lib/pq/oid"

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

// OidToQualifiedBuiltinOverload is a map from builtin function OID to an
// qualified overload.
var OidToQualifiedBuiltinOverload map[oid.Oid]QualifiedOverload

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

// NewFunctionDefinition allocates a function definition corresponding
// to the given built-in definition.
func NewFunctionDefinition(
	name string, props *FunctionProperties, def []Overload,
) *FunctionDefinition {
	overloads := make([]*Overload, len(def))

	for i := range def {
		if def[i].PreferredOverload {
			break
		}
	}

	for i := range def {
		overloads[i] = &def[i]
	}
	return &FunctionDefinition{
		Name:               name,
		Definition:         overloads,
		FunctionProperties: *props,
	}
}

// QualifyBuiltinFunctionDefinition qualified all overloads in a function
// definition with a schema name. Note that this function can only be used for
// builtin function.
func QualifyBuiltinFunctionDefinition(
	def *FunctionDefinition, schema string,
) *ResolvedFunctionDefinition {
	ret := &ResolvedFunctionDefinition{
		Name:      def.Name,
		Overloads: make([]QualifiedOverload, 0, len(def.Definition)),
	}
	for _, o := range def.Definition {
		ret.Overloads = append(
			ret.Overloads,
			MakeQualifiedOverload(schema, o),
		)
	}
	return ret
}

// MakeQualifiedOverload creates a new QualifiedOverload.
func MakeQualifiedOverload(schema string, overload *Overload) QualifiedOverload {
	return QualifiedOverload{Schema: schema, Overload: overload}
}
