package servicedto

import "hona/backend/internal/domain/enums"

type ServiceInfoResponse struct {
	ID          uint              `json:"id"`
	Type        enums.ServiceType `json:"type"`
	Description *string           `json:"description"`
	Price       uint              `json:"price"`
	PetKinds    []enums.PetKind   `json:"petKinds"`
}
