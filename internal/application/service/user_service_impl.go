package service

import (
	"context"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/communication"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	redis "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/redis"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/google/uuid"
)

type UserService struct {
	constants           *bootstrap.Constants
	userRepository      postgres.UserRepository
	userCacheRepository redis.UserCacheRepository
	jwtService          usecase.JwtService
	smsService          communication.SmsService
	otpService          usecase.OtpService
	actionLogService    usecase.ActionLogService
	db                  database.Database
	passwordHasher      usecase.PasswordHasher
}

func NewUserService(
	constants *bootstrap.Constants,
	userRepository postgres.UserRepository,
	userCacheRepository redis.UserCacheRepository,
	jwtService usecase.JwtService,
	smsService communication.SmsService,
	otpService usecase.OtpService,
	actionLogService usecase.ActionLogService,
	db database.Database,
	passwordHasher usecase.PasswordHasher,
) *UserService {
	return &UserService{
		constants:           constants,
		userRepository:      userRepository,
		userCacheRepository: userCacheRepository,
		jwtService:          jwtService,
		smsService:          smsService,
		otpService:          otpService,
		actionLogService:    actionLogService,
		db:                  db,
		passwordHasher:      passwordHasher,
	}
}

func (userService *UserService) GetUserByID(userID uint) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByID(userService.db, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	return user, nil
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) error {
	user, err := userService.userRepository.FindUserByPhone(userService.db, loginInfo.Phone)
	if err != nil {
		return err
	}

	if user == nil {
		user = &entity.User{
			Phone: loginInfo.Phone,
		}
		err = userService.userRepository.CreateUser(userService.db, user)
		if err != nil {
			return err
		}
		patientRole, err := userService.userRepository.FindRoleByName(userService.db, enum.Patient.String())
		if err != nil {
			return err
		}
		if patientRole != nil {
			err = userService.userRepository.AssignRoleToUser(userService.db, user, patientRole)
			if err != nil {
				return err
			}
		}
	}

	otp, exppireMinute, err := userService.otpService.GenerateOTP(loginInfo.Phone)
	if err != nil {
		return err
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(loginInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(exppireMinute)*time.Minute)
	if err != nil {
		return err
	}

	err = userService.smsService.SendOTP(loginInfo.Phone, otp)
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) VerifyOTP(verifyOTPInfo userdto.VerifyOTPRequest) (userdto.LoginResponse, error) {
	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyOTPInfo.Phone)
	err := userService.otpService.VerifyOTP(redisKey, verifyOTPInfo.OTP)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	user, err := userService.userRepository.FindUserByPhone(userService.db, verifyOTPInfo.Phone)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return userdto.LoginResponse{}, notFoundError
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	permissions, err := userService.FindUserPermissions(user)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	roles, err := userService.GetUserRoles(user.ID)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	return userdto.LoginResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Permissions:  permissions,
		Roles:        roles,
	}, nil
}

func (userService *UserService) SetPassword(request userdto.SetPasswordRequest) error {
	user, err := userService.GetUserByID(request.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return notFoundError
	}

	hashedPassword, err := userService.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := userService.userRepository.UpdateUser(userService.db, user); err != nil {
		return err
	}
	return nil
}

func (userService *UserService) LoginWithPassword(loginInfo userdto.LoginRequest) (userdto.LoginResponse, error) {

	user, err := userService.userRepository.FindUserByPhone(userService.db, loginInfo.Phone)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return userdto.LoginResponse{}, notFoundError
	}

	if err := userService.passwordHasher.VerifyPassword(loginInfo.Password, user.Password); err != nil {
		authError := exception.NewInvalidCredentialsError("phone and password not match", nil)
		return userdto.LoginResponse{}, authError
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	permissions, err := userService.FindUserPermissions(user)
	if err != nil {
		return userdto.LoginResponse{}, err
	}

	return userdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Permissions:  permissions,
	}, nil
}

func (userService *UserService) FindUserPermissions(user *entity.User) ([]userdto.PermissionResponse, error) {
	var permissions []userdto.PermissionResponse

	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		return nil, err
	}

	for _, role := range user.Roles {
		rolePermissions, err := userService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}

func (userService *UserService) getRolePermissions(role *entity.Role) ([]userdto.PermissionResponse, error) {
	if err := userService.userRepository.FindRolePermissions(userService.db, role); err != nil {
		return nil, err
	}
	permissions := make([]userdto.PermissionResponse, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissions[i] = userdto.PermissionResponse{
			ID:       permission.ID,
			Name:     permission.Type.String(),
			Category: permission.Category.String(),
		}
	}
	return permissions, nil
}

