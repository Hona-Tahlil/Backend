package user

import "github.com/gin-gonic/gin"

type UserPetController struct {
}

func NewUserPetController() *UserPetController {
	return &UserPetController{}
}

// TODO: Add Pet
func (uc *UserPetController) AddPet(ctx *gin.Context) {
	type AddPetParams struct {
		Name string `json:"name" validate:"required"`
		// TODO: handle pic
	}
}

// TODO: Update Pet

// TODO: Remove Pet

// TODO: Get Pets Basic Data

// TODO: Get Pet's Full Data With ID
