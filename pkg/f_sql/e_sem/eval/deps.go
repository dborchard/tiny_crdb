package eval

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/privilege"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/types"
	"github.com/dborchard/tiny_crdb/pkg/security/username"
)

// Planner is a limited planner that can be used from EvalContext.
type Planner interface {
	DatabaseCatalog
	TypeResolver
	tree.FunctionReferenceResolver

	// QueryRowEx executes the supplied SQL statement and returns a single row, or
	// nil if no row is found, or an error if more that one row is returned.
	//
	// The fields set in session that are set override the respective fields if
	// they have previously been set through SetSessionData().
	QueryRowEx(
		ctx context.Context,
		opName string,
		override sessiondata.InternalExecutorOverride,
		stmt string,
		qargs ...interface{}) (tree.Datums, error)

	// Optimizer returns the optimizer associated with this Planner, if any.
	Optimizer() interface{}

	// AutoCommit indicates whether the Planner has flagged the current statement
	// as eligible for transaction auto-commit.
	AutoCommit() bool
}

// DatabaseCatalog consists of functions that reference the session database
// and is to be used from Context.
type DatabaseCatalog interface {

	// ResolveTableName expands the given table name and
	// makes it point to a valid object.
	// If the database name is not given, it uses the search path to find it, and
	// sets it on the returned TableName.
	// It returns the ID of the resolved table, and an error if the table doesn't exist.
	ResolveTableName(ctx context.Context, tn *tree.TableName) (tree.ID, error)

	// SchemaExists looks up the schema with the given name and determines
	// whether it exists.
	SchemaExists(ctx context.Context, dbName, scName string) (found bool, err error)

	// HasAnyPrivilegeForSpecifier returns whether the current user has privilege
	// to access the given object.
	HasAnyPrivilegeForSpecifier(ctx context.Context, specifier HasPrivilegeSpecifier, user username.SQLUsername, privs []privilege.Privilege) (HasAnyPrivilegeResult, error)
}

// HasPrivilegeSpecifier specifies an object to lookup privilege for.
// Only one of { DatabaseName, DatabaseOID, SchemaName, TableName, TableOID } is filled.
type HasPrivilegeSpecifier struct {
}

// HasAnyPrivilegeResult represents the non-error results of calling HasAnyPrivilege
type HasAnyPrivilegeResult = int8

const (
	// HasPrivilege means at least one of the specified privileges is granted.
	HasPrivilege HasAnyPrivilegeResult = 1
	// HasNoPrivilege means no privileges are granted.
	HasNoPrivilege HasAnyPrivilegeResult = 0
	// ObjectNotFound means the object that privileges are being checked on was not found.
	ObjectNotFound HasAnyPrivilegeResult = -1
)

// TypeResolver is an interface for resolving types and type OIDs.
type TypeResolver interface {
	tree.TypeReferenceResolver

	// ResolveOIDFromString looks up the populated value of the OID with the
	// desired resultType which matches the provided name.
	//
	// The return value is a fresh DOid of the input oid.Oid with name and OID
	// set to the result of the query. If there was not exactly one result to the
	// query, an error will be returned.
	ResolveOIDFromString(
		ctx context.Context, resultType *types.T, toResolve *tree.DString,
	) (_ *tree.DOid, errSafeToIgnore bool, _ error)

	// ResolveOIDFromOID looks up the populated value of the oid with the
	// desired resultType which matches the provided oid.
	//
	// The return value is a fresh DOid of the input oid.Oid with name and OID
	// set to the result of the query. If there was not exactly one result to the
	// query, an error will be returned.
	ResolveOIDFromOID(
		ctx context.Context, resultType *types.T, toResolve *tree.DOid,
	) (_ *tree.DOid, errSafeToIgnore bool, _ error)
}
