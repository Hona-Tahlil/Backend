package entities

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID      uint    `gorm:"not null;index"`
	PetSitterID uint    `gorm:"not null;index"`
	IsApproved  bool    `gorm:"not null;default=false;index"`
	IsPaid      bool    `gorm:"not null=default=false;index"`
	Reserve     Reserve `gorm:"foreignKey:RequestID;not null"`
	Chat        Chat    `gorm:"foreignKey:RequestID;not null"`
}
