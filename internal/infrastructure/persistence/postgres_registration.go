package persistence

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/entities"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var dbOnce sync.Once

type Database interface {
	GetDB() *gorm.DB
	WithTransaction(fn func(Database) error) error
}

type PostgresDatabase struct {
	DB *gorm.DB
}

func (pgx *PostgresDatabase) GetDB() *gorm.DB {
	return pgx.DB
}

func (pgx *PostgresDatabase) WithTransaction(fn func(Database) error) error {
	return pgx.DB.Transaction(func(tx *gorm.DB) error {
		txWrapper := &PostgresDatabase{DB: tx}
		return fn(txWrapper)
	})
}

func NewPostgresDatabase() *PostgresDatabase {
	dbConfig := bootstrap.Run().Env.PrimaryDB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)

	dbOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database"))
		}

		dbInstance = &PostgresDatabase{DB: db}
		// dbInstance = db
		dbInstance.DB.AutoMigrate(
			&entities.User{},
			&entities.User{},
			&entities.Role{},
			&entities.Permission{},
			&entities.Wallet{},
			&entities.Request{},
			&entities.CalendarSlot{},
			&entities.Pet{},
			&entities.PetSitter{},
			&entities.Service{},
			&entities.Chat{},
			&entities.Comment{},
			&entities.Address{},
			&entities.Province{},
			&entities.City{},
			&entities.TextMessage{},
			&entities.Transaction{},
			&entities.Transfer{},
		)
		// user := entities.User{
		// 	Email: "test@email",
		// }
		// db.Create(&user)
		// permission := entities.Permission{
		// 	Type:     enums.RequestPermission,
		// 	Category: enums.ReadPermissionCategory,
		// }
		// db.Create(&permission)

	})

	return dbInstance
}
