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

func (rp *RequestRepository) PreloadFields(request *entities.Request, fields []string) error {
	for _, field := range fields {
		err := rp.db.Preload(field).First(request, request.ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (rp *RequestRepository) SearchRequests(userID uint, options *QueryOptions) ([]entities.Request, int64, error) {
	var requests []entities.Request
	query, total := ApplyModifiers(rp.db.Model(&entities.Request{}).Where("user_id = ?", userID), *options)
	if err := query.Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (rp *RequestRepository) SearchRequestsByPetSitterID(petSitterID uint, options *QueryOptions) ([]entities.Request, int64, error) {
	var requests []entities.Request
	query, total := ApplyModifiers(rp.db.Model(&entities.Request{}).Where("pet_sitter_id = ?", petSitterID), *options)
	if err := query.Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}
