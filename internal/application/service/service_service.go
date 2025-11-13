package service

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/ports"
)

type ServiceService struct {
	unitOfWork ports.UnitOfWork
}

func NewServiceService(unitOfWork ports.UnitOfWork) *ServiceService {
	return &ServiceService{
		unitOfWork: unitOfWork,
	}
}

func (ss *ServiceService) FindServiceByID(id uint) (*entities.Service, error) {
	serviceRepo := ss.unitOfWork.Factory().ServiceRepository()
	service, err := serviceRepo.FindServiceByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, fmt.Errorf("service not found")
	}

	return service, nil
}
