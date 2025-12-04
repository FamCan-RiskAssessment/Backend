package service

import (
	"errors"
	"testing"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	usecaseMocks "github.com/FamCan-RiskAssessment/Backend/mocks/application/usecase"
	externalMocks "github.com/FamCan-RiskAssessment/Backend/mocks/domain/external"
	formRepositoryMocks "github.com/FamCan-RiskAssessment/Backend/mocks/domain/repository/postgres"
	s3Mocks "github.com/FamCan-RiskAssessment/Backend/mocks/domain/storage/s3"
	databaseMocks "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FormServiceTestSuite struct {
	suite.Suite
	constants          *bootstrap.Constants
	formRepository     *formRepositoryMocks.FormRepositoryMock
	userService        *usecaseMocks.UserServiceMock
	actionLogService   *usecaseMocks.ActionLogServiceMock
	s3Storage          *s3Mocks.S3StorageMock
	db                 *databaseMocks.DatabaseMock
	verificationClient *externalMocks.VerificationClientMock
	formService        *FormService
}

func (suite *FormServiceTestSuite) SetupTest() {
	suite.constants = bootstrap.NewConstants()
	suite.formRepository = formRepositoryMocks.NewFormRepositoryMock()
	suite.userService = usecaseMocks.NewUserServiceMock()
	suite.actionLogService = usecaseMocks.NewActionLogServiceMock()
	suite.s3Storage = s3Mocks.NewS3StorageMock()
	suite.db = databaseMocks.NewDatabaseMock()
	suite.verificationClient = externalMocks.NewVerificationClientMock()

	suite.formService = NewFormService(
		suite.constants,
		suite.formRepository,
		suite.userService,
		suite.actionLogService,
		suite.s3Storage,
		suite.db,
		suite.verificationClient,
	)
}

// Test: CreateBasicInfoForm should successfully create form and basic info
func (suite *FormServiceTestSuite) TestCreateBasicInfoForm_Success() {
	// Arrange
	userID := uint(1)
	birthDate, _ := time.Parse("2006-01-02", "1990-01-01")
	request := formdto.CreateBasicFormRequest{
		UserID:               userID,
		Gender:               1,
		BirthDate:            birthDate,
		IsAtba:               false,
		SocialSecurityNumber: "123456789",
		Height:               170,
		Weight:               70,
	}

	user := &entity.User{}
	user.ID = userID
	user.Phone = "09123456789"

	form := &entity.Form{}
	form.ID = 1
	form.UserID = userID
	form.Status = enum.FormStatusPending
	form.FilledByOperatorID = nil
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(user, nil)
	suite.verificationClient.On("VerifyPhoneAndSSN", user.Phone, request.SocialSecurityNumber).Return(true, nil)
	suite.formRepository.On("CreateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.UserID == userID && f.Status == enum.FormStatusPending
	})).Run(func(args mock.Arguments) {
		form := args.Get(1).(*entity.Form)
		form.ID = 1
		form.CreatedAt = time.Now()
		form.UpdatedAt = time.Now()
	}).Return(nil)
	suite.formRepository.On("CreateBasicInfo", suite.db, mock.MatchedBy(func(b *entity.BasicInfo) bool {
		return b.FormID == form.ID
	})).Return(nil)
	suite.formRepository.On("FindAttentionQuestionsByFormID", suite.db, mock.Anything).Return(nil, nil)

	// Act
	response, err := suite.formService.CreateBasicInfoForm(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, response.UserID)
	assert.Equal(suite.T(), enum.FormStatusPending.String(), response.Status)
	suite.userService.AssertExpectations(suite.T())
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: CreateBasicInfoForm should return error when user not found
func (suite *FormServiceTestSuite) TestCreateBasicInfoForm_UserNotFound() {
	// Arrange
	userID := uint(999)
	request := formdto.CreateBasicFormRequest{
		UserID: userID,
	}

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(nil, nil)

	// Act
	response, err := suite.formService.CreateBasicInfoForm(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.BasicFormResponse{}, response)
	notFoundError, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.constants.Field.User, notFoundError.Item)
}

// Test: CreateBasicInfoForm should return error when user service fails
func (suite *FormServiceTestSuite) TestCreateBasicInfoForm_UserServiceError() {
	// Arrange
	userID := uint(1)
	request := formdto.CreateBasicFormRequest{
		UserID: userID,
	}
	serviceError := errors.New("database connection error")

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(nil, serviceError)

	// Act
	response, err := suite.formService.CreateBasicInfoForm(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), response, formdto.BasicFormResponse{})
	assert.Equal(suite.T(), serviceError, err)
}

