package calc

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCalcController struct {
	constants   *bootstrap.Constants
	calcService usecase.CalcService
}

func NewAdminCalcController(
	constants *bootstrap.Constants,
	calcService usecase.CalcService,
) *AdminCalcController {
	return &AdminCalcController{
		constants:   constants,
		calcService: calcService,
	}
}

func (calcController *AdminCalcController) SendFormToCalc(ctx *gin.Context) {
	type SendFormToCalcParams struct {
		FormID uint `json:"formI" validate:"required"`
		CalcID uint `json:"calcID" validate:"required"`
	}
	params := controller.Validate[SendFormToCalcParams](ctx)

	sendFormToCalcRequest := calcdto.SendFormToCalcRequest{
		FormID: params.FormID,
		CalcID: params.CalcID,
	}

	err := calcController.calcService.SendFormToCalc(sendFormToCalcRequest)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, calcController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.sendFormToCalc")
	controller.Response(ctx, 200, message, nil)
}
