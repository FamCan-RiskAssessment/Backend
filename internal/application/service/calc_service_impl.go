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
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	cancerInfo, err := calcService.formRepository.FindCancerByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancerByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	currentAge := uint(2024 - basicInfo.BirthYear) // Adjust year as needed

	request := calcdto.SendFormToPremm5Request{
		Sex:                    mapGenderToPremm5(basicInfo.Gender),
		PersonalCrcOne:         mapPersonalCrcOne(cancerInfo),
		PersonalCrcMultiple:    mapPersonalCrcMultiple(cancerInfo),
		AgeCrcDx:               mapAgeCrcDx(cancerInfo),
		PersonalEndometrial:    mapPersonalEndometrial(cancerInfo),
		AgeEcDx:                mapAgeEcDx(cancerInfo),
		PersonalLsOther:        mapPersonalLsOther(cancerInfo),
		FirstDegreeCrcOne:      mapFirstDegreeCrcOne(familyCancerInfo),
		AgeYoungestRelativeCrc: mapAgeYoungestRelativeCrc(familyCancerInfo),
		CurrentAge:             currentAge,
		FamilyLsOther:          mapFamilyLsOther(familyCancerInfo),
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

	var premm5Response calcdto.Premm5Response
	err = json.Unmarshal(body, &premm5Response)
	if err != nil {
		fmt.Println(err)
		return calcdto.ModelResponse{}, err
	}

	defer resp.Body.Close()

	// Check if PREMM5 result already exists for this form
	existingResult, err := calcService.formRepository.FindPremm5ResultByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	// Create or update PREMM5 result
	if existingResult == nil {
		// Create new result
		premm5Result := &entity.Premm5Result{
			FormID:          form.ID,
			MLH1Probability: premm5Response.GeneProbs["MLH1"],
			MSH2Probability: premm5Response.GeneProbs["MSH2"],
			MSH6Probability: premm5Response.GeneProbs["MSH6"],
			PMS2Probability: premm5Response.GeneProbs["PMS2"],
			PAny:            premm5Response.PAny,
			PNone:           premm5Response.PNone,
		}
		err = calcService.formRepository.CreatePremm5Result(calcService.db, premm5Result)
	} else {
		// Update existing result
		existingResult.MLH1Probability = premm5Response.GeneProbs["MLH1"]
		existingResult.MSH2Probability = premm5Response.GeneProbs["MSH2"]
		existingResult.MSH6Probability = premm5Response.GeneProbs["MSH6"]
		existingResult.PMS2Probability = premm5Response.GeneProbs["PMS2"]
		existingResult.PAny = premm5Response.PAny
		existingResult.PNone = premm5Response.PNone
		err = calcService.formRepository.UpdatePremm5Result(calcService.db, existingResult)
	}

	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	return calcdto.ModelResponse{
		Name:        "PREMM5",
		Probability: premm5Response.PAny,
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

// Map gender: Male=1, Female=0
func mapGenderToPremm5(gender enum.Gender) uint {
	switch gender {
	case enum.GenderMale:
		return 1
	case enum.GenderFemale:
		return 0
	default:
		return 0 // Default to female
	}
}

// Map personal CRC history
func mapPersonalCrcOne(cancerInfo *entity.CancerInfo) uint {
	if cancerInfo != nil && cancerInfo.Cancer &&
		(cancerInfo.CancerType == nil || *cancerInfo.CancerType == enum.CancerTypeColon) {
		return 1
	}
	return 0
}

func mapPersonalCrcMultiple(cancerInfo *entity.CancerInfo) uint {
	// You'd need additional cancer history fields to determine multiple CRCs
	// For now, return 0
	return 0
}

func mapAgeCrcDx(cancerInfo *entity.CancerInfo) uint {
	if cancerInfo != nil && cancerInfo.CancerAge != nil {
		return *cancerInfo.CancerAge
	}
	return 0
}

// Map endometrial cancer
func mapPersonalEndometrial(cancerInfo *entity.CancerInfo) uint {
	if cancerInfo != nil && cancerInfo.Cancer &&
		cancerInfo.CancerType != nil && *cancerInfo.CancerType == enum.CancerTypeEndometrial {
		return 1
	}
	return 0
}

func mapAgeEcDx(cancerInfo *entity.CancerInfo) uint {
	if cancerInfo != nil && cancerInfo.CancerAge != nil &&
		cancerInfo.CancerType != nil && *cancerInfo.CancerType == enum.CancerTypeEndometrial {
		return *cancerInfo.CancerAge
	}
	return 0
}

// Map other LS-associated cancers
func mapPersonalLsOther(cancerInfo *entity.CancerInfo) uint {
	if cancerInfo != nil && cancerInfo.Cancer && cancerInfo.CancerType != nil {
		switch *cancerInfo.CancerType {
		case enum.CancerTypeOvarian, enum.CancerTypeStomach, enum.CancerTypePancreatic:
			return 1
		}
	}
	return 0
}

// Map family history
func mapFirstDegreeCrcOne(familyInfo *entity.FamilyCancerInfo) uint {
	if familyInfo != nil {
		// Check mother, father, siblings for colon cancer
		if (familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
			*familyInfo.MotherCancerType == enum.CancerTypeColon) ||
			(familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
				*familyInfo.FatherCancerType == enum.CancerTypeColon) ||
			(familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
				*familyInfo.SiblingCancerType == enum.CancerTypeColon) {
			return 1
		}
	}
	return 0
}

func mapAgeYoungestRelativeCrc(familyInfo *entity.FamilyCancerInfo) uint {
	if familyInfo != nil {
		var youngestAge uint = 999
		if familyInfo.MotherCancerAge != nil && *familyInfo.MotherCancerAge < youngestAge {
			youngestAge = *familyInfo.MotherCancerAge
		}
		if familyInfo.FatherCancerAge != nil && *familyInfo.FatherCancerAge < youngestAge {
			youngestAge = *familyInfo.FatherCancerAge
		}
		if familyInfo.SiblingCancerAge != nil && *familyInfo.SiblingCancerAge < youngestAge {
			youngestAge = *familyInfo.SiblingCancerAge
		}
		if youngestAge != 999 {
			return youngestAge
		}
	}
	return 0
}

func mapFamilyLsOther(familyInfo *entity.FamilyCancerInfo) uint {
	if familyInfo != nil {
		// Check for other LS-associated cancers in family members
		// LS-associated cancers include: endometrial, ovarian, stomach, pancreatic, brain, small bowel, hepatobiliary, urinary tract

		// Check mother
		if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil {
			switch *familyInfo.MotherCancerType {
			case enum.CancerTypeCervical, enum.CancerTypeOvarian,
				enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}

		// Check father
		if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil {
			switch *familyInfo.FatherCancerType {
			case enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}

		// Check siblings
		if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil {
			switch *familyInfo.SiblingCancerType {
			case enum.CancerTypeCervical, enum.CancerTypeOvarian,
				enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}

		// Check aunt/uncle (AmeAmo/KhaleDaei)
		if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil {
			switch *familyInfo.AmeAmoCancerType {
			case enum.CancerTypeCervical, enum.CancerTypeOvarian,
				enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}

		if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil {
			switch *familyInfo.KhaleDaeiCancerType {
			case enum.CancerTypeCervical, enum.CancerTypeOvarian,
				enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}

		// Check other relatives
		if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
			familyInfo.OtherRelativeCancerType != nil {
			switch *familyInfo.OtherRelativeCancerType {
			case enum.CancerTypeCervical, enum.CancerTypeOvarian,
				enum.CancerTypeStomach, enum.CancerTypePancreatic,
				enum.CancerTypeBrain, enum.CancerTypeLiver:
				return 1
			}
		}
	}
	return 0
}
