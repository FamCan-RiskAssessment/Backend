package jwt

import (
	"crypto/rsa"
	"fmt"
	"sync"

	domainJWT "github.com/FamCan-RiskAssessment/Backend/internal/domain/jwt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTKeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	mutex      sync.RWMutex
	isLoaded   bool
}

func NewJWTKeyManager() domainJWT.KeyManager {
	return &JWTKeyManager{}
}

func (k *JWTKeyManager) LoadKeys(privateKeyPEM, publicKeyPEM string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	k.privateKey = privateKey
	k.publicKey = publicKey
	k.isLoaded = true

	return nil
}

func (k *JWTKeyManager) GetPrivateKey() *rsa.PrivateKey {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.isLoaded {
		return nil
	}
	return k.privateKey
}

func (k *JWTKeyManager) GetPublicKey() *rsa.PublicKey {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.isLoaded {
		return nil
	}
	return k.publicKey
}
