package sql

import kv "github.com/dborchard/tiny_crdb/pkg/g_kv"

// An ExecutorConfig encompasses the auxiliary objects and configuration
// required to create an executor.
// All fields holding a pointer or an interface are required to create
// an Executor; the rest will have sane defaults set if omitted.
type ExecutorConfig struct {
	DB             *kv.DB
	InternalDB     *InternalDB
	DistSQLPlanner *DistSQLPlanner
}
