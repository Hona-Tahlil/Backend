package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type AddressService struct {
	unitOfWork ports.UnitOfWork
}

func NewAddressService(unitOfWork ports.UnitOfWork) *AddressService {
	return &AddressService{
		unitOfWork: unitOfWork,
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
