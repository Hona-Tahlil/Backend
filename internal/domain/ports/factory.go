package ports

import domainpostgres "hona/backend/internal/domain/ports/postgres"

type RepositoryFactory interface {
	UserRepository() domainpostgres.UserRepository
}
