package entities

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance        int
	PendingBalance int
	PaymentInfo    string
	UserID         uint
}
