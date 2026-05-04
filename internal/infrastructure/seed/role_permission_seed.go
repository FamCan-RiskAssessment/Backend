package seed

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	repository "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type RoleSeeder struct {
	superAdmin     *bootstrap.SuperAdmin
	userRepository repository.UserRepository
	db             database.Database
	passwordHasher usecase.PasswordHasher
}

func NewRoleSeeder(
	superAdmin *bootstrap.SuperAdmin,
	userRepository repository.UserRepository,
	db database.Database,
	passwordHasher usecase.PasswordHasher,
) *RoleSeeder {
	return &RoleSeeder{
		superAdmin:     superAdmin,
		userRepository: userRepository,
		db:             db,
		passwordHasher: passwordHasher,
	}
}

func (r *RoleSeeder) SeedRoles() {
	permissions := r.seedPermissions()
	roles := r.seedRolesWithPermissions(permissions)
	superAdminRoles := []*entity.Role{
		roles[enum.SuperAdmin],
	}
	r.seedSuperAdmin(r.superAdmin, superAdminRoles)
}

func (r *RoleSeeder) seedPermissions() map[enum.PermissionType]*entity.Permission {
	permissions := make(map[enum.PermissionType]*entity.Permission)
	permissionTypes := enum.GetAllPermissionTypes()

	for _, permissionType := range permissionTypes {
		permission := r.getOrCreatePermission(permissionType)
		permissions[permissionType] = permission
	}

	return permissions
}

func (roleSeeder *RoleSeeder) getOrCreatePermission(permissionType enum.PermissionType) *entity.Permission {
	permission, err := roleSeeder.userRepository.FindPermissionByType(roleSeeder.db, permissionType)
	if err != nil {
		panic(err)
	}
	if permission != nil {
		return permission
	}

	permission = &entity.Permission{
		Type:     permissionType,
		Category: permissionType.Category(),
	}

	if err := roleSeeder.userRepository.CreatePermission(roleSeeder.db, permission); err != nil {
		panic(err)
	}

	return permission
}

func (roleSeeder *RoleSeeder) seedRolesWithPermissions(permissions map[enum.PermissionType]*entity.Permission) map[enum.RoleName]*entity.Role {
	roles := make(map[enum.RoleName]*entity.Role)
	roleNames := enum.GetAllRoleNames()
	for _, roleName := range roleNames {
		role := roleSeeder.getOrCreateRole(roleName)
		roles[roleName] = role
		roleSeeder.assignPermissionsToRole(role, roleName, permissions)
	}
	return roles
}

func (roleSeeder *RoleSeeder) getOrCreateRole(roleName enum.RoleName) *entity.Role {
	role, err := roleSeeder.userRepository.FindRoleByName(roleSeeder.db, roleName.String())
	if err != nil {
		panic(err)
	}
	if role != nil {
		return role
	}

	role = &entity.Role{
		Name: roleName.String(),
	}

	if err := roleSeeder.userRepository.CreateRole(roleSeeder.db, role); err != nil {
		panic(err)
	}

	return role
}

func (roleSeeder *RoleSeeder) assignPermissionsToRole(role *entity.Role, roleName enum.RoleName, permissions map[enum.PermissionType]*entity.Permission) {
	for _, permissionType := range roleName.Permissions() {
		permission := permissions[permissionType]
		if !roleSeeder.userRepository.RoleHasPermission(roleSeeder.db, role.ID, permission.ID) {
			if err := roleSeeder.userRepository.AssignPermissionToRole(roleSeeder.db, role, permission); err != nil {
				panic(err)
			}
		}
	}
}

func (roleSeeder *RoleSeeder) seedSuperAdmin(adminCred *bootstrap.SuperAdmin, roles []*entity.Role) {
	admin := roleSeeder.getOrCreateAdmin(adminCred)
	for _, role := range roles {
		if !roleSeeder.userRepository.UserHasRole(roleSeeder.db, admin.ID, role.ID) {
			if err := roleSeeder.userRepository.AssignRoleToUser(roleSeeder.db, admin, role); err != nil {
				panic(err)
			}
		}

	}
}

func (roleSeeder *RoleSeeder) getOrCreateAdmin(adminCred *bootstrap.SuperAdmin) *entity.User {
	admin, err := roleSeeder.userRepository.FindUserByPhone(roleSeeder.db, adminCred.Phone)
	if err != nil {
		panic(err)
	}
	if admin == nil {
		hashedPassword, err := roleSeeder.passwordHasher.HashPassword(adminCred.Password)
		if err != nil {
			panic(err)
		}
		admin = &entity.User{
			Phone:    adminCred.Phone,
			Password: hashedPassword,
		}
		err = roleSeeder.userRepository.CreateUser(roleSeeder.db, admin)
		if err != nil {
			panic(err)
		}
	}
	return admin
}
