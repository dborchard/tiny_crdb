package sql

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/c_catalog/colinfo"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	"sync"
)

const iteratorDepthLimit = 1000

var iteratorDepthLimitExceededErr = errors.New("rowsIterator exceeded recursion depth limit")

type rowsIterator struct {
	r     ieResultReader
	depth int64

	rowsAffected int
	resultCols   colinfo.ResultColumns
	lastRow      tree.Datums
	lastErr      error
	done         bool
	wg           *sync.WaitGroup

	mode        ieExecutionMode
	stmtBuf     *StmtBuf
	errCallback func(err error) error
	first       *ieIteratorResult
}

func (r *rowsIterator) Next(ctx context.Context) (_ bool, retErr error) {
	defer func() {
		if r.done {
			_ = r.Close()
		}
		if r.errCallback != nil {
			r.lastErr = r.errCallback(r.lastErr)
			r.errCallback = nil
		}
		retErr = r.lastErr
		r.depth--
	}()

	r.depth++
	if r.depth > iteratorDepthLimit {
		r.lastErr = iteratorDepthLimitExceededErr
		r.done = true
		return false, r.lastErr
	}

	if r.done {
		return false, r.lastErr
	}

	// handleDataObject processes a single object read from ieResultReader and
	// returns the result to be returned by Next. It also might call Next
	// recursively if the object is a piece of metadata.
	handleDataObject := func(data ieIteratorResult) (bool, error) {
		if data.row != nil {
			r.rowsAffected++
			// No need to make a copy because streamingCommandResult does that
			// for us.
			r.lastRow = data.row
			return true, nil
		}
		if data.rowsAffected != nil {
			r.rowsAffected = *data.rowsAffected
			return r.Next(ctx)
		}
		if data.cols != nil {
			if r.mode == rowsAffectedIEExecutionMode {
				// In "rows affected" execution mode we simply ignore the column
				// schema since we always return the number of rows affected
				// (i.e. a single integer column).
				return r.Next(ctx)
			}
			// At this point we don't expect to see the columns - we should only
			// return the rowsIterator to the caller of execInternal after the
			// columns have been determined.
			data.err = errors.New("unexpectedly received non-nil cols in Next: %v")
		}
		if data.err == nil {
			data.err = errors.New("unexpectedly empty ieIteratorResult object")
		}
		r.lastErr = data.err
		r.done = true
		return false, r.lastErr
	}

	if r.first != nil {
		// This is the very first call to Next() and we have already buffered
		// up the first piece of data before returning rowsIterator to the
		// caller.
		first := r.first
		r.first = nil
		return handleDataObject(*first)
	}

	var next ieIteratorResult
	next, r.done, r.lastErr = r.r.nextResult(ctx)
	if r.done || r.lastErr != nil {
		return false, r.lastErr
	}
	return handleDataObject(next)
}
func (r *rowsIterator) Cur() tree.Datums {
	return r.lastRow
}

func (r *rowsIterator) RowsAffected() int {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) Close() error {
	r.stmtBuf.Close()
	if err := r.r.close(); err != nil && r.lastErr == nil {
		r.lastErr = err
	}
	return r.lastErr
}
func (r *rowsIterator) Types() colinfo.ResultColumns {
	//TODO implement me
	panic("implement me")
}

func (r *rowsIterator) HasResults() bool {
	//TODO implement me
	panic("implement me")
}

type ieIteratorResult struct {
	// Exactly one of these 4 fields will be set.
	row          tree.Datums
	rowsAffected *int
	cols         colinfo.ResultColumns
	err          error
}
