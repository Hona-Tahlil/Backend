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
	OTP   string
}
type ForgotPasswordRequest struct {
	Email string
}
