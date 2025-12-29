package usecase

import (
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/domain/entities"
)

type RequestService interface {
	CreateRequest(info request.CreateRequestRequest) error
	GetCreateRequestInfo(info request.GetCreateRequestInfoRequest) (*request.CreateRequestInfoResponse, error)
	EditRequest(info request.EditRequestRequest) error
	CancelRequest(info request.CancelRequestRequest) error
	PayRequest(info request.PayRequestRequest) error
	SearchRequests(info request.SearchRequestsRequest) ([]request.RequestListItemResponse, int64, error)
	SearchPetSitterRequests(info request.SearchPetSitterRequestsRequest) ([]request.RequestListItemResponse, int64, error)
	GetRequestFullData(info request.GetRequestFullDataRequest) (*request.RequestFullDataResponse, error)
	RespondToRequest(info request.RespondToRequestRequest) error
	PreloadFields(request *entities.Request, fields []string) error
	FindRequestByID(id uint) (*entities.Request, error)
	EnsureRequestIsFinished(request *entities.Request) error
}
