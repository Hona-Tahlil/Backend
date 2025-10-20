package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Type        enums.Permission
	Description string
}
