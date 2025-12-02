package petsitter

import (
	"fmt"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type PetSitterRegisterController struct {
	petSitterService usecase.PetSitterService
}

func NewPetSitterRegisterController(petSitterService usecase.PetSitterService) *PetSitterRegisterController {
	return &PetSitterRegisterController{
		petSitterService: petSitterService,
	}
}

func (pc *PetSitterRegisterController) CreateSignupSession(ctx *gin.Context) {
	fmt.Println("✅ PETSITTER REGISTER CONTROLLER")
	UserID := controllers.GetID(ctx)
	fmt.Println("✅ PETSITTER REGISTER CONTROLLER")

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

func (pc *PetSitterRegisterController) SubmitPersonalInfo(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	fmt.Println("✅ 1655400")

	type FirstSubmitParams struct {
		FirstName   string         `json:"first_name"`
		LastName    string         `json:"last_name"`
		Email       string         `json:"email"`
		Gender      enums.Gender   `json:"gender"`
		BirthDate   *time.Time     `json:"birth_date"`
		PhoneNumber string         `json:"phone_number"`
		Province    enums.Province `json:"province"`
		City        enums.City     `json:"city"`
		Address     string         `json:"address"`
		HouseNumber uint           `json:"house_number"`
		Unit        uint           `json:"unit"`
	}
	params := controllers.Receive[FirstSubmitParams](ctx)
	FirstSubmitInfo := petsitter.SubmitPersonalInfoRequest{
		UserID:      UserID,
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		Gender:      params.Gender,
		BirthDate:   params.BirthDate,
		Phone:       params.PhoneNumber,
		Province:    params.Province,
		City:        params.City,
		Address:     params.Address,
		HouseNumber: params.HouseNumber,
		Unit:        params.Unit,
	}
	fmt.Println("✅ 00msbjansltngn")

	err := pc.petSitterService.SubmitPersonalInfo(FirstSubmitInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ 00msbjansltngn")

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterRegisterController) GetPersonalInfo(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	fmt.Println("✅ 00131546")
	res, err := pc.petSitterService.GetPersonalInfo(UserID)
	fmt.Println("✅ 00131547")
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ 00msbjansltngn")

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (pc *PetSitterRegisterController) UploadDocuments(ctx *gin.Context) {
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

func (pc *PetSitterRegisterController) GetDocuments(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetDocuments(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterRegisterController) SubmitSkills(ctx *gin.Context) {
	userID := controllers.GetID(ctx)

	type SkillsParams struct {
		Bio      string              `json:"bio" validate:"required,min=20"`
		PetKinds []enums.PetKind     `json:"pet_kinds" validate:"required,min=1"`
		Services []enums.ServiceType `json:"services" validate:"required,min=1"`
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

func (pc *PetSitterRegisterController) GetPetsitterStatus(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetPetsitterStatus(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
