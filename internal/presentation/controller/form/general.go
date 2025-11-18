package form

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralFormController struct {
	formService usecase.FormService
}

func NewGeneralFormController(
	formService usecase.FormService,
) *GeneralFormController {
	return &GeneralFormController{
		formService: formService,
	}
}

func (fc *GeneralFormController) GetAllCancerTypes(ctx *gin.Context) {
	types, err := fc.formService.GetAllCancerTypes()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", types)
}

func (fc *GeneralFormController) GetAllGenders(ctx *gin.Context) {
	types, err := fc.formService.GetAllGenders()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", types)
}

func (fc *GeneralFormController) GetAllMenopausalStatuses(ctx *gin.Context) {
	types, err := fc.formService.GetAllMenopausalStatuses()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", types)
}

func (fc *GeneralFormController) GetAllFormStatuses(ctx *gin.Context) {
	types, err := fc.formService.GetAllFormStatuses()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", types)
}

func (fc *GeneralFormController) GetAllHyperplasiaInBiopsyStatuses(ctx *gin.Context) {
	statuses, err := fc.formService.GetAllHyperplasiaInBiopsyStatuses()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", statuses)
}

func (fc *GeneralFormController) GetAllLifeStatuses(ctx *gin.Context) {
	statuses, err := fc.formService.GetAllLifeStatuses()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", statuses)
}

func (fc *GeneralFormController) GetAllRelativeTypes(ctx *gin.Context) {
	statuses, err := fc.formService.GetAllRelativeTypes()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", statuses)
}

func (fc *GeneralFormController) GetAllDefiniteAnswers(ctx *gin.Context) {
	statuses, err := fc.formService.GetAllDefiniteAnswers()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", statuses)
}
