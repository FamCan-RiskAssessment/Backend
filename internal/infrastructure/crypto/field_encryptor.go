package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
)

type FieldEncryptor struct {
	key []byte
}

// NewFieldEncryptor creates encryptor with AES-256 key
// Key must be 32 bytes for AES-256
func NewFieldEncryptor(security *bootstrap.Security) (*FieldEncryptor, error) {
	keyBytes := []byte(security.EncryptionKey)
	if len(keyBytes) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes for AES-256")
	}
	return &FieldEncryptor{key: keyBytes}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
// Returns base64-encoded ciphertext
func (fe *FieldEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
// If the input is not valid base64, returns it as-is (backwards compatibility for plaintext data)
func (fe *FieldEncryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// Not valid base64 - assume it's plaintext from before encryption was enabled
		return ciphertext, nil
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		// Data too short to be encrypted - assume plaintext
		return ciphertext, nil
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		// Decryption failed - assume plaintext
		return ciphertext, nil
	}

	return string(plaintext), nil
}
