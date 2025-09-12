package form

import (
	"strconv"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
}

func NewAdminFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
) *AdminFormController {
	return &AdminFormController{
		constants:   constants,
		formService: formService,
	}
}

func (formController *AdminFormController) GetAllForms(ctx *gin.Context) {
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

	response, err := formController.formService.GetAllForms(offset, limit)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *AdminFormController) GetForm(ctx *gin.Context) {
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

func (formController *AdminFormController) DeleteForm(ctx *gin.Context) {
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
