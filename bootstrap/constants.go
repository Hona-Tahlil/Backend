package bootstrap

type Constants struct {
	Context     Context
	JWTKeysPath JWTKeysPath
}

type JWTKeysPath struct {
	PublicKey  string
	PrivateKey string
}

type Context struct {
	Translator string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/infrastructure/jwt/public_key.pem",
			PrivateKey: "./internal/infrastructure/jwt/private_key.pem",
		},
	}
}
