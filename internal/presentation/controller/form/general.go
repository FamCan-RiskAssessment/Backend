package form

import (
	"strconv"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
}

func NewGeneralFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
) *GeneralFormController {
	return &GeneralFormController{
		constants:   constants,
		formService: formService,
	}
}

func (formController *GeneralFormController) CreateForm(ctx *gin.Context) {
	type CreateFormParams struct {
		Name                 string `json:"name" validate:"required"`
		DateOfBirth          string `json:"date_of_birth" validate:"required"`
		Address              string `json:"address" validate:"required"`
		PostalCode           string `json:"postal_code" validate:"required"`
		SocialSecurityNumber string `json:"social_security_number" validate:"required"`
	}

	params := controller.Validate[CreateFormParams](ctx)

	// Get user ID from context (assuming it's set by authentication middleware)
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		panic("User ID not found in context")
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		panic("Invalid user ID type")
	}

	request := formdto.CreateFormRequest{
		Name:                 params.Name,
		DateOfBirth:          params.DateOfBirth,
		Address:              params.Address,
		PostalCode:           params.PostalCode,
		SocialSecurityNumber: params.SocialSecurityNumber,
	}

	response, err := formController.formService.CreateForm(userID, request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 201, "Form created successfully", response)
}

func (formController *GeneralFormController) GetUserForms(ctx *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		panic("User ID not found in context")
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		panic("Invalid user ID type")
	}

	// Get pagination parameters
	offsetStr := ctx.DefaultQuery("offset", "0")
	limitStr := ctx.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	request := formdto.GetUserFormsRequest{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}

	response, err := formController.formService.GetUserForms(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *GeneralFormController) GetForm(ctx *gin.Context) {
	type GetFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[GetFormParams](ctx)

	response, err := formController.formService.GetForm(params.FormID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *GeneralFormController) UpdateForm(ctx *gin.Context) {
	type UpdateFormParams struct {
		FormID               uint    `json:"form_id" validate:"required"`
		Name                 *string `json:"name,omitempty"`
		DateOfBirth          *string `json:"date_of_birth,omitempty"`
		Address              *string `json:"address,omitempty"`
		PostalCode           *string `json:"postal_code,omitempty"`
		SocialSecurityNumber *string `json:"social_security_number,omitempty"`
	}

	params := controller.Validate[UpdateFormParams](ctx)

	request := formdto.UpdateFormRequest{
		FormID:               params.FormID,
		Name:                 params.Name,
		DateOfBirth:          params.DateOfBirth,
		Address:              params.Address,
		PostalCode:           params.PostalCode,
		SocialSecurityNumber: params.SocialSecurityNumber,
	}

	response, err := formController.formService.UpdateForm(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "Form updated successfully", response)
}

func (formController *GeneralFormController) DeleteForm(ctx *gin.Context) {
	type DeleteFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[DeleteFormParams](ctx)

	response, err := formController.formService.DeleteForm(params.FormID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "Form deleted successfully", response)
}
