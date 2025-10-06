package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormRepository interface {
	CreateForm(db database.Database, form *entity.Form) error
	CreateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error
	FindBasicInfoByFormID(db database.Database, formID uint) (*entity.BasicInfo, error)
	FindFormByID(db database.Database, id uint) (*entity.Form, error)
	FindFormsByUserID(db database.Database, userID uint, options *QueryOptions) ([]*entity.Form, error)
	CountFormsByUserID(db database.Database, userID uint) (int64, error)
	UpdateForm(db database.Database, form *entity.Form) error
	UpdateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error
	DeleteForm(db database.Database, id uint) error
	FindAllForms(db database.Database, offset, limit int) ([]*entity.Form, error)
	CountAllForms(db database.Database) (int64, error)

	FindGeneralHealthByFormID(db database.Database, formID uint) (*entity.GeneralHealthInfo, error)
	CreateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error
	UpdateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error

	FindMamographyByFormID(db database.Database, formID uint) (*entity.MamoGraphyInfo, error)
	CreateMamography(db database.Database, info *entity.MamoGraphyInfo) error
	UpdateMamography(db database.Database, info *entity.MamoGraphyInfo) error

	FindCancerByFormID(db database.Database, formID uint) (*entity.CancerInfo, error)
	CreateCancer(db database.Database, info *entity.CancerInfo) error
	UpdateCancer(db database.Database, info *entity.CancerInfo) error

	FindFamilyCancerByFormID(db database.Database, formID uint) (*entity.FamilyCancerInfo, error)
	CreateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error
	UpdateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error

	FindContactByFormID(db database.Database, formID uint) (*entity.ContactInfo, error)
	CreateContact(db database.Database, info *entity.ContactInfo) error
	UpdateContact(db database.Database, info *entity.ContactInfo) error

	FindLungCancerByFormID(db database.Database, formID uint) (*entity.LungCancerInfo, error)
	CreateLungCancer(db database.Database, info *entity.LungCancerInfo) error
	UpdateLungCancer(db database.Database, info *entity.LungCancerInfo) error
}
