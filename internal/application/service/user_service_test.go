package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	postgresRepository "github.com/FamCan-RiskAssessment/Backend/mocks/domain/repository/postgres"
	redisRepository "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/repository/redis"
	communicationMocks "github.com/FamCan-RiskAssessment/Backend/mocks/domain/communication"
	usecaseMocks "github.com/FamCan-RiskAssessment/Backend/mocks/application/usecase"
	databaseMocks "github.com/FamCan-RiskAssessment/Backend/mocks/infrastructure/database"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	userService         *UserService
	userRepository      *postgresRepository.UserRepositoryMock
	userCacheRepository *redisRepository.UserCacheRepositoryMock
	jwtService          *usecaseMocks.JwtServiceMock
	smsService          *communicationMocks.SmsServiceMock
	otpService          *usecaseMocks.OtpServiceMock
	actionLogService    *usecaseMocks.ActionLogServiceMock
	db                  *databaseMocks.DatabaseMock
	constants           *bootstrap.Constants
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepository = postgresRepository.NewUserRepositoryMock()
	suite.userCacheRepository = redisRepository.NewUserCacheRepositoryMock()
	suite.jwtService = usecaseMocks.NewJwtServiceMock()
	suite.smsService = communicationMocks.NewSmsServiceMock()
	suite.otpService = usecaseMocks.NewOtpServiceMock()
	suite.actionLogService = usecaseMocks.NewActionLogServiceMock()
	suite.db = databaseMocks.NewDatabaseMock()

	// Create a simple constants structure for testing
	suite.constants = &bootstrap.Constants{
		Field: bootstrap.Field{
			User: "کاربر",
			Role: "نقش",
		},
		RedisKey: bootstrap.RedisKey{},
		Tag:      bootstrap.Tag{},
	}

	suite.userService = NewUserService(
		suite.constants,
		suite.userRepository,
		suite.userCacheRepository,
		suite.jwtService,
		suite.smsService,
		suite.otpService,
		suite.actionLogService,
		suite.db,
	)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// ==================== GetUserByID Tests ====================

// Test: GetUserByID should return user when found
func (suite *UserServiceTestSuite) TestGetUserByID_Success() {
	// Arrange
	userID := uint(1)
	expectedUser := &entity.User{
		Phone: "09123456789",
	}
	expectedUser.ID = userID

	suite.userRepository.On("FindUserByID", suite.db, userID).Return(expectedUser, nil)

	// Act
	user, err := suite.userService.GetUserByID(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser, user)
	suite.userRepository.AssertCalled(suite.T(), "FindUserByID", suite.db, userID)
}

