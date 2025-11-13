package usecase

import "hona/backend/internal/domain/entities"

type ServiceService interface {
	FindServiceByID(id uint) (*entities.Service, error)
}
