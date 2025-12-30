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

func (wc *UserWalletController) GetWallet(ctx *gin.Context) {
	userID := controllers.GetID(ctx)

	res, err := wc.walletService.GetWallet(userID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (wc *UserWalletController) ListTransfers(ctx *gin.Context) {
	type Params struct {
		Page  int `form:"page" validate:"min=1"`
		Count int `form:"count" validate:"min=1,max=100"`
	}
	params := controllers.Receive[Params](ctx)
	userID := controllers.GetID(ctx)
	offset, limit := controllers.GetOffsetLimit(params.Page, params.Count)

	info := wallet.HistoryRequest{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}

	items, total, err := wc.walletService.ListTransfers(info)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(items, total, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}

func (wc *UserWalletController) ListTransactions(ctx *gin.Context) {
	type Params struct {
		Page  int `form:"page" validate:"min=1"`
		Count int `form:"count" validate:"min=1,max=100"`
	}
	params := controllers.Receive[Params](ctx)
	userID := controllers.GetID(ctx)
	offset, limit := controllers.GetOffsetLimit(params.Page, params.Count)

	info := wallet.HistoryRequest{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}

	items, total, err := wc.walletService.ListTransactions(info)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(items, total, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}
