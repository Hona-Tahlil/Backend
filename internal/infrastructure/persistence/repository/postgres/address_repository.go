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

func (ar *AddressRepository) FindUserAddressByUserID(id uint) (*entities.Address, error) {
	var foundAddress entities.Address

	err := ar.db.First(&foundAddress, "refer = ? AND type = ?", id, "User").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &foundAddress, nil
}

func (ar *AddressRepository) Create(address *entities.Address) error {
	if err := ar.db.Create(address).Error; err != nil {
		return err
	}
	return nil
}
func (ar *AddressRepository) Update(address *entities.Address) error {
	if err := ar.db.Save(address).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AddressRepository) FindUserRequestAddressesByID(id uint) ([]entities.Address, error) {
	var addresses []entities.Address

	err := ar.db.Where("refer = ? AND type = ?", id, "Request").Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (ar *AddressRepository) FindRequestAddressByID(requestID uint) (*entities.Address, error) {
	var address entities.Address

	err := ar.db.Where("refer = ? AND type = ?", requestID, "Request").First(&address).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &address, nil
}
