package entities

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	ReceiverWalletID uint `gorm:"index"`
	SenderWalletID   uint `gorm:"index"`
	Amount           uint
	Request          Request `gorm:"foreignKey:TransferID"`
}