// Test: CreateBasicInfoForm with operator should log action
func (suite *FormServiceTestSuite) TestCreateBasicInfoForm_WithOperator_LogsAction() {
	// Arrange
	userID := uint(1)
	operatorID := uint(2)
	birthDate, _ := time.Parse("2006-01-02", "1990-01-01")
	request := formdto.CreateBasicFormRequest{
		UserID:             userID,
		FilledByOperatorID: &operatorID,
		Gender:             1,
		BirthDate:          birthDate,
		Height:             170,
		Weight:             70,
	}

	user := &entity.User{}
	user.ID = userID
	user.Phone = "09123456789"
	operator := &entity.User{}
	operator.ID = operatorID

	form := &entity.Form{}
	form.ID = 1
	form.UserID = userID
	form.Status = enum.FormStatusPending
	form.FilledByOperatorID = &operatorID
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(user, nil)
	suite.verificationClient.On("VerifyPhoneAndSSN", user.Phone, mock.Anything).Return(true, nil)
	suite.userService.On("ValidateUserForFormCreation", operatorID, userID).Return(nil)
	suite.userService.On("GetUserByID", operatorID).Return(operator, nil)
	suite.formRepository.On("CreateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.UserID == userID
	})).Run(func(args mock.Arguments) {
		f := args.Get(1).(*entity.Form)
		f.ID = 1
		f.CreatedAt = time.Now()
		f.UpdatedAt = time.Now()
	}).Return(nil)
	suite.formRepository.On("CreateBasicInfo", suite.db, mock.MatchedBy(func(b *entity.BasicInfo) bool {
		return b.FormID == form.ID
	})).Return(nil)
	suite.formRepository.On("FindAttentionQuestionsByFormID", suite.db, mock.Anything).Return(nil, nil)
	suite.actionLogService.On("LogAction", mock.MatchedBy(func(log actionlogdto.LogAction) bool {
		return log.ActorID == operatorID && log.TargetID != nil && *log.TargetID == userID
	})).Return(nil)

	// Act
	response, err := suite.formService.CreateBasicInfoForm(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, response.UserID)
	suite.actionLogService.AssertExpectations(suite.T())
}

// Test: CreateBasicInfoForm should return error when form creation fails
func (suite *FormServiceTestSuite) TestCreateBasicInfoForm_FormCreationFails() {
	// Arrange
	userID := uint(1)
	request := formdto.CreateBasicFormRequest{
		UserID: userID,
	}
	user := &entity.User{}
	user.ID = userID
	user.Phone = "09123456789"
	repoError := errors.New("database error")

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(user, nil)
	suite.verificationClient.On("VerifyPhoneAndSSN", user.Phone, mock.Anything).Return(true, nil)
	suite.formRepository.On("CreateForm", suite.db, mock.Anything).Return(repoError)

	// Act
	response, err := suite.formService.CreateBasicInfoForm(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), response, formdto.BasicFormResponse{})
	assert.Equal(suite.T(), repoError, err)
}

