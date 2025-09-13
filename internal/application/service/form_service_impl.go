package service

import (
	"fmt"

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

func (formService *FormService) CreateForm(userID uint, request formdto.CreateFormRequest) (formdto.CreateFormResponse, error) {
	user, err := formService.userService.GetUserByID(userID)
	if err != nil {
		return formdto.CreateFormResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return formdto.CreateFormResponse{}, notFoundError
	}

	form := &entity.Form{
		UserID:               userID,
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
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
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

func (formService *FormService) GetUserForms(request formdto.GetUserFormsRequest) (formdto.GetUserFormsResponse, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	fmt.Println("user", user, err)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return formdto.GetUserFormsResponse{}, notFoundError
	}

	offset := request.Offset
	limit := request.Limit
	if limit <= 0 {
		limit = 10
	}

	forms, err := formService.formRepository.FindFormsByUserID(formService.db, request.UserID, offset, limit)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}

	total, err := formService.formRepository.CountFormsByUserID(formService.db, request.UserID)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
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

	response := formdto.GetUserFormsResponse{
		Forms: formResponses,
		Total: int(total),
	}

	return response, nil
}

func (formService *FormService) UpdateForm(request formdto.UpdateFormRequest) (formdto.UpdateFormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.UpdateFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.UpdateFormResponse{}, notFoundError
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
		return formdto.UpdateFormResponse{}, err
	}

	response := formdto.UpdateFormResponse{
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
		Message: "Form updated successfully",
	}

	return response, nil
}

func (formService *FormService) DeleteForm(formID uint) (formdto.DeleteFormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.DeleteFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.DeleteFormResponse{}, notFoundError
	}

	err = formService.formRepository.DeleteForm(formService.db, formID)
	if err != nil {
		return formdto.DeleteFormResponse{}, err
	}

	response := formdto.DeleteFormResponse{
		Message: "Form deleted successfully",
	}

	return response, nil
}

func (formService *FormService) GetAllForms(offset, limit int) (formdto.GetUserFormsResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	forms, err := formService.formRepository.FindAllForms(formService.db, offset, limit)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}

	total, err := formService.formRepository.CountAllForms(formService.db)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
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

	response := formdto.GetUserFormsResponse{
		Forms: formResponses,
		Total: int(total),
	}

	return response, nil
}
