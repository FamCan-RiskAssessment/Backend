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
			formGroup.PUT("/basic", app.Controllers.Customer.FormController.UpdateBasicInfo)
			formGroup.PUT("/generalhealth", app.Controllers.Customer.FormController.UpsertGeneralHealth)
			formGroup.PUT("/mamography", app.Controllers.Customer.FormController.UpsertMamography)
			// formGroup.PUT("/familycancer", app.Controllers.Customer.FormController.UpsertFamilyCancer)
			formGroup.PUT("/contact", app.Controllers.Customer.FormController.UpsertContact)
			formGroup.PUT("/lungcancer", app.Controllers.Customer.FormController.UpsertLungCancer)
			formGroup.PUT("/status", app.Controllers.Customer.FormController.ChangeFormStatus)

			formGroup.GET("/basic", app.Controllers.Customer.FormController.GetBasicForm)
			formGroup.GET("/generalhealth", app.Controllers.Customer.FormController.GetGeneralHealth)
			formGroup.GET("/mamography", app.Controllers.Customer.FormController.GetMamography)
			formGroup.GET("/familycancer", app.Controllers.Customer.FormController.GetFamilyCancer)
			formGroup.GET("/contact", app.Controllers.Customer.FormController.GetContact)
			formGroup.GET("/lungcancer", app.Controllers.Customer.FormController.GetLungCancer)

			// Single cancer operations (more specific routes first)
			formGroup.PUT("/cancer/:cancerID", app.Controllers.Customer.FormController.UpdateCancer)
			formGroup.DELETE("/cancer/:cancerID", app.Controllers.Customer.FormController.DeleteCancer)
			formGroup.POST("/cancer", app.Controllers.Customer.FormController.CreateCancer)
			formGroup.GET("/cancer", app.Controllers.Customer.FormController.GetAllCancers)

			// Single family cancer operations (more specific routes first)
			formGroup.POST("/familycancer", app.Controllers.Customer.FormController.CreateFamilyCancer)
			formGroup.PUT("/familycancer/:familyCancerID", app.Controllers.Customer.FormController.UpdateFamilyCancer)
			formGroup.DELETE("/familycancer/:familyCancerID", app.Controllers.Customer.FormController.DeleteFamilyCancer)

			formGroup.DELETE("", app.Controllers.Customer.FormController.DeleteForm)
		}

		form.GET("", app.Controllers.Customer.FormController.GetUserForms)
	}
}