// Test: isOperator should return true when user has operator role
func (suite *FormServiceTestSuite) TestIsOperator_True() {
	// Arrange
	userID := uint(1)
	roles := []userdto.RoleResponse{
		{
			ID:   1,
			Name: enum.Operator.String(),
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(roles, nil)

	// Act
	isOp, err := suite.formService.isOperator(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isOp)
	suite.userService.AssertExpectations(suite.T())
}

// Test: isOperator should return false when user doesn't have operator role
func (suite *FormServiceTestSuite) TestIsOperator_False() {
	// Arrange
	userID := uint(1)
	roles := []userdto.RoleResponse{
		{
			ID:   1,
			Name: enum.Patient.String(),
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(roles, nil)

	// Act
	isOp, err := suite.formService.isOperator(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isOp)
}

// Test: isOperator should return error when service fails
func (suite *FormServiceTestSuite) TestIsOperator_ServiceError() {
	// Arrange
	userID := uint(1)
	serviceError := errors.New("database error")

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(nil, serviceError)

	// Act
	isOp, err := suite.formService.isOperator(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.False(suite.T(), isOp)
	assert.Equal(suite.T(), serviceError, err)
}

// Test: isSupervisor should return true when user has supervisor role
func (suite *FormServiceTestSuite) TestIsSupervisor_True() {
	// Arrange
	userID := uint(1)
	roles := []userdto.RoleResponse{
		{
			ID:   1,
			Name: enum.Supervisor.String(),
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(roles, nil)

	// Act
	isSup, err := suite.formService.isSupervisor(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isSup)
}

// Test: isSupervisor should return false when user doesn't have supervisor role
func (suite *FormServiceTestSuite) TestIsSupervisor_False() {
	// Arrange
	userID := uint(1)
	roles := []userdto.RoleResponse{
		{
			ID:   1,
			Name: enum.Operator.String(),
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(roles, nil)

	// Act
	isSup, err := suite.formService.isSupervisor(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isSup)
}

// Test: canOperatorEditForm should return true when operator is FilledByOperatorID
func (suite *FormServiceTestSuite) TestCanOperatorEditForm_FilledByOperatorMatch() {
	// Arrange
	operatorID := uint(2)
	form := &entity.Form{}
	form.ID = 1
	form.UserID = 1
	form.FilledByOperatorID = &operatorID

	// Act
	canEdit, err := suite.formService.canOperatorEditForm(form, operatorID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), canEdit)
}

// Test: canOperatorEditForm should return true when operator is OperatorID
func (suite *FormServiceTestSuite) TestCanOperatorEditForm_OperatorIDMatch() {
	// Arrange
	operatorID := uint(2)
	form := &entity.Form{}
	form.ID = 1
	form.UserID = 1
	form.OperatorID = &operatorID

	// Act
	canEdit, err := suite.formService.canOperatorEditForm(form, operatorID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), canEdit)
}

// Test: canOperatorEditForm should return false when operator doesn't match
func (suite *FormServiceTestSuite) TestCanOperatorEditForm_NoMatch() {
	// Arrange
	operatorID := uint(2)
	otherOperatorID := uint(3)
	form := &entity.Form{}
	form.ID = 1
	form.UserID = 1
	form.OperatorID = &otherOperatorID

	// Act
	canEdit, err := suite.formService.canOperatorEditForm(form, operatorID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), canEdit)
}

// Test: UpsertGeneralHealth should create when not exists
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_Create() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertGeneralHealthRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID
	form.Status = enum.FormStatusPending

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateGeneralHealth", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: UpsertGeneralHealth should return error when form not found
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.UpsertGeneralHealthRequest{
		UserID: userID,
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.Error(suite.T(), err)
	notFoundError, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.constants.Field.Form, notFoundError.Item)
}

// Test: UpsertGeneralHealth should return forbidden when operator can't edit
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_OperatorForbidden() {
	// Arrange
	userID := uint(1)
	operatorID := uint(2)
	formID := uint(1)
	request := formdto.UpsertGeneralHealthRequest{
		UserID: operatorID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID
	form.OperatorID = nil // Different operator

	operatorRole := userdto.RoleResponse{
		ID:   1,
		Name: enum.Operator.String(),
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", operatorID).Return([]userdto.RoleResponse{operatorRole}, nil)
	// canOperatorEditForm will return false because form.OperatorID is nil

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.Error(suite.T(), err)
	forbiddenError, ok := err.(exception.ForbiddenError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.constants.Field.Form, forbiddenError.Resource)
}

// Test: DeleteForm should successfully delete form
func (suite *FormServiceTestSuite) TestDeleteForm_Success() {
	// Arrange
	formID := uint(1)
	request := formdto.DeleteFormRequest{
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = 1

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("DeleteForm", suite.db, formID).Return(nil)

	// Act
	err := suite.formService.DeleteForm(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: DeleteForm should return error when form not found
func (suite *FormServiceTestSuite) TestDeleteForm_FormNotFound() {
	// Arrange
	formID := uint(999)
	request := formdto.DeleteFormRequest{
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.DeleteForm(request)

	// Assert
	assert.Error(suite.T(), err)
	_, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
}

// Test: GetBasicForm should return basic form info
func (suite *FormServiceTestSuite) TestGetBasicForm_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	basicInfo := &entity.BasicInfo{}
	basicInfo.ID = 1
	basicInfo.FormID = formID
	basicInfo.Gender = enum.Gender(1)

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindBasicInfoByFormID", suite.db, formID).Return(basicInfo, nil)

	// Act
	response, err := suite.formService.GetBasicForm(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), basicInfo.ID, response.ID)
	assert.Equal(suite.T(), basicInfo.Gender, response.Gender)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: GetBasicForm should return error when form not found
func (suite *FormServiceTestSuite) TestGetBasicForm_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.GetBasicForm(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.GetBasicFormResponse{}, response)
}

// Test: UpsertGeneralHealth should update when exists
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_Update() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertGeneralHealthRequest{
		UserID:                    userID,
		FormID:                    formID,
		DrinksAlcohol:             boolPtr(true),
		CupsPerWeek:               strPtr("5"),
		LastMonthSabzijatMeal:     "high",
		LastMonthSabzijatWeight:   "2kg",
		MediumActivityMonthInYear: 3,
		MediumActivityHourInWeek:  "3",
		HardActivityMonthInYear:   1,
		HardActivityHourInWeek:    "1",
		SmokingNow:                true,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID
	form.Status = enum.FormStatusPending

	existingHealth := &entity.GeneralHealthInfo{}
	existingHealth.ID = 5
	existingHealth.FormID = formID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(existingHealth, nil)
	suite.formRepository.On("UpdateGeneralHealth", suite.db, mock.MatchedBy(func(g *entity.GeneralHealthInfo) bool {
		return g.ID == 5 && g.SmokingNow == true
	})).Return(nil)

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: UpsertGeneralHealth should log action when operator performs update
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_OperatorLogsAction() {
	// Arrange
	userID := uint(1)
	operatorID := uint(2)
	formID := uint(1)
	request := formdto.UpsertGeneralHealthRequest{
		UserID: operatorID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID
	form.OperatorID = &operatorID

	operatorRole := userdto.RoleResponse{
		ID:   1,
		Name: enum.Operator.String(),
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", operatorID).Return([]userdto.RoleResponse{operatorRole}, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateGeneralHealth", suite.db, mock.Anything).Return(nil)
	suite.actionLogService.On("LogAction", mock.MatchedBy(func(log actionlogdto.LogAction) bool {
		return log.ActorID == operatorID && log.ResourceID != nil && *log.ResourceID == formID
	})).Return(nil)

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.actionLogService.AssertExpectations(suite.T())
}

// Test: ChangeFormStatus should update form status
func (suite *FormServiceTestSuite) TestChangeFormStatus_Success() {
	// Arrange
	formID := uint(1)
	request := formdto.ChangeFormStatusRequest{
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = 1
	form.Status = enum.FormStatusPending
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.Status == enum.FormStatusReady
	})).Return(nil)

	// Act
	response, err := suite.formService.ChangeFormStatus(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), enum.FormStatusReady.String(), response.Form.Status)
	assert.Equal(suite.T(), formID, response.Form.FormID)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: ChangeFormStatus should return error when form not found
func (suite *FormServiceTestSuite) TestChangeFormStatus_FormNotFound() {
	// Arrange
	formID := uint(999)
	request := formdto.ChangeFormStatusRequest{
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.ChangeFormStatus(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.ChangeFormStatusResponse{}, response)
	notFoundErr, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.constants.Field.Form, notFoundErr.Item)
}

// Test: ChangeFormStatus should return error when update fails
func (suite *FormServiceTestSuite) TestChangeFormStatus_UpdateFails() {
	// Arrange
	formID := uint(1)
	request := formdto.ChangeFormStatusRequest{
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = 1
	form.Status = enum.FormStatusPending

	updateError := errors.New("database update error")

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.Anything).Return(updateError)

	// Act
	response, err := suite.formService.ChangeFormStatus(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), updateError, err)
	assert.Equal(suite.T(), formdto.ChangeFormStatusResponse{}, response)
}

// Test: GetBasicForm should return error when basic info not found
func (suite *FormServiceTestSuite) TestGetBasicForm_BasicInfoNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindBasicInfoByFormID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.GetBasicForm(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.GetBasicFormResponse{}, response)
	_, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
}

// Test: GetGeneralHealth should return health info
func (suite *FormServiceTestSuite) TestGetGeneralHealth_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	healthInfo := &entity.GeneralHealthInfo{}
	healthInfo.ID = 5
	healthInfo.FormID = formID
	healthInfo.DrinksAlcohol = boolPtr(true)
	healthInfo.SmokingNow = true

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(healthInfo, nil)

	// Act
	response, err := suite.formService.GetGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), healthInfo.ID, response.ID)
	assert.Equal(suite.T(), *healthInfo.DrinksAlcohol, *response.DrinksAlcohol)
	assert.Equal(suite.T(), healthInfo.SmokingNow, response.SmokingNow)
}

// Test: GetGeneralHealth should return error when form not found
func (suite *FormServiceTestSuite) TestGetGeneralHealth_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.GetGeneralHealth(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.GetGeneralHealthResponse{}, response)
}

// Test: GetGeneralHealth should return error when health info not found
func (suite *FormServiceTestSuite) TestGetGeneralHealth_InfoNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.GetGeneralHealth(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.GetGeneralHealthResponse{}, response)
}

// Test: GetMamography should return mamography info
func (suite *FormServiceTestSuite) TestGetMamography_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	mamoInfo := &entity.MamoGraphyInfo{}
	mamoInfo.ID = 10
	mamoInfo.FormID = formID
	mamoInfo.GhaedeAge = 45
	mamoInfo.HasChildren = true

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindMamographyByFormID", suite.db, formID).Return(mamoInfo, nil)

	// Act
	response, err := suite.formService.GetMamography(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mamoInfo.ID, response.ID)
	assert.Equal(suite.T(), mamoInfo.GhaedeAge, response.GhaedeAge)
	assert.Equal(suite.T(), mamoInfo.HasChildren, response.HasChildren)
}

// Test: GetMamography should return error when form not found
func (suite *FormServiceTestSuite) TestGetMamography_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.GetMamography(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.GetMamographyResponse{}, response)
}

