package tree

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
)

// NameParts is the array of strings that composes the path in an
// UnresolvedName.
type NameParts = [4]string

// UnresolvedName corresponds to an unresolved qualified name.
type UnresolvedName struct {
	NumParts int
	Star     bool
	Parts    NameParts
}

func (u *UnresolvedName) functionReference() {
	//TODO implement me
	panic("implement me")
}

var _ Expr = new(UnresolvedName)

func (u *UnresolvedName) String() string {
	//TODO implement me
	panic("implement me")
}

func (u *UnresolvedName) Format(ctx *FmtCtx) {
	//TODO implement me
	panic("implement me")
}

func (u *UnresolvedName) Walk(visitor Visitor) Expr {
	//TODO implement me
	panic("implement me")
}

func (u *UnresolvedName) TypeCheck(ctx context.Context, semaCtx *SemaContext, desired *types.T) (TypedExpr, error) {
	//TODO implement me
	panic("implement me")
}

// ToRoutineName converts an UnresolvedName to a RoutineName.
func (u *UnresolvedName) ToRoutineName() (RoutineName, error) {
	return RoutineName{}, nil
}

// UnresolvedRoutineName is an unresolved function or procedure name. The two
// implementations of this interface are used to differentiate between the two
// types of routines for things like error messages.
type UnresolvedRoutineName interface {
	UnresolvedName() *UnresolvedName
	isUnresolvedRoutineName()
}

// UnresolvedObjectName is an unresolved qualified name for a database object
// (table, view, etc). It is like UnresolvedName but more restrictive.
// It should only be constructed via NewUnresolvedObjectName.
type UnresolvedObjectName struct {
	// NumParts indicates the number of name parts specified; always 1 or greater.
	NumParts int

	// Parts are the name components, in reverse order.
	// There are at most 3: object name, schema, catalog/db.
	//
	// Note: Parts has a fixed size so that we avoid a heap allocation for the
	// slice every time we construct an UnresolvedObjectName. It does imply
	// however that Parts does not have a meaningful "length"; its actual length
	// (the number of parts specified) is populated in NumParts above.
	Parts [3]string
}

// UnresolvedFunctionName is an unresolved function name.
type UnresolvedFunctionName struct {
	u *UnresolvedName
}

// MakeUnresolvedFunctionName returns a new UnresolvedFunctionName containing
// the give UnresolvedName.
func MakeUnresolvedFunctionName(u *UnresolvedName) UnresolvedFunctionName {
	return UnresolvedFunctionName{u: u}
}
func (u UnresolvedFunctionName) UnresolvedName() *UnresolvedName {
	//TODO implement me
	panic("implement me")
}

func (u UnresolvedFunctionName) isUnresolvedRoutineName() {
	//TODO implement me
	panic("implement me")
}
