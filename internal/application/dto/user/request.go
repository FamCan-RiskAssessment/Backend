package userdto

type LoginRequest struct {
	Phone    string
	Password string
}

type VerifyOTPRequest struct {
	Phone string
	OTP   string
}

type NewRoleRequest struct {
	Name          string
	PermissionIDs []uint
}

type UpdateRoleRequest struct {
	RoleID        uint
	Name          *string
	PermissionIDs []uint
}

type UpdateUserRolesRequest struct {
	ActorID uint
	UserID  uint
	RoleIDs []uint
}

type GetUsersListRequest struct {
	Offset int
	Limit  int
}

type GetPermissionRolesRequest struct {
	PermissionID uint
	Offset       int
	Limit        int
}

type SetPasswordRequest struct {
	UserID   uint
	Password string
}

type RequestUserValidationOTPRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11"`
}

type VerifyUserValidationOTPRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11"`
	OTP   string `json:"otp" validate:"required,min=6,max=6"`
}

type SubmitUserProfileRequest struct {
	UserID               uint
	Name                 string
	LastName             string
	HealthCenter         string
	SocialSecurityNumber string
}
