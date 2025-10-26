package usecase

import (
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
)

type FormService interface {
	CreateBasicInfoForm(request formdto.CreateBasicFormRequest) (formdto.BasicFormResponse, error)
	GetBasicForm(request formdto.GetPartialFormRequest) (formdto.GetBasicFormResponse, error)
	GetGeneralHealth(request formdto.GetPartialFormRequest) (formdto.GetGeneralHealthResponse, error)
	GetMamography(request formdto.GetPartialFormRequest) (formdto.GetMamographyResponse, error)
	GetCancer(request formdto.GetPartialFormRequest) (formdto.GetCancerResponse, error)
	GetFamilyCancer(request formdto.GetPartialFormRequest) (formdto.GetFamilyCancerResponse, error)
	GetContact(request formdto.GetPartialFormRequest) (formdto.GetContactResponse, error)
	GetLungCancer(request formdto.GetPartialFormRequest) (formdto.GetLungCancerResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error)
	UpdateForm(request formdto.UpdateBasicFormRequest) error
	UpdateBasicInfo(request formdto.UpdateBasicFormRequest) error
	UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error
	UpsertMamography(request formdto.UpsertMamographyRequest) error
	UpsertCancer(request formdto.UpsertCancerRequest) error
	UpsertFamilyCancer(request formdto.UpsertFamilyCancerRequest) error
	UpsertContact(request formdto.UpsertContactRequest) error
	UpsertLungCancer(request formdto.UpsertLungCancerRequest) error

	UpdateGeneralHealth(request formdto.UpdateGeneralHealthRequest) error
	UpdateMamography(request formdto.UpdateMamographyRequest) error
	UpdateCancer(request formdto.UpdateCancerRequest) error
	UpdateFamilyCancer(request formdto.UpdateFamilyCancerRequest) error
	UpdateContact(request formdto.UpdateContactRequest) error
	UpdateLungCancer(request formdto.UpdateLungCancerRequest) error

	ChangeFormStatus(request formdto.ChangeFormStatusRequest) (formdto.ChangeFormStatusResponse, error)
	DeleteForm(request formdto.DeleteFormRequest) error
	GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.BasicFormResponse, int64, error)
	GetAllOperatorForms(offset, limit int, filters *postgres.OperatorFormFilters) ([]formdto.BasicFormResponse, int64, error)
	AcceptForm(formID uint) error
	RejectForm(formID uint) error
	AssignOperator(request formdto.AssignOperatorRequest) error
	UnassignOperator(request formdto.UnassignOperatorRequest) error

	GetAllCancerTypes() ([]generaldto.EnumResponse, error)
	GetAllGenders() ([]generaldto.EnumResponse, error)
	GetAllMenopausalStatuses() ([]generaldto.EnumResponse, error)
	GetAllFormStatuses() ([]generaldto.EnumResponse, error)
}