// Test: hasPermission should return true when user has permission
func (suite *FormServiceTestSuite) TestHasPermission_True() {
	// Arrange
	userID := uint(1)
	permission := enum.PermissionCreateFormForUser

	role := userdto.RoleResponse{
		ID:   1,
		Name: enum.Operator.String(),
		Permissions: []userdto.PermissionResponse{
			{
				Name: permission.String(),
			},
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{role}, nil)

	// Act
	hasPerm, err := suite.formService.hasPermission(userID, permission)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPerm)
}

// Test: hasPermission should return true when user has PermissionAll
func (suite *FormServiceTestSuite) TestHasPermission_WithPermissionAll() {
	// Arrange
	userID := uint(1)
	permission := enum.PermissionCreateFormForUser

	role := userdto.RoleResponse{
		ID:   1,
		Name: enum.SuperAdmin.String(),
		Permissions: []userdto.PermissionResponse{
			{
				Name: enum.PermissionAll.String(),
			},
		},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{role}, nil)

	// Act
	hasPerm, err := suite.formService.hasPermission(userID, permission)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPerm)
}

// Test: hasPermission should return false when user lacks permission
func (suite *FormServiceTestSuite) TestHasPermission_False() {
	// Arrange
	userID := uint(1)
	permission := enum.PermissionCreateFormForUser

	role := userdto.RoleResponse{
		ID:          1,
		Name:        enum.Patient.String(),
		Permissions: []userdto.PermissionResponse{},
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{role}, nil)

	// Act
	hasPerm, err := suite.formService.hasPermission(userID, permission)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPerm)
}

// Test: hasPermission should return error when service fails
func (suite *FormServiceTestSuite) TestHasPermission_ServiceError() {
	// Arrange
	userID := uint(1)
	permission := enum.PermissionCreateFormForUser
	serviceError := errors.New("database error")

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return(nil, serviceError)

	// Act
	hasPerm, err := suite.formService.hasPermission(userID, permission)

	// Assert
	assert.Error(suite.T(), err)
	assert.False(suite.T(), hasPerm)
	assert.Equal(suite.T(), serviceError, err)
}

