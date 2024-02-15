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
