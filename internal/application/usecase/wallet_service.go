package usecase

import (
	"hona/backend/internal/application/dto/wallet"
	"hona/backend/internal/domain/ports"
)

type WalletService interface {
	TopUp(info wallet.TopUpRequest) (*wallet.WalletResponse, error)
	Withdraw(info wallet.WithdrawRequest) (*wallet.WalletResponse, error)
	TransferInTransaction(rf ports.RepositoryFactory, senderUserID, receiverUserID, amount uint) (*wallet.WalletResponse, *wallet.WalletResponse, error)
}
