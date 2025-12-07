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
	err := ar.db.Preload("Province.Cities").First(&foundAddress, foundAddress.ID).Error
	if err != nil {
		return nil, err
	}
	err = ar.db.Preload("City").First(&foundAddress, foundAddress.ID).Error
	if err != nil {
		return nil, err
	}

	return &foundAddress, nil
}

func (ar *AddressRepository) FindAddressesByUserID(id uint) ([]entities.Address, error) {
	addresses := make([]entities.Address, 0)

	var mainAddress entities.Address

	if err := ar.db.First(&mainAddress, "refer = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return addresses, nil
		}
		return nil, err
	}

	err := ar.db.Preload("Province.Cities").First(&mainAddress, mainAddress.ID).Error
	if err != nil {
		return nil, err
	}
	err = ar.db.Preload("City").First(&mainAddress, mainAddress.ID).Error
	if err != nil {
		return nil, err
	}

	err = ar.db.Find(&addresses, "refer = ? AND type = ?", id, "Request").Error
	if err != nil {
		return nil, err
	}

	flag := false
	for _, address := range addresses {
		err = ar.db.Preload("Province.Cities").First(&address, address.ID).Error
		if err != nil {
			return nil, err
		}
		err = ar.db.Preload("City").First(&address, address.ID).Error
		if err != nil {
			return nil, err
		}

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

func (ar *AddressRepository) FinduserAddressByUserID(id uint) (*entities.Address, error) {
	var foundAddress entities.Address

	err := ar.db.First(&foundAddress, "refer = ? AND type = ?", id, "User").Error
	if err != nil {
		return nil, err
	}

	err = ar.db.Preload("Province.Cities").First(&foundAddress, foundAddress.ID).Error
	if err != nil {
		return nil, err
	}
	err = ar.db.Preload("City").First(&foundAddress, foundAddress.ID).Error
	if err != nil {
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
