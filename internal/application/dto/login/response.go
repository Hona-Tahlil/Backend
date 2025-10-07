package login

type LoginResponse struct {
	JWTToken string `json:jwt_token`
}
