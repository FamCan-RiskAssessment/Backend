package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	accessManagement := routerGroup.Group("")
	// accessManagement.Use(app.Middlewares.Auth.AuthRequired)
	{
		permissions := accessManagement.Group("/permission")
		{
			permissions.GET("", app.Controllers.Admin.UserController.GetPermissionsList)
			permissions.GET("/:permissionID/roles", app.Controllers.Admin.UserController.GetPermissionRoles)
		}

		roles := accessManagement.Group("/role")
		{
			roles.GET("", app.Controllers.Admin.UserController.GetRolesList)
			roles.POST("", app.Controllers.Admin.UserController.CreateRole)

			rolesSubGroup := roles.Group("/:roleID")
			{
				rolesSubGroup.GET("", app.Controllers.Admin.UserController.GetRoleDetails)
				rolesSubGroup.GET("/owners", app.Controllers.Admin.UserController.GetRoleOwners)
				rolesSubGroup.PUT("", app.Controllers.Admin.UserController.UpdateRole)
				rolesSubGroup.DELETE("", app.Controllers.Admin.UserController.DeleteRole)
			}
		}

		userRoles := accessManagement.Group("/user/:userID/role")
		{
			userRoles.GET("", app.Controllers.Admin.UserController.GetUserRoles)
			userRoles.PUT("", app.Controllers.Admin.UserController.UpdateUserRoles)
		}

		password := accessManagement.Group("/user/password")
		{
			password.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionSetPassword}))
			password.PUT("", app.Controllers.Admin.UserController.SetPassword)
		}
	}

	userManagement := routerGroup.Group("/user")
  {
		userManagement.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.CategoryPatientManagement)}))
		userManagement.GET("", app.Controllers.Admin.UserController.GetUsers)
    userManagement.GET("/forms", app.Controllers.Customer.FormController.GetUserForms)
	}

	forms := routerGroup.Group("/form")
	// forms.Use(app.Middlewares.Auth.AuthRequired)
	{
		forms.GET("", app.Controllers.Admin.FormController.GetAllForms)
		forms.DELETE("/:formID", app.Controllers.Admin.FormController.DeleteForm)

		forms.GET("/basic", app.Controllers.Customer.FormController.GetBasicForm)
		forms.GET("/generalhealth", app.Controllers.Customer.FormController.GetGeneralHealth)
		forms.GET("/mamography", app.Controllers.Customer.FormController.GetMamography)
		forms.GET("/cancer", app.Controllers.Customer.FormController.GetCancer)
		forms.GET("/familycancer", app.Controllers.Customer.FormController.GetFamilyCancer)
		forms.GET("/contact", app.Controllers.Customer.FormController.GetContact)
		forms.GET("/lungcancer", app.Controllers.Customer.FormController.GetLungCancer)
	}

	formManagement := routerGroup.Group("/form")
	{
		formManagement.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.CategoryFormManagement)}))
		formManagement.PUT("/:formID/accept", app.Controllers.Admin.FormController.AcceptForm)
		formManagement.PUT("/:formID/reject", app.Controllers.Admin.FormController.RejectForm)
		formManagement.PUT("/basic", app.Controllers.Customer.FormController.UpdateBasicInfo)
		formManagement.PUT("/generalhealth", app.Controllers.Customer.FormController.UpsertGeneralHealth)
		formManagement.PUT("/mamography", app.Controllers.Customer.FormController.UpsertMamography)
		formManagement.PUT("/cancer", app.Controllers.Customer.FormController.UpsertCancer)
		formManagement.PUT("/familycancer", app.Controllers.Customer.FormController.UpsertFamilyCancer)
		formManagement.PUT("/contact", app.Controllers.Customer.FormController.UpsertContact)
		formManagement.PUT("/lungcancer", app.Controllers.Customer.FormController.UpsertLungCancer)
		formManagement.PUT("/status", app.Controllers.Customer.FormController.ChangeFormStatus)
	}
}
