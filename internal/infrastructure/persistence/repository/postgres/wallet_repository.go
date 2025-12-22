package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (wr *WalletRepository) FindWalletByID(id uint) (*entities.Wallet, error) {
	var wallet entities.Wallet

	if result := wr.db.First(&wallet, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &wallet, nil
}

func (wr *WalletRepository) FindWalletByUserID(userID uint) (*entities.Wallet, error) {
	var wallet entities.Wallet

	if result := wr.db.First(&wallet, "user_id = ?", userID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &wallet, nil
}

func (wr *WalletRepository) SaveWallet(wallet *entities.Wallet) error {
	return wr.db.Save(wallet).Error
}
