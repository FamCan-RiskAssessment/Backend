package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormService struct {
	constants      *bootstrap.Constants
	formRepository postgres.FormRepository
	userService    usecase.UserService
	db             database.Database
}

func NewFormService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	userService usecase.UserService,
	db database.Database,
) *FormService {
	return &FormService{
		constants:      constants,
		formRepository: formRepository,
		userService:    userService,
		db:             db,
	}
}

func (formService *FormService) CreateForm(request formdto.CreateFormRequest) (formdto.CreateFormResponse, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return formdto.CreateFormResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return formdto.CreateFormResponse{}, notFoundError
	}

	form := &entity.Form{
		UserID:               request.UserID,
		Name:                 request.Name,
		DateOfBirth:          request.DateOfBirth,
		Address:              request.Address,
		PostalCode:           request.PostalCode,
		SocialSecurityNumber: request.SocialSecurityNumber,
	}

	err = formService.formRepository.CreateForm(formService.db, form)
	if err != nil {
		return formdto.CreateFormResponse{}, err
	}

	response := formdto.CreateFormResponse{
		Form: formdto.FormResponse{
			ID:                   form.ID,
			UserID:               form.UserID,
			Name:                 form.Name,
			DateOfBirth:          form.DateOfBirth,
			Address:              form.Address,
			PostalCode:           form.PostalCode,
			SocialSecurityNumber: form.SocialSecurityNumber,
			CreatedAt:            form.CreatedAt,
			UpdatedAt:            form.UpdatedAt,
		},
		Message: "Form created successfully",
	}

	return response, nil
}

func (formService *FormService) GetForm(formID uint) (formdto.FormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.FormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.FormResponse{}, notFoundError
	}

	response := formdto.FormResponse{
		ID:                   form.ID,
		UserID:               form.UserID,
		Name:                 form.Name,
		DateOfBirth:          form.DateOfBirth,
		Address:              form.Address,
		PostalCode:           form.PostalCode,
		SocialSecurityNumber: form.SocialSecurityNumber,
		CreatedAt:            form.CreatedAt,
		UpdatedAt:            form.UpdatedAt,
	}

	return response, nil
}

func (formService *FormService) GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.FormResponse, int64, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return nil, 0, notFoundError
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	forms, err := formService.formRepository.FindFormsByUserID(formService.db, request.UserID, options)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountFormsByUserID(formService.db, request.UserID)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.FormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.FormResponse{
			ID:                   form.ID,
			UserID:               form.UserID,
			Name:                 form.Name,
			DateOfBirth:          form.DateOfBirth,
			Address:              form.Address,
			PostalCode:           form.PostalCode,
			SocialSecurityNumber: form.SocialSecurityNumber,
			CreatedAt:            form.CreatedAt,
			UpdatedAt:            form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) UpdateForm(request formdto.UpdateFormRequest) (formdto.FormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.FormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.FormResponse{}, notFoundError
	}

	if request.Name != nil {
		form.Name = *request.Name
	}
	if request.DateOfBirth != nil {
		form.DateOfBirth = *request.DateOfBirth
	}
	if request.Address != nil {
		form.Address = *request.Address
	}
	if request.PostalCode != nil {
		form.PostalCode = *request.PostalCode
	}
	if request.SocialSecurityNumber != nil {
		form.SocialSecurityNumber = *request.SocialSecurityNumber
	}

	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return formdto.FormResponse{}, err
	}

	response := formdto.FormResponse{
		ID:                   form.ID,
		UserID:               form.UserID,
		Name:                 form.Name,
		DateOfBirth:          form.DateOfBirth,
		Address:              form.Address,
		PostalCode:           form.PostalCode,
		SocialSecurityNumber: form.SocialSecurityNumber,
		CreatedAt:            form.CreatedAt,
		UpdatedAt:            form.UpdatedAt,
	}

	return response, nil
}

func (formService *FormService) DeleteForm(formID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	err = formService.formRepository.DeleteForm(formService.db, formID)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) GetAllForms(offset, limit int) ([]formdto.FormResponse, int64, error) {

	forms, err := formService.formRepository.FindAllForms(formService.db, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllForms(formService.db)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.FormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.FormResponse{
			ID:                   form.ID,
			UserID:               form.UserID,
			Name:                 form.Name,
			DateOfBirth:          form.DateOfBirth,
			Address:              form.Address,
			PostalCode:           form.PostalCode,
			SocialSecurityNumber: form.SocialSecurityNumber,
			CreatedAt:            form.CreatedAt,
			UpdatedAt:            form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}
