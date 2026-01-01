package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type FormRepositoryMock struct {
	mock.Mock
}

func NewFormRepositoryMock() *FormRepositoryMock {
	return &FormRepositoryMock{}
}

func (f *FormRepositoryMock) CreateForm(db database.Database, form *entity.Form) error {
	args := f.Called(db, form)
	return args.Error(0)
}

func (f *FormRepositoryMock) CreateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error {
	args := f.Called(db, basicInfo)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindBasicInfoByFormID(db database.Database, formID uint) (*entity.BasicInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BasicInfo), args.Error(1)
}

func (f *FormRepositoryMock) FindFormByID(db database.Database, id uint) (*entity.Form, error) {
	args := f.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Form), args.Error(1)
}

func (f *FormRepositoryMock) FindFormsByUserID(db database.Database, userID uint, options *postgres.QueryOptions) ([]*entity.Form, error) {
	args := f.Called(db, userID, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Form), args.Error(1)
}

func (f *FormRepositoryMock) CountFormsByUserID(db database.Database, userID uint) (int64, error) {
	args := f.Called(db, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (f *FormRepositoryMock) UpdateForm(db database.Database, form *entity.Form) error {
	args := f.Called(db, form)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateBasicInfo(db database.Database, basicInfo *entity.BasicInfo) error {
	args := f.Called(db, basicInfo)
	return args.Error(0)
}

func (f *FormRepositoryMock) DeleteForm(db database.Database, id uint) error {
	args := f.Called(db, id)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindAllForms(db database.Database, offset, limit int, filters *postgres.FormFilters) ([]*entity.Form, error) {
	args := f.Called(db, offset, limit, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Form), args.Error(1)
}

func (f *FormRepositoryMock) CountAllForms(db database.Database, filters *postgres.FormFilters) (int64, error) {
	args := f.Called(db, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (f *FormRepositoryMock) FindAllOperatorForms(db database.Database, offset, limit int, filters *postgres.OperatorFormFilters) ([]*entity.Form, error) {
	args := f.Called(db, offset, limit, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Form), args.Error(1)
}

func (f *FormRepositoryMock) CountAllOperatorForms(db database.Database, filters *postgres.OperatorFormFilters) (int64, error) {
	args := f.Called(db, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (f *FormRepositoryMock) FindGeneralHealthByFormID(db database.Database, formID uint) (*entity.GeneralHealthInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.GeneralHealthInfo), args.Error(1)
}

func (f *FormRepositoryMock) CreateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateGeneralHealth(db database.Database, info *entity.GeneralHealthInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindMamographyByFormID(db database.Database, formID uint) (*entity.MamoGraphyInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MamoGraphyInfo), args.Error(1)
}

func (f *FormRepositoryMock) CreateMamography(db database.Database, info *entity.MamoGraphyInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateMamography(db database.Database, info *entity.MamoGraphyInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindCancersByFormID(db database.Database, formID uint) ([]*entity.CancerInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.CancerInfo), args.Error(1)
}

func (f *FormRepositoryMock) FindCancerByID(db database.Database, cancerID uint) (*entity.CancerInfo, error) {
	args := f.Called(db, cancerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CancerInfo), args.Error(1)
}

func (f *FormRepositoryMock) DeleteCancersByFormID(db database.Database, formID uint) error {
	args := f.Called(db, formID)
	return args.Error(0)
}

func (f *FormRepositoryMock) DeleteCancerByID(db database.Database, cancerID uint) error {
	args := f.Called(db, cancerID)
	return args.Error(0)
}

func (f *FormRepositoryMock) CreateCancer(db database.Database, info *entity.CancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateCancer(db database.Database, info *entity.CancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindFamilyCancersByFormID(db database.Database, formID uint) ([]*entity.FamilyCancerInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FamilyCancerInfo), args.Error(1)
}

func (f *FormRepositoryMock) FindFamilyCancerByID(db database.Database, familyCancerID uint) (*entity.FamilyCancerInfo, error) {
	args := f.Called(db, familyCancerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyCancerInfo), args.Error(1)
}

func (f *FormRepositoryMock) DeleteFamilyCancersByFormID(db database.Database, formID uint) error {
	args := f.Called(db, formID)
	return args.Error(0)
}

func (f *FormRepositoryMock) DeleteFamilyCancerByID(db database.Database, familyCancerID uint) error {
	args := f.Called(db, familyCancerID)
	return args.Error(0)
}

func (f *FormRepositoryMock) CreateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateFamilyCancer(db database.Database, info *entity.FamilyCancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindContactByFormID(db database.Database, formID uint) (*entity.ContactInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactInfo), args.Error(1)
}

func (f *FormRepositoryMock) CreateContact(db database.Database, info *entity.ContactInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateContact(db database.Database, info *entity.ContactInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindLungCancerByFormID(db database.Database, formID uint) (*entity.LungCancerInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LungCancerInfo), args.Error(1)
}

func (f *FormRepositoryMock) CreateLungCancer(db database.Database, info *entity.LungCancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateLungCancer(db database.Database, info *entity.LungCancerInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindNavidInfoByFormID(db database.Database, formID uint) (*entity.NavidInfo, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.NavidInfo), args.Error(1)
}

func (f *FormRepositoryMock) CreateNavidInfo(db database.Database, info *entity.NavidInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateNavidInfo(db database.Database, info *entity.NavidInfo) error {
	args := f.Called(db, info)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindPremm5ResultByFormID(db database.Database, formID uint) (*entity.Premm5Result, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Premm5Result), args.Error(1)
}

func (f *FormRepositoryMock) CreatePremm5Result(db database.Database, result *entity.Premm5Result) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdatePremm5Result(db database.Database, result *entity.Premm5Result) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindBCRAResultByFormID(db database.Database, formID uint) (*entity.BCRAResult, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BCRAResult), args.Error(1)
}

func (f *FormRepositoryMock) CreateBCRAResult(db database.Database, result *entity.BCRAResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateBCRAResult(db database.Database, result *entity.BCRAResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindGailResultByFormID(db database.Database, formID uint) (*entity.GailResult, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.GailResult), args.Error(1)
}

func (f *FormRepositoryMock) CreateGailResult(db database.Database, result *entity.GailResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateGailResult(db database.Database, result *entity.GailResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindPLCOResultByFormID(db database.Database, formID uint) (*entity.PLCOResult, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PLCOResult), args.Error(1)
}

func (f *FormRepositoryMock) CreatePLCOResult(db database.Database, result *entity.PLCOResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdatePLCOResult(db database.Database, result *entity.PLCOResult) error {
	args := f.Called(db, result)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindAttentionQuestionsByFormID(db database.Database, formID uint) (*entity.AttentionQuestions, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AttentionQuestions), args.Error(1)
}

func (f *FormRepositoryMock) CreateAttentionQuestions(db database.Database, questions *entity.AttentionQuestions) error {
	args := f.Called(db, questions)
	return args.Error(0)
}

func (f *FormRepositoryMock) UpdateAttentionQuestions(db database.Database, questions *entity.AttentionQuestions) error {
	args := f.Called(db, questions)
	return args.Error(0)
}

func (f *FormRepositoryMock) CreateCalculationHistory(db database.Database, history *entity.CalculationHistory) error {
	args := f.Called(db, history)
	return args.Error(0)
}

func (f *FormRepositoryMock) FindCalculationHistoryByFormID(db database.Database, formID uint) ([]*entity.CalculationHistory, error) {
	args := f.Called(db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.CalculationHistory), args.Error(1)
}
