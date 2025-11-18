package petsitter

import (
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/service"
	"hona/backend/internal/domain/enums"
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

func (pc *PetSitterController) CreateSignupSession(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	petSitterInfo := petsitter.GetPetSitterRequest{
		UserID: UserID,
	}
	res, err := pc.petSitterService.CreateSignupSession(petSitterInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterController) SubmitPersonalInfo(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	type FirstSubmitParams struct {
		FirstName   string `json:"firstName" validate:"required"`
		LastName    string `json:"lastName" validate:"required"`
		Email       string `json:"email" validate:"required"`
		PhoneNumber string
	}
	params := controllers.Receive[FirstSubmitParams](ctx)
	FirstSubmitInfo := petsitter.SubmitPersonalInfoRequest{
		UserID:    UserID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Phone:     params.PhoneNumber,
	}
	err := pc.petSitterService.SubmitPersonalInfo(FirstSubmitInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterController) GetPersonalInfo(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetPersonalInfo(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterController) UploadDocuments(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	form, err := ctx.MultipartForm()
	if err != nil {
		form = nil
		panic(err)
	}
	CertificateFiles := form.File["certificateFiles"]
	Files := form.File["files"]
	UploadDocumentsInfo := petsitter.UploadDocumentsRequest{
		UserID:           UserID,
		CertificateFiles: CertificateFiles,
		Files:            Files,
	}
	err = pc.petSitterService.UploadDocuments(UploadDocumentsInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterController) GetDocuments(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetDocuments(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterController) SubmitSkills(ctx *gin.Context) {
	userID := controllers.GetID(ctx)

	type SkillsParams struct {
		Bio      string `json:"bio" validate:"required,min=20"`
		PetKinds []enums.PetKind
		Services []enums.ServiceType
	}

	params := controllers.Receive[SkillsParams](ctx)
	dto := petsitter.SubmitSkillsRequest{
		UserID:   userID,
		Bio:      params.Bio,
		Petkind:  params.PetKinds,
		Services: params.Services,
	}

	err := pc.petSitterService.SubmitSkills(dto)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterController) GetPetsitterStatus(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetPetsitterStatus(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
