package repository

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (userRepository *UserRepository) GetUserByPhone(db database.Database, phone string) (*entity.User, error) {
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
