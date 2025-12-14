package domainpostgres

import "hona/backend/internal/domain/entities"

type AddressRepository interface {
	FindAddressByID(id uint) (*entities.Address, error)
	FindUserAddressByUserID(id uint) (*entities.Address, error)
	Create(address *entities.Address) error
	Update(address *entities.Address) error
	FindUserRequestAddressesByID(id uint) ([]entities.Address, error)
	FindRequestAddressByID(requestID uint) (*entities.Address, error)
}
