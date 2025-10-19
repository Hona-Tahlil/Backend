package entities

import "gorm.io/gorm"

type TextMessage struct {
	gorm.Model
	Content  string
	SenderID uint
	ChatID   uint
}
