package persistence

import (
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

func (f *RepositoryFactory) UserRepository() *postgres.UserRepository {
	return postgres.NewUserRepository(f.db)
}
