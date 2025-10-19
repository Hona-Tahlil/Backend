package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type        string
	Permissions []Permission `gorm:"many2many:role_permission"`
}
