package entities

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance        uint `gorm:"default=0;not null"`
	PendingBalance uint `gorm:"default=0;not null"`
	PaymentInfo    *string
	UserID         uint
}
