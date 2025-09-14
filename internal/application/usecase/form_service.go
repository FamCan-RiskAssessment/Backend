package usecase

import formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"

type FormService interface {
	CreateForm(request formdto.CreateFormRequest) (formdto.CreateFormResponse, error)
	GetForm(formID uint) (formdto.FormResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.FormResponse, int64, error)
	UpdateForm(request formdto.UpdateFormRequest) (formdto.FormResponse, error)
	DeleteForm(formID uint) error
	GetAllForms(offset, limit int) ([]formdto.FormResponse, int64, error)
}
