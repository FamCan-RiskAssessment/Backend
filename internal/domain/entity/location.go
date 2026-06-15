package entity

import "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"

type Province struct {
	database.Model
	Name   string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Cities []City `gorm:"foreignKey:ProvinceID"`
}

type City struct {
	database.Model
	Name       string   `gorm:"type:varchar(100);not null;uniqueIndex:idx_city_name_province"`
	ProvinceID uint     `gorm:"not null;index;uniqueIndex:idx_city_name_province"`
	Province   Province `gorm:"foreignKey:ProvinceID;constraint:OnDelete:CASCADE"`
}
