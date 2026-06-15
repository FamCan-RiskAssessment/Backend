package seed

import (
	"fmt"

	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type SuperAdminPasswordChanger struct {
	db             database.Database
	userRepository postgres.UserRepository
	passwordHasher usecase.PasswordHasher
}

func NewSuperAdminPasswordChanger(
	db database.Database,
	userRepository postgres.UserRepository,
	passwordHasher usecase.PasswordHasher,
) *SuperAdminPasswordChanger {
	return &SuperAdminPasswordChanger{
		db:             db,
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (c *SuperAdminPasswordChanger) ChangePassword(phone, password string) error {
	if phone == "" {
		return fmt.Errorf("super admin phone is required")
	}
	if password == "" {
		return fmt.Errorf("new password is required")
	}

	user, err := c.userRepository.FindUserByPhone(c.db, phone)
	if err != nil {
		return fmt.Errorf("find user by phone: %w", err)
	}
	if user == nil {
		return fmt.Errorf("no user found with phone %s", phone)
	}

	if err := c.userRepository.FindUserRoles(c.db, user); err != nil {
		return fmt.Errorf("find user roles: %w", err)
	}
	if !hasSuperAdminRole(user) {
		return fmt.Errorf("user with phone %s is not a super admin", phone)
	}

	hashedPassword, err := c.passwordHasher.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.Password = hashedPassword
	if err := c.userRepository.UpdateUser(c.db, user); err != nil {
		return fmt.Errorf("update user password: %w", err)
	}

	return nil
}

func hasSuperAdminRole(user *entity.User) bool {
	for _, role := range user.Roles {
		if role.Name == enum.SuperAdmin.String() {
			return true
		}
	}
	return false
}
