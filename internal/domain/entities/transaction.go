package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Type     enums.TransactionType `gorm:"index"`
	WalletID uint                  `gorm:"index"`
	Amount   uint
}