func (userService *UserService) GetAllPermissions() ([]userdto.PermissionResponse, error) {
	permissions, err := userService.userRepository.FindAllPermissions(userService.db)
	if err != nil {
		return nil, err
	}
	permissionsResponse := make([]userdto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionsResponse[i] = userdto.PermissionResponse{
			ID:       permission.ID,
			Name:     permission.Type.String(),
			Category: permission.Category.String(),
		}
	}
	return permissionsResponse, nil
}

func (userService *UserService) GetAllRoles() ([]userdto.RoleResponse, error) {
	roles, err := userService.userRepository.FindAllRoles(userService.db)
	if err != nil {
		return nil, err
	}
	rolesResponse := make([]userdto.RoleResponse, len(roles))
	for i, role := range roles {
		permissions, err := userService.getRolePermissions(role)
		if err != nil {
			return nil, err
		}
		rolesResponse[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}
	return rolesResponse, nil
}

func (userService *UserService) getPermission(permissionID uint) (*entity.Permission, error) {
	permission, err := userService.userRepository.FindPermissionByID(userService.db, permissionID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
		return nil, notFoundError
	}
	return permission, nil
}

func (userService *UserService) getRole(roleID uint) (*entity.Role, error) {
	role, err := userService.userRepository.FindRoleByID(userService.db, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		return nil, notFoundError
	}
	return role, nil
}

func (userService *UserService) CreateRole(newRoleRequest userdto.NewRoleRequest) error {
	existingRole, err := userService.userRepository.FindRoleByName(userService.db, newRoleRequest.Name)
	if err != nil {
		return err
	}
	if existingRole != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Role, userService.constants.Tag.AlreadyExist)
		return conflictErrors
	}
	err = userService.db.WithTransaction(func(tx database.Database) error {
		role := &entity.Role{
			Name: newRoleRequest.Name,
		}
		err = userService.userRepository.CreateRole(tx, role)
		if err != nil {
			return err
		}

		existingPermissions := make(map[uint]bool)
		for _, permissionID := range newRoleRequest.PermissionIDs {
			if existingPermissions[permissionID] {
				continue
			}

			permission, err := userService.getPermission(permissionID)
			if err != nil {
				return err
			}

			if err := userService.userRepository.AssignPermissionToRole(tx, role, permission); err != nil {
				return err
			}
			existingPermissions[permissionID] = true
		}

		return nil
	})
	return nil
}

