package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type CalcService struct {
	constants      *bootstrap.Constants
	formRepository postgres.FormRepository
	db             database.Database
}

func NewCalcService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	db database.Database,
) *CalcService {
	return &CalcService{
		constants:      constants,
		formRepository: formRepository,
		db:             db,
	}
}

func (calcService *CalcService) SendFormToCalc(request calcdto.SendFormToCalcRequest) error {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return notFoundError
	}

	switch enum.Calc(request.CalcID) {
	case enum.CalcPremm5:
		err = calcService.sendFormToPremm5(form)
		if err != nil {
			return err
		}
	case enum.CalcBCRA:
		err = calcService.sendFormToBCRA(form)
		if err != nil {
			return err
		}
	case enum.CalcGBR:
		err = calcService.sendFormToGBR(form)
		if err != nil {
			return err
		}
	}
	return nil
}

func (calcService *CalcService) sendFormToPremm5(form *entity.Form) error {
	return nil
}

func (calcService *CalcService) sendFormToBCRA(form *entity.Form) error {
	return nil
}

func (calcService *CalcService) sendFormToGBR(form *entity.Form) error {
	return nil
}
