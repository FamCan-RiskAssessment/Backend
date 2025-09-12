package service

import (
	"strconv"
	"time"

	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormService struct {
	formRepository postgres.FormRepository
	userRepository postgres.UserRepository
	db             database.Database
}

func NewFormService(
	formRepository postgres.FormRepository,
	userRepository postgres.UserRepository,
	db database.Database,
) *FormService {
	return &FormService{
		formRepository: formRepository,
		userRepository: userRepository,
		db:             db,
	}
}

func (s *FormService) CreateForm(userID uint, request formdto.CreateFormRequest) (formdto.CreateFormResponse, error) {
	// Verify user exists
	user, err := s.userRepository.FindUserByID(s.db, userID)
	if err != nil {
		return formdto.CreateFormResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: strconv.Itoa(int(userID))}
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

	err = s.formRepository.CreateForm(s.db, form)
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
			CreatedAt:            form.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            form.UpdatedAt.Format(time.RFC3339),
		},
		Message: "Form created successfully",
	}

	return response, nil
}

func (s *FormService) GetForm(formID uint) (formdto.FormResponse, error) {
	form, err := s.formRepository.FindFormByID(s.db, formID)
	if err != nil {
		return formdto.FormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: strconv.Itoa(int(formID))}
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
		CreatedAt:            form.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            form.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *FormService) GetUserForms(request formdto.GetUserFormsRequest) (formdto.GetUserFormsResponse, error) {
	// Verify user exists
	user, err := s.userRepository.FindUserByID(s.db, request.UserID)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: strconv.Itoa(int(request.UserID))}
		return formdto.GetUserFormsResponse{}, notFoundError
	}

	// Set default pagination values
	offset := request.Offset
	limit := request.Limit
	if limit <= 0 {
		limit = 10 // default limit
	}

	forms, err := s.formRepository.FindFormsByUserID(s.db, request.UserID, offset, limit)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}

	total, err := s.formRepository.CountFormsByUserID(s.db, request.UserID)
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
			CreatedAt:            form.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            form.UpdatedAt.Format(time.RFC3339),
		}
	}

	response := formdto.GetUserFormsResponse{
		Forms: formResponses,
		Total: int(total),
	}

	return response, nil
}

func (s *FormService) UpdateForm(request formdto.UpdateFormRequest) (formdto.UpdateFormResponse, error) {
	form, err := s.formRepository.FindFormByID(s.db, request.FormID)
	if err != nil {
		return formdto.UpdateFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: strconv.Itoa(int(request.FormID))}
		return formdto.UpdateFormResponse{}, notFoundError
	}

	// Update fields if provided
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

	err = s.formRepository.UpdateForm(s.db, form)
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
			CreatedAt:            form.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            form.UpdatedAt.Format(time.RFC3339),
		},
		Message: "Form updated successfully",
	}

	return response, nil
}

func (s *FormService) DeleteForm(formID uint) (formdto.DeleteFormResponse, error) {
	// Check if form exists
	form, err := s.formRepository.FindFormByID(s.db, formID)
	if err != nil {
		return formdto.DeleteFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: strconv.Itoa(int(formID))}
		return formdto.DeleteFormResponse{}, notFoundError
	}

	err = s.formRepository.DeleteForm(s.db, formID)
	if err != nil {
		return formdto.DeleteFormResponse{}, err
	}

	response := formdto.DeleteFormResponse{
		Message: "Form deleted successfully",
	}

	return response, nil
}

func (s *FormService) GetAllForms(offset, limit int) (formdto.GetUserFormsResponse, error) {
	// Set default pagination values
	if limit <= 0 {
		limit = 10 // default limit
	}

	forms, err := s.formRepository.FindAllForms(s.db, offset, limit)
	if err != nil {
		return formdto.GetUserFormsResponse{}, err
	}

	total, err := s.formRepository.CountAllForms(s.db)
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
			CreatedAt:            form.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            form.UpdatedAt.Format(time.RFC3339),
		}
	}

	response := formdto.GetUserFormsResponse{
		Forms: formResponses,
		Total: int(total),
	}

	return response, nil
}
