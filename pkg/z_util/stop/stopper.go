package stop

import "context"

// A Stopper provides control over the lifecycle of goroutines started
// through it via its RunTask, RunAsyncTask, and other similar methods.
//
// # When Stop is invoked, the Stopper
//
//   - it invokes Quiesce, which causes the Stopper to refuse new work
//     (that is, its Run* family of methods starts returning ErrUnavailable),
//     closes the channel returned by ShouldQuiesce, and blocks until
//     until no more tasks are tracked, then
//   - it runs all of the methods supplied to AddCloser, then
//   - closes the IsStopped channel.
//
// When ErrUnavailable is returned from a task, the caller needs
// to handle it appropriately by terminating any work that it had
// hoped to defer to the task (which is guaranteed to never have been
// invoked). A simple example of this can be seen in the below snippet:
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//	if err := s.RunAsyncTask("foo", func(ctx context.Context) {
//	  defer wg.Done()
//	}); err != nil {
//	  // Task never ran.
//	  wg.Done()
//	}
//
// To ensure that tasks that do get started are sensitive to Quiesce,
// they need to observe the ShouldQuiesce channel similar to how they
// are expected to observe context cancellation:
//
//	func x() {
//	  select {
//	  case <-s.ShouldQuiesce:
//	    return
//	  case <-ctx.Done():
//	    return
//	  case <-someChan:
//	    // Do work.
//	  }
//	}
//
// TODO(tbg): many improvements here are possible:
//   - propagate quiescing via context cancellation
//   - better API around refused tasks
//   - all the other things mentioned in:
//     https://github.com/cockroachdb/cockroach/issues/58164
type Stopper struct {
	quiescer chan struct{} // Closed when quiescing
	stopped  chan struct{} // Closed when stopped completely
}

// ShouldQuiesce returns a channel which will be closed when Stop() has been
// invoked and outstanding tasks should begin to quiesce.
func (s *Stopper) ShouldQuiesce() <-chan struct{} {
	if s == nil {
		// A nil stopper will never signal ShouldQuiesce, but will also never panic.
		return nil
	}
	return s.quiescer
}

func (s *Stopper) RunAsyncTaskEx(ctx context.Context, f func(ctx context.Context)) error {
	return nil
}

// NewStopper returns an instance of Stopper.
func NewStopper() *Stopper {
	s := &Stopper{
		quiescer: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	register(s)
	return s
}

func register(s *Stopper) {
}

// Stop signals all live workers to stop and then waits for each to
// confirm it has stopped.
//
// Stop is idempotent; concurrent calls will block on each other.
func (s *Stopper) Stop(ctx context.Context) {

}
