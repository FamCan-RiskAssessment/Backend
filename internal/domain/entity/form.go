package entity

import "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"

type Form struct {
	database.Model
	UserID               uint   `gorm:"not null;index"`
	User                 User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name                 string `gorm:"type:varchar(255);not null"`
	DateOfBirth          string `gorm:"type:date;not null"`
	Address              string `gorm:"type:text;not null"`
	PostalCode           string `gorm:"type:varchar(20);not null"`
	SocialSecurityNumber string `gorm:"type:varchar(20);not null"`
}
