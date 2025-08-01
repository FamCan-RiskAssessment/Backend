package repository

import "github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"

type UserRepository interface {
	GetUserByPhone(phone string) (*entity.User, error)
	CreateUser(phone string) error
}
