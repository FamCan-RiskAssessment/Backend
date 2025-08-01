package service

import (
	"context"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/communication"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	redis "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/redis"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type UserService struct {
	constants           *bootstrap.Constants
	userRepository      postgres.UserRepository
	userCacheRepository redis.UserCacheRepository
	jwtService          usecase.JwtService
	smsService          communication.SmsService
	otpService          usecase.OtpService
	db                  database.Database
}

func NewUserService(
	constants *bootstrap.Constants,
	userRepository postgres.UserRepository,
	userCacheRepository redis.UserCacheRepository,
	jwtService usecase.JwtService,
	smsService communication.SmsService,
	otpService usecase.OtpService,
	db database.Database,
) *UserService {
	return &UserService{
		constants:           constants,
		userRepository:      userRepository,
		userCacheRepository: userCacheRepository,
		jwtService:          jwtService,
		smsService:          smsService,
		otpService:          otpService,
		db:                  db,
	}
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) error {
	user, err := userService.userRepository.GetUserByPhone(userService.db, loginInfo.Phone)
	if err != nil {
		return err
	}

	if user == nil {
		user = &entity.User{
			Phone: loginInfo.Phone,
		}
		err = userService.userRepository.CreateUser(userService.db, user)
		if err != nil {
			return err
		}
	}

	otp, exppireMinute, err := userService.otpService.GenerateOTP(loginInfo.Phone)
	if err != nil {
		return err
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(loginInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(exppireMinute)*time.Minute)
	if err != nil {
		return err
	}

	// err = userService.smsService.SendOTP(loginInfo.Phone, otp)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (userService *UserService) VerifyOTP(verifyOTPInfo userdto.VerifyOTPRequest) (userdto.LoginResponse, error) {
	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyOTPInfo.Phone)
	err := userService.otpService.VerifyOTP(redisKey, verifyOTPInfo.OTP)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	user, err := userService.userRepository.GetUserByPhone(userService.db, verifyOTPInfo.Phone)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return userdto.LoginResponse{}, notFoundError
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	return userdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
