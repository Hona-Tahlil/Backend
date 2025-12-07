package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type AddressService struct {
	unitOfWork      ports.UnitOfWork
	userService     usecase.UserService
	provinceService usecase.ProvinceService
}

func NewAddressService(unitOfWork ports.UnitOfWork, provinceService usecase.ProvinceService, userService usecase.UserService) *AddressService {
	return &AddressService{
		unitOfWork:      unitOfWork,
		provinceService: provinceService,
		userService:     userService,
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
	user, err := as.userService.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	err = as.userService.PreloadFields(user, []string{"Address"})
	if err != nil {
		return nil, err
	}
	var mainAddress *entities.Address
	if user.Address != nil {
		mainAddress = user.Address
	}
	addresses := make([]entities.Address, 0)
	flag := false
	if user.Requests != nil {
		for _, request := range user.Requests {
			addresses = append(addresses, request.Address)
			if mainAddress != nil && request.Address.City == mainAddress.City && request.Address.Province == mainAddress.Province && request.Address.StreetAddress == mainAddress.StreetAddress {
				flag = true
			}
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

func (as *AddressService) CreateAddress(addressInfo address.AddressInfo) (*entities.Address, error) {
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