func (userService *UserService) GetRoleDetails(roleID uint) (userdto.RoleResponse, error) {
	role, err := userService.getRole(roleID)
	if err != nil {
		return userdto.RoleResponse{}, err
	}

	permissions, err := userService.getRolePermissions(role)
	if err != nil {
		return userdto.RoleResponse{}, err
	}

	return userdto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (userService *UserService) GetRoleOwners(roleID uint) ([]userdto.UserResponse, error) {
	_, err := userService.getRole(roleID)
	if err != nil {
		return nil, err
	}

	users, err := userService.userRepository.FindProfilesByRoleID(userService.db, roleID)
	if err != nil {
		return nil, err
	}

	userCreds := make([]userdto.UserResponse, len(users))
	for i, user := range users {
		var Name, LastName, HealthCenter, SocialSecurityNumber string = "-", "-", "-", "-"
		if user.UserProfile != nil {
			Name = user.UserProfile.Name
			LastName = user.UserProfile.LastName
			HealthCenter = user.UserProfile.HealthCenter
			SocialSecurityNumber = user.UserProfile.SocialSecurityNumber
		}
		userCreds[i] = userdto.UserResponse{
			ID:                   user.ID,
			Phone:                user.Phone,
			Name:                 &Name,
			LastName:             &LastName,
			HealthCenter:         &HealthCenter,
			SocialSecurityNumber: &SocialSecurityNumber,
		}
	}
	return userCreds, nil
}

func (userService *UserService) GetUserRoles(userID uint) ([]userdto.RoleResponse, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		return nil, err
	}
	roles := make([]userdto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		permissions, err := userService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		roles[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}
	return roles, nil
}

func (userService *UserService) DeleteRole(roleID uint) error {
	_, err := userService.getRole(roleID)
	if err != nil {
		return err
	}

	if err := userService.userRepository.DeleteRole(userService.db, roleID); err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UpdateRole(newRoleRequest userdto.UpdateRoleRequest) error {
	role, err := userService.getRole(newRoleRequest.RoleID)
	if err != nil {
		return err
	}

	existingPermissions := make(map[uint]bool)
	var permissions []entity.Permission
	for _, permissionID := range newRoleRequest.PermissionIDs {
		if existingPermissions[permissionID] {
			continue
		}

		permission, err := userService.getPermission(permissionID)
		if err != nil {
			return err
		}

		permissions = append(permissions, *permission)
		existingPermissions[permissionID] = true
	}

	err = userService.db.WithTransaction(func(tx database.Database) error {
		if newRoleRequest.Name != nil {
			role.Name = *newRoleRequest.Name
			if err := userService.userRepository.UpdateRole(tx, role); err != nil {
				return err
			}
		}

		if err := userService.userRepository.ReplaceRolePermissions(tx, role, permissions); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (userService *UserService) UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) error {
	user, err := userService.GetUserByID(userRolesRequest.UserID)
	if err != nil {
		return err
	}

	existingRoles := make(map[uint]bool)
	var roles []entity.Role
	for _, roleID := range userRolesRequest.RoleIDs {
		if existingRoles[roleID] {
			continue
		}

		role, err := userService.getRole(roleID)
		if err != nil {
			return err
		}

		roles = append(roles, *role)
		existingRoles[roleID] = true
	}

	if err := userService.userRepository.ReplaceUserRoles(userService.db, user, roles); err != nil {
		return err
	}

	log := actionlogdto.LogAction{
		ActorID:  userRolesRequest.ActorID,
		TargetID: &userRolesRequest.UserID,
		Action:   enum.ActionTypeUserRoleUpdated,
		Details:  "رول های یوزر بروزرسانی شد",
	}

	userService.actionLogService.LogAction(log)

	return nil
}

var userAllowedSortColumns = []string{"id", "phone", "created_at"}

func (userService *UserService) GetUsers(request userdto.GetUsersListRequest) ([]userdto.UserResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	sortCol := "id"
	if request.SortBy != nil {
		sortCol = postgres.ValidateSortColumn(*request.SortBy, userAllowedSortColumns, "id")
	}
	asc := true
	if request.SortOrder != nil && *request.SortOrder == "desc" {
		asc = false
	}
	options.WithSorting(sortCol, asc)

	if request.Search != nil && *request.Search != "" {
		options.WithSearch(*request.Search, []string{"phone"})
	}

	filters := &postgres.UserFilters{
		RoleID: request.RoleID,
	}

	users, count, err := userService.userRepository.FindUsers(userService.db, options, filters)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]userdto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = userdto.UserResponse{
			ID:    user.ID,
			Phone: user.Phone,
		}
	}
	return userResponses, count, nil
}

func (userService *UserService) GetPermissionRoles(request userdto.GetPermissionRolesRequest) ([]userdto.RoleResponse, error) {
	permission, err := userService.userRepository.FindPermissionByID(userService.db, request.PermissionID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
		return nil, notFoundError
	}

	roles, err := userService.userRepository.FindRolesByPermission(userService.db, request.PermissionID)
	if err != nil {
		return nil, err
	}
	rolesResponse := make([]userdto.RoleResponse, len(roles))
	for i, role := range roles {
		rolesResponse[i] = userdto.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}
	return rolesResponse, nil
}

func (userService *UserService) RequestUserValidationOTP(operatorID uint, request userdto.RequestUserValidationOTPRequest) error {
	// Verify operator has operator role
	userRoles, err := userService.GetUserRoles(operatorID)
	if err != nil {
		return err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		return exception.ForbiddenError{Resource: userService.constants.Field.Role}
	}

	// Find or create user by phone
	user, err := userService.userRepository.FindUserByPhone(userService.db, request.Phone)
	if err != nil {
		return err
	}
	if user == nil {
		// Create user if doesn't exist
		user = &entity.User{
			Phone: request.Phone,
		}
		err = userService.userRepository.CreateUser(userService.db, user)
		if err != nil {
			return err
		}
		patientRole, err := userService.userRepository.FindRoleByName(userService.db, enum.Patient.String())
		if err != nil {
			return err
		}
		if patientRole != nil {
			err = userService.userRepository.AssignRoleToUser(userService.db, user, patientRole)
			if err != nil {
				return err
			}
		}
	}

	// Generate OTP
	otp, expireMinute, err := userService.otpService.GenerateOTP(request.Phone)
	if err != nil {
		return err
	}

	// Store OTP with operator-specific key
	redisKey := userService.constants.RedisKey.GenerateOperatorValidationOTPKey(operatorID, request.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expireMinute)*time.Minute)
	if err != nil {
		return err
	}

	// Send OTP via SMS (uncomment when ready)
	// err = userService.smsService.SendOTP(request.Phone, otp)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (userService *UserService) VerifyUserValidationOTP(operatorID uint, request userdto.VerifyUserValidationOTPRequest) (userdto.UserValidationResponse, error) {
	// Verify operator has operator role
	userRoles, err := userService.GetUserRoles(operatorID)
	if err != nil {
		return userdto.UserValidationResponse{}, err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		return userdto.UserValidationResponse{}, exception.ForbiddenError{Resource: userService.constants.Field.Role}
	}

	// Verify OTP
	redisKey := userService.constants.RedisKey.GenerateOperatorValidationOTPKey(operatorID, request.Phone)
	err = userService.otpService.VerifyOTP(redisKey, request.OTP)
	if err != nil {
		return userdto.UserValidationResponse{}, err
	}

	// Find user by phone
	user, err := userService.userRepository.FindUserByPhone(userService.db, request.Phone)
	if err != nil {
		return userdto.UserValidationResponse{}, err
	}
	if user == nil {
		return userdto.UserValidationResponse{}, exception.NotFoundError{Item: userService.constants.Field.User}
	}

	// Generate validation token and store it
	validationToken := uuid.New().String()
	validationKey := userService.constants.RedisKey.GenerateOperatorValidationTokenKey(operatorID, user.ID)
	// Store for 30 minutes (operator can create forms for this user)
	expiration := 30 * time.Minute
	err = userService.userCacheRepository.SetValidationToken(context.Background(), validationKey, expiration)
	if err != nil {
		return userdto.UserValidationResponse{}, err
	}

	// Delete OTP after successful verification
	_ = userService.userCacheRepository.Delete(context.Background(), redisKey)

	return userdto.UserValidationResponse{
		ValidationToken: validationToken,
		ExpiresIn:       int(expiration.Seconds()),
		UserID:          user.ID,
	}, nil
}

func (userService *UserService) ValidateUserForFormCreation(operatorID, userID uint) error {
	validationKey := userService.constants.RedisKey.GenerateOperatorValidationTokenKey(operatorID, userID)
	valid, err := userService.userCacheRepository.GetValidationToken(context.Background(), validationKey)
	if err != nil {
		return err
	}
	if !valid {
		return exception.ForbiddenError{Resource: userService.constants.Field.User}
	}
	return nil
}

func (userService *UserService) GetUserProfile(operatorID uint) (userdto.UserProfileResponse, error) {
	user, err := userService.userRepository.FindProfileByUserID(userService.db, operatorID)
	if err != nil {
		return userdto.UserProfileResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return userdto.UserProfileResponse{}, notFoundError
	}
	if user.UserProfile == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Profile}
		return userdto.UserProfileResponse{}, notFoundError
	}
	return userdto.UserProfileResponse{
		Name:                 user.UserProfile.Name,
		LastName:             user.UserProfile.LastName,
		Phone:                user.Phone,
		HealthCenter:         user.UserProfile.HealthCenter,
		SocialSecurityNumber: user.UserProfile.SocialSecurityNumber,
	}, nil
}

