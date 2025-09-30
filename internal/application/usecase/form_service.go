package usecase

import formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"

type FormService interface {
	CreateBasicInfoForm(request formdto.CreateBasicFormRequest) (formdto.BasicFormResponse, error)
	GetForm(formID uint) (formdto.BasicFormResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error)
	UpdateForm(request formdto.UpdateBasicFormRequest) error
	UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error
	UpsertMamography(request formdto.UpsertMamographyRequest) error
	UpsertCancer(request formdto.UpsertCancerRequest) error
	UpsertFamilyCancer(request formdto.UpsertFamilyCancerRequest) error
	UpsertContact(request formdto.UpsertContactRequest) error
	UpsertLungCancer(request formdto.UpsertLungCancerRequest) error
	DeleteForm(formID uint) error
	GetAllForms(offset, limit int) ([]formdto.BasicFormResponse, int64, error)
	AcceptForm(formID uint) error
	RejectForm(formID uint) error
}
