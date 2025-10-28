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

func (r *FormRepository) CreateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error {
	return db.GetDB().Create(basicInfo).Error
}

func (r *FormRepository) FindBasicInfoByFormID(db database.Database, formID uint) (*entity.BasicInfo, error) {
	var info entity.BasicInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
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

func (r *FormRepository) UpdateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error {
	return db.GetDB().Save(basicInfo).Error
}

func (r *FormRepository) DeleteForm(db database.Database, id uint) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("form_id = ?", id).Delete(&entity.BasicInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.GeneralHealthInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.MamoGraphyInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.CancerInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.FamilyCancerInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.ContactInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.LungCancerInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.Premm5Result{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entity.Form{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *FormRepository) FindAllForms(db database.Database, offset, limit int, filters *postgres.FormFilters) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")
	query = ApplyFormFilters(query, filters)

	query = query.Order("id ASC")

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

func (r *FormRepository) FindAllOperatorForms(db database.Database, offset, limit int, filters *postgres.OperatorFormFilters) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")
	query = ApplyOperatorFormFilters(query, filters)

	query = query.Order("id ASC")

	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountAllOperatorForms(db database.Database, filters *postgres.OperatorFormFilters) (int64, error) {
	var count int64
	query := db.GetDB().Model(&entity.Form{})
	query = ApplyOperatorFormFilters(query, filters)
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
func ApplyOperatorFormFilters(query *gorm.DB, filters *postgres.OperatorFormFilters) *gorm.DB {
	if filters == nil {
		return query
	}

	query = query.Where("operator_id = ?", filters.OperatorID)
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
func (r *FormRepository) FindGeneralHealthByFormID(db database.Database, formID uint) (*entity.GeneralHealthInfo, error) {
	var info entity.GeneralHealthInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindMamographyByFormID(db database.Database, formID uint) (*entity.MamoGraphyInfo, error) {
	var info entity.MamoGraphyInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateMamography(db database.Database, info *entity.MamoGraphyInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateMamography(db database.Database, info *entity.MamoGraphyInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindCancerByFormID(db database.Database, formID uint) (*entity.CancerInfo, error) {
	var info entity.CancerInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) FindAllCancersByFormID(db database.Database, formID uint) ([]*entity.NewCancerInfo, error) {
	var info []*entity.NewCancerInfo
	err := db.GetDB().
		Preload("Cancer").
		Where("form_id = ?", formID).
		Find(&info).Error

	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *FormRepository) CreateCancer(db database.Database, info *entity.CancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) CreateNewCancer(db database.Database, info *entity.NewCancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateCancer(db database.Database, info *entity.CancerInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) UpdateNewCancer(db database.Database, info *entity.NewCancerInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindFamilyCancerByFormID(db database.Database, formID uint) (*entity.FamilyCancerInfo, error) {
	var info entity.FamilyCancerInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) DeleteAllCancersByFormID(db database.Database, formID uint) error {
	return db.GetDB().Where("form_id = ?", formID).Delete(&entity.NewCancerInfo{}).Error
}

func (r *FormRepository) CreateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindContactByFormID(db database.Database, formID uint) (*entity.ContactInfo, error) {
	var info entity.ContactInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateContact(db database.Database, info *entity.ContactInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateContact(db database.Database, info *entity.ContactInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindLungCancerByFormID(db database.Database, formID uint) (*entity.LungCancerInfo, error) {
	var info entity.LungCancerInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateLungCancer(db database.Database, info *entity.LungCancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateLungCancer(db database.Database, info *entity.LungCancerInfo) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindPremm5ResultByFormID(db database.Database, formID uint) (*entity.Premm5Result, error) {
	var result entity.Premm5Result
	err := db.GetDB().Where("form_id = ?", formID).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *FormRepository) CreatePremm5Result(db database.Database, result *entity.Premm5Result) error {
	return db.GetDB().Create(result).Error
}

func (r *FormRepository) UpdatePremm5Result(db database.Database, result *entity.Premm5Result) error {
	return db.GetDB().Save(result).Error
}

func (r *FormRepository) FindBCRAResultByFormID(db database.Database, formID uint) (*entity.BCRAResult, error) {
	var result entity.BCRAResult
	err := db.GetDB().Where("form_id = ?", formID).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *FormRepository) CreateBCRAResult(db database.Database, result *entity.BCRAResult) error {
	return db.GetDB().Create(result).Error
}

func (r *FormRepository) UpdateBCRAResult(db database.Database, result *entity.BCRAResult) error {
	return db.GetDB().Save(result).Error
}
