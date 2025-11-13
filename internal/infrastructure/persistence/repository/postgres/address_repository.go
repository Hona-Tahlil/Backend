package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{
		db: db,
	}
}

func (ar *AddressRepository) FindAddressByID(id uint) (*entities.Address, error) {
	var foundAddress entities.Address

	if result := ar.db.First(&foundAddress, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundAddress, nil
}
