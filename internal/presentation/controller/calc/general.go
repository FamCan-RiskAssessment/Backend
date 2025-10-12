package calc

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralCalcController struct {
	calcService usecase.CalcService
}

func NewGeneralCalcController(
	calcService usecase.CalcService,
) *GeneralCalcController {
	return &GeneralCalcController{
		calcService: calcService,
	}
}

func (calcController *GeneralCalcController) GetAllModelTypes(ctx *gin.Context) {
	response := calcController.calcService.GetAllModelTypes
	controller.Response(ctx, 200, "", response)
}
