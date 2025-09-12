package usecase

import formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"

type FormService interface {
	CreateForm(userID uint, request formdto.CreateFormRequest) (formdto.CreateFormResponse, error)
	GetForm(formID uint) (formdto.FormResponse, error)
	GetUserForms(request formdto.GetUserFormsRequest) (formdto.GetUserFormsResponse, error)
	UpdateForm(request formdto.UpdateFormRequest) (formdto.UpdateFormResponse, error)
	DeleteForm(formID uint) (formdto.DeleteFormResponse, error)
	GetAllForms(offset, limit int) (formdto.GetUserFormsResponse, error)
}
