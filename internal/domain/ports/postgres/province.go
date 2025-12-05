package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type ProvinceRepository interface {
	FindProvinceByName(name enums.Province) (*entities.Province, error)
	CreateProvince(province *entities.Province) error
	GetAllProvinces() ([]entities.Province, error)
}

