package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	form := routerGroup.Group("/form")
	form.Use(app.Middlewares.Auth.AuthRequired)
	{
		form.POST("/basic", app.Controllers.Customer.FormController.CreateForm)
		formGroup := form.Group("/:formID")
		{
			formGroup.PUT("/basic", app.Controllers.Customer.FormController.UpdateForm)
			formGroup.PUT("/generalhealth", app.Controllers.Customer.FormController.UpsertGeneralHealth)
			formGroup.PUT("/mamography", app.Controllers.Customer.FormController.UpsertMamography)
			formGroup.PUT("/cancer", app.Controllers.Customer.FormController.UpsertCancer)
			formGroup.PUT("/familycancer", app.Controllers.Customer.FormController.UpsertFamilyCancer)
			formGroup.PUT("/contact", app.Controllers.Customer.FormController.UpsertContact)
			formGroup.PUT("/lungcancer", app.Controllers.Customer.FormController.UpsertLungCancer)

			formGroup.GET("", app.Controllers.Customer.FormController.GetForm)
			formGroup.DELETE("", app.Controllers.Customer.FormController.DeleteForm)
		}

		// To list all forms of the authenticated user
		form.GET("", app.Controllers.Customer.FormController.GetUserForms)
	}
}
