package usecase

import (
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (u *UserServiceMock) GetUserByID(userID uint) (*entity.User, error) {
	args := u.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserServiceMock) Login(loginInfo userdto.LoginRequest) error {
	args := u.Called(loginInfo)
	return args.Error(0)
}

func (u *UserServiceMock) SetPassword(request userdto.SetPasswordRequest) error {
	args := u.Called(request)
	return args.Error(0)
}

func (u *UserServiceMock) VerifyOTP(verifyOTPInfo userdto.VerifyOTPRequest) (userdto.LoginResponse, error) {
	args := u.Called(verifyOTPInfo)
	if args.Get(0) == nil {
		return userdto.LoginResponse{}, args.Error(1)
	}
	return args.Get(0).(userdto.LoginResponse), args.Error(1)
}

func (u *UserServiceMock) GetAllPermissions() ([]userdto.PermissionResponse, error) {
	args := u.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.PermissionResponse), args.Error(1)
}

func (u *UserServiceMock) GetPermissionRoles(request userdto.GetPermissionRolesRequest) ([]userdto.RoleResponse, error) {
	args := u.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.RoleResponse), args.Error(1)
}

func (u *UserServiceMock) GetAllRoles() ([]userdto.RoleResponse, error) {
	args := u.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.RoleResponse), args.Error(1)
}

func (u *UserServiceMock) CreateRole(request userdto.NewRoleRequest) error {
	args := u.Called(request)
	return args.Error(0)
}

func (u *UserServiceMock) GetRoleDetails(roleID uint) (userdto.RoleResponse, error) {
	args := u.Called(roleID)
	if args.Get(0) == nil {
		return userdto.RoleResponse{}, args.Error(1)
	}
	return args.Get(0).(userdto.RoleResponse), args.Error(1)
}

func (u *UserServiceMock) GetRoleOwners(roleID uint) ([]userdto.UserResponse, error) {
	args := u.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.UserResponse), args.Error(1)
}

func (u *UserServiceMock) UpdateRole(request userdto.UpdateRoleRequest) error {
	args := u.Called(request)
	return args.Error(0)
}

func (u *UserServiceMock) DeleteRole(roleID uint) error {
	args := u.Called(roleID)
	return args.Error(0)
}

func (u *UserServiceMock) GetUserRoles(userID uint) ([]userdto.RoleResponse, error) {
	args := u.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.RoleResponse), args.Error(1)
}

func (u *UserServiceMock) UpdateUserRoles(request userdto.UpdateUserRolesRequest) error {
	args := u.Called(request)
	return args.Error(0)
}

func (u *UserServiceMock) GetUsers(request userdto.GetUsersListRequest) ([]userdto.UserResponse, int64, error) {
	args := u.Called(request)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]userdto.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (u *UserServiceMock) LoginWithPassword(loginInfo userdto.LoginRequest) (userdto.LoginResponse, error) {
	args := u.Called(loginInfo)
	if args.Get(0) == nil {
		return userdto.LoginResponse{}, args.Error(1)
	}
	return args.Get(0).(userdto.LoginResponse), args.Error(1)
}

func (u *UserServiceMock) RequestUserValidationOTP(operatorID uint, request userdto.RequestUserValidationOTPRequest) error {
	args := u.Called(operatorID, request)
	return args.Error(0)
}

func (u *UserServiceMock) VerifyUserValidationOTP(operatorID uint, request userdto.VerifyUserValidationOTPRequest) (userdto.UserValidationResponse, error) {
	args := u.Called(operatorID, request)
	if args.Get(0) == nil {
		return userdto.UserValidationResponse{}, args.Error(1)
	}
	return args.Get(0).(userdto.UserValidationResponse), args.Error(1)
}

func (u *UserServiceMock) ValidateUserForFormCreation(operatorID, userID uint) error {
	args := u.Called(operatorID, userID)
	return args.Error(0)
}

func (u *UserServiceMock) GetUserProfile(operatorID uint) (userdto.UserProfileResponse, error) {
	args := u.Called(operatorID)

	var resp userdto.UserProfileResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(userdto.UserProfileResponse)
	}

	return resp, args.Error(1)
}

func (u *UserServiceMock) SubmitUserProfile(request userdto.SubmitUserProfileRequest) (userdto.UserProfileResponse, error) {
	args := u.Called(request)

	var resp userdto.UserProfileResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(userdto.UserProfileResponse)
	}

	return resp, args.Error(1)
}
