package sql

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvserver/concurrency/isolation"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/fsm"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
	"time"
)

type stateNoTxn struct{}

var _ fsm.State = &stateNoTxn{}

func (stateNoTxn) State() {}

// advanceCode is part of advanceInfo; it instructs the module managing the
// statements buffer on what action to take.
type advanceCode int

const (
	advanceUnknown advanceCode = iota
	stayInPlace
	advanceOne
	skipBatch
	rewind
)

// advanceInfo represents instructions for the connExecutor about what statement
// to execute next (how to move its cursor over the input statements) and how
// to handle the results produced so far - can they be delivered to the client
// ASAP or not. advanceInfo is the "output" of performing a state transition.
type advanceInfo struct {
	code advanceCode
}
type eventTxnStart struct {
	ImplicitTxn fsm.Bool
}

func (e eventTxnStart) Event() {
	//TODO implement me
	panic("implement me")
}

var eventStartImplicitTxn fsm.Event = eventTxnStart{ImplicitTxn: fsm.True}
var eventStartExplicitTxn fsm.Event = eventTxnStart{ImplicitTxn: fsm.False}

type eventTxnStartPayload struct {
	pri                 roachpb.UserPriority
	txnSQLTimestamp     time.Time
	readOnly            tree.ReadWriteMode
	historicalTimestamp *hlc.Timestamp
	isoLevel            isolation.Level
}
