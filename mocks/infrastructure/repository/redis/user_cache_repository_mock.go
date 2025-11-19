package redis

import (
	"context"
	"time"

	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/stretchr/testify/mock"
)

type UserCacheRepositoryMock struct {
	mock.Mock
}

func NewUserCacheRepositoryMock() *UserCacheRepositoryMock {
	return &UserCacheRepositoryMock{}
}

func (u *UserCacheRepositoryMock) Get(ctx context.Context, key string) (*userdto.OTPData, error) {
	args := u.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdto.OTPData), args.Error(1)
}

func (u *UserCacheRepositoryMock) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	args := u.Called(ctx, key, otp, expiration)
	return args.Error(0)
}

func (u *UserCacheRepositoryMock) SetValidationToken(ctx context.Context, key string, expiration time.Duration) error {
	args := u.Called(ctx, key, expiration)
	return args.Error(0)
}

func (u *UserCacheRepositoryMock) GetValidationToken(ctx context.Context, key string) (bool, error) {
	args := u.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (u *UserCacheRepositoryMock) Delete(ctx context.Context, key string) error {
	args := u.Called(ctx, key)
	return args.Error(0)
}
