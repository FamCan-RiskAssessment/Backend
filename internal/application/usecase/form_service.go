package usecase

import (
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
)

type FormService interface {
	CreateBasicInfoForm(request formdto.CreateBasicFormRequest) (formdto.BasicFormResponse, error)
	GetBasicForm(formID uint) (formdto.GetBasicFormResponse, error)
	GetGeneralHealth(formID uint) (formdto.GetGeneralHealthResponse, error)
	GetMamography(formID uint) (formdto.GetMamographyResponse, error)
	GetCancer(formID uint) (formdto.GetCancerResponse, error)
	GetFamilyCancer(formID uint) (formdto.GetFamilyCancerResponse, error)
	GetContact(formID uint) (formdto.GetContactResponse, error)
	GetLungCancer(formID uint) (formdto.GetLungCancerResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error)
	UpdateForm(request formdto.UpdateBasicFormRequest) error
	UpdateBasicInfo(request formdto.UpdateBasicFormRequest) error
	UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error
	UpsertMamography(request formdto.UpsertMamographyRequest) error
	UpsertCancer(request formdto.UpsertCancerRequest) error
	UpsertFamilyCancer(request formdto.UpsertFamilyCancerRequest) error
	UpsertContact(request formdto.UpsertContactRequest) error
	UpsertLungCancer(request formdto.UpsertLungCancerRequest) error
	ChangeFormStatus(request formdto.ChangeFormStatusRequest) (formdto.ChangeFormStatusResponse, error)
	DeleteForm(formID uint) error
	GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.FormResponse, int64, error)
	AcceptForm(formID uint) error
	RejectForm(formID uint) error
}
