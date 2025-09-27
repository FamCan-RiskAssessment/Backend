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
