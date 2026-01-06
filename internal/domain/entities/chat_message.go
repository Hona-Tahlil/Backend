package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	RoomID           uint `gorm:"index"`
	SenderID         uint `gorm:"index"`
	MessageType      enums.ChatMessageType `gorm:"type:varchar(16);default:'TEXT'"`
	Content          string
	MediaKey         *string `gorm:"index"`
	IsEdited         bool    `gorm:"default:false"`
	EditedAt         *time.Time `gorm:"index"`
	IsDeleted        bool    `gorm:"default:false"`
	DeletedAtTime    *time.Time `gorm:"index"`
	ReplyToMessageID *uint   `gorm:"index"`
	// ClientMessageID  *string `gorm:"uniqueIndex:ux_sender_client_msg"`
	// Type             MessageType `gorm:"not null;default:'TEXT'"` // TEXT | IMAGE | FILE | SYSTEM
	// EditedAt         *time.Time
	// DeletedAt        *time.Time
	// CreatedAt        time.Time
}
