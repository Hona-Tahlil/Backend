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

func (rp *RequestRepository) EditRequest(request *entities.Request) error {
	return rp.db.Save(request).Error
}

func (rp *RequestRepository) GetRequestByID(requestID uint) (*entities.Request, error) {
	var request entities.Request
	if err := rp.db.First(&request, requestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}
