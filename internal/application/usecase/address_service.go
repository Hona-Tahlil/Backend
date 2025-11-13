package usecase

import "hona/backend/internal/domain/entities"

type AddressService interface {
	FindAddressByID(id uint) (*entities.Address, error)
}
