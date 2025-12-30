package wallet

type TopUpRequest struct {
	UserID uint
	Amount uint
}

type WithdrawRequest struct {
	UserID uint
	Amount uint
}

type HistoryRequest struct {
	UserID uint
	Offset int
	Limit  int
}
