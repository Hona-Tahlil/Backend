package usecase

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/dto/provincecity"

	// "hona/backend/internal/application/dto/request"
	"hona/backend/internal/domain/entities"
)

type AddressService interface {
	FindAddressByID(id uint) (*entities.Address, error)
	GetUserAddressesInfo(id uint) ([]address.AddressInfoResponse, error)
	GetUserAddressInfo(addressEntity *entities.Address) address.AddressInfoResponse
	CreateAddressEntity(addressInfo address.AddressInfo) (*entities.Address, error)
	GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error)
	GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error)
}
