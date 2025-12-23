package wallet

import "time"

type WalletResponse struct {
	ID             uint    `json:"id"`
	Balance        uint    `json:"balance"`
	PendingBalance uint    `json:"pendingBalance"`
	PaymentInfo    *string `json:"paymentInfo"`
	UserID         uint    `json:"userID"`
}

type TransferHistoryItemResponse struct {
	ID               uint      `json:"id"`
	Amount           uint      `json:"amount"`
	Direction        string    `json:"direction"`
	SenderWalletID   uint      `json:"senderWalletID"`
	ReceiverWalletID uint      `json:"receiverWalletID"`
	CreatedAt        time.Time `json:"createdAt"`
}

type TransactionHistoryItemResponse struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Amount    uint      `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}
