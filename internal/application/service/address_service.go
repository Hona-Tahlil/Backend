package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/dto/provincecity"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type AddressService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
}

func NewAddressService(unitOfWork ports.UnitOfWork, userService usecase.UserService) *AddressService {
	return &AddressService{
		unitOfWork:  unitOfWork,
		userService: userService,
	}
}

func (as *AddressService) FindAddressByID(id uint) (*entities.Address, error) {
	addressRepo := as.unitOfWork.Factory().AddressRepository()
	address, err := addressRepo.FindAddressByID(id)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Address)
	}

	return address, nil
}

func (as *AddressService) GetUserAddressesInfo(id uint) ([]address.AddressInfoResponse, error) {
	addressRepo := as.unitOfWork.Factory().AddressRepository()
	mainAddress, err := addressRepo.FindUserAddressByUserID(id)
	if err != nil {
		return nil, err
	}
	addresses, err := addressRepo.FindUserRequestAddressesByID(id)
	if err != nil {
		return nil, err
	}

	flag := false
	for _, address := range addresses {
		if mainAddress != nil && address.City == mainAddress.City && address.Province == mainAddress.Province && address.StreetAddress == mainAddress.StreetAddress {
			flag = true
		}
	}

	if !flag {
		addresses = append(addresses, *mainAddress)
	}

	r := make([]address.AddressInfoResponse, 0)
	for _, address := range addresses {
		r = append(r, as.GetUserAddressInfo(&address))
	}

	return r, nil
}

func (as *AddressService) GetUserAddressInfo(addressEntity *entities.Address) address.AddressInfoResponse {
	return address.AddressInfoResponse{
		ID:            addressEntity.ID,
		ProvinceName:  addressEntity.Province.String(),
		CityName:      addressEntity.City.String(),
		StreetAddress: addressEntity.StreetAddress,
		HouseNumber:   addressEntity.HouseNumber,
		Unit:          addressEntity.Unit,
		PostalCode:    addressEntity.PostalCode,
	}
}

func (as *AddressService) CreateAddressEntity(addressInfo address.AddressInfo) (*entities.Address, error) {
	cities := enums.ProvinceWithCities[addressInfo.ProvinceName]
	var foundCity enums.City
	for _, city := range cities {
		if city == addressInfo.CityName {
			foundCity = city
			break
		}
	}
	if foundCity == 0 {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.City)
	}

	address := &entities.Address{
		Province:      addressInfo.ProvinceName,
		City:          foundCity,
		StreetAddress: addressInfo.StreetAddress,
		HouseNumber:   addressInfo.HouseNumber,
		Unit:          addressInfo.Unit,
		PostalCode:    addressInfo.PostalCode,
	}

	return address, nil
}

func (as *AddressService) GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error) {
	provinces := enums.GetAllProvinces()

	provinceResponses := make([]provincecity.ProvinceResponse, 0)
	for _, province := range provinces {
		provinceResponses = append(provinceResponses, provincecity.ProvinceResponse{
			Num:  province,
			Name: province.String(),
		})
	}

	return provinceResponses, nil
}

func (as *AddressService) GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error) {
	cities := enums.ProvinceWithCities[info.ProvinceNum]
	cityResponses := make([]provincecity.CityResponse, 0)
	for _, city := range cities {
		cityResponses = append(cityResponses, provincecity.CityResponse{
			Num:  city,
			Name: city.String(),
		})
	}

	return cityResponses, nil
}

func (as *AddressService) FindRequestAddressByID(requestID uint) (*entities.Address, error) {
	addressRepo := as.unitOfWork.Factory().AddressRepository()
	address, err := addressRepo.FindRequestAddressByID(requestID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Address)
	}

	return address, nil
}
