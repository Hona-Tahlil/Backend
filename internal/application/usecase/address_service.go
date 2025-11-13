package usecase

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/domain/entities"
)

type AddressService interface {
	FindAddressByID(id uint) (*entities.Address, error)
	GetUserAddressesInfo(id uint) ([]address.AddressInfoResponse, error)
	GetUserAddressInfo(addressEntity *entities.Address) address.AddressInfoResponse
}
