package entities

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID      uint
	PetSitterID uint
	IsApproved  bool
	IsPaid      bool
	Reserve     Reserve `gorm:"foreignKey:RequestID"`
	Chat        Chat    `gorm:"foreignKey:RequestID"`
}
