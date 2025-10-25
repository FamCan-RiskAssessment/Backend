package entity

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type ActionLog struct {
	database.Model
	ActorID    uint            `gorm:"not null"`
	Actor      User            `gorm:"foreignKey:ActorID"`
	TargetID   *uint           `gorm:"index"`
	Target     *User           `gorm:"foreignKey:TargetID"`
	Action     enum.ActionType `gorm:"not null;index"`
	Resource   *string         `gorm:"type:varchar(100)"`
	ResourceID *uint           `gorm:"index"`
	Details    string          `gorm:"type:text"`
}
