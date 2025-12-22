package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

func (tr *TransferRepository) CreateTransfer(transfer *entities.Transfer) error {
	return tr.db.Create(transfer).Error
}
