package exceptions

import (
	"fmt"
	"hona/backend/bootstrap"
)

var errTags = bootstrap.Run().Constants.ErrorTags

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
		Type:    errTags.InvalidAuthCredentials,
		Message: message,
	}
}

func NewExpiredTokenError() *AuthError {
	return &AuthError{
		Type:    errTags.ExpiredAuthToken,
		Message: "Authentication token has expired",
	}
}

func NewInvalidTokenError() *AuthError {
	return &AuthError{
		Type:    errTags.InvalidAuthToken,
		Message: "Invalid authentication token",
	}
}

func NewUnauthorizedError(message string) *AuthError {
	if message == "" {
		message = "Unauthorized access"
	}

	return &AuthError{
		Type:    errTags.Unauthorized,
		Message: message,
	}
}

func NewAccessDeniedError(message string) *AuthError {
	if message == "" {
		message = "Access Denied"
	}

	return &AuthError{
		Type:    errTags.AccessDenied,
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
