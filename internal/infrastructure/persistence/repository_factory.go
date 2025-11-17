package persistence

import (
	domainpostgres "hona/backend/internal/domain/ports/postgres"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"

	"gorm.io/gorm"
)

type RepositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

func (f *RepositoryFactory) UserRepository() domainpostgres.UserRepository {
	return postgres.NewUserRepository(f.db)
}

func (f *RepositoryFactory) RBACRepository() domainpostgres.RBACRepository {
	return postgres.NewRBACRepository(f.db)
}

func (f *RepositoryFactory) PetSitterRepository() domainpostgres.PetSitterRepository {
	return postgres.NewPetSitterRepository(f.db)
}

func (f *RepositoryFactory) AddressRepository() domainpostgres.AddressRepository {
	return postgres.NewAddressRepository(f.db)
}