package repository

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (userRepository *UserRepository) FindUserByPhone(db database.Database, phone string) (*entity.User, error) {
	var user entity.User
	result := db.GetDB().Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (userRepository *UserRepository) CreateUser(db database.Database, user *entity.User) error {
	return db.GetDB().Create(user).Error
}

func (userRepository *UserRepository) FindUserRoles(db database.Database, user *entity.User) error {
	return db.GetDB().Preload("Roles").First(&user).Error
}

func (userRepository *UserRepository) FindRolePermissions(db database.Database, role *entity.Role) error {
	return db.GetDB().Preload("Permissions").First(&role).Error
}

func (userRepository *UserRepository) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	var role entity.Role
	result := db.GetDB().Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (userRepository *UserRepository) CreateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Create(role).Error
}

func (userRepository *UserRepository) RoleHasPermission(db database.Database, roleID, permissionID uint) bool {
	var count int64
	db.GetDB().
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count)
	return count > 0
}

func (userRepository *UserRepository) AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error {
	return db.GetDB().Model(role).Association("Permissions").Append(permission)
}

func (userRepository *UserRepository) AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error {
	return db.GetDB().Model(user).Association("Roles").Append(role)
}

func (userRepository *UserRepository) UserHasRole(db database.Database, userID, roleID uint) bool {
	var count int64
	db.GetDB().
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count)
	return count > 0
}

func (userRepository *UserRepository) FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, error) {
	var permission entity.Permission
	result := db.GetDB().Where("type = ?", permissionType).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &permission, nil
}

func (userRepository *UserRepository) CreatePermission(db database.Database, permission *entity.Permission) error {
	return db.GetDB().Create(permission).Error
}
