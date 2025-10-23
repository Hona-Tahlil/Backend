package jwt

import (
	"errors"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: interface

type JWTService struct {
	keyManager *JWTKeyManager
}

func NewJWTService(keyManager *JWTKeyManager) *JWTService {
	return &JWTService{
		keyManager: keyManager,
	}
}

func (js *JWTService) GenerateTokens(userID uint, rememberMe bool) (accessTokenString string, refreshTokenString string, expireTime int) {
	accessTokenClaims, refreshTokenClaims, expireTime := js.generateClaims(userID, rememberMe)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)

	var err error
	accessTokenString, err = accessToken.SignedString(js.keyManager.GetPrivateKey())
	if err != nil {
		panic(err)
	}
	refreshTokenString, err = refreshToken.SignedString(js.keyManager.GetPrivateKey())
	if err != nil {
		panic(err)
	}

	return
}

func (js *JWTService) generateClaims(userID uint, rememberMe bool) (accessTokenClaims jwt.MapClaims, refreshTokenClaims jwt.MapClaims, expireTime int) {
	if rememberMe {
		expireTime = int(time.Now().Add(time.Hour * time.Duration(bootstrap.Run().Env.TokenExpires.LongRefreshHours)).Unix())
	} else {
		expireTime = int(time.Now().Add(time.Hour * time.Duration(bootstrap.Run().Env.TokenExpires.ShortRefreshHours)).Unix())
	}
	accessTokenClaims = jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Minute * time.Duration(bootstrap.Run().Env.TokenExpires.AccessMinutes)).Unix(),
		"iat":  time.Now().Unix(),
		"type": bootstrap.Run().Constants.JWTConstants.AccessTokenType,
	}

	refreshTokenClaims = jwt.MapClaims{
		"sub":  userID,
		"exp":  expireTime,
		"iat":  time.Now().Unix(),
		"type": bootstrap.Run().Constants.JWTConstants.RefreshTokenType,
	}

	return
}

func (js *JWTService) ValidateToken(tokenString string, tokenType string) uint {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, exceptions.NewInvalidTokenError()
		}
		return js.keyManager.GetPublicKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			expiredTokenErr := exceptions.NewExpiredTokenError()
			panic(expiredTokenErr)
		}
		invalidTokenErr := exceptions.NewInvalidTokenError()
		panic(invalidTokenErr)
	}

	if !token.Valid {
		invalidTokenErr := exceptions.NewInvalidTokenError()
		panic(invalidTokenErr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		invalidTokenErr := exceptions.NewInvalidTokenError()
		panic(invalidTokenErr)
	}

	if claims["type"] != tokenType {
		invalidTokenErr := exceptions.NewInvalidTokenError()
		panic(invalidTokenErr)
	}

	userID := uint(claims["sub"].(float64))
	return userID
}

func (js *JWTService) RefreshTokens(refreshTokenString string) (accessTokenString string, newRefreshTokenString string, userID uint, expireTime int) {
	userID = js.ValidateToken(refreshTokenString, bootstrap.Run().Constants.JWTConstants.RefreshTokenType)

	accessTokenString, newRefreshTokenString, expireTime = js.GenerateTokens(userID, false)
	return
}
