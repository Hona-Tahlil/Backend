package domainpostgres

import "hona/backend/internal/domain/entities"

type RequestRepository interface {
	CreateRequest(request *entities.Request) error
}
