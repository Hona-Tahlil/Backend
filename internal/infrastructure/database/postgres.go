package database

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/entity"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *PostgresDatabase
var dbOnce sync.Once

type Database interface {
	GetDB() *gorm.DB
	WithTransaction(fn func(Database) error) error
}

type PostgresDatabase struct {
	DB *gorm.DB
}

func (pgx *PostgresDatabase) GetDB() *gorm.DB {
	// TODO: ? was dbInstance.DB - which is better?
	return pgx.DB
}

func (pgx *PostgresDatabase) WithTransaction(fn func(Database) error) error {
	return pgx.DB.Transaction(func(tx *gorm.DB) error {
		txWrapper := &PostgresDatabase{DB: tx}
		return fn(txWrapper)
	})
}

func NewPostgresDatabase() *PostgresDatabase {
	dsn := bootstrap.ProjectConfig.Env.DSN

	dbOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database"))
		}

		dbInstance = &PostgresDatabase{DB: db}

		dbInstance.DB.AutoMigrate(&entity.User{})

	})

	return dbInstance
}
