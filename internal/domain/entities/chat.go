package entities

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Messages           []TextMessage `gorm:"foreignKey:ChatID"`
	IsUserBlocked      bool
	IsPetSitterBlocked bool
	IsAccepted         bool
	RequestID          uint
}
