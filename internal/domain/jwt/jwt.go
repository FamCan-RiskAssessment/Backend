package jwt

import (
	"crypto/rsa"
)

type KeyManager interface {
	LoadKeys(privateKeyPEM, publicKeyPEM string) error
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
}
