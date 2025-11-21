package domainpostgres

import "hona/backend/internal/domain/entities"

type RequestRepository interface {
	CreateRequest(request *entities.Request) error
	GetRequestByID(requestID uint) (*entities.Request, error)
	EditRequest(request *entities.Request) error
	PreloadFields(request *entities.Request, fields []string) error
}
