package domainpostgres

import "hona/backend/internal/domain/entities"

type TransactionRepository interface {
	CreateTransaction(transaction *entities.Transaction) error
}
