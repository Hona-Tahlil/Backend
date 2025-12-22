package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (tr *TransactionRepository) CreateTransaction(transaction *entities.Transaction) error {
	return tr.db.Create(transaction).Error
}
