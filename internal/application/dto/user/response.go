package user

type LoginResponse struct {
	JWTToken string `json:jwt_token`
}

type MLData struct {
	Token      string `json:"ml"`
}


