package domainpostgres

import "hona/backend/internal/domain/entities"

type ServiceRepository interface {
	FindServiceByID(id uint) (*entities.Service, error)
}
