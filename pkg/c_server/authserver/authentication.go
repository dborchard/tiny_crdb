package authserver

import (
	sql "github.com/dborchard/tiny_crdb/pkg/i_sql"
	"github.com/dborchard/tiny_crdb/pkg/i_sql/isql"
)

type authenticationServer struct {
	sqlServer SQLServerInterface
}

func NewServer(sqlServer SQLServerInterface) Server {
	return &authenticationServer{
		sqlServer: sqlServer,
	}
}

type SQLServerInterface interface {
	ExecutorConfig() *sql.ExecutorConfig
	InternalExecutor() isql.Executor
}
