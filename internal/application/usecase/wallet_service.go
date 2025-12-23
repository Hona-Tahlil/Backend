package usecase

import (
	"hona/backend/internal/application/dto/wallet"
	"hona/backend/internal/domain/ports"
)

type WalletService interface {
	GetWallet(userID uint) (*wallet.WalletResponse, error)
	TopUp(info wallet.TopUpRequest) (*wallet.WalletResponse, error)
	Withdraw(info wallet.WithdrawRequest) (*wallet.WalletResponse, error)
	ListTransfers(info wallet.HistoryRequest) ([]wallet.TransferHistoryItemResponse, int64, error)
	ListTransactions(info wallet.HistoryRequest) ([]wallet.TransactionHistoryItemResponse, int64, error)
	TransferInTransaction(rf ports.RepositoryFactory, senderUserID, receiverUserID, amount uint) (*wallet.WalletResponse, *wallet.WalletResponse, error)
}
