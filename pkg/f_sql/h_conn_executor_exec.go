package sql

import (
	"github.com/dborchard/tiny_crdb/pkg/f_sql/b_parser/statements"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/g_kv/kvserver/concurrency/isolation"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
	"github.com/dborchard/tiny_crdb/pkg/z_util/fsm"
	"github.com/dborchard/tiny_crdb/pkg/z_util/hlc"
	"time"
)

func (ex *connExecutor) execStmtInNoTxnState(parserStmt statements.Statement[tree.Statement]) (_ fsm.Event, payload fsm.EventPayload) {

	ast := parserStmt.AST
	switch ast.(type) {
	default:
		mode, sqlTs, historicalTs, err := ex.beginTransactionTimestampsAndReadMode()
		if err != nil {
			return nil, nil
		}
		return eventStartImplicitTxn, makeEventTxnStartPayload(roachpb.NormalUserPriority, mode, sqlTs, historicalTs, isolation.Serializable)
	}
}

func (ex *connExecutor) beginTransactionTimestampsAndReadMode() (rwMode tree.ReadWriteMode, txnSQLTimestamp time.Time, historicalTimestamp *hlc.Timestamp, err error) {
	return tree.ReadOnly, time.Time{}, nil, nil
}

func makeEventTxnStartPayload(pri roachpb.UserPriority, readOnly tree.ReadWriteMode, txnSQLTimestamp time.Time, historicalTimestamp *hlc.Timestamp, isoLevel isolation.Level) eventTxnStartPayload {
	return eventTxnStartPayload{
		pri:                 pri,
		readOnly:            readOnly,
		txnSQLTimestamp:     txnSQLTimestamp,
		historicalTimestamp: historicalTimestamp,
		isoLevel:            isoLevel,
	}
}
