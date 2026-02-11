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
	GetCancers(request formdto.GetPartialFormRequest) (formdto.GetCancersResponse, error)
	GetFamilyCancer(request formdto.GetPartialFormRequest) (formdto.GetFamilyCancerResponse, error)
	GetFamilyCancerList(request formdto.GetPartialFormRequest) (formdto.GetFamilyCancerListResponse, error)
	GetContact(request formdto.GetPartialFormRequest) (formdto.GetContactResponse, error)
	GetLungCancer(request formdto.GetPartialFormRequest) (formdto.GetLungCancerResponse, error)
	GetNavidForm(request formdto.GetPartialFormRequest) (formdto.GetNavidFormResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error)
	UpdateForm(request formdto.UpdateBasicFormRequest) error
	UpdateBasicInfo(request formdto.UpdateBasicFormRequest) error
	UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error
	UpsertMamography(request formdto.UpsertMamographyRequest) error
	CreateCancer(request formdto.CreateCancerRequest) error
	UpdateCancer(request formdto.UpdateCancerRequest) (formdto.UpdateCancerResponse, error)
	DeleteCancer(request formdto.DeleteCancerRequest) error
	CreateFamilyCancer(request formdto.CreateFamilyCancerRequest) (formdto.CreateFamilyCancerResponse, error)
	UpdateFamilyCancer(request formdto.UpdateFamilyCancerRequest) (formdto.UpdateFamilyCancerResponse, error)
	DeleteFamilyCancer(request formdto.DeleteFamilyCancerRequest) error
	UpsertContact(request formdto.UpsertContactRequest) error
	UpsertLungCancer(request formdto.UpsertLungCancerRequest) error
	UpsertNavidForm(request formdto.UpsertNavidFormRequest) error

	UpdateGeneralHealth(request formdto.UpdateGeneralHealthRequest) error
	UpdateMamography(request formdto.UpdateMamographyRequest) error
	UpdateContact(request formdto.UpdateContactRequest) error
	UpdateLungCancer(request formdto.UpdateLungCancerRequest) error
	UpdateNavidForm(request formdto.UpdateNavidFormRequest) error

	ChangeFormStatus(request formdto.ChangeFormStatusRequest) (formdto.ChangeFormStatusResponse, error)
	DeleteForm(request formdto.DeleteFormRequest) error
	GetAllForms(offset, limit int, filters *postgres.FormFilters, sortBy, sortOrder, search *string) ([]formdto.BasicFormResponse, int64, error)
	GetAllOperatorForms(offset, limit int, filters *postgres.OperatorFormFilters, sortBy, sortOrder, search *string) ([]formdto.BasicFormResponse, int64, error)
	AcceptForm(formID uint, userID uint) error
	RejectForm(formID uint, userID uint) error
	RequestPatientResponse(formID uint, userID uint) error
	RequestDocuments(formID uint, userID uint) error
	ResubmitRejectedForm(formID uint, userID uint) error
	SubmitDocuments(formID uint, userID uint) error
	AssignOperator(request formdto.AssignOperatorRequest) error
	UnassignOperator(request formdto.UnassignOperatorRequest) error

	GetAllCancerTypes() ([]generaldto.EnumResponse, error)
	GetAllGenders() ([]generaldto.EnumResponse, error)
	GetAllMenopausalStatuses() ([]generaldto.EnumResponse, error)
	GetAllFormStatuses() ([]generaldto.EnumResponse, error)
	GetAllHyperplasiaInBiopsyStatuses() ([]generaldto.EnumResponse, error)
	GetAllLifeStatuses() ([]generaldto.EnumResponse, error)
	GetAllRelativeTypes() ([]generaldto.EnumResponse, error)
	GetAllFormTypes() ([]generaldto.EnumResponse, error)
	GetAllAnswers() ([]generaldto.EnumResponse, error)
}
