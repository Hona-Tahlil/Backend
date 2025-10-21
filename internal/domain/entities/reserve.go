package entities

import "gorm.io/gorm"

type Reserve struct {
	gorm.Model
	RequestID     uint           `gorm:"index"`
	IsFinished    bool           `gorm:"default=false"`
	CalenderSlots []CalenderSlot `gorm:"foreignKey:Refer"`
	Pets          []Pet          `gorm:"many2many:reserve_pet"`
	TotalPrice    uint
	Notes         *string
	Comment       *Comment
}
