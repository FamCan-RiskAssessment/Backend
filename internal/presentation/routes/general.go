package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login", app.Controllers.General.UserController.Login)
		auth.POST("/verify-otp", app.Controllers.General.UserController.VerifyOTP)
		auth.POST("/login-with-password", app.Controllers.Admin.UserController.LoginWithPassword)
	}
}
