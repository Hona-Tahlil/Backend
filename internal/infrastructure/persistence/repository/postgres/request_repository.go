package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{
		db: db,
	}
}

func (rp *RequestRepository) CreateRequest(request *entities.Request) error {
	return rp.db.Create(request).Error
}
