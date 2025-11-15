package middleware

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (am *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	authHeader := ctx.GetHeader(bootstrap.Run().Constants.Context.Authorization)
	if authHeader == "" {
		unauthorizedError := exceptions.NewUnauthorizedError("empty auth header")
		panic(unauthorizedError)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError := exceptions.NewUnauthorizedError("invalid token format")
		panic(unauthorizedError)
	}

	tokenString := parts[1]
	if tokenString == "" {
		unauthorizedError := exceptions.NewUnauthorizedError("empty token")
		panic(unauthorizedError)
	}

	userID := am.jwtService.ValidateToken(tokenString, bootstrap.Run().Constants.JWTConstants.AccessTokenType)

	ctx.Set(bootstrap.Run().Constants.Context.ID, userID)

	ctx.Next()
}
