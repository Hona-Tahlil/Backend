package usecase

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type CityService interface {
	FindCityByNameInProvince(name enums.City, province *entities.Province) (*entities.City, error)
}
