package formdto

import "time"

type BasicFormResponse struct {
	FormID uint   `json:"id"`
	Status string `json:"status"`
	UserID uint   `json:"user_id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateFormResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpdateFormResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type GetUserFormsResponse struct {
	Forms []BasicFormResponse `json:"forms"` // TODO: CHANGE TO FORMRESPONSE
	Total int                 `json:"total"`
}

type DeleteFormResponse struct {
	Message string `json:"message"`
}

type UpsertGeneralHealthResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpsertMamographyResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpsertCancerResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpsertFamilyCancerResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpsertContactResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}

type UpsertLungCancerResponse struct {
	Form    BasicFormResponse `json:"form"`
	Message string            `json:"message"`
}