func (userService *UserService) SubmitUserProfile(request userdto.SubmitUserProfileRequest) (userdto.UserProfileResponse, error) {
	user, err := userService.userRepository.FindProfileByUserID(userService.db, request.UserID)
	if err != nil {
		return userdto.UserProfileResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return userdto.UserProfileResponse{}, notFoundError
	}
	profile := &entity.UserProfile{
		UserID:               user.ID,
		Name:                 request.Name,
		LastName:             request.LastName,
		HealthCenter:         request.HealthCenter,
		SocialSecurityNumber: request.SocialSecurityNumber,
	}
	userProfileResponse := userdto.UserProfileResponse{
		Name:                 profile.Name,
		Phone:                user.Phone,
		LastName:             profile.LastName,
		HealthCenter:         profile.HealthCenter,
		SocialSecurityNumber: profile.SocialSecurityNumber,
	}
	if user.UserProfile != nil {
		// If the UserProfile Exists, then update the current one
		if err := userService.userRepository.UpdateProfile(userService.db, profile); err != nil {
			return userProfileResponse, nil
		}
		return userProfileResponse, nil
	}
	if err := userService.userRepository.CreateProfile(userService.db, profile); err != nil {
		return userdto.UserProfileResponse{}, nil
	}
	return userProfileResponse, nil
}
