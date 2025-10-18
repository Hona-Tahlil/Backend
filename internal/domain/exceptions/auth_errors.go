package exceptions

import (
	"fmt"
)

// TODO: should be in constants
const (
	ErrorTypeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrorTypeExpiredToken       = "EXPIRED_TOKEN"
	ErrorTypeInvalidToken       = "INVALID_TOKEN"
	ErrorTypeUnauthorized       = "UNAUTHORIZED"
	ErrorTypeAccessDenied       = "ACCESS_DENIED"
)

type AuthError struct {
	Type    string
	Message string
}

func (e AuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewInvalidCredentialsError(message string) *AuthError {
	if message == "" {
		message = "username and password not match"
	}
	return &AuthError{
		Type:    ErrorTypeInvalidCredentials,
		Message: message,
	}
}

func NewExpiredTokenError() *AuthError {
	return &AuthError{
		Type:    ErrorTypeExpiredToken,
		Message: "Authentication token has expired",
	}
}

func NewInvalidTokenError() *AuthError {
	return &AuthError{
		Type:    ErrorTypeInvalidToken,
		Message: "Invalid authentication token",
	}
}

func NewUnauthorizedError(message string) *AuthError {
	if message == "" {
		message = "Unauthorized access"
	}

	return &AuthError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
	}
}

func NewAccessDeniedError(message string) *AuthError {
	if message == "" {
		message = "Access Denied"
	}

	return &AuthError{
		Type:    ErrorTypeAccessDenied,
		Message: message,
	}
}

func IsAuthError(err error) bool {
	_, ok := err.(*AuthError)
	return ok
}

func GetAuthErrorType(err error) string {
	if authErr, ok := err.(*AuthError); ok {
		return authErr.Type
	}
	return ""
}
