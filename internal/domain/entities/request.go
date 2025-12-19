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
	Chat          ChatRoom `gorm:"foreignKey:RequestID;not null"`
	TransferID    *uint
	CalendarSlots []CalendarSlot `gorm:"foreignKey:Refer;not null"`
	Pets          []Pet          `gorm:"foreignKey:RequestID;not null"`
	TotalPrice    uint
	Notes         *string
	Comment       *Comment `gorm:"foreignKey:RequestID"`
	Address       Address  `gorm:"polymorphicType:Type;polymorphicId:Refer;polymorphicValue:Request"`
	Service       Service  `gorm:"foreignKey:RequestID;not null"`
}