// Test: GetUserByID should return error when user not found
func (suite *UserServiceTestSuite) TestGetUserByID_NotFound() {
	// Arrange
	userID := uint(999)
	suite.userRepository.On("FindUserByID", suite.db, userID).Return(nil, nil)

	// Act
	user, err := suite.userService.GetUserByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

// Test: GetUserByID should return error from repository
func (suite *UserServiceTestSuite) TestGetUserByID_RepositoryError() {
	// Arrange
	userID := uint(1)
	expectedErr := errors.New("database error")
	suite.userRepository.On("FindUserByID", suite.db, userID).Return(nil, expectedErr)

	// Act
	user, err := suite.userService.GetUserByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
	assert.Nil(suite.T(), user)
}

// ==================== Login Tests ====================

// Test: Login should create new user and send OTP when user doesn't exist
func (suite *UserServiceTestSuite) TestLogin_NewUser_Success() {
	// Arrange
	phone := "09123456789"
	request := userdto.LoginRequest{Phone: phone}
	otp := "123456"
	expireMinute := 5

	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(nil, nil)
	suite.userRepository.On("CreateUser", suite.db, mock.Anything).Return(nil)

	patientRole := &entity.Role{
		Name: enum.Patient.String(),
	}
	patientRole.ID = 1
	suite.userRepository.On("FindRoleByName", suite.db, enum.Patient.String()).Return(patientRole, nil)
	suite.userRepository.On("AssignRoleToUser", suite.db, mock.Anything, mock.Anything).Return(nil)
	suite.otpService.On("GenerateOTP", phone).Return(otp, expireMinute, nil)
	suite.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, time.Duration(expireMinute)*time.Minute).Return(nil)
	suite.smsService.On("SendOTP", phone, otp).Return(nil)

	// Act
	err := suite.userService.Login(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.userRepository.AssertCalled(suite.T(), "CreateUser", suite.db, mock.Anything)
	suite.smsService.AssertCalled(suite.T(), "SendOTP", phone, otp)
}

// Test: Login should send OTP when user exists
func (suite *UserServiceTestSuite) TestLogin_ExistingUser_Success() {
	// Arrange
	phone := "09123456789"
	request := userdto.LoginRequest{Phone: phone}
	otp := "123456"
	expireMinute := 5

	existingUser := &entity.User{
		Phone: phone,
	}
	existingUser.ID = 1

	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(existingUser, nil)
	suite.otpService.On("GenerateOTP", phone).Return(otp, expireMinute, nil)
	suite.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, time.Duration(expireMinute)*time.Minute).Return(nil)
	suite.smsService.On("SendOTP", phone, otp).Return(nil)

	// Act
	err := suite.userService.Login(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.userRepository.AssertNotCalled(suite.T(), "CreateUser", suite.db, mock.Anything)
	suite.smsService.AssertCalled(suite.T(), "SendOTP", phone, otp)
}

// Test: Login should return error when OTP generation fails
func (suite *UserServiceTestSuite) TestLogin_OTPGenerationError() {
	// Arrange
	phone := "09123456789"
	request := userdto.LoginRequest{Phone: phone}
	expectedErr := errors.New("otp generation failed")

	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(nil, nil)
	suite.userRepository.On("CreateUser", suite.db, mock.Anything).Return(nil)
	suite.userRepository.On("FindRoleByName", suite.db, enum.Patient.String()).Return(&entity.Role{}, nil)
	suite.userRepository.On("AssignRoleToUser", suite.db, mock.Anything, mock.Anything).Return(nil)
	suite.otpService.On("GenerateOTP", phone).Return("", 0, expectedErr)

	// Act
	err := suite.userService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

// Test: Login should return error when SMS sending fails
func (suite *UserServiceTestSuite) TestLogin_SMSSendError() {
	// Arrange
	phone := "09123456789"
	request := userdto.LoginRequest{Phone: phone}
	otp := "123456"
	expireMinute := 5
	expectedErr := errors.New("sms send failed")

	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(nil, nil)
	suite.userRepository.On("CreateUser", suite.db, mock.Anything).Return(nil)
	suite.userRepository.On("FindRoleByName", suite.db, enum.Patient.String()).Return(&entity.Role{}, nil)
	suite.userRepository.On("AssignRoleToUser", suite.db, mock.Anything, mock.Anything).Return(nil)
	suite.otpService.On("GenerateOTP", phone).Return(otp, expireMinute, nil)
	suite.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, time.Duration(expireMinute)*time.Minute).Return(nil)
	suite.smsService.On("SendOTP", phone, otp).Return(expectedErr)

	// Act
	err := suite.userService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

// ==================== VerifyOTP Tests ====================

// Test: VerifyOTP should return error when OTP verification fails
func (suite *UserServiceTestSuite) TestVerifyOTP_InvalidOTP() {
	// Arrange
	phone := "09123456789"
	otp := "000000"
	request := userdto.VerifyOTPRequest{Phone: phone, OTP: otp}
	expectedErr := errors.New("invalid otp")

	suite.otpService.On("VerifyOTP", mock.Anything, otp).Return(expectedErr)

	// Act
	response, err := suite.userService.VerifyOTP(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
	assert.Equal(suite.T(), userdto.LoginResponse{}, response)
}

// Test: VerifyOTP should return error when user not found after OTP verification
func (suite *UserServiceTestSuite) TestVerifyOTP_UserNotFound() {
	// Arrange
	phone := "09123456789"
	otp := "123456"
	request := userdto.VerifyOTPRequest{Phone: phone, OTP: otp}

	suite.otpService.On("VerifyOTP", mock.Anything, otp).Return(nil)
	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(nil, nil)

	// Act
	response, err := suite.userService.VerifyOTP(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), userdto.LoginResponse{}, response)
}

// ==================== SetPassword Tests ====================

// Test: SetPassword should update user password
func (suite *UserServiceTestSuite) TestSetPassword_Success() {
	// Arrange
	userID := uint(1)
	password := "newPassword123"
	request := userdto.SetPasswordRequest{
		UserID:   userID,
		Password: password,
	}

	user := &entity.User{
		Phone: "09123456789",
	}
	user.ID = userID

	suite.userRepository.On("FindUserByID", suite.db, userID).Return(user, nil)
	suite.userRepository.On("UpdateUser", suite.db, mock.Anything).Return(nil)

	// Act
	err := suite.userService.SetPassword(request)

	// Assert
	assert.NoError(suite.T(), err)
	suite.userRepository.AssertCalled(suite.T(), "UpdateUser", suite.db, mock.Anything)
}

// Test: SetPassword should return error when user not found
func (suite *UserServiceTestSuite) TestSetPassword_UserNotFound() {
	// Arrange
	userID := uint(999)
	request := userdto.SetPasswordRequest{
		UserID:   userID,
		Password: "newPassword123",
	}

	suite.userRepository.On("FindUserByID", suite.db, userID).Return(nil, nil)

	// Act
	err := suite.userService.SetPassword(request)

	// Assert
	assert.Error(suite.T(), err)
}

// ==================== LoginWithPassword Tests ====================

// Test: LoginWithPassword should return error when password is invalid
func (suite *UserServiceTestSuite) TestLoginWithPassword_InvalidPassword() {
	// Arrange
	phone := "09123456789"
	password := "wrongPassword"
	request := userdto.LoginRequest{Phone: phone, Password: password}

	user := &entity.User{
		Phone:    phone,
		Password: "correctPassword",
	}
	user.ID = 1

	suite.userRepository.On("FindUserByPhone", suite.db, phone).Return(user, nil)

	// Act
	response, err := suite.userService.LoginWithPassword(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), userdto.LoginResponse{}, response)
}

// ==================== GetAllPermissions Tests ====================

// Test: GetAllPermissions should return all permissions
func (suite *UserServiceTestSuite) TestGetAllPermissions_Success() {
	// Arrange
	permissions := []*entity.Permission{
		{Name: "read"},
		{Name: "write"},
		{Name: "delete"},
	}
	permissions[0].ID = 1
	permissions[1].ID = 2
	permissions[2].ID = 3

	suite.userRepository.On("FindAllPermissions", suite.db).Return(permissions, nil)

	// Act
	response, err := suite.userService.GetAllPermissions()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(permissions), len(response))
}

