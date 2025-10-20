package user

type LoginResponse struct {
	JWTToken string `json:jwt_token`
}

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}