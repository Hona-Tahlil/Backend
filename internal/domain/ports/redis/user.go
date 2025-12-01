package domainredis

import (
	"context"
	"hona/backend/internal/application/dto/user"
	"time"
)

type UserCacheRepository interface {
	Set(ctx context.Context, key, token string, expiration time.Duration) error
	Get(ctx context.Context, key string) (*user.MLData, error)
}
