package postgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"

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

func (cr *CityRepository) FindCityByName(name enums.City) (*entities.City, error) {
	var foundCity entities.City

	if result := cr.db.First(&foundCity, "name = ?", name); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundCity, nil
}
