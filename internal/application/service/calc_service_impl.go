package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (calcService *CalcService) SendFormToCalc(request calcdto.SendFormToCalcRequest) (calcdto.ModelResponse, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return calcdto.ModelResponse{}, notFoundError
	}
	var response calcdto.ModelResponse
	switch enum.Calc(request.CalcID) {
	case enum.CalcPremm5:
		response, err = calcService.sendFormToPremm5(form)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	case enum.CalcBCRA:
		response, err = calcService.sendFormToBCRA(form)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	}
	return response, nil
}

func (calcService *CalcService) sendFormToPremm5(form *entity.Form) (calcdto.ModelResponse, error) {
	request := calcdto.SendFormToPremm5Request{
		Sex:                    0,
		PersonalCrcOne:         1,
		PersonalCrcMultiple:    1,
		AgeCrcDx:               20,
		PersonalEndometrial:    1,
		AgeEcDx:                20,
		PersonalLsOther:        1,
		FirstDegreeCrcOne:      1,
		AgeYoungestRelativeCrc: 20,
		CurrentAge:             20,
		FamilyLsOther:          1,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.Premm5)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return calcdto.ModelResponse{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return calcdto.ModelResponse{}, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return calcdto.ModelResponse{}, err
	}

	defer resp.Body.Close()
	return calcdto.ModelResponse{
		Name:        "PREMM5",
		Probability: response["p_any"].(float64),
	}, nil
}

func (calcService *CalcService) sendFormToBCRA(form *entity.Form) (calcdto.ModelResponse, error) {
	request := calcdto.SendFormToBCRARequest{}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.BCRA)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	req.Header.Set("content-tytpe", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	defer resp.Body.Close()
	return calcdto.ModelResponse{}, nil
}

func (calcService *CalcService) GetAllModelTypes() []calcdto.CalcEnumResponse {
	modelTypes := enum.GetAllCalcs()
	response := make([]calcdto.CalcEnumResponse, len(modelTypes))

	for i, modelType := range modelTypes {
		response[i] = calcdto.CalcEnumResponse{
			ID:   uint(modelType),
			Name: modelType.String(),
		}
	}
	return response
}
