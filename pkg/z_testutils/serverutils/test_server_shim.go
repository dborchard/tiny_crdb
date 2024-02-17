package serverutils

import (
	gosql "database/sql"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
)

// StartServer creates and starts a test server.
// The returned server should be stopped by calling
// server.Stopper().Stop().
//
// The second and third return values are equivalent to
// .ApplicationLayer().SQLConn() and .ApplicationLayer().DB(),
// respectively. If your test does not need them, consider
// using StartServerOnly() instead.
func StartServer() (TestServerInterface, *gosql.DB, *kv.DB) {
	return nil, nil, nil
}
