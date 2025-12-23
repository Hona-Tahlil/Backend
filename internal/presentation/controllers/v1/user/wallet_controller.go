package user

import (
	"hona/backend/internal/application/dto/wallet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type UserWalletController struct {
	walletService usecase.WalletService
}

func NewUserWalletController(walletService usecase.WalletService) *UserWalletController {
	return &UserWalletController{
		walletService: walletService,
	}
}

func (wc *UserWalletController) TopUp(ctx *gin.Context) {
	type Params struct {
		Amount uint `json:"amount" validate:"required,min=1"`
	}
	params := controllers.Receive[Params](ctx)
	userID := controllers.GetID(ctx)

	info := wallet.TopUpRequest{
		UserID: userID,
		Amount: params.Amount,
	}
	res, err := wc.walletService.TopUp(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
