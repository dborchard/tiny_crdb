package sql

// Server is the top level singleton for handling SQL connections. It creates
// connExecutors to server every incoming connection.
type Server struct {
	cfg *ExecutorConfig
}
