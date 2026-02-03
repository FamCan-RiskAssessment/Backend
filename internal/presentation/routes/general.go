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
		admin := auth.Group("/admin")
		{
			admin.POST("/login", app.Controllers.Admin.UserController.LoginWithPassword)
		}
	}

	enums := routerGroup.Group("/enum")
	{
		enums.GET("/cancer-types", app.Controllers.General.FormController.GetAllCancerTypes)
		enums.GET("/genders", app.Controllers.General.FormController.GetAllGenders)
		enums.GET("/menopausal-statuses", app.Controllers.General.FormController.GetAllMenopausalStatuses)
		enums.GET("/form-statuses", app.Controllers.General.FormController.GetAllFormStatuses)
		enums.GET("/hyperplasiaInBiopsy-statuses", app.Controllers.General.FormController.GetAllHyperplasiaInBiopsyStatuses)
		enums.GET("/life-statuses", app.Controllers.General.FormController.GetAllLifeStatuses)
		enums.GET("/relatives", app.Controllers.General.FormController.GetAllRelativeTypes)
		enums.GET("/form-types", app.Controllers.General.FormController.GetAllFormTypes)
		enums.GET("/answers", app.Controllers.General.FormController.GetAllAnswers)

	}
}
