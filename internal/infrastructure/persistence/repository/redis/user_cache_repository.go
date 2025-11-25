package redis

import (
	"context"
	"encoding/json"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/infrastructure/persistence"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	rdb persistence.Cache
}

func NewUserCacheRepository(rdb persistence.Cache) *UserCacheRepository {
	return &UserCacheRepository{
		rdb: rdb,
	}
}

func (userCache *UserCacheRepository) Get(ctx context.Context, key string) (*user.MLData, error) {
	value, err := userCache.rdb.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var mlData user.MLData
	if err = json.Unmarshal([]byte(value), &mlData); err != nil {
		return nil, err
	}

	return &mlData, nil

}

func (uc *UserCacheRepository) Set(ctx context.Context, key, token string, expiration time.Duration) error {
	mldata := user.MLData{
		Token: token,
	}
	value, err := json.Marshal(mldata)
	if err != nil {
		return err
	}
	err = uc.rdb.GetRDB().Set(ctx, key, string(value), expiration).Err()
	if err != nil {

		return err
	}
	return nil
}

