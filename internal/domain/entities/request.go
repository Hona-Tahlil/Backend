package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	UserID        uint `gorm:"index"`
	PetSitterID   uint `gorm:"index"`
	Status        enums.RequestStatus
	Chat          Chat `gorm:"foreignKey:RequestID;not null"`
	TransferID    *uint
	CalendarSlots []CalendarSlot `gorm:"foreignKey:Refer;not null"`
	Pets          []RequestPet   `gorm:"foreignKey:RequestID;not null"`
	TotalPrice    uint
	Notes         *string
	Comment       *Comment         `gorm:"foreignKey:RequestID"`
	Address       RequestAddress   `gorm:"foreignKey:RequestID;not null"`
	Services      []RequestService `gorm:"foreignKey:RequestID;not null"`
}
