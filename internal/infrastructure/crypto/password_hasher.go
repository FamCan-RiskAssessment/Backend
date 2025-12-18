package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is bcrypt cost 12 (good balance of security and performance)
	// Each increment doubles computation time
	DefaultCost = 12
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: DefaultCost,
	}
}

// HashPassword takes a plaintext password and returns bcrypt hash
func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword compares plaintext password with hash
func (ph *PasswordHasher) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// NeedsRehash checks if password was hashed with different cost
// Useful for upgrading hash strength over time
func (ph *PasswordHasher) NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return false
	}
	return cost != ph.cost
}
