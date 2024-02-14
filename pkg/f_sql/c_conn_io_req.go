package sql

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/z_util/ring"
	"sync"
)

type CmdPos int64

type StmtBuf struct {
	mu struct {
		sync.Mutex
		closed bool

		// cond is signaled when new commands are pushed.
		cond *sync.Cond

		// data contains the elements of the buffer.
		data ring.Buffer[Command]

		// startPos indicates the index of the first command currently in data
		// relative to the start of the connection.
		startPos CmdPos
		// curPos is the current position of the cursor going through the commands.
		// At any time, curPos indicates the position of the command to be returned
		// by CurCmd().
		curPos CmdPos
		// lastPos indicates the position of the last command that was pushed into
		// the buffer.
		lastPos CmdPos
	}
}

func NewStmtBuf() *StmtBuf {
	var buf StmtBuf
	buf.Init()
	return &buf
}

// Init initializes a StmtBuf. It exists to avoid the allocation imposed by
// NewStmtBuf.
func (buf *StmtBuf) Init() {
	buf.mu.lastPos = -1
	buf.mu.cond = sync.NewCond(&buf.mu.Mutex)
	buf.mu.data = ring.NewBuffer[Command]()
}

func (buf *StmtBuf) Close() {

}

// Push adds a Command to the end of the buffer. If a CurCmd() call was blocked
// waiting for this command to arrive, it will be woken up.
//
// An error is returned if the buffer has been closed.
func (buf *StmtBuf) Push(ctx context.Context, cmd Command) error {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	if buf.mu.closed {
		return errors.New("buffer is closed")
	}
	buf.mu.data.AddLast(cmd)
	buf.mu.lastPos++

	buf.mu.cond.Signal()
	return nil
}
