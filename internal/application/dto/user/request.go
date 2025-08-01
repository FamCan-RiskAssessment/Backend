package userdto

type LoginRequest struct {
	Phone string
}

type VerifyOTPRequest struct {
	Phone string
	OTP   string
}
