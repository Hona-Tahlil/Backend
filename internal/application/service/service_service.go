package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
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
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Service, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	return service, nil
}

func (ss *ServiceService) GetServiceResponse(serviceEntity *entities.Service) servicedto.ServiceInfoResponse {
	return servicedto.ServiceInfoResponse{
		ID:          serviceEntity.ID,
		Type:        serviceEntity.Type.String(),
		Description: serviceEntity.Description,
		Price:       serviceEntity.Price,
		PetKinds:    serviceEntity.PetKinds,
	}
}
