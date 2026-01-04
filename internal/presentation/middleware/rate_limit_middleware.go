package middleware

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}

func NewRateLimit() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		mu:       sync.Mutex{},
		limiters: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimitMiddleware) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit := bootstrap.Run().Env.RateLimit.Limit
	burst := bootstrap.Run().Env.RateLimit.Burst

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(limit), burst)
		rl.limiters[ip] = limiter
	}
	return limiter
}

func (rl *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.RemoteIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			panic(exceptions.NewRequestRateLimitError())
		}
		ctx.Next()
	}
}
