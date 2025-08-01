package usecase

type OtpService interface {
	GenerateOTP(phone string) (string, int, error)
	VerifyOTP(redisKey, otp string) error
}
