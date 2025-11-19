package communication

import (
	"github.com/stretchr/testify/mock"
)

type SmsServiceMock struct {
	mock.Mock
}

func NewSmsServiceMock() *SmsServiceMock {
	return &SmsServiceMock{}
}

func (s *SmsServiceMock) SendOTP(phone, otp string) error {
	args := s.Called(phone, otp)
	return args.Error(0)
}
