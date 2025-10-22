package userdto

type LoginResponse struct {
	UserID       uint                 `json:"id"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	Permissions  []PermissionResponse `json:"permissions"`
	Roles        []RoleResponse       `json:"roles"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
}

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type PermissionResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
}
