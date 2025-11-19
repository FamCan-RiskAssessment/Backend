package jwt

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/stretchr/testify/mock"
)

type KeyManagerMock struct {
	mock.Mock
}

func NewKeyManagerMock() *KeyManagerMock {
	return &KeyManagerMock{}
}

func (k *KeyManagerMock) LoadKeys(privateKeyPath, publicKeyPath string) error {
	args := k.Called(privateKeyPath, publicKeyPath)
	return args.Error(0)
}

func (k *KeyManagerMock) GetPrivateKey() *rsa.PrivateKey {
	args := k.Called()
	if args.Get(0) == nil {
		// Generate a test RSA key pair
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return privateKey
	}
	return args.Get(0).(*rsa.PrivateKey)
}

func (k *KeyManagerMock) GetPublicKey() *rsa.PublicKey {
	args := k.Called()
	if args.Get(0) == nil {
		// Generate a test RSA key pair
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return &privateKey.PublicKey
	}
	return args.Get(0).(*rsa.PublicKey)
}
