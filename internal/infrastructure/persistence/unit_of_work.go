package persistence

import (
	"hona/backend/internal/domain/ports"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (u *UnitOfWork) Factory() ports.RepositoryFactory {
	return NewRepositoryFactory(u.db)
}

func (u *UnitOfWork) WithTransaction(fn func(ports.RepositoryFactory) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		factory := NewRepositoryFactory(tx)
		return fn(factory)
	})
}
