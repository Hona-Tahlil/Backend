package servicedto

type ServiceInfoResponse struct {
	ID          uint    `json:"id"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
	Price       uint    `json:"price"`
	// PetKinds    []enums.PetKind `json:"petKinds"`
}
