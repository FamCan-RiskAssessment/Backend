package form

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
	pagination  *bootstrap.Pagination
}

func NewAdminFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
	pagination *bootstrap.Pagination,
) *AdminFormController {
	return &AdminFormController{
		constants:   constants,
		formService: formService,
		pagination:  pagination,
	}
}

func (formController *AdminFormController) GetAllForms(ctx *gin.Context) {
	type GetAllFormsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	params := controller.Validate[GetAllFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	forms, count, err := formController.formService.GetAllForms(offset, limit)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(forms, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (formController *AdminFormController) DeleteForm(ctx *gin.Context) {
	type DeleteFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[DeleteFormParams](ctx)

	err := formController.formService.DeleteForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) AcceptForm(ctx *gin.Context) {
	type AcceptFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[AcceptFormParams](ctx)

	err := formController.formService.AcceptForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.acceptForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) RejectForm(ctx *gin.Context) {
	type RejectFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[RejectFormParams](ctx)

	err := formController.formService.RejectForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) GetBasicForm(ctx *gin.Context) {
	type GetBasicFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetBasicFormParams](ctx)

	response, err := formController.formService.GetBasicForm(params.FormID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
