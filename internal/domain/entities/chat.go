package entities

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Messages           []TextMessage `gorm:"foreignKey:ChatID"`
	IsUserBlocked      bool          `gorm:"default=false"`
	IsPetSitterBlocked bool          `gorm:"default=false"`
	IsAccepted         bool          `gorm:"default=false"`
	RequestID          uint          `gorm:"index"`
}
