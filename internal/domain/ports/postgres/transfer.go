package domainpostgres

import "hona/backend/internal/domain/entities"

type TransferRepository interface {
	CreateTransfer(transfer *entities.Transfer) error
}
