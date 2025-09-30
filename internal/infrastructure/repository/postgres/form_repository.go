package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type FormRepository struct{}

func NewFormRepository() *FormRepository {
	return &FormRepository{}
}

func (r *FormRepository) CreateForm(db database.Database, form *entity.Form) error {
	return db.GetDB().Create(form).Error
}

func (r *FormRepository) FindFormByID(db database.Database, id uint) (*entity.Form, error) {
	var form entity.Form
	err := db.GetDB().Preload("User").First(&form, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &form, nil
}

func (r *FormRepository) FindFormsByUserID(db database.Database, userID uint, options *postgres.QueryOptions) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Where("user_id = ?", userID).Preload("User")
	query = applyQueryOptions(query, options)
	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountFormsByUserID(db database.Database, userID uint) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.Form{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *FormRepository) UpdateForm(db database.Database, form *entity.Form) error {
	return db.GetDB().Save(form).Error
}

func (r *FormRepository) DeleteForm(db database.Database, id uint) error {
	return db.GetDB().Delete(&entity.Form{}, id).Error
}

func (r *FormRepository) FindAllForms(db database.Database, offset, limit int, filters *postgres.FormFilters) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")
	query = ApplyFormFilters(query, filters)

	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountAllForms(db database.Database, filters *postgres.FormFilters) (int64, error) {
	var count int64
	query := db.GetDB().Model(&entity.Form{})
	query = ApplyFormFilters(query, filters)
	err := query.Count(&count).Error
	return count, err
}

func ApplyFormFilters(query *gorm.DB, filters *postgres.FormFilters) *gorm.DB {
	if filters == nil {
		return query
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Gender != nil {
		query = query.Where("gender = ?", *filters.Gender)
	}
	if filters.BirthYear != nil {
		query = query.Where("birth_year = ?", *filters.BirthYear)
	}
	if filters.DrinksAlcohol != nil {
		query = query.Where("drinks_alcohol = ?", *filters.DrinksAlcohol)
	}
	if filters.SmokingNow != nil {
		query = query.Where("smoking_now = ?", *filters.SmokingNow)
	}
	if filters.Cancer != nil {
		query = query.Where("cancer = ?", *filters.Cancer)
	}

	return query
}
