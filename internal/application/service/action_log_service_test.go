package service

import (
	"errors"
	"testing"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	databaseMocks "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/database"
	repositoryMocks "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ActionLogServiceTestSuite struct {
	suite.Suite
	constants           *bootstrap.Constants
	db                  *databaseMocks.DatabaseMock
	actionLogRepository *repositoryMocks.ActionLogRepositoryMock
	actionLogService    *ActionLogService
}

func (suite *ActionLogServiceTestSuite) SetupTest() {
	suite.constants = bootstrap.NewConstants()
	suite.db = databaseMocks.NewDatabaseMock()
	suite.actionLogRepository = repositoryMocks.NewActionLogRepositoryMock()

	suite.actionLogService = NewActionLogService(
		suite.constants,
		suite.actionLogRepository,
		suite.db,
	)
}

// Test: LogAction should successfully create an action log
func (suite *ActionLogServiceTestSuite) TestLogAction_Success() {
	// Arrange
	actorID := uint(1)
	resource := "form"
	resourceID := uint(100)
	request := actionlogdto.LogAction{
		ActorID:    actorID,
		TargetID:   nil,
		Action:     enum.ActionTypeOperatorCreatedForm,
		Resource:   &resource,
		ResourceID: &resourceID,
		Details:    "Created a new form",
	}

	suite.actionLogRepository.On("CreateActionLog", suite.db, mock.MatchedBy(func(actionLog *entity.ActionLog) bool {
		return actionLog.ActorID == actorID &&
			actionLog.Action == enum.ActionTypeOperatorCreatedForm &&
			actionLog.Details == "Created a new form"
	})).Return(nil)

	// Act
	err := suite.actionLogService.LogAction(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.actionLogRepository.AssertCalled(suite.T(), "CreateActionLog", suite.db, mock.Anything)
}

// Test: LogAction should handle repository errors
func (suite *ActionLogServiceTestSuite) TestLogAction_RepositoryError() {
	// Arrange
	request := actionlogdto.LogAction{
		ActorID: 1,
		Action:  enum.ActionTypeOperatorCreatedForm,
		Details: "Test action",
	}

	repositoryError := errors.New("database connection error")
	suite.actionLogRepository.On("CreateActionLog", suite.db, mock.Anything).Return(repositoryError)

	// Act
	err := suite.actionLogService.LogAction(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
}

// Test: LogAction should handle target user ID correctly
func (suite *ActionLogServiceTestSuite) TestLogAction_WithTargetID() {
	// Arrange
	actorID := uint(1)
	targetID := uint(2)
	request := actionlogdto.LogAction{
		ActorID:  actorID,
		TargetID: &targetID,
		Action:   enum.ActionTypeUserRoleUpdated,
		Details:  "Updated user",
	}

	suite.actionLogRepository.On("CreateActionLog", suite.db, mock.MatchedBy(func(actionLog *entity.ActionLog) bool {
		return actionLog.ActorID == actorID &&
			actionLog.TargetID != nil &&
			*actionLog.TargetID == targetID
	})).Return(nil)

	// Act
	err := suite.actionLogService.LogAction(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: GetAllActionLogs should return paginated results
func (suite *ActionLogServiceTestSuite) TestGetAllActionLogs_Success() {
	// Arrange
	offset := 0
	limit := 10
	now := time.Now()
	actionLogs := []*entity.ActionLog{
		{
			Model: database.Model{ID: 1, CreatedAt: now},
			ActorID: 1,
			Action:  enum.ActionTypeOperatorCreatedForm,
			Details: "Created form",
		},
		{
			Model: database.Model{ID: 2, CreatedAt: now},
			ActorID: 2,
			Action:  enum.ActionTypeOperatorUpdatedForm,
			Details: "Updated form",
		},
	}

	suite.actionLogRepository.On("FindAllActionLogs", suite.db, mock.MatchedBy(func(opts interface{}) bool {
		return true
	})).Return(actionLogs, nil)

	suite.actionLogRepository.On("CountAllActionLogs", suite.db, mock.MatchedBy(func(opts interface{}) bool {
		return true
	})).Return(int64(2), nil)

	// Act
	response, count, err := suite.actionLogService.GetAllActionLogs(offset, limit)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
	assert.Len(suite.T(), response, 2)
	assert.Equal(suite.T(), enum.ActionTypeOperatorCreatedForm.String(), response[0].Action)
	assert.Equal(suite.T(), enum.ActionTypeOperatorUpdatedForm.String(), response[1].Action)
}

// Test: GetAllActionLogs should handle repository errors
func (suite *ActionLogServiceTestSuite) TestGetAllActionLogs_FindError() {
	// Arrange
	offset := 0
	limit := 10
	findError := errors.New("database error")

	suite.actionLogRepository.On("FindAllActionLogs", suite.db, mock.Anything).Return(nil, findError)

	// Act
	response, count, err := suite.actionLogService.GetAllActionLogs(offset, limit)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), int64(0), count)
	assert.Equal(suite.T(), findError, err)
}

// Test: GetAllActionLogs should handle count errors
func (suite *ActionLogServiceTestSuite) TestGetAllActionLogs_CountError() {
	// Arrange
	offset := 0
	limit := 10
	now := time.Now()
	actionLogs := []*entity.ActionLog{
		{
			Model: database.Model{ID: 1, CreatedAt: now},
			ActorID: 1,
			Action:  enum.ActionTypeOperatorCreatedForm,
		},
	}
	countError := errors.New("count error")

	suite.actionLogRepository.On("FindAllActionLogs", suite.db, mock.Anything).Return(actionLogs, nil)
	suite.actionLogRepository.On("CountAllActionLogs", suite.db, mock.Anything).Return(int64(0), countError)

	// Act
	response, count, err := suite.actionLogService.GetAllActionLogs(offset, limit)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), int64(0), count)
	assert.Equal(suite.T(), countError, err)
}

// Test: GetAllActionLogs should return empty list when no logs exist
func (suite *ActionLogServiceTestSuite) TestGetAllActionLogs_Empty() {
	// Arrange
	offset := 0
	limit := 10
	emptyLogs := []*entity.ActionLog{}

	suite.actionLogRepository.On("FindAllActionLogs", suite.db, mock.Anything).Return(emptyLogs, nil)
	suite.actionLogRepository.On("CountAllActionLogs", suite.db, mock.Anything).Return(int64(0), nil)

	// Act
	response, count, err := suite.actionLogService.GetAllActionLogs(offset, limit)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
	assert.Len(suite.T(), response, 0)
}

// Test: GetAllActionTypes should return all action types
func (suite *ActionLogServiceTestSuite) TestGetAllActionTypes_Success() {
	// Act
	response, err := suite.actionLogService.GetAllActionTypes()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Greater(suite.T(), len(response), 0, "Should return at least one action type")

	// Verify response structure
	for _, actionType := range response {
		assert.Greater(suite.T(), actionType.ID, uint(0), "Action type ID should be greater than 0")
		assert.NotEmpty(suite.T(), actionType.Name, "Action type name should not be empty")
	}
}

func TestActionLogServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ActionLogServiceTestSuite))
}
