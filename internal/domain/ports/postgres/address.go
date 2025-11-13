package domainpostgres

import "hona/backend/internal/domain/entities"

type AddressRepository interface {
	FindAddressByID(id uint) (*entities.Address, error)
	FindAddressesByUserID(id uint) ([]entities.Address, error)
}
