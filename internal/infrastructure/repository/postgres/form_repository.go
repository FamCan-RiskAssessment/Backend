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

		if err := tx.Delete(&entity.Form{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *FormRepository) FindAllForms(db database.Database, offset, limit int) ([]*entity.Form, error) {
	var forms []*entity.Form
	query := db.GetDB().Preload("User")

	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&forms).Error
	return forms, err
}

func (r *FormRepository) CountAllForms(db database.Database) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.Form{}).Count(&count).Error
	return count, err
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

func (r *FormRepository) CreateCancer(db database.Database, info *entity.CancerInfo) error {
	return db.GetDB().Create(info).Error
}

func (r *FormRepository) UpdateCancer(db database.Database, info *entity.CancerInfo) error {
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