// Test: UpsertMamography should create when not exists
func (suite *FormServiceTestSuite) TestUpsertMamography_Create() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertMamographyRequest{
		UserID:           userID,
		FormID:           formID,
		GhaedeAge:        45,
		HasChildren:      true,
		MenopausalStatus: uint(enum.MenopausalStatusPreMenopausal),
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindMamographyByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateMamography", suite.db, mock.MatchedBy(func(m *entity.MamoGraphyInfo) bool {
		return m.FormID == formID && m.GhaedeAge == 45
	})).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpsertMamography(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: CreateCancer should create cancer info
func (suite *FormServiceTestSuite) TestCreateCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.CreateCancerRequest{
		UserID:     userID,
		FormID:     formID,
		CancerAge:  35,
		CancerType: uint(enum.CancerTypeBreast),
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindCancersByFormID", suite.db, formID).Return([]*entity.CancerInfo{}, nil)
	suite.formRepository.On("CreateCancer", suite.db, mock.MatchedBy(func(c *entity.CancerInfo) bool {
		return c.FormID == formID && c.CancerAge == 35
	})).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.CreateCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.formRepository.AssertExpectations(suite.T())
}

// Test: CreateCancer should return error when form not found
func (suite *FormServiceTestSuite) TestCreateCancer_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.CreateCancerRequest{
		UserID:     userID,
		FormID:     formID,
		CancerAge:  35,
		CancerType: uint(enum.CancerTypeBreast),
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.CreateCancer(request)

	// Assert
	assert.Error(suite.T(), err)
	notFoundErr, ok := err.(exception.NotFoundError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.constants.Field.Form, notFoundErr.Item)
}

// Test: UpdateCancer should update existing cancer
func (suite *FormServiceTestSuite) TestUpdateCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	cancerID := uint(5)
	request := formdto.UpdateCancerRequest{
		CancerID:   cancerID,
		FormID:     formID,
		UserID:     userID,
		CancerAge:  40,
		CancerType: uint(enum.CancerTypeBreast),
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingCancer := &entity.CancerInfo{}
	existingCancer.ID = cancerID
	existingCancer.FormID = formID
	existingCancer.CancerAge = 45
	existingCancer.CancerType = enum.CancerTypeOvarian

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindCancerByID", suite.db, cancerID).Return(existingCancer, nil)
	suite.formRepository.On("FindCancersByFormID", suite.db, formID).Return([]*entity.CancerInfo{existingCancer}, nil)
	suite.formRepository.On("UpdateCancer", suite.db, mock.MatchedBy(func(c *entity.CancerInfo) bool {
		return c.ID == cancerID && c.CancerAge == 40
	})).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	response, err := suite.formService.UpdateCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cancerID, response.Cancer.ID)
}

// Test: UpdateCancer should return error when form not found
func (suite *FormServiceTestSuite) TestUpdateCancer_FormNotFound() {
	// Arrange
	formID := uint(999)
	request := formdto.UpdateCancerRequest{
		CancerID: 1,
		FormID:   formID,
		UserID:   1,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	response, err := suite.formService.UpdateCancer(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), formdto.UpdateCancerResponse{}, response)
}

// Test: DeleteCancer should delete cancer record
func (suite *FormServiceTestSuite) TestDeleteCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	cancerID := uint(5)
	request := formdto.DeleteCancerRequest{
		CancerID: cancerID,
		FormID:   formID,
		UserID:   userID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	cancer := &entity.CancerInfo{}
	cancer.ID = cancerID
	cancer.FormID = formID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindCancerByID", suite.db, cancerID).Return(cancer, nil)
	suite.formRepository.On("DeleteCancerByID", suite.db, cancerID).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.DeleteCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: CreateFamilyCancer should create family cancer record
func (suite *FormServiceTestSuite) TestCreateFamilyCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.CreateFamilyCancerRequest{
		UserID:     userID,
		FormID:     formID,
		Relative:   enum.PaternalAunt,
		CancerAge:  45,
		CancerType: uint(enum.CancerTypeBreast),
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("CreateFamilyCancer", suite.db, mock.MatchedBy(func(f *entity.FamilyCancerInfo) bool {
		return f.FormID == formID && f.CancerAge == 45
	})).Run(func(args mock.Arguments) {
		f := args.Get(1).(*entity.FamilyCancerInfo)
		f.ID = 5
	}).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	response, err := suite.formService.CreateFamilyCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), response.FamilyCancer.ID)
}

// Test: UpdateFamilyCancer should update family cancer
func (suite *FormServiceTestSuite) TestUpdateFamilyCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	familyCancerID := uint(10)
	request := formdto.UpdateFamilyCancerRequest{
		FamilyCancerID: familyCancerID,
		FormID:         formID,
		UserID:         userID,
		CancerAge:      50,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingFC := &entity.FamilyCancerInfo{}
	existingFC.ID = familyCancerID
	existingFC.FormID = formID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFamilyCancerByID", suite.db, familyCancerID).Return(existingFC, nil)
	suite.formRepository.On("UpdateFamilyCancer", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	response, err := suite.formService.UpdateFamilyCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), familyCancerID, response.FamilyCancer.ID)
}

