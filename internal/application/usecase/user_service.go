package usecase

import userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"

type UserService interface {
	Login(loginInfo userdto.LoginRequest) error
	VerifyOTP(verifyOTPInfo userdto.VerifyOTPRequest) (userdto.LoginResponse, error)
}
