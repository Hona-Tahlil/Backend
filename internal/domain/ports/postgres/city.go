package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type CityRepository interface {
	FindCityByName(name enums.City) (*entities.City, error)
}
