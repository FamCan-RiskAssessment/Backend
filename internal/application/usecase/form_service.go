package usecase

import (
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
)

type FormService interface {
	CreateForm(request formdto.CreateFormRequest) error
	GetForm(formID uint) (formdto.FormResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.FormResponse, int64, error)
	UpdateForm(request formdto.UpdateFormRequest) error
	DeleteForm(formID uint) error
	GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.FormResponse, int64, error)
	AcceptForm(formID uint) error
	RejectForm(formID uint) error
}
