package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{
		db: db,
	}
}

func (cr *CityRepository) CreateCity(city *entities.City) error {
	return cr.db.Create(city).Error
}
