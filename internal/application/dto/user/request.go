package user

import (
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"time"
)

type LoginRequest struct {
	Email      string
	Password   string
	RememberMe bool
}
type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type VerifyEmailRequest struct {
	Email string
	Token string
}
type ForgotPasswordRequest struct {
	Email string
}

type ResetPasswordRequest struct {
	Email    string
	Password string
	Token    string
}

type SendVerificationEmailRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type GetProfileRequest struct {
	UserID uint
}

type GetIdentityRequest struct {
	UserID uint
}

type UpdateProfileRequest struct {
	UserID        uint
	FirstName     string
	LastName      string
	Phone         *string
	Gender        enums.Gender
	BirthDate     *time.Time
	Province      enums.Province
	City          enums.City
	StreetAddress string
	HouseNumber   uint
	Unit          uint
	PostalCode    *string
	ProfilePic    *multipart.FileHeader
}
