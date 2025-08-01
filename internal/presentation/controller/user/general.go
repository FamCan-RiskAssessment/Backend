package user

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
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

func (userController *GeneralUserController) Login(ctx *gin.Context) {
	type LoginParams struct {
		Phone string `json:"phone" validate:"required,min=11,max=11"`
	}

	params := controller.Validate[LoginParams](ctx)

	loginInfo := userdto.LoginRequest{
		Phone: params.Phone,
	}

	err := userController.userService.Login(loginInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "کد پیامکی ارسال شد", nil)
}

func (userController *GeneralUserController) VerifyOTP(ctx *gin.Context) {
	type VerifyCodeParams struct {
		Phone string `json:"phone" validate:"required,min=11,max=11"`
		OTP   string `json:"otp" validate:"required,min=6,max=6"`
	}

	params := controller.Validate[VerifyCodeParams](ctx)

	verifyOTPInfo := userdto.VerifyOTPRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}

	userInfo, err := userController.userService.VerifyOTP(verifyOTPInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "شما با موفقیت وارد شدید", userInfo)

}
