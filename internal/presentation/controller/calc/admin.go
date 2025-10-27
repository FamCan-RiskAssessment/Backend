package calc

import (
	"fmt"

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
		fmt.Printf("ERROR in SendFormToCalc: %v\n", err)
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
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (calcController *AdminCalcController) GetAllModelTypes(ctx *gin.Context) {
	response := calcController.calcService.GetAllModelTypes()
	controller.Response(ctx, 200, "", response)
}
