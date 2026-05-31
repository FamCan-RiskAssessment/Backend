package calc

import (
	"log"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCalcController struct {
	constants   *bootstrap.Constants
	calcService usecase.CalcService
	pagination  *bootstrap.Pagination
}

func NewAdminCalcController(
	constants *bootstrap.Constants,
	calcService usecase.CalcService,
	pagination *bootstrap.Pagination,
) *AdminCalcController {
	return &AdminCalcController{
		constants:   constants,
		calcService: calcService,
		pagination:  pagination,
	}
}

func (calcController *AdminCalcController) GetCalcBrowse(ctx *gin.Context) {
	type GetCalcBrowseParams struct {
		Page      int     `form:"page"`
		PageSize  int     `form:"pageSize"`
		Status    *uint   `form:"status"`
		SortBy    *string `form:"sortBy"`
		SortOrder *string `form:"sortOrder"`
		Search    *string `form:"search"`
	}

	params := controller.Validate[GetCalcBrowseParams](ctx)
	offset, limit := controller.GetOffsetLimit(
		params.Page,
		params.PageSize,
		calcController.pagination.DefaultPage,
		calcController.pagination.DefaultPageSize,
	)

	filters := &postgres.FormFilters{
		Status: params.Status,
	}

	request := calcdto.GetCalcBrowseRequest{
		Offset:    offset,
		Limit:     limit,
		Filters:   filters,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		Search:    params.Search,
	}

	items, count, err := calcController.calcService.GetCalcBrowse(request)
	if err != nil {
		log.Printf("[CALC ERROR] GetCalcBrowse failed - Error: %v\n", err)
		panic(err)
	}

	data := controller.NewPaginatedResponse(items, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (calcController *AdminCalcController) SendFormToCalc(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `json:"formID" validate:"required"`
		CalcID uint `json:"calcID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	userID, _ := ctx.Get(calcController.constants.Context.ID)
	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
		CalcID: params.CalcID,
	}

	response, err := calcController.calcService.SendFormToCalc(sendFormToCalcRequest)
	if err != nil {
		log.Printf("[CALC ERROR] SendFormToCalc failed - UserID: %d, FormID: %d, CalcID: %d, Error: %v\n",
			sendFormToCalcRequest.UserID, params.FormID, params.CalcID, err)
		panic(err)
	}

	trans := controller.GetTranslator(ctx, calcController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.sendFormToCalc")
	controller.Response(ctx, 200, message, response)
}

func (calcController *AdminCalcController) GetPremm5Results(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	userID, _ := ctx.Get(calcController.constants.Context.ID)
	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := calcController.calcService.GetPremm5Results(sendFormToCalcRequest)
	if err != nil {
		log.Printf("[CALC ERROR] GetPremm5Results failed - UserID: %d, FormID: %d, Error: %v\n",
			sendFormToCalcRequest.UserID, params.FormID, err)
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (calcController *AdminCalcController) GetBCRAResults(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	userID, _ := ctx.Get(calcController.constants.Context.ID)
	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := calcController.calcService.GetBCRAResults(sendFormToCalcRequest)
	if err != nil {
		log.Printf("[CALC ERROR] GetBCRAResults failed - UserID: %d, FormID: %d, Error: %v\n",
			sendFormToCalcRequest.UserID, params.FormID, err)
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (calcController *AdminCalcController) GetGailResults(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	userID, _ := ctx.Get(calcController.constants.Context.ID)
	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := calcController.calcService.GetGailResults(sendFormToCalcRequest)
	if err != nil {
		log.Printf("[CALC ERROR] GetGailResults failed - UserID: %d, FormID: %d, Error: %v\n",
			sendFormToCalcRequest.UserID, params.FormID, err)
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (calcController *AdminCalcController) GetPLCOResults(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	userID, _ := ctx.Get(calcController.constants.Context.ID)
	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := calcController.calcService.GetPLCOResults(sendFormToCalcRequest)
	if err != nil {
		log.Printf("[CALC ERROR] GetPLCOResults failed - UserID: %d, FormID: %d, Error: %v\n",
			sendFormToCalcRequest.UserID, params.FormID, err)
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (calcController *AdminCalcController) GetAllModelTypes(ctx *gin.Context) {
	response := calcController.calcService.GetAllModelTypes()
	controller.Response(ctx, 200, "", response)
}
