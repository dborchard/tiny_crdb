package tree

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/volatility"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	"github.com/lib/pq/oid"
)

// TypeList is a list of types representing a function parameter list.
type TypeList interface {
	Match(types []*types.T) bool
	GetAt(i int) *types.T
	Length() int
	Types() []*types.T
	String() string
}

// FunctionClass specifies the class of the builtin function.
type FunctionClass int

const (
	// NormalClass is a standard builtin function.
	NormalClass FunctionClass = iota
	// AggregateClass is a builtin aggregate function.
	AggregateClass
	// WindowClass is a builtin window function.
	WindowClass
	// GeneratorClass is a builtin generator function.
	GeneratorClass
	// SQLClass is a builtin function that executes a SQL statement as a side
	// effect of the function call.
	//
	// For example, AddGeometryColumn is a SQLClass function that executes an
	// ALTER TABLE ... ADD COLUMN statement to add a geometry column to an
	// existing table. It returns metadata about the column added.
	//
	// All builtin functions of this class should include a definition for
	// Overload.SQLFn, which returns the SQL statement to be executed. They
	// should also include a definition for Overload.Fn, which is executed
	// like a NormalClass function and returns a Datum.
	SQLClass
)

// ReturnTyper defines the type-level function in which a builtin function's return type
// is determined. ReturnTypers should make sure to return unknownReturnType when necessary.
type ReturnTyper func(args []TypedExpr) *types.T
type FnOverload interface{}

// Overload is one of the overloads of a built-in function.
// Each FunctionDefinition may contain one or more overloads.
type Overload struct {
	Types             TypeList
	ReturnType        ReturnTyper
	PreferredOverload bool
	Info              string
	Class             FunctionClass
	Fn                FnOverload
	Oid               oid.Oid
	ReturnSet         bool
	Volatility        volatility.V
}

// ParamTypes is a list of function parameter names and their types.
type ParamTypes []ParamType

// Match is part of the TypeList interface.
func (p ParamTypes) Match(types []*types.T) bool {
	if len(types) != len(p) {
		return false
	}
	for i := range types {
		if !p.MatchAt(types[i], i) {
			return false
		}
	}
	return true
}

// MatchAt is part of the TypeList interface.
func (p ParamTypes) MatchAt(typ *types.T, i int) bool {
	return i < len(p) && (typ.Family() == types.UnknownFamily || p[i].Typ.Equivalent(typ))
}

func (p ParamTypes) GetAt(i int) *types.T {
	//TODO implement me
	panic("implement me")
}

func (p ParamTypes) Length() int {
	//TODO implement me
	panic("implement me")
}

func (p ParamTypes) Types() []*types.T {
	//TODO implement me
	panic("implement me")
}

func (p ParamTypes) String() string {
	//TODO implement me
	panic("implement me")
}

// ParamType encapsulate a function parameter name and type.
type ParamType struct {
	Name string
	Typ  *types.T
}

// FixedReturnType functions simply return a fixed type, independent of argument types.
func FixedReturnType(typ *types.T) ReturnTyper {
	return func(args []TypedExpr) *types.T { return typ }
}
