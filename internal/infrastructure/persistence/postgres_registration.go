package persistence

import (
	"fmt"
	"hona/backend/bootstrap"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var dbOnce sync.Once

func NewPostgresDatabase() *gorm.DB {
	dbConfig := bootstrap.ProjectConfig.Env.PrimaryDB
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
		dbInstance = db
	})

	return dbInstance
}
