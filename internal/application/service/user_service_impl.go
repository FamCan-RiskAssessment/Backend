package service

import "github.com/FamCan-RiskAssessment/Backend/bootstrap"

type UserService struct {
	constants *bootstrap.Constants
}

func NewUserService(
	constants *bootstrap.Constants,
) *UserService {
	return &UserService{
		constants: constants,
	}
}
