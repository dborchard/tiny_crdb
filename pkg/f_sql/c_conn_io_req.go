package sql

import (
	"context"
	"errors"
	"github.com/dborchard/tiny_crdb/pkg/z_util/ring"
	"io"
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

func (buf *StmtBuf) CurCmd() (Command, CmdPos, error) {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	for {
		if buf.mu.closed {
			return nil, 0, io.EOF
		}
		curPos := buf.mu.curPos
		length := buf.mu.data.Len()
		cmdIdx := int(curPos - buf.mu.startPos)
		if cmdIdx < length {
			return buf.mu.data.Get(cmdIdx), curPos, nil
		}
		if cmdIdx != length {
			return nil, 0, errors.New("can only wait for next command; corrupt cursor: %d")
		}
		// Wait for the next Command to arrive to the buffer.
		buf.mu.cond.Wait()
	}
}

// AdvanceOne advances the cursor one Command over. The command over which
// the cursor will be positioned when this returns may not be in the buffer
// yet. The previous CmdPos is returned.
func (buf *StmtBuf) AdvanceOne() CmdPos {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	prev := buf.mu.curPos
	buf.mu.curPos++
	return prev
}
