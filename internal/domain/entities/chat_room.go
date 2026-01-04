package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	UserID                     uint                 `gorm:"index"`
	PetSitterID                uint                 `gorm:"index"`
	User                       User                 `gorm:"foreignKey:UserID"`
	PetSitter                  PetSitter            `gorm:"foreignKey:PetSitterID"`
	RequestID                  *uint                `gorm:"index;unique;null"`
	Status                     enums.ChatRoomStatus `gorm:"type:integer;index"` // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy                  *enums.ChatBlockedBy `gorm:"index;null"`         // USER | PET_SITTER | ADMIN
	LastMessageID              *uint                `gorm:"index"`
	UserLastReadMessageID      *uint                `gorm:"null"`
	PetSitterLastReadMessageID *uint                `gorm:"null"`
	// CreatedAt                  time.Time
	// BlockedAt                  time.Time
}
