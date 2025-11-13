package service

import (
	"fmt"
	"hona/backend/internal/domain/entities"
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
		return nil, fmt.Errorf("invalid address")
	}

	return address, nil
}
