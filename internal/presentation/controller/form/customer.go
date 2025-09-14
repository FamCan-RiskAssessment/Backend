package form

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
	pagination  *bootstrap.Pagination
}

func NewCustomerFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
	pagination *bootstrap.Pagination,
) *CustomerFormController {
	return &CustomerFormController{
		constants:   constants,
		formService: formService,
		pagination:  pagination,
	}
}

func (formController *CustomerFormController) CreateForm(ctx *gin.Context) {
	type CreateFormParams struct {
		Name                 string `json:"name" validate:"required"`
		DateOfBirth          string `json:"date_of_birth" validate:"required"`
		Address              string `json:"address" validate:"required"`
		PostalCode           string `json:"postal_code" validate:"required"`
		SocialSecurityNumber string `json:"social_security_number" validate:"required"`
	}

	params := controller.Validate[CreateFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.CreateFormRequest{
		UserID:               userID.(uint),
		Name:                 params.Name,
		DateOfBirth:          params.DateOfBirth,
		Address:              params.Address,
		PostalCode:           params.PostalCode,
		SocialSecurityNumber: params.SocialSecurityNumber,
	}

	response, err := formController.formService.CreateForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, response)
}

func (formController *CustomerFormController) GetUserForms(ctx *gin.Context) {
	type GetUserFormsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	params := controller.Validate[GetUserFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetUserFormsRequest{
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}

	response, err := formController.formService.GetUserForms(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) GetForm(ctx *gin.Context) {
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

func (formController *CustomerFormController) UpdateForm(ctx *gin.Context) {
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

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, response)
}

func (formController *CustomerFormController) DeleteForm(ctx *gin.Context) {
	type DeleteFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[DeleteFormParams](ctx)

	response, err := formController.formService.DeleteForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, response)
}
