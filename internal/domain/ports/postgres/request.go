package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type RequestRepository interface {
	CreateRequest(request *entities.Request) error
	GetRequestByID(requestID uint) (*entities.Request, error)
	EditRequest(request *entities.Request) error
	PreloadFields(request *entities.Request, fields []string) error
	SearchRequests(userID uint, options *postgres.QueryOptions) ([]entities.Request, int64, error)
	SearchRequestsByPetSitterID(petSitterID uint, options *postgres.QueryOptions) ([]entities.Request, int64, error)
	FindRequestsByUserAndPetSitter(userID uint, petSitterID uint) ([]entities.Request, error)
}
