package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormRepository interface {
	CreateForm(db database.Database, form *entity.Form) error
	FindFormByID(db database.Database, id uint) (*entity.Form, error)
	FindFormsByUserID(db database.Database, userID uint, options *QueryOptions) ([]*entity.Form, error)
	CountFormsByUserID(db database.Database, userID uint) (int64, error)
	UpdateForm(db database.Database, form *entity.Form) error
	DeleteForm(db database.Database, id uint) error
	FindAllForms(db database.Database, offset, limit int, filters *FormFilters) ([]*entity.Form, error)
	CountAllForms(db database.Database, filters *FormFilters) (int64, error)
}

type FormFilters struct {
	Status        *uint
	Gender        *string
	BirthYear     *uint
	DrinksAlcohol *bool
	SmokingNow    *bool
	Cancer        *bool
}
