package repository

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUserByPhone(db database.Database, phone string) (*entity.User, error)
	CreateUser(db database.Database, user *entity.User) error
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
}
