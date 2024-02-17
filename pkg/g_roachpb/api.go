package roachpb

type Transaction struct {
}

type UserPriority int32

const (
	NormalUserPriority UserPriority = 0
	MinUserPriority    UserPriority = 25
	MaxUserPriority    UserPriority = 100
)

type LeafTxnInputState struct {
}

type TransactionStatus int32

type Locality struct {
}
