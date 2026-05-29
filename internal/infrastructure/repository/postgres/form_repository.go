package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
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
		if err := tx.Where("form_id = ?", id).Delete(&entity.BCRAResult{}).Error; err != nil {
			return err
		}
		if err := tx.Where("form_id = ?", id).Delete(&entity.GailResult{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entity.Form{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *FormRepository) FindAllForms(db database.Database, options *postgres.QueryOptions, filters *postgres.FormFilters) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")

	if options != nil && options.HasSearch() {
		query = query.Joins("JOIN users ON users.id = forms.user_id")
	}

	query = ApplyFormFilters(query, filters)
	query = applyQueryOptions(query, options)

	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountAllForms(db database.Database, options *postgres.QueryOptions, filters *postgres.FormFilters) (int64, error) {
	var count int64
	query := db.GetDB().Model(&entity.Form{})

	if options != nil && options.HasSearch() {
		query = query.Joins("JOIN users ON users.id = forms.user_id")
	}

	query = ApplyFormFilters(query, filters)
	query = applySearchOnly(query, options)

	err := query.Count(&count).Error
	return count, err
}

func (r *FormRepository) FindAllOperatorForms(db database.Database, options *postgres.QueryOptions, filters *postgres.OperatorFormFilters) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")

	if options != nil && options.HasSearch() {
		query = query.Joins("JOIN users ON users.id = forms.user_id")
	}

	query = ApplyOperatorFormFilters(query, filters)
	query = applyQueryOptions(query, options)

	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountAllOperatorForms(db database.Database, options *postgres.QueryOptions, filters *postgres.OperatorFormFilters) (int64, error) {
	var count int64
	query := db.GetDB().Model(&entity.Form{})

	if options != nil && options.HasSearch() {
		query = query.Joins("JOIN users ON users.id = forms.user_id")
	}

	query = ApplyOperatorFormFilters(query, filters)
	query = applySearchOnly(query, options)

	err := query.Count(&count).Error
	return count, err
}

type formFilterJoins struct {
	basicInfo     bool
	generalHealth bool
	cancerVisit   bool
}

func (j *formFilterJoins) joinBasicInfo(query *gorm.DB) *gorm.DB {
	if !j.basicInfo {
		query = query.Joins("JOIN basic_infos ON basic_infos.form_id = forms.id")
		j.basicInfo = true
	}
	return query
}

func (j *formFilterJoins) joinGeneralHealth(query *gorm.DB) *gorm.DB {
	if !j.generalHealth {
		query = query.Joins("JOIN general_health_infos ON general_health_infos.form_id = forms.id")
		j.generalHealth = true
	}
	return query
}

func (j *formFilterJoins) joinCancerVisit(query *gorm.DB) *gorm.DB {
	if !j.cancerVisit {
		query = query.Joins("JOIN cancers ON cancers.form_id = forms.id")
		j.cancerVisit = true
	}
	return query
}

func applyCrossTableFormFilters(
	query *gorm.DB,
	gender *string,
	birthYear *uint,
	drinksAlcohol, smokingNow, cancer *bool,
	ssnHash *string,
) *gorm.DB {
	var joins formFilterJoins

	if gender != nil {
		if g, ok := enum.ParseGenderFilter(*gender); ok {
			query = joins.joinBasicInfo(query)
			query = query.Where("basic_infos.gender = ?", g)
		}
	}
	if birthYear != nil {
		year := int(*birthYear)
		query = joins.joinBasicInfo(query)
		// Frontend may send Jalali year; map to overlapping Gregorian years.
		if year < 1700 {
			query = query.Where("EXTRACT(YEAR FROM basic_infos.birth_date)::int IN (?, ?)", year+621, year+622)
		} else {
			query = query.Where("EXTRACT(YEAR FROM basic_infos.birth_date)::int = ?", year)
		}
	}
	if drinksAlcohol != nil {
		query = joins.joinGeneralHealth(query)
		answer := enum.AnswerNo
		if *drinksAlcohol {
			answer = enum.AnswerYes
		}
		query = query.Where("general_health_infos.drinks_alcohol = ?", answer)
	}
	if smokingNow != nil {
		query = joins.joinGeneralHealth(query)
		answer := enum.AnswerNo
		if *smokingNow {
			answer = enum.AnswerYes
		}
		query = query.Where("general_health_infos.smoking_now = ?", answer)
	}
	if cancer != nil {
		query = joins.joinCancerVisit(query)
		query = query.Where("cancers.cancer = ?", *cancer)
	}
	if ssnHash != nil {
		query = joins.joinBasicInfo(query)
		query = query.Where("basic_infos.social_security_number_hash = ?", *ssnHash)
	}

	return query
}

func ApplyFormFilters(query *gorm.DB, filters *postgres.FormFilters) *gorm.DB {
	if filters == nil {
		return query
	}

	if filters.FormType != nil {
		query = query.Where("forms.form_type = ?", *filters.FormType)
	}
	if filters.Status != nil {
		query = query.Where("forms.status = ?", *filters.Status)
	}
	if filters.FilledByOperatorID != nil {
		query = query.Where("forms.filled_by_operator_id = ?", *filters.FilledByOperatorID)
	}

	query = applyCrossTableFormFilters(
		query,
		filters.Gender,
		filters.BirthYear,
		filters.DrinksAlcohol,
		filters.SmokingNow,
		filters.Cancer,
		filters.SSNHash,
	)

	return query
}

func ApplyOperatorFormFilters(query *gorm.DB, filters *postgres.OperatorFormFilters) *gorm.DB {
	if filters == nil {
		return query
	}

	query = query.Where("forms.operator_id = ?", filters.OperatorID)
	if filters.FormType != nil {
		query = query.Where("forms.form_type = ?", *filters.FormType)
	}
	if filters.Status != nil {
		query = query.Where("forms.status = ?", *filters.Status)
	}

	query = applyCrossTableFormFilters(
		query,
		filters.Gender,
		filters.BirthYear,
		filters.DrinksAlcohol,
		filters.SmokingNow,
		filters.Cancer,
		filters.SSNHash,
	)

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

func (r *FormRepository) FindCancersByFormID(db database.Database, formID uint) ([]*entity.CancerInfo, error) {
	var info []*entity.CancerInfo
	err := db.GetDB().
		Where("form_id = ?", formID).
		Find(&info).Error

	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *FormRepository) FindCancerByID(db database.Database, cancerID uint) (*entity.CancerInfo, error) {
	var info entity.CancerInfo
	err := db.GetDB().First(&info, cancerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateCancer(db database.Database, info *entity.CancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) CreateCancerVisit(db database.Database, info *entity.Cancer) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateCancerVisit(db database.Database, info *entity.Cancer) error {
	return db.GetDB().Save(info).Error
}

func (r *FormRepository) FindCancerVisit(db database.Database, formID uint) (*entity.Cancer, error) {
	var info entity.Cancer
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) UpdateCancer(db database.Database, info *entity.CancerInfo) error {
	return db.GetDB().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(info).Error
}

func (r *FormRepository) FindFamilyCancersByFormID(db database.Database, formID uint) ([]*entity.FamilyCancerInfo, error) {
	var info []*entity.FamilyCancerInfo
	err := db.GetDB().
		Where("form_id = ?", formID).
		Order("relative ASC").
		Find(&info).Error

	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *FormRepository) FindFamilyCancerByID(db database.Database, familyCancerID uint) (*entity.FamilyCancerInfo, error) {
	var info entity.FamilyCancerInfo
	err := db.GetDB().First(&info, familyCancerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) DeleteCancersByFormID(db database.Database, formID uint) error {
	return db.GetDB().
		Where("form_id = ?", formID).
		Delete(&entity.CancerInfo{}).Error
}

func (r *FormRepository) DeleteCancerByID(db database.Database, cancerID uint) error {
	return db.GetDB().
		Delete(&entity.CancerInfo{}, cancerID).Error
}

func (r *FormRepository) DeleteFamilyCancersByFormID(db database.Database, formID uint) error {
	return db.GetDB().
		Where("form_id = ?", formID).
		Delete(&entity.FamilyCancerInfo{}).Error
}

func (r *FormRepository) DeleteFamilyCancerByID(db database.Database, familyCancerID uint) error {
	return db.GetDB().
		Delete(&entity.FamilyCancerInfo{}, familyCancerID).Error
}

func (r *FormRepository) CreateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	return db.GetDB().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(info).Error
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

func (r *FormRepository) FindNavidInfoByFormID(db database.Database, formID uint) (*entity.NavidInfo, error) {
	var info entity.NavidInfo
	err := db.GetDB().Where("form_id = ?", formID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *FormRepository) CreateNavidInfo(db database.Database, info *entity.NavidInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateNavidInfo(db database.Database, info *entity.NavidInfo) error {
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

func (r *FormRepository) FindGailResultByFormID(db database.Database, formID uint) (*entity.GailResult, error) {
	var result entity.GailResult
	err := db.GetDB().Where("form_id = ?", formID).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *FormRepository) CreateGailResult(db database.Database, result *entity.GailResult) error {
	return db.GetDB().Create(result).Error
}

func (r *FormRepository) UpdateGailResult(db database.Database, result *entity.GailResult) error {
	return db.GetDB().Save(result).Error
}

func (r *FormRepository) FindPLCOResultByFormID(db database.Database, formID uint) (*entity.PLCOResult, error) {
	var result entity.PLCOResult
	err := db.GetDB().Where("form_id = ?", formID).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *FormRepository) CreatePLCOResult(db database.Database, result *entity.PLCOResult) error {
	return db.GetDB().Create(result).Error
}

func (r *FormRepository) UpdatePLCOResult(db database.Database, result *entity.PLCOResult) error {
	return db.GetDB().Save(result).Error
}

func (r *FormRepository) FindAttentionQuestionsByFormID(db database.Database, formID uint) (*entity.AttentionQuestions, error) {
	var result entity.AttentionQuestions
	err := db.GetDB().Where("form_id = ?", formID).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *FormRepository) CreateAttentionQuestions(db database.Database, questions *entity.AttentionQuestions) error {
	return db.GetDB().Create(questions).Error
}

func (r *FormRepository) UpdateAttentionQuestions(db database.Database, questions *entity.AttentionQuestions) error {
	return db.GetDB().Save(questions).Error
}
