package usecase

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type ProvinceService interface {
	FindProvinceByName(name enums.Province) (*entities.Province, error)
}
