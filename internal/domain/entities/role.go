package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type        enums.Role
	Permissions []Permission `gorm:"many2many:role_permission"`
}
