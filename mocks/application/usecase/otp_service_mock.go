package usecase

import (
	"github.com/stretchr/testify/mock"
)

type OtpServiceMock struct {
	mock.Mock
}

func NewOtpServiceMock() *OtpServiceMock {
	return &OtpServiceMock{}
}

func (o *OtpServiceMock) GenerateOTP(phone string) (string, int, error) {
	args := o.Called(phone)
	return args.String(0), args.Int(1), args.Error(2)
}

func (o *OtpServiceMock) VerifyOTP(redisKey, otp string) error {
	args := o.Called(redisKey, otp)
	return args.Error(0)
}
