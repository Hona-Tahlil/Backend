package petsitter

import (
	"hona/backend/internal/application/dto/wallet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterWalletController struct {
	walletService usecase.WalletService
}

func NewPetSitterWalletController(walletService usecase.WalletService) *PetSitterWalletController {
	return &PetSitterWalletController{
		walletService: walletService,
	}
}

func (wc *PetSitterWalletController) Withdraw(ctx *gin.Context) {
	type Params struct {
		Amount uint `json:"amount" validate:"required,min=1"`
	}
	params := controllers.Receive[Params](ctx)
	userID := controllers.GetID(ctx)

	info := wallet.WithdrawRequest{
		UserID: userID,
		Amount: params.Amount,
	}
	res, err := wc.walletService.Withdraw(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
