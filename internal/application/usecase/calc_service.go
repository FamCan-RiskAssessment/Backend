package usecase

import calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"

type CalcService interface {
	SendFormToCalc(request calcdto.SendFormToCalcRequest) (calcdto.ModelResponse, error)
	GetAllModelTypes() []calcdto.CalcEnumResponse
}
