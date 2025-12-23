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

func (tr *TransactionRepository) GetTransactionsByWalletID(walletID uint, options *QueryOptions) ([]entities.Transaction, int64, error) {
	if options == nil {
		options = NewQueryOptions()
	}

	query := tr.db.Model(&entities.Transaction{}).Where("wallet_id = ?", walletID)
	query, total := ApplyModifiers(query, *options)

	var transactions []entities.Transaction
	if err := query.Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
