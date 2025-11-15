package usecase

import "hona/backend/internal/application/dto/request"

type RequestService interface {
	CreateRequest(info request.CreateRequestRequest) error
	GetCreateRequestInfo(info request.GetCreateRequestInfoRequest) (*request.CreateRequestInfoResponse, error)
	EditRequest(info request.EditRequestRequest) error
	CancelRequest(info request.CancelRequestRequest) error
	GetRequestFullData(info request.GetRequestFullDataRequest) (*request.RequestFullDataResponse, error)
	RespondToRequest(info request.RespondToRequestRequest) error
}
