package authserver

import (
	sql "github.com/dborchard/tiny_crdb/pkg/i_sql"
	"github.com/dborchard/tiny_crdb/pkg/i_sql/isql"
)

type authenticationServer struct {
	sqlServer SQLServerInterface
}
type SQLServerInterface interface {
	ExecutorConfig() *sql.ExecutorConfig
	InternalExecutor() isql.Executor
}
