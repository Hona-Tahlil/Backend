package persistence

import (
	"context"
	"fmt"
	"hona/backend/bootstrap"
	"log"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetRDB() redis.Client
}

type RedisDatabase struct {
	RDB *redis.Client
}

var (
	rdbOnce     sync.Once
	rdbInstance *RedisDatabase
)

func NewRedisDatabase() *RedisDatabase {
	rdbOnce.Do(func() {
		rdbNumber, _ := strconv.Atoi(bootstrap.Run().Env.PrimaryRedis.RDBNumber)
		address := fmt.Sprintf("%s:%s", bootstrap.Run().Env.PrimaryRedis.Address, bootstrap.Run().Env.PrimaryRedis.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: bootstrap.Run().Env.PrimaryRedis.Password,
			DB:       rdbNumber,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal("Error connecting to Redis:", err)
		}
		rdbInstance = &RedisDatabase{RDB: rdb}
	})

	return rdbInstance
}

func (rdb *RedisDatabase) GetRDB() redis.Client {
	return *rdbInstance.RDB
}
