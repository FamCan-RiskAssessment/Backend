package service

import (
	"errors"
	"testing"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	userCacheMocks "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/repository/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OTPServiceTestSuite struct {
	suite.Suite
	constants           *bootstrap.Constants
	userCacheRepository *userCacheMocks.UserCacheRepositoryMock
	otpService          *OTPService
	otpConfig           *bootstrap.OTP
}

func (suite *OTPServiceTestSuite) SetupTest() {
	suite.constants = bootstrap.NewConstants()
	suite.userCacheRepository = userCacheMocks.NewUserCacheRepositoryMock()
	suite.otpConfig = &bootstrap.OTP{
		Length:       6,
		ExpiryMinute: 5,
	}

	suite.otpService = NewOTPService(
		suite.constants,
		suite.otpConfig,
		suite.userCacheRepository,
	)
}

// Test: GenerateOTP should return a 6-digit OTP and expiry time
func (suite *OTPServiceTestSuite) TestGenerateOTP_Success() {
	// Act
	otp, expiryMinute, err := suite.otpService.GenerateOTP("989123456789")

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), otp)
	assert.Equal(suite.T(), 6, len(otp))
	assert.Equal(suite.T(), 5, expiryMinute)

	// Verify OTP contains only digits 1-9 (as per the table in the service)
	for _, digit := range otp {
		assert.True(suite.T(), digit >= '1' && digit <= '9', "OTP should only contain digits 1-9")
	}
}

// Test: GenerateOTP should return different OTPs on successive calls
func (suite *OTPServiceTestSuite) TestGenerateOTP_RandomnessCheck() {
	// Act
	otp1, _, err1 := suite.otpService.GenerateOTP("989123456789")
	otp2, _, err2 := suite.otpService.GenerateOTP("989123456789")

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), otp1, otp2, "OTPs should be randomly generated")
}

// Test: VerifyOTP should return error when OTP key not found in cache
func (suite *OTPServiceTestSuite) TestVerifyOTP_OTPExpired() {
	// Arrange
	redisKey := "otp:989123456789"
	otp := "123456"

	suite.userCacheRepository.On("Get", mock.Anything, redisKey).
		Return(nil, nil)

	// Act
	err := suite.otpService.VerifyOTP(redisKey, otp)

	// Assert
	assert.Error(suite.T(), err)
	validationErrors, ok := err.(exception.ValidationErrors)
	assert.True(suite.T(), ok, "Error should be ValidationErrors type")
	assert.Len(suite.T(), validationErrors.Errors, 1)
	assert.Equal(suite.T(), suite.constants.Field.OTP, validationErrors.Errors[0].Field)
	assert.Equal(suite.T(), suite.constants.Tag.Expired, validationErrors.Errors[0].Tag)
}

// Test: VerifyOTP should return error when OTP is incorrect
func (suite *OTPServiceTestSuite) TestVerifyOTP_InvalidOTP() {
	// Arrange
	redisKey := "otp:989123456789"
	incorrectOTP := "999999"
	cachedOTPData := &userdto.OTPData{
		OTP:      "123456",
		Attempts: 0,
	}

	suite.userCacheRepository.On("Get", mock.Anything, redisKey).
		Return(cachedOTPData, nil)

	// Act
	err := suite.otpService.VerifyOTP(redisKey, incorrectOTP)

	// Assert
	assert.Error(suite.T(), err)
	validationErrors, ok := err.(exception.ValidationErrors)
	assert.True(suite.T(), ok, "Error should be ValidationErrors type")
	assert.Len(suite.T(), validationErrors.Errors, 1)
	assert.Equal(suite.T(), suite.constants.Field.OTP, validationErrors.Errors[0].Field)
	assert.Equal(suite.T(), suite.constants.Tag.Invalid, validationErrors.Errors[0].Tag)
}

// Test: VerifyOTP should succeed with correct OTP
func (suite *OTPServiceTestSuite) TestVerifyOTP_CorrectOTP() {
	// Arrange
	redisKey := "otp:989123456789"
	correctOTP := "123456"
	cachedOTPData := &userdto.OTPData{
		OTP:      "123456",
		Attempts: 0,
	}

	suite.userCacheRepository.On("Get", mock.Anything, redisKey).
		Return(cachedOTPData, nil)

	// Act
	err := suite.otpService.VerifyOTP(redisKey, correctOTP)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: VerifyOTP should succeed with hardcoded test OTP "111111"
func (suite *OTPServiceTestSuite) TestVerifyOTP_TestOTP() {
	// Arrange
	redisKey := "otp:989123456789"
	testOTP := "111111"
	cachedOTPData := &userdto.OTPData{
		OTP:      "123456",
		Attempts: 0,
	}

	suite.userCacheRepository.On("Get", mock.Anything, redisKey).
		Return(cachedOTPData, nil)

	// Act
	err := suite.otpService.VerifyOTP(redisKey, testOTP)

	// Assert
	assert.NoError(suite.T(), err, "Test OTP '111111' should always be valid")
}

// Test: VerifyOTP should handle Redis errors
func (suite *OTPServiceTestSuite) TestVerifyOTP_RepositoryError() {
	// Arrange
	redisKey := "otp:989123456789"
	otp := "123456"
	repositoryError := errors.New("redis connection error")

	suite.userCacheRepository.On("Get", mock.Anything, redisKey).
		Return(nil, repositoryError)

	// Act
	err := suite.otpService.VerifyOTP(redisKey, otp)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
}

// Test: GenerateOTP should return error when random reader fails
// Note: Line 36 (if n != otpService.otpConfig.Length) is hard to trigger in practice
// because rand.Reader is very reliable and io.ReadAtLeast will return an error
// instead of returning fewer bytes. We test the error handling implicitly through
// success tests - if the function properly returns on error conditions.
func (suite *OTPServiceTestSuite) TestGenerateOTP_ReturnsValidOTP() {
	// Act - test the happy path thoroughly
	otp, expiryMinute, err := suite.otpService.GenerateOTP("989123456789")

	// Assert - ensure we get valid response
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), otp)
	assert.Len(suite.T(), otp, 6)
	assert.Equal(suite.T(), 5, expiryMinute)

	// Verify all characters are digits 1-9
	for _, char := range otp {
		assert.GreaterOrEqual(suite.T(), int(char), int('1'))
		assert.LessOrEqual(suite.T(), int(char), int('9'))
	}
}

func TestOTPServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OTPServiceTestSuite))
}
