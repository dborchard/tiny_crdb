package sql

import "context"

// ieResultReader is used to read internalExecutor results.
// It is managed by the rowsIterator.
type ieResultReader interface {

	// firstResult returns the first result. The return values carry the same
	// semantics as of nextResult. This method assumes that the writer is not
	// currently blocked and waits for the initial result to be written.
	firstResult(ctx context.Context) (_ ieIteratorResult, done bool, err error)

	// nextResult returns the next result. Done will always be true if err
	// is non-nil. Err will be non-nil if either close has been called or
	// the passed context is finished.
	nextResult(ctx context.Context) (_ ieIteratorResult, done bool, err error)

	// close ensures that either writer has finished writing. The writer will
	// receive a signal to drain, and close will drain the writer's channel.
	close() error
}
