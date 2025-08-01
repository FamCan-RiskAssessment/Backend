package usecase

type OtpService interface {
	GenerateOTP(phone string) (string, int, error)
	VerifyOTP(phone string, otp string) error
}
