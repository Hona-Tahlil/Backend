package enums

type TransactionType uint

const (
	PayIn TransactionType = iota + 1
	Withdraw
)

func (transactionType TransactionType) String() string {
	switch transactionType {
	case TransactionType(PayIn):
		return "pay in"
	case TransactionType(Withdraw):
		return "withdraw"
	}
	return ""
}
