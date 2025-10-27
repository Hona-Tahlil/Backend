package entities

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID        uint `gorm:"index"`
	PetSitterID   uint `gorm:"index"`
	IsApproved    bool `gorm:"default=false;index"`
	IsPaid        bool `gorm:"default=false;index"`
	Chat          Chat `gorm:"foreignKey:RequestID;not null"`
	TransferID    uint
	IsCanceled    bool           `gorm:"default=false"`
	CalenderSlots []CalenderSlot `gorm:"foreignKey:Refer"`
	Pets          []Pet          `gorm:"many2many:reserve_pet"`
	TotalPrice    uint
	Notes         *string
	Comment       *Comment `gorm:"foreignKey:RequestID"`
	Address       Address  `gorm:"foreignKey:AddressID;not null"`
	AddressID     uint
	Services      []Service `gorm:"many2many:reserve_service"`
}
