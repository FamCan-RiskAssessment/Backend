package redis

import (
	"context"
	"encoding/json"
	"time"

	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	rdb database.Cache
}

func NewUserCacheRepository(rdb database.Cache) *UserCacheRepository {
	return &UserCacheRepository{
		rdb: rdb,
	}
}

func (userCache *UserCacheRepository) Get(ctx context.Context, key string) (*userdto.OTPData, error) {
	value, err := userCache.rdb.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var otpData userdto.OTPData
	if err = json.Unmarshal([]byte(value), &otpData); err != nil {
		return nil, err
	}

	return &otpData, nil

}

func (userCache *UserCacheRepository) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	otpData := userdto.OTPData{
		OTP:      otp,
		Attempts: 0,
	}
	value, err := json.Marshal(otpData)
	if err != nil {
		return err
	}
	err = userCache.rdb.GetRDB().Set(ctx, key, string(value), expiration).Err()
	if err != nil {

		return err
	}
	return nil
}

func (userCache *UserCacheRepository) SetValidationToken(ctx context.Context, key string, expiration time.Duration) error {
	err := userCache.rdb.GetRDB().Set(ctx, key, "validated", expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (userCache *UserCacheRepository) GetValidationToken(ctx context.Context, key string) (bool, error) {
	_, err := userCache.rdb.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (userCache *UserCacheRepository) Delete(ctx context.Context, key string) error {
	return userCache.rdb.GetRDB().Del(ctx, key).Err()
}