// ==================== Role Management Tests ====================

// Test: CreateRole should return error when role name already exists
func (suite *UserServiceTestSuite) TestCreateRole_DuplicateRole() {
	// Arrange
	roleName := "existing_role"
	request := userdto.NewRoleRequest{
		Name:          roleName,
		PermissionIDs: []uint{1},
	}

	existingRole := &entity.Role{
		Name: roleName,
	}
	existingRole.ID = 1

	suite.userRepository.On("FindRoleByName", suite.db, roleName).Return(existingRole, nil)

	// Act
	err := suite.userService.CreateRole(request)

	// Assert
	assert.Error(suite.T(), err)
}

// Test: GetRoleDetails should return error when role not found
func (suite *UserServiceTestSuite) TestGetRoleDetails_NotFound() {
	// Arrange
	roleID := uint(999)
	suite.userRepository.On("FindRoleByID", suite.db, roleID).Return(nil, nil)

	// Act
	response, err := suite.userService.GetRoleDetails(roleID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), userdto.RoleResponse{}, response)
}

// Test: GetRoleOwners should return users assigned to role
func (suite *UserServiceTestSuite) TestGetRoleOwners_Success() {
	// Arrange
	roleID := uint(1)
	users := []*entity.User{
		{Phone: "09123456789"},
		{Phone: "09987654321"},
	}
	users[0].ID = 1
	users[1].ID = 2

	role := &entity.Role{
		Name: enum.Patient.String(),
	}
	role.ID = roleID

	suite.userRepository.On("FindRoleByID", suite.db, roleID).Return(role, nil)
	suite.userRepository.On("FindUsersByRoleID", suite.db, roleID).Return(users, nil)

	// Act
	response, err := suite.userService.GetRoleOwners(roleID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(response))
}

