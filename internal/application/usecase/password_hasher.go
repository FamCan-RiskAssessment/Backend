package usecase

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
	NeedsRehash(hash string) bool
}
