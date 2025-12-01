package domainpostgres

import "hona/backend/internal/domain/entities"

type AddressRepository interface {
	FindAddressByID(id uint) (*entities.Address, error)
	FindAddressesByUserID(id uint) ([]entities.Address, error)
<<<<<<< HEAD
	FinduserAddressByUserID(id uint) (*entities.Address, error)
	Create(address *entities.Address) error
	Update(address *entities.Address) error
}
=======
}
>>>>>>> dev
