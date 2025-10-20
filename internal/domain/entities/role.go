package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type        string       // TODO: should be enum? don't think so
	Permissions []Permission `gorm:"many2many:role_permission"`
}
