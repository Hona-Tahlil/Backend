package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entities.Transaction) error
	GetTransactionsByWalletID(walletID uint, options *postgres.QueryOptions) ([]entities.Transaction, int64, error)
}
