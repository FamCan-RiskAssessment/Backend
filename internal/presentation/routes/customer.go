package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	form := routerGroup.Group("/form")
	form.Use(app.Middlewares.Auth.AuthRequired)
	{
		form.POST("", app.Controllers.Customer.FormController.CreateForm)
		form.GET("", app.Controllers.Customer.FormController.GetUserForms)
		form.GET("/:formID", app.Controllers.Customer.FormController.GetForm)
		form.PUT("", app.Controllers.Customer.FormController.UpdateForm)
		form.DELETE("/:formID", app.Controllers.Customer.FormController.DeleteForm)
	}
}
