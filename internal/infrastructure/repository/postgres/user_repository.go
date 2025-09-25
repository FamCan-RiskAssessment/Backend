package postgres

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

func (repo *UserRepository) FindUserByID(db database.Database, id uint) (*entity.User, error) {
	var user entity.User
	result := db.GetDB().First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
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

func (userRepository *UserRepository) UpdateUser(db database.Database, user *entity.User) error {
	return db.GetDB().Save(&user).Error
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

func (repo *UserRepository) FindAllPermissions(db database.Database) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	result := db.GetDB().Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (repo *UserRepository) FindAllRoles(db database.Database) ([]*entity.Role, error) {
	var roles []*entity.Role
	result := db.GetDB().Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (repo *UserRepository) FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error) {
	var permission entity.Permission
	result := db.GetDB().First(&permission, permissionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &permission, nil
}

func (repo *UserRepository) FindRoleByID(db database.Database, roleID uint) (*entity.Role, error) {
	var role entity.Role
	result := db.GetDB().First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (repo *UserRepository) FindUsersByRoleID(db database.Database, roleID uint) ([]*entity.User, error) {
	var users []*entity.User
	result := db.GetDB().
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *UserRepository) FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	var users []*entity.User

	result := db.GetDB().
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("permissions.type IN ?", permissionTypes).
		Distinct().
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo *UserRepository) DeleteRole(db database.Database, roleID uint) error {
	return db.GetDB().Unscoped().Delete(&entity.Role{}, roleID).Error
}

func (repo *UserRepository) UpdateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Save(&role).Error
}

func (repo *UserRepository) ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error {
	return db.GetDB().Model(&role).Association("Permissions").Replace(permissions)

}

func (repo *UserRepository) ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error {
	return db.GetDB().Model(&user).Association("Roles").Replace(roles)

}

func (repo *UserRepository) FindRolesByPermission(db database.Database, permissionID uint) ([]*entity.Role, error) {
	var roles []*entity.Role
	query := db.GetDB().
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Where("role_permissions.permission_id = ?", permissionID)

	result := query.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
