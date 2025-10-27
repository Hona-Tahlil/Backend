package domainjwt

type JWTService interface {
	GenerateTokens(userID uint, rememberMe bool) (accessTokenString string, refreshTokenString string, expireTime int)
	ValidateToken(tokenString string, tokenType string) uint
	RefreshTokens(refreshTokenString string) (accessTokenString string, newRefreshTokenString string, userID uint, expireTime int)
}
