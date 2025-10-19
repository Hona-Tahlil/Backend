package entities

import "gorm.io/gorm"

type Reserve struct {
	gorm.Model
	RequestID     uint `gorm:"not null;index"`
	TotalPrice    uint `gorm:"not null;default=0"`
	IsFinished    bool `gorm:"not null;default=false"`
	Notes         *string
	CalenderSlots []CalenderSlot `gorm:"foreignKey:Refer"`
	Comment       *Comment
}
