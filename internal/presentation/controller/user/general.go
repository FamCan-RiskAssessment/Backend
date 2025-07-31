package user

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
)

type GeneralUserController struct {
	constants   *bootstrap.Constants
	userService usecase.UserService
}

func NewGeneralUserController(
	constants *bootstrap.Constants,
	userService usecase.UserService,
) *GeneralUserController {
	return &GeneralUserController{
		constants:   constants,
		userService: userService,
	}
}
