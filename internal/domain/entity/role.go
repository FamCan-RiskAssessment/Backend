package entity

import "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"

type Role struct {
	database.Model
	Name        string       `gorm:"not null"`
	Users       []User       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}
