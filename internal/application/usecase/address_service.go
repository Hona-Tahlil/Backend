package usecase

import (
	"hona/backend/internal/application/dto/address"
	// "hona/backend/internal/application/dto/request"
	"hona/backend/internal/domain/entities"
)

type AddressService interface {
	FindAddressByID(id uint) (*entities.Address, error)
	GetUserAddressesInfo(id uint) ([]address.AddressInfoResponse, error)
	GetUserAddressInfo(addressEntity *entities.Address) address.AddressInfoResponse
	// CreateAddress(addressInfo request.AddressInfoRequest) (*entities.Address, error)
	CreateAddress(addressInfo address.AddressInfo) (*entities.Address, error)
}
