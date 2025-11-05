package redis

import (
	"context"
	"time"

	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
)

type UserCacheRepository interface {
	Get(ctx context.Context, key string) (*userdto.OTPData, error)
	Set(ctx context.Context, key, otp string, expiration time.Duration) error
	SetValidationToken(ctx context.Context, key string, expiration time.Duration) error
	GetValidationToken(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
