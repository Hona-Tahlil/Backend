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

func (userCache *UserCacheRepository) Get(ctx context.Context, key string) (*user.OTPData, error) {
	value, err := userCache.rdb.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var otpData user.OTPData
	if err = json.Unmarshal([]byte(value), &otpData); err != nil {
		return nil, err
	}

	return &otpData, nil

}

func (userCache *UserCacheRepository) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	otpData := user.OTPData{
		OTP:      otp,
		Attempts: 0,
	}
	value, err := json.Marshal(otpData)
	if err != nil {
		return err
	}
	err = userCache.rdb.GetRDB().Set(ctx, key, string(value), expiration).Err()
	if err != nil {

		return err
	}
	return nil
}
