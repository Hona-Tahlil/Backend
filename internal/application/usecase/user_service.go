package usecase

import "hona/backend/internal/application/dto/user"

type UserService interface {
	Login(loginInfo user.LoginRequest) user.LoginResponse
	Register(registerInfo user.RegisterRequest) error
	ResetPassword(resetPasswordInfo user.ResetPasswordRequest) error
	ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error
	VerifyEmail(info user.VerifyEmailRequest) error
	SendVerificationEmail(info user.SendVerificationEmailRequest) error
}
