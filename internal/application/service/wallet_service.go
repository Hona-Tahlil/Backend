package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/wallet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainpostgres "hona/backend/internal/domain/ports/postgres"
)

type WalletService struct {
	unitOfWork       ports.UnitOfWork
	petSitterService usecase.PetSitterService
}

func NewWalletService(unitOfWork ports.UnitOfWork, petSitterService usecase.PetSitterService) *WalletService {
	return &WalletService{
		unitOfWork:       unitOfWork,
		petSitterService: petSitterService,
	}
}

func (ws *WalletService) TopUp(info wallet.TopUpRequest) (*wallet.WalletResponse, error) {
	var response *wallet.WalletResponse
	err := ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		walletRepo := rf.WalletRepository()
		transactionRepo := rf.TransactionRepository()

		foundWallet, err := ws.loadWalletByUserID(walletRepo, info.UserID)
		if err != nil {
			return err
		}

		foundWallet.Balance += info.Amount
		if err := walletRepo.SaveWallet(foundWallet); err != nil {
			return err
		}

		transaction := &entities.Transaction{
			Type:     enums.PayIn,
			WalletID: foundWallet.ID,
			Amount:   info.Amount,
		}
		if err := transactionRepo.CreateTransaction(transaction); err != nil {
			return err
		}

		response = ws.toWalletResponse(foundWallet)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ws *WalletService) Withdraw(info wallet.WithdrawRequest) (*wallet.WalletResponse, error) {
	var response *wallet.WalletResponse
	if _, err := ws.petSitterService.GetPetSitterByUserID(info.UserID); err != nil {
		return nil, err
	}
	err := ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		walletRepo := rf.WalletRepository()
		transactionRepo := rf.TransactionRepository()

		foundWallet, err := ws.loadWalletByUserID(walletRepo, info.UserID)
		if err != nil {
			return err
		}

		if err := ws.ensureSufficientBalance(foundWallet, info.Amount); err != nil {
			return err
		}

		foundWallet.Balance -= info.Amount
		if err := walletRepo.SaveWallet(foundWallet); err != nil {
			return err
		}

		transaction := &entities.Transaction{
			Type:     enums.Withdraw,
			WalletID: foundWallet.ID,
			Amount:   info.Amount,
		}
		if err := transactionRepo.CreateTransaction(transaction); err != nil {
			return err
		}

		response = ws.toWalletResponse(foundWallet)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ws *WalletService) TransferInTransaction(rf ports.RepositoryFactory, senderUserID, receiverUserID, amount uint) (*wallet.WalletResponse, *wallet.WalletResponse, error) {
	walletRepo := rf.WalletRepository()

	senderWallet, err := ws.loadWalletByUserID(walletRepo, senderUserID)
	if err != nil {
		return nil, nil, err
	}

	receiverWallet, err := ws.loadWalletByUserID(walletRepo, receiverUserID)
	if err != nil {
		return nil, nil, err
	}

	if err := ws.ensureSufficientBalance(senderWallet, amount); err != nil {
		return nil, nil, err
	}

	senderWallet.Balance -= amount
	receiverWallet.Balance += amount

	if err := walletRepo.SaveWallet(senderWallet); err != nil {
		return nil, nil, err
	}
	if err := walletRepo.SaveWallet(receiverWallet); err != nil {
		return nil, nil, err
	}

	return ws.toWalletResponse(senderWallet), ws.toWalletResponse(receiverWallet), nil
}

func (ws *WalletService) loadWalletByUserID(walletRepo domainpostgres.WalletRepository, userID uint) (*entities.Wallet, error) {
	foundWallet, err := walletRepo.FindWalletByUserID(userID)
	if err != nil {
		return nil, err
	}
	if foundWallet == nil {
		foundWallet = &entities.Wallet{
			UserID: userID,
		}
		if err := walletRepo.SaveWallet(foundWallet); err != nil {
			return nil, err
		}
	}
	return foundWallet, nil
}

func (ws *WalletService) ensureSufficientBalance(wallet *entities.Wallet, amount uint) error {
	if wallet.Balance >= amount {
		return nil
	}
	var ve exceptions.ValidationErrors
	ve.AddError(bootstrap.Run().Constants.ErrorFields.Wallet, bootstrap.Run().Constants.ErrorTags.InsufficientBalance)
	return &ve
}

func (ws *WalletService) toWalletResponse(walletEntity *entities.Wallet) *wallet.WalletResponse {
	return &wallet.WalletResponse{
		ID:             walletEntity.ID,
		Balance:        walletEntity.Balance,
		PendingBalance: walletEntity.PendingBalance,
		PaymentInfo:    walletEntity.PaymentInfo,
		UserID:         walletEntity.UserID,
	}
}
