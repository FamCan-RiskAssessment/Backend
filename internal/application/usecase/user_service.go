package usecase

import userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"

type UserService interface {
	Login(loginInfo userdto.LoginRequest) error
	VerifyOTP(verifyOTPInfo userdto.VerifyOTPRequest) (userdto.LoginResponse, error)
	GetAllPermissions() ([]userdto.PermissionResponse, error)
	GetPermissionRoles(request userdto.GetPermissionRolesRequest) ([]userdto.RoleResponse, error)
	GetAllRoles() ([]userdto.RoleResponse, error)
	CreateRole(request userdto.NewRoleRequest) error
	GetRoleDetails(roleID uint) (userdto.RoleResponse, error)
	GetRoleOwners(roleID uint) ([]userdto.UserResponse, error)
	UpdateRole(request userdto.UpdateRoleRequest) error
	DeleteRole(roleID uint) error
	GetUserRoles(userID uint) ([]userdto.RoleResponse, error)
	UpdateUserRoles(request userdto.UpdateUserRolesRequest) error
}
