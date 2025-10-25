package actionlog

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type ActionLogController struct {
	actionLogService usecase.ActionLogService
	pagination       *bootstrap.Pagination
}

func NewActionLogController(
	actionLogService usecase.ActionLogService,
	pagination *bootstrap.Pagination,

) *ActionLogController {
	return &ActionLogController{
		actionLogService: actionLogService,
		pagination:       pagination,
	}
}

func (alc *ActionLogController) GetAllActionLogs(ctx *gin.Context) {
	type GetAllActionLogs struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	params := controller.Validate[GetAllActionLogs](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, alc.pagination.DefaultPage, alc.pagination.DefaultPageSize)

	actionLogs, count, err := alc.actionLogService.GetAllActionLogs(offset, limit)
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
