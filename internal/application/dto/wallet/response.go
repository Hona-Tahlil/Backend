package wallet

type WalletResponse struct {
	ID             uint    `json:"id"`
	Balance        uint    `json:"balance"`
	PendingBalance uint    `json:"pendingBalance"`
	PaymentInfo    *string `json:"paymentInfo"`
	UserID         uint    `json:"userID"`
}
