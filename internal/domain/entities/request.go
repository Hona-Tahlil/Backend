package entities

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID      uint    `gorm:"index"`
	PetSitterID uint    `gorm:"index"`
	IsApproved  bool    `gorm:"default=false;index"`
	IsPaid      bool    `gorm:"default=false;index"`
	Reserve     Reserve `gorm:"foreignKey:RequestID;not null"`
	Chat        Chat    `gorm:"foreignKey:RequestID;not null"`
	TransferID  uint
}
