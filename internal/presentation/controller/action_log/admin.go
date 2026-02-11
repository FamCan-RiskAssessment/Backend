package actionlog

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type ActionLogController struct {
	actionLogService usecase.ActionLogService
	formService      usecase.FormService
	userService      usecase.UserService
	pagination       *bootstrap.Pagination
}

func NewActionLogController(
	actionLogService usecase.ActionLogService,
	formService usecase.FormService,
	userService usecase.UserService,
	pagination *bootstrap.Pagination,

) *ActionLogController {
	return &ActionLogController{
		actionLogService: actionLogService,
		formService:      formService,
		userService:      userService,
		pagination:       pagination,
	}
}

func (alc *ActionLogController) GetAllActionLogs(ctx *gin.Context) {
	type GetAllActionLogsParams struct {
		Page      int              `form:"page"`
		PageSize  int              `form:"pageSize"`
		SortBy    *string          `form:"sortBy"`
		SortOrder *string          `form:"sortOrder"`
		Search    *string          `form:"search"`
		Action    *enum.ActionType `form:"action"`
		ActorID   *uint            `form:"actorId"`
		DateFrom  *string          `form:"dateFrom"`
		DateTo    *string          `form:"dateTo"`
	}

	params := controller.Validate[GetAllActionLogsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, alc.pagination.DefaultPage, alc.pagination.DefaultPageSize)

	request := actionlogdto.GetAllActionLogsRequest{
		Offset:    offset,
		Limit:     limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		Search:    params.Search,
		Action:    params.Action,
		ActorID:   params.ActorID,
		DateFrom:  params.DateFrom,
		DateTo:    params.DateTo,
	}

	actionLogs, count, err := alc.actionLogService.GetAllActionLogs(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(actionLogs, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (alc *ActionLogController) GetAllActionTypes(ctx *gin.Context) {
	types, err := alc.actionLogService.GetAllActionTypes()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", types)
}
