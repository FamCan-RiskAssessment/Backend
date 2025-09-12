package repository

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormRepository interface {
	CreateForm(db database.Database, form *entity.Form) error
	FindFormByID(db database.Database, id uint) (*entity.Form, error)
	FindFormsByUserID(db database.Database, userID uint, offset, limit int) ([]*entity.Form, error)
	CountFormsByUserID(db database.Database, userID uint) (int64, error)
	UpdateForm(db database.Database, form *entity.Form) error
	DeleteForm(db database.Database, id uint) error
	FindAllForms(db database.Database, offset, limit int) ([]*entity.Form, error)
	CountAllForms(db database.Database) (int64, error)
}
