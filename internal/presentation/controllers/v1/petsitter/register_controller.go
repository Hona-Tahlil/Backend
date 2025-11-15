package petsitter

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterController struct {
	petSitterService *service.PetSitterService
}

func NewPetsitterController(petSitterService *service.PetSitterService) *PetSitterController {
	return &PetSitterController{
		petSitterService: petSitterService,
	}
}

// func (pc *PetSitterController) CreateSignupSession(ctx *gin.Context) {
// 	userID := controllers.GetID(ctx)

// 	res, err := pc.petSitterService.CreateSignupSession(userID)
// 	if err != nil {
// 		controllers.RespondError(ctx, err)
// 		return
// 	}

// 	controllers.Respond(ctx, 200, controllers.Message{Message: "Signup session created"}, res)
// }

func (pc *PetSitterController) CreateSignupSession(ctx *gin.Context) {
	// userID, _ := ctx.Get(bootstrap.Run().Constants.Context.ID)
	// UserID, _ := userID.(uint)
	UserID := controllers.GetID(ctx)
	petSitterInfo := petsitter.GetPetSitterRequest{
		UserID: UserID,
	}
	res, err := pc.petSitterService.GetAndUpdateUser(petSitterInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterController) GetPersonalInfo(ctx *gin.Context) {
	// userID, _ := ctx.Get(bootstrap.Run().Constants.Context.ID)
	// UserID, _ := userID.(uint)
	UserID := controllers.GetID(ctx)
	petSitterInfo := petsitter.GetPetSitterRequest{
		UserID: UserID,
	}
	res, err := pc.petSitterService.GetPersonalInfo(petSitterInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (pc *PetSitterController) SubmitPersonalInfo(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	type FirstSubmitParams struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required"`
	}
	params := controllers.Receive[FirstSubmitParams](ctx)
	FirstSubmitInfo := petsitter.FirstSubmit{
		UserID:    UserID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
	}
	if err := pc.petSitterService.SubmitPersonalInfo(FirstSubmitInfo); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterController) UploadDocuments(ctx *gin.Context) {
	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.PetProfilePic)
	if err != nil {
		file = nil
	}

	if err := pc.petSitterService.UploadDocuments(SecondSubmitInfo); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterController) GetDocuments(ctx *gin.Context) {

}

func (pc *PetSitterController) SubmitSkills(ctx *gin.Context) {
	userID := controllers.GetID(ctx)

	type SkillsParams struct {
		Bio      string   `json:"bio" validate:"required"`
		Skills   []string `json:"skills" validate:"required,min=1"`
		Services []string `json:"services" validate:"required,min=1"`
	}

	params := controllers.Receive[SkillsParams](ctx)

	dto := petsitter.SkillsRequest{
		UserID:   userID,
		Bio:      params.Bio,
		Skills:   params.Skills,
		Services: params.Services,
	}

	res, err := pc.petSitterService.SubmitSkills(dto)
	if err != nil {
		controllers.RespondError(ctx, err)
		return
	}

	controllers.Respond(ctx, 200, controllers.Message{Message: "skills updated"}, res)
}

func (pc *PetSitterController) GetSkills(ctx *gin.Context) {
	userID := controllers.GetID(ctx)

	res, err := pc.petSitterService.GetSkills(userID)
	if err != nil {
		controllers.RespondError(ctx, err)
		return
	}

	controllers.Respond(ctx, 200, controllers.Message{}, res)
}

func (pc *PetSitterController) GetSignupStatus(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	petSitterInfo := petsitter.GetPetSitterRequest{
		UserID: UserID,
	}
	res, err := pc.petSitterService.GetAndUpdateUser(petSitterInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
