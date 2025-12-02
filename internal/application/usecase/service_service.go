package usecase

import (
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/domain/entities"
)

type ServiceService interface {
	FindServiceByID(id uint) (*entities.Service, error)
	GetServiceResponse(serviceEntity *entities.Service) servicedto.ServiceInfoResponse
}
