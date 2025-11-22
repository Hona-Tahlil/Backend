package domainpostgres

import "hona/backend/internal/domain/entities"

type CityRepository interface {
	CreateCity(city *entities.City) error
}