// Test: DeleteFamilyCancer should delete family cancer
func (suite *FormServiceTestSuite) TestDeleteFamilyCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	familyCancerID := uint(10)
	request := formdto.DeleteFamilyCancerRequest{
		FamilyCancerID: familyCancerID,
		FormID:         formID,
		UserID:         userID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	fc := &entity.FamilyCancerInfo{}
	fc.ID = familyCancerID
	fc.FormID = formID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFamilyCancerByID", suite.db, familyCancerID).Return(fc, nil)
	suite.formRepository.On("DeleteFamilyCancerByID", suite.db, familyCancerID).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.DeleteFamilyCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: GetCancers should return list of cancers
func (suite *FormServiceTestSuite) TestGetCancers_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	cancers := []*entity.CancerInfo{
		{
			FormID:    formID,
			CancerAge: 35,
		},
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindCancersByFormID", suite.db, formID).Return(cancers, nil)

	// Act
	response, err := suite.formService.GetCancers(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.Cancers, 1)
}

// Test: GetFamilyCancer should return family cancers
func (suite *FormServiceTestSuite) TestGetFamilyCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	familyCancers := []*entity.FamilyCancerInfo{
		{
			FormID:    formID,
			CancerAge: 50,
		},
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindFamilyCancersByFormID", suite.db, formID).Return(familyCancers, nil)

	// Act
	response, err := suite.formService.GetFamilyCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.FamilyCancers, 1)
}

// Test: GetContact should return contact info
func (suite *FormServiceTestSuite) TestGetContact_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	contact := &entity.ContactInfo{
		FormID: formID,
		Name:   "John Doe",
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindContactByFormID", suite.db, formID).Return(contact, nil)

	// Act
	response, err := suite.formService.GetContact(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), contact.Name, response.Name)
}

// Test: GetLungCancer should return lung cancer info
func (suite *FormServiceTestSuite) TestGetLungCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.GetPartialFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	lungCancer := &entity.LungCancerInfo{
		FormID:         formID,
		CurrentSmoking: true,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindLungCancerByFormID", suite.db, formID).Return(lungCancer, nil)

	// Act
	response, err := suite.formService.GetLungCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), lungCancer.CurrentSmoking, response.CurrentSmoking)
}

// Test: GetUserForms should return user forms
func (suite *FormServiceTestSuite) TestGetUserForms_Success() {
	// Arrange
	userID := uint(1)
	request := formdto.GetUserFormsRequest{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	user := &entity.User{}
	user.ID = userID

	forms := []*entity.Form{
		{
			UserID: userID,
			Status: enum.FormStatusPending,
		},
	}

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(user, nil)
	suite.formRepository.On("FindFormsByUserID", suite.db, userID, mock.Anything).Return(forms, nil)
	suite.formRepository.On("CountFormsByUserID", suite.db, userID).Return(int64(1), nil)

	// Act
	responses, count, err := suite.formService.GetUserForms(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)
	assert.Len(suite.T(), responses, 1)
}

// Test: GetUserForms should return error when user not found
func (suite *FormServiceTestSuite) TestGetUserForms_UserNotFound() {
	// Arrange
	userID := uint(999)
	request := formdto.GetUserFormsRequest{
		UserID: userID,
	}

	// Setup expectations
	suite.userService.On("GetUserByID", userID).Return(nil, nil)

	// Act
	responses, count, err := suite.formService.GetUserForms(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), responses)
	assert.Equal(suite.T(), int64(0), count)
}

