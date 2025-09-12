package formdto

type FormResponse struct {
	ID                   uint   `json:"id"`
	UserID               uint   `json:"user_id"`
	Name                 string `json:"name"`
	DateOfBirth          string `json:"date_of_birth"`
	Address              string `json:"address"`
	PostalCode           string `json:"postal_code"`
	SocialSecurityNumber string `json:"social_security_number"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

type CreateFormResponse struct {
	Form    FormResponse `json:"form"`
	Message string       `json:"message"`
}

type UpdateFormResponse struct {
	Form    FormResponse `json:"form"`
	Message string       `json:"message"`
}

type GetUserFormsResponse struct {
	Forms []FormResponse `json:"forms"`
	Total int            `json:"total"`
}

type DeleteFormResponse struct {
	Message string `json:"message"`
}
