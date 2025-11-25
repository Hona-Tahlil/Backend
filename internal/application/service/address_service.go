package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type AddressService struct {
	unitOfWork      ports.UnitOfWork
	provinceService usecase.ProvinceService
}

func NewAddressService(unitOfWork ports.UnitOfWork, provinceService usecase.ProvinceService) *AddressService {
	return &AddressService{
		unitOfWork:      unitOfWork,
		provinceService: provinceService,
	}
}

func (as *AddressService) FindAddressByID(id uint) (*entities.Address, error) {
	addressRepo := as.unitOfWork.Factory().AddressRepository()
	address, err := addressRepo.FindAddressByID(id)
	if err != nil {
		return nil, err
	}
	if address == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Address, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	return address, nil
}

func (as *AddressService) GetUserAddressesInfo(id uint) ([]address.AddressInfoResponse, error) {
	addressRepo := as.unitOfWork.Factory().AddressRepository()
	addresses, err := addressRepo.FindAddressesByUserID(id)
	if err != nil {
		return nil, err
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
		ProvinceName:  addressEntity.Province.Name.String(),
		CityName:      addressEntity.City.Name.String(),
		StreetAddress: addressEntity.StreetAddress,
		HouseNumber:   addressEntity.HouseNumber,
		Unit:          addressEntity.Unit,
		PostalCode:    addressEntity.PostalCode,
	}
}

func (as *AddressService) CreateAddress(addressInfo address.AddressInfo) (*entities.Address, error) {
	province, err := as.provinceService.FindProvinceByName(addressInfo.ProvinceName)
	if err != nil {
		return nil, err
	}
	var foundCity *entities.City
	for _, city := range province.Cities {
		if city.Name == addressInfo.CityName {
			foundCity = &city
			break
		}
	}
	if foundCity == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.City)
	}

	address := &entities.Address{
		ProvinceID:    province.ID,
		CityID:        foundCity.ID,
		StreetAddress: addressInfo.StreetAddress,
		HouseNumber:   addressInfo.HouseNumber,
		Unit:          addressInfo.Unit,
		PostalCode:    addressInfo.PostalCode,
	}

	return address, nil
}
