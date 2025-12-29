package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type TransferRepository interface {
	CreateTransfer(transfer *entities.Transfer) error
	GetTransfersByWalletID(walletID uint, options *postgres.QueryOptions) ([]entities.Transfer, int64, error)
}
