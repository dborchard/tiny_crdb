package sql

import (
	"context"
	"errors"
	"sync"
)

// ieExecutionMode determines how the internal executor consumes the results of
// the statement evaluation.
type ieExecutionMode int

const (
	// defaultIEExecutionMode is the execution mode in which the results of the
	// statement evaluation are consumed according to the statement's type.
	defaultIEExecutionMode ieExecutionMode = iota
	// rowsAffectedIEExecutionMode is the execution mode in which the internal
	// executor is only interested in the number of rows affected, regardless of
	// the statement's type.
	//
	// With this mode, if a stmt encounters a retry error, the internal executor
	// will proceed to transparently reset the number of rows affected (if any
	// have been seen by the rowsIterator) and retry the corresponding command.
	// Such behavior makes sense given that in production code at most one
	// command in the StmtBuf results in "rows affected".
	rowsAffectedIEExecutionMode
)

// ieResultChannel is used to coordinate passing results from an
// internalExecutor to its corresponding iterator. It can be constructed to
// ensure that there is no concurrency between the reader and writer.
type ieResultChannel struct {

	// dataCh is the channel on which the connExecutor goroutine sends the rows
	// (in addResult) and, in the synchronous case, will block on waitCh after
	// each send. The iterator goroutine blocks on dataCh until there is
	// something to receive (rows or other metadata) and will return the data to
	// the caller. On the next call to Next(), the iterator goroutine unblocks
	// the producer and will block itself again. dataCh will be closed (in
	// finish()) when the connExecutor goroutine exits its run() loop whereas
	// waitCh is closed when closing the iterator.
	dataCh chan ieIteratorResult

	// waitCh is nil for async ieResultChannels. It is never closed. In all places
	// where the caller may interact with it the doneCh is also used. This policy
	// is in place to make it safe to unblock both the reader and the writer
	// without any hazards of a blocked reader attempting to send on a closed
	// channel.
	waitCh chan struct{}

	// doneCh is used to indicate that the ieResultReader has been closed and is
	// closed under the doneOnce, the writer will transition to draining. This
	// is crucial to ensure that a synchronous writer does not attempt to
	// continue to operate after the reader has called close.
	doneCh   chan struct{}
	doneErr  error
	doneOnce sync.Once
}

func (i *ieResultChannel) finish() {
	close(i.dataCh)
}

func (i *ieResultChannel) firstResult(
	ctx context.Context,
) (_ ieIteratorResult, done bool, err error) {
	select {
	case <-ctx.Done():
		return ieIteratorResult{}, true, errors.New("context done")
	case <-i.doneCh:
		return ieIteratorResult{}, true, errors.New("done")
	case res, ok := <-i.dataCh:
		if !ok {
			return ieIteratorResult{}, true, errors.New("channel closed")
		}
		return res, false, nil
	}
}

func (i *ieResultChannel) nextResult(ctx context.Context) (_ ieIteratorResult, done bool, err error) {
	if done, err = i.maybeUnblockWriter(ctx); done {
		return ieIteratorResult{}, done, err
	}
	return i.firstResult(ctx)
}

func (i *ieResultChannel) close() error {
	return i.doneErr
}

func (i *ieResultChannel) maybeUnblockWriter(ctx context.Context) (done bool, err error) {
	select {
	case <-ctx.Done():
		return true, ctx.Err()
	case <-i.doneCh:
		return true, ctx.Err()
	case i.waitCh <- struct{}{}:
		return false, nil
	}
}

func (i *ieResultChannel) addResult(ctx context.Context, result ieIteratorResult) error {
	select {
	case <-ctx.Done():
		return errors.New("context done")
	case <-i.doneCh:
		if ctxErr := ctx.Err(); ctxErr != nil {
			return errors.New("context done")
		}
		return errors.New("context done")
	case i.dataCh <- result:
	}
	return nil
}

// newAsyncIEResultChannel returns an ieResultChannel which does not attempt to
// synchronize the writer with the reader.
func newAsyncIEResultChannel() *ieResultChannel {
	return &ieResultChannel{
		dataCh: make(chan ieIteratorResult, 1),
		doneCh: make(chan struct{}),
	}
}

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

// ieResultWriter is used by the internalExecutor to write results to an
// iterator.
type ieResultWriter interface {

	// addResult adds a result. It may block until the next result is requested
	// by the reader, depending on the synchronization strategy.
	addResult(ctx context.Context, result ieIteratorResult) error

	// finish is used to indicate that the writer is done writing rows.
	finish()
}
