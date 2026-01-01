package request

import (
	"hona/backend/internal/application/dto/address"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/servicedto"
	"time"
)

type CreateRequestInfoResponse struct {
	Services           []servicedto.ServiceInfoResponse        `json:"services"`
	Addresses          []address.AddressInfoResponse           `json:"addresses"`
	Pets               []pet.PetBasicDataResponse              `json:"pets"`
	FreeCalendarSlots  []calendarslot.CalendarSlotInfoResponse `json:"freeCalendarSlots"`
	PetSitterFirstName string                                  `json:"petSitterFirstName"`
	PetSitterLastName  string                                  `json:"petSitterLastName"`
}

type RequestFullDataResponse struct {
	RequestID          uint                                    `json:"requestID"`
	PetSitterUserID    uint                                    `json:"petSitterUserID"`
	PetSitterFirstName string                                  `json:"petSitterFirstName"`
	PetSitterLastName  string                                  `json:"petSitterLastName"`
	UserFirstName      string                                  `json:"userFirstName"`
	UserLastName       string                                  `json:"userLastName"`
	Service            servicedto.ServiceInfoResponse          `json:"service"`
	Pets               []pet.PetBasicDataResponse              `json:"pets"`
	CalendarSlots      []calendarslot.CalendarSlotInfoResponse `json:"calendarSlots"`
	Notes              *string                                 `json:"notes"`
	TotalPrice         uint                                    `json:"totalPrice"`
	Comment            comment.CommentResponse                 `json:"comment"`
	Address            address.AddressInfoResponse             `json:"address"`
	Status             string                                  `json:"status"`
	TransferID         *uint                                   `json:"transferID"`
	UpdatedAt          time.Time                               `json:"updatedAt"`
}

type RequestListItemResponse struct {
	RequestID          uint                           `json:"requestID"`
	PetSitterUserID    uint                           `json:"petSitterUserID"`
	PetSitterFirstName string                         `json:"petSitterFirstName"`
	PetSitterLastName  string                         `json:"petSitterLastName"`
	Service            servicedto.ServiceInfoResponse `json:"service"`
	TotalPrice         uint                           `json:"totalPrice"`
	Status             string                         `json:"status"`
	UpdatedAt          time.Time                      `json:"updatedAt"`
}
