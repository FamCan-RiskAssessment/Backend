package userdto

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}
