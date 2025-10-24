package entities

import "gorm.io/gorm"

type Reserve struct {
	gorm.Model
	RequestID uint `gorm:"index"`
	// remove
	IsCanceled    bool           `gorm:"default=false"`
	CalenderSlots []CalenderSlot `gorm:"foreignKey:Refer"`
	Pets          []Pet          `gorm:"many2many:reserve_pet"`
	TotalPrice    uint
	Notes         *string
	Comment       *Comment
	Address       Address `gorm:"foreignKey:AddressID"`
	AddressID     uint
	Services      []Service `gorm:"many2many:reserve_service"`
}
