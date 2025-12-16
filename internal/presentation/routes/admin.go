package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	accessManagement := routerGroup.Group("")
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
	}

	supervisorFormManagement := routerGroup.Group("/form")
	supervisorFormManagement.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.PermissionHandleOperators)}))
	{

	}

	operators := routerGroup.Group("/operator")
	operators.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionCreateFormForUser}))
	{
		operators.POST("/validate-user/request", app.Controllers.Admin.FormController.RequestUserValidationOTP)
		operators.POST("/validate-user/verify", app.Controllers.Admin.FormController.VerifyUserValidationOTP)
		operators.POST("/form", app.Controllers.Admin.FormController.CreateFormForUser)

		operatorFormGroup := operators.Group("/form/:formID")
		{
			operatorFormGroup.PUT("/accept", app.Controllers.Admin.FormController.AcceptForm)
			operatorFormGroup.PUT("/reject", app.Controllers.Admin.FormController.RejectForm)
			operatorFormGroup.PATCH("/basic", app.Controllers.Admin.FormController.UpdateBasicInfo)
			operatorFormGroup.PATCH("/generalhealth", app.Controllers.Admin.FormController.UpdateGeneralHealth)
			operatorFormGroup.PATCH("/mamography", app.Controllers.Admin.FormController.UpdateMamography)
			// operatorFormGroup.PATCH("/familycancer", app.Controllers.Admin.FormController.UpdateFamilyCancer)
			operatorFormGroup.PATCH("/contact", app.Controllers.Admin.FormController.UpdateContact)
			operatorFormGroup.PATCH("/lungcancer", app.Controllers.Admin.FormController.UpdateLungCancer)

			// navid form
			operatorFormGroup.PATCH("/navid", app.Controllers.Admin.FormController.UpdateNavidForm)

			operatorFormGroup.GET("/basic", app.Controllers.Admin.FormController.GetBasicForm)
			operatorFormGroup.GET("/generalhealth", app.Controllers.Admin.FormController.GetGeneralHealth)
			operatorFormGroup.GET("/mamography", app.Controllers.Admin.FormController.GetMamography)
			operatorFormGroup.GET("/familycancer", app.Controllers.Admin.FormController.GetFamilyCancer)
			operatorFormGroup.GET("/contact", app.Controllers.Admin.FormController.GetContact)
			operatorFormGroup.GET("/lungcancer", app.Controllers.Admin.FormController.GetLungCancer)

			// navid form
			operatorFormGroup.GET("/navid", app.Controllers.Admin.FormController.GetNavidForm)

			// Single cancer operations (more specific routes first)
			operatorFormGroup.PATCH("/cancer/:cancerID", app.Controllers.Admin.FormController.UpdateCancer)
			operatorFormGroup.DELETE("/cancer/:cancerID", app.Controllers.Admin.FormController.DeleteCancer)
			operatorFormGroup.POST("/cancer", app.Controllers.Admin.FormController.CreateCancer)
			operatorFormGroup.GET("/cancer", app.Controllers.Admin.FormController.GetAllCancers)

			// Single family cancer operations (more specific routes first)
			operatorFormGroup.PATCH("/familycancer/:familyCancerID", app.Controllers.Admin.FormController.UpdateFamilyCancer)
			operatorFormGroup.DELETE("/familycancer/:familyCancerID", app.Controllers.Admin.FormController.DeleteFamilyCancer)
			operatorFormGroup.POST("/familycancer", app.Controllers.Admin.FormController.CreateFamilyCancer)
		}
	}

	supervisor := routerGroup.Group("/supervisor/operator/")
	supervisor.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.PermissionHandleOperators)}))
	{
		// supervisor.GET("", app.Controllers.Admin.FormController.GetAllForms)
	}

	forms := routerGroup.Group("/form")
	forms.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.CategoryFormManagement)}))
	{
		forms.GET("", app.Controllers.Admin.FormController.GetAllForms)
		forms.DELETE("/:formID", app.Controllers.Admin.FormController.DeleteForm)

		supervisorFormManagement := forms.Group("/:formID")
		supervisorFormManagement.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.PermissionHandleOperators)}))
		{
			supervisorFormManagement.PUT("/operator", app.Controllers.Admin.FormController.AssignOperator)
			supervisorFormManagement.DELETE("/operator", app.Controllers.Admin.FormController.UnassignOperator)
			supervisorFormManagement.PUT("/accept", app.Controllers.Admin.FormController.AcceptForm)
			supervisorFormManagement.PUT("/reject", app.Controllers.Admin.FormController.RejectForm)
		}

		formManagement := forms.Group("/:formID")
		{
			formManagement.PATCH("/basic", app.Controllers.Admin.FormController.UpdateBasicInfo)
			formManagement.PATCH("/generalhealth", app.Controllers.Admin.FormController.UpdateGeneralHealth)
			formManagement.PATCH("/mamography", app.Controllers.Admin.FormController.UpdateMamography)
			// formManagement.PATCH("/familycancer", app.Controllers.Admin.FormController.UpdateFamilyCancer)
			formManagement.PATCH("/contact", app.Controllers.Admin.FormController.UpdateContact)
			formManagement.PATCH("/lungcancer", app.Controllers.Admin.FormController.UpdateLungCancer)

			// navid form
			formManagement.PATCH("/navid", app.Controllers.Admin.FormController.UpdateNavidForm)

			formManagement.GET("/basic", app.Controllers.Admin.FormController.GetBasicForm)
			formManagement.GET("/generalhealth", app.Controllers.Admin.FormController.GetGeneralHealth)
			formManagement.GET("/mamography", app.Controllers.Admin.FormController.GetMamography)
			formManagement.GET("/contact", app.Controllers.Admin.FormController.GetContact)
			formManagement.GET("/lungcancer", app.Controllers.Admin.FormController.GetLungCancer)

			// navid form
			formManagement.GET("/navid", app.Controllers.Admin.FormController.GetNavidForm)

			// Single cancer operations (more specific routes first)
			formManagement.PATCH("/cancer/:cancerID", app.Controllers.Admin.FormController.UpdateCancer)
			formManagement.DELETE("/cancer/:cancerID", app.Controllers.Admin.FormController.DeleteCancer)
			formManagement.POST("/cancer", app.Controllers.Admin.FormController.CreateCancer)
			formManagement.GET("/cancer", app.Controllers.Admin.FormController.GetAllCancers)

			// Single family cancer operations (more specific routes first)
			formManagement.GET("/familycancer", app.Controllers.Admin.FormController.GetFamilyCancer)
			formManagement.POST("/familycancer", app.Controllers.Admin.FormController.CreateFamilyCancer)
			formManagement.PATCH("/familycancer/:familyCancerID", app.Controllers.Admin.FormController.UpdateFamilyCancer)
			formManagement.DELETE("/familycancer/:familyCancerID", app.Controllers.Admin.FormController.DeleteFamilyCancer)
		}
	}

	operatorForms := routerGroup.Group("/operator-form")
	operatorForms.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.CategoryFormManagement)}))
	{
		operatorForms.GET("", app.Controllers.Admin.FormController.GetAllOperatorForms)
	}

	actionLog := routerGroup.Group("/log")
	actionLog.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionType(enum.CategoryLogManagement)}))
	{
		actionLog.GET("", app.Controllers.Admin.ActionLogController.GetAllActionLogs)
		actionLog.GET("/types", app.Controllers.Admin.ActionLogController.GetAllActionTypes)
	}

	calc := routerGroup.Group("/calc")
	{
		calc.POST("/model", app.Controllers.Admin.CalcController.SendFormToCalc)
		calc.GET("/premm5/:formID", app.Controllers.Admin.CalcController.GetPremm5Results)
		calc.GET("/bcra/:formID", app.Controllers.Admin.CalcController.GetBCRAResults)
		calc.GET("/gail/:formID", app.Controllers.Admin.CalcController.GetGailResults)
		calc.GET("/plco/:formID", app.Controllers.Admin.CalcController.GetPLCOResults)
		calc.GET("/all-models", app.Controllers.Admin.CalcController.GetAllModelTypes)
	}
}
