package user

type LoginResponse struct {
	JWTToken string `json:jwt_token`
}
