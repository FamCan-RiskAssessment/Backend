package usecase

import "github.com/stretchr/testify/mock"

type PasswordHasherMock struct {
	mock.Mock
}

func NewPasswordHasherMock() *PasswordHasherMock {
	return &PasswordHasherMock{}
}

func (m *PasswordHasherMock) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *PasswordHasherMock) VerifyPassword(password, hash string) error {
	args := m.Called(password, hash)
	return args.Error(0)
}

func (m *PasswordHasherMock) NeedsRehash(hash string) bool {
	args := m.Called(hash)
	return args.Bool(0)
}
