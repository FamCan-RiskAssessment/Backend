package user

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	constants   *bootstrap.Constants
	userService usecase.UserService
	pagination  *bootstrap.Pagination
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	userService usecase.UserService,
	pagination *bootstrap.Pagination,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		userService: userService,
		pagination:  pagination,
	}
}

func (userController *AdminUserController) GetUsers(ctx *gin.Context) {
	type usersParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := controller.Validate[usersParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetUsersListRequest{
		Offset: offset,
		Limit:  limit,
	}

	users, count, err := userController.userService.GetUsers(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(users, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetPermissionsList(ctx *gin.Context) {
	permissions, err := userController.userService.GetAllPermissions()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", permissions)
}

func (userController *AdminUserController) GetPermissionRoles(ctx *gin.Context) {
	type getPermissionRolesParams struct {
		PermissionID uint `uri:"permissionID" validate:"required"`
	}
	params := controller.Validate[getPermissionRolesParams](ctx)

	request := userdto.GetPermissionRolesRequest{
		PermissionID: params.PermissionID,
	}

	roles, err := userController.userService.GetPermissionRoles(request)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roles)
}

func (userController *AdminUserController) GetRolesList(ctx *gin.Context) {
	roles, err := userController.userService.GetAllRoles()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roles)
}

func (userController *AdminUserController) CreateRole(ctx *gin.Context) {
	type newRoleParams struct {
		Name          string `json:"name" validate:"required"`
		PermissionIDs []uint `json:"permissionIDs"`
	}
	params := controller.Validate[newRoleParams](ctx)

	newRoleRequest := userdto.NewRoleRequest{
		Name:          params.Name,
		PermissionIDs: params.PermissionIDs,
	}
	if err := userController.userService.CreateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetRoleDetails(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validate[getRoleParams](ctx)

	role, err := userController.userService.GetRoleDetails(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", role)
}

func (userController *AdminUserController) GetRoleOwners(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validate[getRoleParams](ctx)
	roleOwners, err := userController.userService.GetRoleOwners(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roleOwners)
}

func (userController *AdminUserController) UpdateRole(ctx *gin.Context) {
	type updateRoleParams struct {
		RoleID        uint    `uri:"roleID" validate:"required"`
		Name          *string `json:"name"`
		PermissionIDs []uint  `json:"permissionIDs"`
	}
	params := controller.Validate[updateRoleParams](ctx)

	newRoleRequest := userdto.UpdateRoleRequest{
		RoleID:        params.RoleID,
		Name:          params.Name,
		PermissionIDs: params.PermissionIDs,
	}
	if err := userController.userService.UpdateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) DeleteRole(ctx *gin.Context) {
	type deleteRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validate[deleteRoleParams](ctx)

	if err := userController.userService.DeleteRole(params.RoleID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetUserRoles(ctx *gin.Context) {
	type getRolesParams struct {
		UserID uint `uri:"userID" validate:"required"`
	}
	params := controller.Validate[getRolesParams](ctx)
	roles, err := userController.userService.GetUserRoles(params.UserID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roles)
}

func (userController *AdminUserController) UpdateUserRoles(ctx *gin.Context) {
	type updateUserRolesParams struct {
		UserID  uint   `uri:"userID" validate:"required"`
		RoleIDs []uint `json:"roleIDs"`
	}
	params := controller.Validate[updateUserRolesParams](ctx)

	userRolesRequest := userdto.UpdateUserRolesRequest{
		UserID:  params.UserID,
		RoleIDs: params.RoleIDs,
	}
	if err := userController.userService.UpdateUserRoles(userRolesRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateUserRoles")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) SetPassword(ctx *gin.Context) {
	type setPasswordParams struct {
		Password string `json:"password"`
	}
	params := controller.Validate[setPasswordParams](ctx)

	userID, _ := ctx.Get(userController.constants.Context.ID)
	userPasswordRequest := userdto.SetPasswordRequest{
		UserID:   userID.(uint),
		Password: params.Password,
	}
	if err := userController.userService.SetPassword(userPasswordRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setPassword")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) LoginWithPassword(ctx *gin.Context) {
	type loginWithPasswordParams struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	params := controller.Validate[loginWithPasswordParams](ctx)

	loginInfo := userdto.LoginRequest{
		Phone:    params.Phone,
		Password: params.Password,
	}

	response, err := userController.userService.LoginWithPassword(loginInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.login")
	controller.Response(ctx, 200, message, response)

}
