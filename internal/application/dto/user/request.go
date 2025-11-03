package user

type LoginRequest struct {
	Email    string
	Password string
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
	Email string
}
