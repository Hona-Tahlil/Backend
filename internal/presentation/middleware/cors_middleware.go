package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSMiddleware struct {
}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (c *CORSMiddleware) CORS() gin.HandlerFunc {
	return cors.New(c.config())
}

func (c *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := c.config()
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if config.AllowAllOrigins {
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Add("Vary", "Origin")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}

		if config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if len(config.AllowMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if len(config.AllowHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", "*")
		}
		if len(config.ExposeHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		}

		next.ServeHTTP(w, r)
	})
}

func (c *CORSMiddleware) config() cors.Config {
	return cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
