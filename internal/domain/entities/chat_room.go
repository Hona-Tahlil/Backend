package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	UserID                     uint                 `gorm:"index"`
	PetSitterID                uint                 `gorm:"index"`
	RequestID                  uint                 `gorm:"index;unique"` 
	Status                     enums.ChatRoomStatus `gorm:"index"` // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy                  enums.ChatBlockedBy  `gorm:"type:text"`
	LastMessageID              *uint                `gorm:"index"`
	UserLastReadMessageID      *uint                `gorm:"null"`
	PetsitterLastReadMessageID *uint                `gorm:"null"`
	// CreatedAt                  time.Time
	// BlockedAt                  time.Time
}
