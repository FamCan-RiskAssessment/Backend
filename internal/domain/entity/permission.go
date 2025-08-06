package entity

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type Permission struct {
	database.Model
	Name     string                  `gorm:"not null"`
	Category enum.PermissionCategory `gorm:"not null"`
	Roles    []Role                  `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}
