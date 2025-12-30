package domainpostgres

import "hona/backend/internal/domain/entities"

type WalletRepository interface {
	FindWalletByID(id uint) (*entities.Wallet, error)
	FindWalletByUserID(userID uint) (*entities.Wallet, error)
	SaveWallet(wallet *entities.Wallet) error
}
