package usecase

type JwtService interface {
	GenerateToken(userID string) (string, string, error)
}
