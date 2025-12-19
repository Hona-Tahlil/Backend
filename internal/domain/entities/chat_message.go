package entities

import (

	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	RoomID           uint `gorm:"index"`
	SenderID         uint `gorm:"index"`
	Content          string
	ReplyToMessageID *uint   `gorm:"index"`
	// ClientMessageID  *string `gorm:"uniqueIndex:ux_sender_client_msg"`
	// Type             MessageType `gorm:"not null;default:'TEXT'"` // TEXT | IMAGE | FILE | SYSTEM
	// EditedAt         *time.Time
	// DeletedAt        *time.Time
	// CreatedAt        time.Time
}
