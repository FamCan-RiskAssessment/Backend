package communication

type SmsService interface {
	SendOTP(receptor, token string) error
}
