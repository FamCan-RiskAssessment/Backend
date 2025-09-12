package formdto

type CreateFormRequest struct {
	Name                 string `json:"name" binding:"required"`
	DateOfBirth          string `json:"date_of_birth" binding:"required"`
	Address              string `json:"address" binding:"required"`
	PostalCode           string `json:"postal_code" binding:"required"`
	SocialSecurityNumber string `json:"social_security_number" binding:"required"`
}

type UpdateFormRequest struct {
	FormID               uint    `json:"form_id" binding:"required"`
	Name                 *string `json:"name,omitempty"`
	DateOfBirth          *string `json:"date_of_birth,omitempty"`
	Address              *string `json:"address,omitempty"`
	PostalCode           *string `json:"postal_code,omitempty"`
	SocialSecurityNumber *string `json:"social_security_number,omitempty"`
}

type GetUserFormsRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Offset int  `json:"offset"`
	Limit  int  `json:"limit"`
}

type GetFormRequest struct {
	FormID uint `json:"form_id" binding:"required"`
}

type DeleteFormRequest struct {
	FormID uint `json:"form_id" binding:"required"`
}
