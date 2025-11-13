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

func (ar *AddressRepository) FindAddressesByUserID(id uint) ([]entities.Address, error) {
	addresses := make([]entities.Address, 0)

	var mainAddress entities.Address

	err := ar.db.First(&mainAddress, "refer = ?", id).Error
	if err != nil {
		return nil, err
	}

	err = ar.db.Find(&addresses, "refer = ? AND type = ?", id, "Request").Error
	if err != nil {
		return nil, err
	}

	flag := false
	for _, address := range addresses {
		if address.City == mainAddress.City && address.Province.Name == mainAddress.Province.Name && address.StreetAddress == mainAddress.StreetAddress {
			flag = true
			break
		}
	}
	if !flag {
		addresses = append(addresses, mainAddress)
	}

	return addresses, nil
}
