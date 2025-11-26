package general

import (
	"hona/backend/internal/application/dto/provincecity"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralProvinceController struct {
	provinceService usecase.ProvinceService
}

func NewGeneralProvinceController(provinceService usecase.ProvinceService) *GeneralProvinceController {
	return &GeneralProvinceController{
		provinceService: provinceService,
	}
}

func (gpc *GeneralProvinceController) GetAllProvinces(ctx *gin.Context) {
	res, err := gpc.provinceService.GetAllProvincesResponse()
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (gpc *GeneralProvinceController) GetCitiesByProvinceName(ctx *gin.Context) {
	type Params struct {
		ProvinceNum enums.Province `uri:"province_num" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	info := provincecity.GetProvinceCitiesRequest{
		ProvinceNum: params.ProvinceNum,
	}
	res, err := gpc.provinceService.GetCitiesByProvinceName(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