// Test: UpdateBasicInfo should update basic form info
func (suite *FormServiceTestSuite) TestUpdateBasicInfo_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	genderVal := uint(1)
	request := formdto.UpdateBasicFormRequest{
		UserID: userID,
		FormID: formID,
		Gender: &genderVal,
	}

	form := &entity.Form{}
	form.ID = formID

	basicInfo := &entity.BasicInfo{}
	basicInfo.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindBasicInfoByFormID", suite.db, formID).Return(basicInfo, nil)
	suite.formRepository.On("UpdateBasicInfo", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.formService.UpdateBasicInfo(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateBasicInfo should return error when form not found
func (suite *FormServiceTestSuite) TestUpdateBasicInfo_FormNotFound() {
	// Arrange
	formID := uint(999)
	request := formdto.UpdateBasicFormRequest{
		FormID: formID,
	}

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.UpdateBasicInfo(request)

	// Assert
	assert.Error(suite.T(), err)
}

// Test: UpsertContact should create contact info
func (suite *FormServiceTestSuite) TestUpsertContact_Create() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertContactRequest{
		UserID:  userID,
		FormID:  formID,
		Name:    "John Doe",
		Address: "123 Main St",
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindContactByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateContact", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.formService.UpsertContact(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpsertLungCancer should create lung cancer info
func (suite *FormServiceTestSuite) TestUpsertLungCancer_Create() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertLungCancerRequest{
		UserID:         userID,
		FormID:         formID,
		CurrentSmoking: true,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindLungCancerByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateLungCancer", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.formService.UpsertLungCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: GetAllForms should return forms with filters
func (suite *FormServiceTestSuite) TestGetAllForms_Success() {
	// Arrange
	form := &entity.Form{}
	form.ID = 1
	form.UserID = 1
	form.Status = enum.FormStatusPending

	forms := []*entity.Form{form}
	filters := &postgres.FormFilters{}

	// Setup expectations
	suite.formRepository.On("FindAllForms", suite.db, 0, 10, filters).Return(forms, nil)
	suite.formRepository.On("CountAllForms", suite.db, filters).Return(int64(1), nil)

	// Act
	responses, count, err := suite.formService.GetAllForms(0, 10, filters)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)
	assert.Len(suite.T(), responses, 1)
}

// Test: GetAllOperatorForms should return operator forms
func (suite *FormServiceTestSuite) TestGetAllOperatorForms_Success() {
	// Arrange
	operatorID := uint(1)
	form := &entity.Form{}
	form.ID = 1
	form.UserID = 1
	form.Status = enum.FormStatusPending

	forms := []*entity.Form{form}
	filters := &postgres.OperatorFormFilters{OperatorID: operatorID}

	operator := &entity.User{}
	operator.ID = operatorID

	// Setup expectations
	suite.userService.On("GetUserByID", operatorID).Return(operator, nil)
	suite.userService.On("GetUserRoles", operatorID).Return([]userdto.RoleResponse{
		{Name: enum.Operator.String()},
	}, nil)
	suite.formRepository.On("FindAllOperatorForms", suite.db, 0, 10, filters).Return(forms, nil)
	suite.formRepository.On("CountAllOperatorForms", suite.db, filters).Return(int64(1), nil)

	// Act
	responses, count, err := suite.formService.GetAllOperatorForms(0, 10, filters)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)
	assert.Len(suite.T(), responses, 1)
}

// Test: AcceptForm should accept form
func (suite *FormServiceTestSuite) TestAcceptForm_Success() {
	// Arrange
	formID := uint(1)
	userID := uint(1)

	form := &entity.Form{}
	form.ID = formID
	form.Status = enum.FormStatusPending

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{
		{
			Name: enum.Supervisor.String(),
			Permissions: []userdto.PermissionResponse{
				{Name: enum.PermissionHandleOperators.String()},
			},
		},
	}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.Status == enum.FormStatusApproved
	})).Return(nil)

	// Act
	err := suite.formService.AcceptForm(formID, userID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: RejectForm should reject form
func (suite *FormServiceTestSuite) TestRejectForm_Success() {
	// Arrange
	formID := uint(1)
	userID := uint(1)

	form := &entity.Form{}
	form.ID = formID
	form.Status = enum.FormStatusPending

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{
		{
			Name: enum.Supervisor.String(),
			Permissions: []userdto.PermissionResponse{
				{Name: enum.PermissionHandleOperators.String()},
			},
		},
	}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.Status == enum.FormStatusRejected
	})).Return(nil)

	// Act
	err := suite.formService.RejectForm(formID, userID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: AssignOperator should assign operator to form
func (suite *FormServiceTestSuite) TestAssignOperator_Success() {
	// Arrange
	formID := uint(1)
	operatorID := uint(2)
	request := formdto.AssignOperatorRequest{
		FormID:     formID,
		OperatorID: operatorID,
	}

	form := &entity.Form{}
	form.ID = formID

	operator := &entity.User{}
	operator.ID = operatorID

	// Setup expectations
	suite.userService.On("GetUserByID", operatorID).Return(operator, nil)
	suite.userService.On("GetUserRoles", operatorID).Return([]userdto.RoleResponse{
		{Name: enum.Operator.String()},
	}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.OperatorID != nil && *f.OperatorID == operatorID
	})).Return(nil)
	suite.actionLogService.On("LogAction", mock.MatchedBy(func(log actionlogdto.LogAction) bool {
		return log.TargetID != nil && *log.TargetID == operatorID
	})).Return(nil)

	// Act
	err := suite.formService.AssignOperator(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UnassignOperator should unassign operator from form
func (suite *FormServiceTestSuite) TestUnassignOperator_Success() {
	// Arrange
	formID := uint(1)
	userID := uint(1)
	operatorID := uint(2)
	request := formdto.UnassignOperatorRequest{
		FormID: formID,
		UserID: userID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.OperatorID = &operatorID

	// Setup expectations
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.MatchedBy(func(f *entity.Form) bool {
		return f.OperatorID == nil
	})).Return(nil)
	suite.actionLogService.On("LogAction", mock.MatchedBy(func(log actionlogdto.LogAction) bool {
		return log.ActorID == userID
	})).Return(nil)

	// Act
	err := suite.formService.UnassignOperator(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: Enum getters should return all enum values
func (suite *FormServiceTestSuite) TestGetAllCancerTypes_Success() {
	// Act
	response, err := suite.formService.GetAllCancerTypes()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllGenders_Success() {
	// Act
	response, err := suite.formService.GetAllGenders()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllMenopausalStatuses_Success() {
	// Act
	response, err := suite.formService.GetAllMenopausalStatuses()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllFormStatuses_Success() {
	// Act
	response, err := suite.formService.GetAllFormStatuses()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllHyperplasiaInBiopsyStatuses_Success() {
	// Act
	response, err := suite.formService.GetAllHyperplasiaInBiopsyStatuses()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllLifeStatuses_Success() {
	// Act
	response, err := suite.formService.GetAllLifeStatuses()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

func (suite *FormServiceTestSuite) TestGetAllRelativeTypes_Success() {
	// Act
	response, err := suite.formService.GetAllRelativeTypes()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
}

// Test: UpdateForm should update form
func (suite *FormServiceTestSuite) TestUpdateForm_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateBasicFormRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("UpdateForm", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.formService.UpdateForm(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateGeneralHealth should update general health
func (suite *FormServiceTestSuite) TestUpdateGeneralHealth_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateGeneralHealthRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	generalHealth := &entity.GeneralHealthInfo{}
	generalHealth.ID = 1
	generalHealth.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(generalHealth, nil)
	suite.formRepository.On("UpdateGeneralHealth", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateMamography should update mamography
func (suite *FormServiceTestSuite) TestUpdateMamography_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateMamographyRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	mamography := &entity.MamoGraphyInfo{}
	mamography.ID = 1
	mamography.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindMamographyByFormID", suite.db, formID).Return(mamography, nil)
	suite.formRepository.On("UpdateMamography", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateMamography(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateContact should update contact
func (suite *FormServiceTestSuite) TestUpdateContact_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateContactRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	contact := &entity.ContactInfo{}
	contact.ID = 1
	contact.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindContactByFormID", suite.db, formID).Return(contact, nil)
	suite.formRepository.On("UpdateContact", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateContact(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateLungCancer should update lung cancer
func (suite *FormServiceTestSuite) TestUpdateLungCancer_Success() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateLungCancerRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	lungCancer := &entity.LungCancerInfo{}
	lungCancer.ID = 1
	lungCancer.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindLungCancerByFormID", suite.db, formID).Return(lungCancer, nil)
	suite.formRepository.On("UpdateLungCancer", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateLungCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpsertMamography should update when exists
func (suite *FormServiceTestSuite) TestUpsertMamography_Update() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertMamographyRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingMamography := &entity.MamoGraphyInfo{}
	existingMamography.ID = 1
	existingMamography.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindMamographyByFormID", suite.db, formID).Return(existingMamography, nil)
	suite.formRepository.On("UpdateMamography", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpsertMamography(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpsertGeneralHealth should update when exists (additional path)
func (suite *FormServiceTestSuite) TestUpsertGeneralHealth_UpdatePath() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertGeneralHealthRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingHealth := &entity.GeneralHealthInfo{}
	existingHealth.ID = 1
	existingHealth.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(existingHealth, nil)
	suite.formRepository.On("UpdateGeneralHealth", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpsertGeneralHealth(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpsertContact should update when exists
func (suite *FormServiceTestSuite) TestUpsertContact_Update() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertContactRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingContact := &entity.ContactInfo{}
	existingContact.ID = 1
	existingContact.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindContactByFormID", suite.db, formID).Return(existingContact, nil)
	suite.formRepository.On("UpdateContact", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpsertContact(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpsertLungCancer should update when exists
func (suite *FormServiceTestSuite) TestUpsertLungCancer_Update() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpsertLungCancerRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	existingLungCancer := &entity.LungCancerInfo{}
	existingLungCancer.ID = 1
	existingLungCancer.FormID = formID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindLungCancerByFormID", suite.db, formID).Return(existingLungCancer, nil)
	suite.formRepository.On("UpdateLungCancer", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpsertLungCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateForm should return error when form not found
func (suite *FormServiceTestSuite) TestUpdateForm_FormNotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(999)
	request := formdto.UpdateBasicFormRequest{
		UserID: userID,
		FormID: formID,
	}

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.UpdateForm(request)

	// Assert
	assert.Error(suite.T(), err)
}

// Test: UpdateGeneralHealth should return error when general health not found
func (suite *FormServiceTestSuite) TestUpdateGeneralHealth_NotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateGeneralHealthRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindGeneralHealthByFormID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.UpdateGeneralHealth(request)

	// Assert
	assert.Error(suite.T(), err)
}

// Test: UpdateMamography should return error when mamography not found
func (suite *FormServiceTestSuite) TestUpdateMamography_NotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateMamographyRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindMamographyByFormID", suite.db, formID).Return(nil, nil)

	// Act
	err := suite.formService.UpdateMamography(request)

	// Assert
	assert.Error(suite.T(), err)
}

// Test: UpdateContact should create new contact when not found
func (suite *FormServiceTestSuite) TestUpdateContact_NotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateContactRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindContactByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateContact", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateContact(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: UpdateLungCancer should create new lung cancer when not found
func (suite *FormServiceTestSuite) TestUpdateLungCancer_NotFound() {
	// Arrange
	userID := uint(1)
	formID := uint(1)
	request := formdto.UpdateLungCancerRequest{
		UserID: userID,
		FormID: formID,
	}

	form := &entity.Form{}
	form.ID = formID
	form.UserID = userID

	// Setup expectations
	suite.userService.On("GetUserRoles", userID).Return([]userdto.RoleResponse{}, nil)
	suite.formRepository.On("FindFormByID", suite.db, formID).Return(form, nil)
	suite.formRepository.On("FindLungCancerByFormID", suite.db, formID).Return(nil, nil)
	suite.formRepository.On("CreateLungCancer", suite.db, mock.Anything).Return(nil)
	suite.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(database.Database) error)
		fn(suite.db)
	}).Return(nil)

	// Act
	err := suite.formService.UpdateLungCancer(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Helper functions for pointers
func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func TestFormServiceTestSuite(t *testing.T) {
	suite.Run(t, new(FormServiceTestSuite))
}
