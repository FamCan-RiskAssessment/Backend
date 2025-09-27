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
		permissions := accessManagement.Group("/permissions")
		{
			permissions.GET("", app.Controllers.Admin.UserController.GetPermissionsList)
			permissions.GET("/:permissionID/roles", app.Controllers.Admin.UserController.GetPermissionRoles)
		}

		roles := accessManagement.Group("/roles")
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

		userRoles := accessManagement.Group("/users/:userID/roles")
		{
			userRoles.GET("", app.Controllers.Admin.UserController.GetUserRoles)
			userRoles.PUT("", app.Controllers.Admin.UserController.UpdateUserRoles)
		}

		password := accessManagement.Group("/users/password")
		{
			password.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionSetPassword}))
			password.PUT("", app.Controllers.Admin.UserController.SetPassword)
		}
	}

	userManagement := routerGroup.Group("/users")
	{
		userManagement.GET("", app.Controllers.Admin.UserController.GetUsers)
	}
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login-with-password", app.Controllers.Admin.UserController.LoginWithPassword)
	}

	forms := routerGroup.Group("/forms")
	// forms.Use(app.Middlewares.Auth.AuthRequired)
	{
		forms.GET("", app.Controllers.Admin.FormController.GetAllForms)
		forms.GET("/:formID", app.Controllers.Admin.FormController.GetForm)
		forms.DELETE("/:formID", app.Controllers.Admin.FormController.DeleteForm)
	}
}
