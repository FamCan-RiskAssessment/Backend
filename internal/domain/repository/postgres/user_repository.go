package repository

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	GetUserByPhone(db database.Database, phone string) (*entity.User, error)
	CreateUser(db database.Database, user *entity.User) error
	FindUserRoles(db database.Database, user *entity.User) error
	FindRolePermissions(db database.Database, role *entity.Role) error
}
