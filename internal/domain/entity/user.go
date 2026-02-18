package entity

import "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"

type User struct {
	database.Model
	Phone    string `gorm:"type:varchar(11);unique"`
	Password string `gorm:"type:varchar(255)"`
	Roles    []Role `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`

	UserProfile *UserProfile `gorm:"constraint:OnDelete:CASCADE"`
}

type UserProfile struct {
	database.Model
	UserID               uint `gorm:"uniqueIndex"`
	Name                 string
	LastName             string
	HealthCenter         string
	SocialSecurityNumber string
}
