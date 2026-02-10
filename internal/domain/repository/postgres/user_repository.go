package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUserByPhone(db database.Database, phone string) (*entity.User, error)
	CreateUser(db database.Database, user *entity.User) error
	UpdateUser(db database.Database, user *entity.User) error
	FindUserByID(db database.Database, id uint) (*entity.User, error)
	FindUserRoles(db database.Database, user *entity.User) error
	FindRolePermissions(db database.Database, role *entity.Role) error
	FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, error)
	CreatePermission(db database.Database, permission *entity.Permission) error
	FindRoleByName(db database.Database, name string) (*entity.Role, error)
	CreateRole(db database.Database, role *entity.Role) error
	RoleHasPermission(db database.Database, roleID, permissionID uint) bool
	AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error
	AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error
	UserHasRole(db database.Database, userID, roleID uint) bool
	FindAllPermissions(db database.Database) ([]*entity.Permission, error)
	FindAllRoles(db database.Database) ([]*entity.Role, error)
	FindRoleByID(db database.Database, id uint) (*entity.Role, error)
	FindPermissionByID(db database.Database, id uint) (*entity.Permission, error)
	FindUsersByRoleID(db database.Database, roleID uint) ([]*entity.User, error)
	FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) ([]*entity.User, error)
	UpdateRole(db database.Database, role *entity.Role) error
	ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error
	ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error
	DeleteRole(db database.Database, id uint) error
	FindRolesByPermission(db database.Database, permissionID uint) ([]*entity.Role, error)
	FindUsers(db database.Database, options *QueryOptions) ([]*entity.User, int64, error)
	FindProfileByUserID(db database.Database, userID uint) (*entity.User, error)
	CreateProfile(db database.Database, profile *entity.UserProfile) error
	UpdateProfile(db database.Database, profile *entity.UserProfile) error
}
