package entities

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	ReceiverWalletID uint
	SenderWalletID   uint
	Amount           uint
	Request          Request `gorm:"foreignKey:TransferID"`
}
