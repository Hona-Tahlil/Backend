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

func (tr *TransferRepository) GetTransfersByWalletID(walletID uint, options *QueryOptions) ([]entities.Transfer, int64, error) {
	if options == nil {
		options = NewQueryOptions()
	}

	query := tr.db.Model(&entities.Transfer{}).
		Where("sender_wallet_id = ? OR receiver_wallet_id = ?", walletID, walletID)
	query, total := ApplyModifiers(query, *options)

	var transfers []entities.Transfer
	if err := query.Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}
