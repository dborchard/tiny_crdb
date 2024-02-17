package kvcoord

import (
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	"github.com/dborchard/tiny_crdb/pkg/g_roachpb"
)

// TxnCoordSenderFactory implements kv.TxnSenderFactory.
type TxnCoordSenderFactory struct {
	wrapped kv.Sender
}

// NewTxnCoordSenderFactory creates a new TxnCoordSenderFactory. The
// factory creates new instances of TxnCoordSenders.
func NewTxnCoordSenderFactory(wrapped kv.Sender) *TxnCoordSenderFactory {
	tcf := &TxnCoordSenderFactory{
		wrapped: wrapped,
	}
	return tcf
}

func (t TxnCoordSenderFactory) RootTransactionalSender(txn *roachpb.Transaction, pri roachpb.UserPriority) kv.TxnSender {
	//TODO implement me
	panic("implement me")
}

func (t TxnCoordSenderFactory) LeafTransactionalSender(tis *roachpb.LeafTxnInputState) kv.TxnSender {
	//TODO implement me
	panic("implement me")
}

func (t TxnCoordSenderFactory) NonTransactionalSender() kv.Sender {
	//TODO implement me
	panic("implement me")
}
