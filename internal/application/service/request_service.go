package service

import "hona/backend/internal/application/dto/request"

type RequestService struct {
}

func NewRequestService() *RequestService {
	return &RequestService{}
}
func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	return nil
}
