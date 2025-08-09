package entity

import "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"

type User struct {
	database.Model
	Phone string `gorm:"type:varchar(11);unique"`
	Roles []Role `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
}
