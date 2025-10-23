package entities

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance        uint `gorm:"default=0"`
	PendingBalance uint `gorm:"default=0;index"`
	PaymentInfo    *string
	UserID         uint          `gorm:"index"`
	InTransfers    []Transfer    `gorm:"foreignKey:ReceiverWalletID"`
	OutTransfers   []Transfer    `gorm:"foreignKey:SenderWalletID"`
	Transactions   []Transaction `gorm:"foreignKey:WalletID"`
}
