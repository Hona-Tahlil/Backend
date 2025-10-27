package domainjwt

import "crypto/rsa"

type JWTKeyManager interface {
	LoadKeys()
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
}
