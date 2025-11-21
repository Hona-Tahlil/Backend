package ports

import (
	domainpostgres "hona/backend/internal/domain/ports/postgres"
)

type RepositoryFactory interface {
	UserRepository() domainpostgres.UserRepository
	RBACRepository() domainpostgres.RBACRepository
	PetRepository() domainpostgres.PetRepository
	ProvinceRepository() domainpostgres.ProvinceRepository
	AddressRepository() domainpostgres.AddressRepository
	ServiceRepository() domainpostgres.ServiceRepository
	RequestRepository() domainpostgres.RequestRepository
	PetSitterRepository() domainpostgres.PetSitterRepository
}
