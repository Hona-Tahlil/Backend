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
type VerifyEamilRequest struct {
	Email  string
	OTP    string
}
type ForgotPasswordRequest struct {
	Email string
}