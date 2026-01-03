package middleware

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
}

func NewRateLimit() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}

func (rl *RateLimitMiddleware) RateLimit(c *gin.Context) {
	limit := bootstrap.Run().Env.RateLimit.Limit
	burst := bootstrap.Run().Env.RateLimit.Burst
	limiter := rate.NewLimiter(rate.Limit(limit), burst)
	if !limiter.Allow() {
		rateLimitError := exceptions.NewRequestRateLimitError()
		panic(rateLimitError)
	}
	c.Next()
}