// Test: DeleteRole should delete role
func (suite *UserServiceTestSuite) TestDeleteRole_Success() {
	// Arrange
	roleID := uint(1)
	role := &entity.Role{
		Name: "test_role",
	}
	role.ID = roleID

	suite.userRepository.On("FindRoleByID", suite.db, roleID).Return(role, nil)
	suite.userRepository.On("DeleteRole", suite.db, roleID).Return(nil)

	// Act
	err := suite.userService.DeleteRole(roleID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test: DeleteRole should return error when role not found
func (suite *UserServiceTestSuite) TestDeleteRole_NotFound() {
	// Arrange
	roleID := uint(999)
	suite.userRepository.On("FindRoleByID", suite.db, roleID).Return(nil, nil)

	// Act
	err := suite.userService.DeleteRole(roleID)

	// Assert
	assert.Error(suite.T(), err)
}

// ==================== UpdateUserRoles Tests ====================

// Test: UpdateUserRoles should update user's roles
func (suite *UserServiceTestSuite) TestUpdateUserRoles_Success() {
	// Arrange
	actorID := uint(1)
	userID := uint(2)
	roleIDs := []uint{1, 2}
	request := userdto.UpdateUserRolesRequest{
		ActorID: actorID,
		UserID:  userID,
		RoleIDs: roleIDs,
	}

	user := &entity.User{
		Phone: "09123456789",
	}
	user.ID = userID

	role1 := &entity.Role{}
	role1.ID = 1
	role2 := &entity.Role{}
	role2.ID = 2

	suite.userRepository.On("FindUserByID", suite.db, userID).Return(user, nil)
	suite.userRepository.On("FindRoleByID", suite.db, uint(1)).Return(role1, nil)
	suite.userRepository.On("FindRoleByID", suite.db, uint(2)).Return(role2, nil)
	suite.userRepository.On("ReplaceUserRoles", suite.db, user, mock.Anything).Return(nil)
	suite.actionLogService.On("LogAction", mock.Anything).Return(nil)

	// Act
	err := suite.userService.UpdateUserRoles(request)

	// Assert
	assert.NoError(suite.T(), err)
}

// ==================== GetUsers Tests ====================

// Test: GetUsers should return paginated user list
func (suite *UserServiceTestSuite) TestGetUsers_Success() {
	// Arrange
	request := userdto.GetUsersListRequest{
		Offset: 0,
		Limit:  10,
	}

	users := []*entity.User{
		{Phone: "09123456789"},
		{Phone: "09987654321"},
	}
	users[0].ID = 1
	users[1].ID = 2

	suite.userRepository.On("FindUsers", suite.db, mock.Anything).Return(users, int64(2), nil)

	// Act
	response, count, err := suite.userService.GetUsers(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
	assert.Equal(suite.T(), 2, len(response))
}

// ==================== RequestUserValidationOTP Tests ====================

// Test: RequestUserValidationOTP should return error when operator not found
func (suite *UserServiceTestSuite) TestRequestUserValidationOTP_OperatorNotFound() {
	// Arrange
	operatorID := uint(999)
	request := userdto.RequestUserValidationOTPRequest{Phone: "09123456789"}

	suite.userRepository.On("FindUserByID", suite.db, operatorID).Return(nil, nil)

	// Act
	err := suite.userService.RequestUserValidationOTP(operatorID, request)

	// Assert
	assert.Error(suite.T(), err)
}
