package pet

import (
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"time"
)

type AddPetRequest struct {
	UserID     uint
	Name       string
	Kind       enums.PetKind
	Species    enums.Species
	BirthDate  *time.Time
	IsAdult    bool
	Gender     enums.PetGender
	Weight     *float32
	AboutPet   *string
	ProfilePic *multipart.FileHeader
}
