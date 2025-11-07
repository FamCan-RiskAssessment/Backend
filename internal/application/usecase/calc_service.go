package usecase

import calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"

type CalcService interface {
	SendFormToCalc(request calcdto.SendFormToCalcRequest) (calcdto.ModelResponse, error)
	GetPremm5Results(request calcdto.SendFormToCalcRequest) (calcdto.Premm5Response, error)
	GetBCRAResults(request calcdto.SendFormToCalcRequest) (calcdto.BCRAResponse, error)
	GetGailResults(request calcdto.SendFormToCalcRequest) (calcdto.GailResponse, error)
	GetAllModelTypes() []calcdto.CalcEnumResponse
}
