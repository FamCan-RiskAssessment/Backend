package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
	calcURL        *bootstrap.CalcURL
}

func NewCalcService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	db database.Database,
	calcURL *bootstrap.CalcURL,
) *CalcService {
	return &CalcService{
		constants:      constants,
		formRepository: formRepository,
		db:             db,
		calcURL:        calcURL,
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
	request := calcdto.SendFormToPremm5Request{}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.Premm5)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (calcService *CalcService) sendFormToBCRA(form *entity.Form) error {
	return nil
}

func (calcService *CalcService) sendFormToGBR(form *entity.Form) error {
	return nil
}
