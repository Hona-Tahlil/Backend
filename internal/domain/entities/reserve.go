package entities

import "gorm.io/gorm"

type Reserve struct {
	gorm.Model
	RequestID     uint
	TotalPrice    int
	IsFinished    bool
	Notes         string
	CalenderSlots []CalenderSlot `gorm:"foreignKey:Refer"`
	Comment       Comment
}
