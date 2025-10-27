package user

import (
	"hona/backend/internal/domain/enums"
	"time"

	"github.com/gin-gonic/gin"
)

type UserPetController struct {
}

func NewUserPetController() *UserPetController {
	return &UserPetController{}
}

// TODO: Add Pet
func (uc *UserPetController) AddPet(ctx *gin.Context) {
	type AddPetParams struct {
		Name      string          `json:"name" validate:"required,min=1,max=100"`
		Kind      enums.PetKind   `json:"kind" validate:"required"`
		Species   enums.Species   `json:"species" validate:"required,min=1,max=100"` // TODO: exact number for max
		BirthDate *time.Time      `json:"birthDate" validate:"omitempty,datetime"`
		IsAdult   *bool           `json:"isAdult" validate:"omitempty"`
		Gender    enums.PetGender `json:"gender" validate:"omitempty,min=1,max=3"`
		Weight    *float32        `json:"weight" validate:"omitempty,min=0.1,max=500"`
		AboutPet  *string         `json:"aboutPet" validate:"omitempty,max=10000"`
		// TODO: handle pic
	}

}

// TODO: Update Pet

// TODO: Remove Pet

// TODO: Get Pets Basic Data

// TODO: Get Pet's Full Data With ID
