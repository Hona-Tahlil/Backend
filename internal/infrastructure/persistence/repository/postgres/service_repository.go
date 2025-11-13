package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (sr *ServiceRepository) FindServiceByID(id uint) (*entities.Service, error) {
	var foundService entities.Service

	if result := sr.db.First(&foundService, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundService, nil
}
