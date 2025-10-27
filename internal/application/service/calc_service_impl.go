package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type CalcService struct {
	constants        *bootstrap.Constants
	formRepository   postgres.FormRepository
	actionLogService usecase.ActionLogService
	db               database.Database
	calcURL          *bootstrap.CalcURL
}

func NewCalcService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	actionLogService usecase.ActionLogService,
	db database.Database,
	calcURL *bootstrap.CalcURL,
) *CalcService {
	return &CalcService{
		constants:        constants,
		formRepository:   formRepository,
		actionLogService: actionLogService,
		db:               db,
		calcURL:          calcURL,
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
		response, err = calcService.sendFormToPremm5(form, request.UserID)
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

func (calcService *CalcService) GetPremm5Results(request calcdto.SendFormToCalcRequest) (calcdto.Premm5Response, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.Premm5Response{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return calcdto.Premm5Response{}, notFoundError
	}

	result, err := calcService.formRepository.FindPremm5ResultByFormID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.Premm5Response{}, err
	}

	if result == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Premm5Result}
		return calcdto.Premm5Response{}, notFoundError
	}

	result, err = calcService.formRepository.FindPremm5ResultByFormID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.Premm5Response{}, err
	}

	response := calcdto.Premm5Response{
		GeneProbs: map[string]float64{
			"MLH1": result.MLH1Probability,
			"MSH2": result.MSH2Probability,
			"MSH6": result.MSH6Probability,
			"PMS2": result.PMS2Probability,
		},
		PAny:  result.PAny,
		PNone: result.PNone,
	}

	return response, nil

}

func (calcService *CalcService) GetBCRAResults(request calcdto.SendFormToCalcRequest) (calcdto.BCRAResponse, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.BCRAResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return calcdto.BCRAResponse{}, notFoundError
	}

	result, err := calcService.formRepository.FindBCRAResultByFormID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.BCRAResponse{}, err
	}

	if result == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.BCRAResult}
		return calcdto.BCRAResponse{}, notFoundError
	}

	response := calcdto.BCRAResponse{
		AbsRisk:    result.AbsRisk,
		AbsRiskAvg: result.AbsRiskAvg,
		RRStar1:    result.RRStar1,
		RRStar2:    result.RRStar2,
		ProjIntvl:  result.ProjIntvl,
	}

	return response, nil
}

func (calcService *CalcService) sendFormToPremm5(form *entity.Form, userID uint) (calcdto.ModelResponse, error) {
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

	currentAge := uint(2024 - basicInfo.BirthYear)

	request := calcdto.SendFormToPremm5Request{
		Sex:        mapGenderToPremm5(basicInfo.Gender),
		CurrentAge: int(currentAge),

		PersonalCrcCount:       mapPersonalCrcCount(cancerInfo),
		PersonalCrcYoungestAge: mapAgeCrcDx(cancerInfo),
		PersonalEc:             mapPersonalEndometrial(cancerInfo),
		PersonalEcAge:          mapAgeEcDx(cancerInfo),
		PersonalOtherLs:        mapPersonalLsOther(cancerInfo),

		NumFdrCrc:         mapNumFdrCrc(familyCancerInfo),
		YoungestFdrCrcAge: mapYoungestFdrCrcAge(familyCancerInfo),
		NumFdrEc:          mapNumFdrEc(familyCancerInfo),
		YoungestFdrEcAge:  mapYoungestFdrEcAge(familyCancerInfo),

		NumSdrCrc:         mapNumSdrCrc(familyCancerInfo),
		YoungestSdrCrcAge: mapYoungestSdrCrcAge(familyCancerInfo),
		NumSdrEc:          mapNumSdrEc(familyCancerInfo),
		YoungestSdrEcAge:  mapYoungestSdrEcAge(familyCancerInfo),

		HasFdrOtherLs: mapHasFdrOtherLs(familyCancerInfo),
		HasSdrOtherLs: mapHasSdrOtherLs(familyCancerInfo),
	}

	// Make API call
	premm5Response, err := calcService.callPremm5API(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	// Save result to database
	err = calcService.savePremm5Result(form.ID, premm5Response)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	log := actionlogdto.LogAction{
		ActorID:    userID,
		Action:     enum.ActionTypeFormSentToPremm5,
		ResourceID: &form.ID,
	}
	calcService.actionLogService.LogAction(log)

	return calcdto.ModelResponse{
		Name:        "PREMM5",
		Probability: premm5Response.PAny,
	}, nil
}

func (calcService *CalcService) sendFormToBCRA(form *entity.Form) (calcdto.ModelResponse, error) {
	// Load all required data
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	mamographyInfo, err := calcService.formRepository.FindMamographyByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	if mamographyInfo == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.MamoGraphyInfo}
		return calcdto.ModelResponse{}, notFoundError
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancerByFormID(calcService.db, form.ID)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	if familyCancerInfo == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.FamilyCancerInfo}
		return calcdto.ModelResponse{}, notFoundError
	}

	// Calculate current age
	currentAge := float64(2024 - basicInfo.BirthYear)
	projectionAge := currentAge + 5.0

	// Build BCRA request
	request := calcdto.SendFormToBCRARequest{
		T1:      currentAge,
		T2:      projectionAge,
		N_Biop:  mapBiopsyCount(mamographyInfo),
		HypPlas: mapHyperplasiaStatus(mamographyInfo),
		AgeMen:  mapAgeAtMenarche(mamographyInfo),
		Age1st:  mapAgeAtFirstBirth(mamographyInfo),
		N_Rels:  mapFirstDegreeBreastCancerRelatives(familyCancerInfo),
		Race:    1, // Default to White
	}

	// Make API call
	bcraResponse, err := calcService.callBCRAAPI(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	// Save result to database
	err = calcService.saveBCRAResult(form.ID, bcraResponse)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	return calcdto.ModelResponse{
		Name:        "BCRA",
		Probability: bcraResponse.AbsRisk,
	}, nil
}

// callBCRAAPI makes the HTTP request to the BCRA calculator API
func (calcService *CalcService) callBCRAAPI(request calcdto.SendFormToBCRARequest) (calcdto.BCRAResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return calcdto.BCRAResponse{}, err
	}

	// Log the request for debugging
	fmt.Printf("BCRA Request: T1=%.1f, T2=%.1f, N_Biop=%d, HypPlas=%d, AgeMen=%d, Age1st=%d, N_Rels=%d, Race=%d\n",
		request.T1, request.T2, request.N_Biop, request.HypPlas, request.AgeMen, request.Age1st, request.N_Rels, request.Race)

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.BCRA)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return calcdto.BCRAResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return calcdto.BCRAResponse{}, fmt.Errorf("failed to call BCRA API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return calcdto.BCRAResponse{}, fmt.Errorf("failed to read BCRA response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return calcdto.BCRAResponse{}, fmt.Errorf("BCRA API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var bcraResponse calcdto.BCRAResponse
	err = json.Unmarshal(body, &bcraResponse)
	if err != nil {
		return calcdto.BCRAResponse{}, fmt.Errorf("failed to parse BCRA response: %w", err)
	}

	return bcraResponse, nil
}

// saveBCRAResult creates or updates the BCRA result in the database
func (calcService *CalcService) saveBCRAResult(formID uint, bcraResponse calcdto.BCRAResponse) error {
	existingResult, err := calcService.formRepository.FindBCRAResultByFormID(calcService.db, formID)
	if err != nil {
		return err
	}

	if existingResult == nil {
		bcraResult := &entity.BCRAResult{
			FormID:     formID,
			AbsRisk:    bcraResponse.AbsRisk,
			AbsRiskAvg: bcraResponse.AbsRiskAvg,
			RRStar1:    bcraResponse.RRStar1,
			RRStar2:    bcraResponse.RRStar2,
			ProjIntvl:  bcraResponse.ProjIntvl,
		}
		return calcService.formRepository.CreateBCRAResult(calcService.db, bcraResult)
	}

	existingResult.AbsRisk = bcraResponse.AbsRisk
	existingResult.AbsRiskAvg = bcraResponse.AbsRiskAvg
	existingResult.RRStar1 = bcraResponse.RRStar1
	existingResult.RRStar2 = bcraResponse.RRStar2
	existingResult.ProjIntvl = bcraResponse.ProjIntvl

	return calcService.formRepository.UpdateBCRAResult(calcService.db, existingResult)
}

// mapHyperplasiaStatus maps hyperplasia status to BCRA code
func mapHyperplasiaStatus(mamographyInfo *entity.MamoGraphyInfo) int {
	if mamographyInfo.HyperplasiaInBiopsy == nil {
		return 99 // Unknown
	}

	switch *mamographyInfo.HyperplasiaInBiopsy {
	case enum.NoHyperplasiaInBiopsy:
		return 0
	case enum.HasHyperplasiaInBiopsy:
		return 1
	case enum.UnkownHyperplasiaInBiopsy:
		return 99
	default:
		return 99
	}
}

// mapBiopsyCount maps number of biopsies to BCRA code
func mapBiopsyCount(mamographyInfo *entity.MamoGraphyInfo) int {
	if mamographyInfo.NumberOfBreastBiopsies == nil {
		return 99 // Unknown
	}

	count := *mamographyInfo.NumberOfBreastBiopsies
	if count >= 2 {
		return 2
	} else if count == 1 {
		return 1
	}
	return 0
}

// mapAgeAtMenarche maps GhaedeAge to BCRA AgeMen
func mapAgeAtMenarche(mamographyInfo *entity.MamoGraphyInfo) int {
	if mamographyInfo.GhaedeAge > 0 {
		return int(mamographyInfo.GhaedeAge)
	}
	return 99 // Unknown
}

// mapAgeAtFirstBirth maps age at first birth to BCRA code with validation
func mapAgeAtFirstBirth(mamographyInfo *entity.MamoGraphyInfo) int {
	// First check if patient has children
	if !mamographyInfo.HasChildren {
		return 98 // Nulliparous
	}

	// If has children but no age of first birth provided, return unknown
	if mamographyInfo.AgeOfFirstBirth == nil {
		return 99 // Unknown
	}

	age := int(*mamographyInfo.AgeOfFirstBirth)
	ageMen := mapAgeAtMenarche(mamographyInfo)

	// Validate: first birth age must be >= menarche age
	if ageMen != 99 && age < ageMen {
		return 99 // Unknown due to invalid data
	}

	// Map age to category
	var category int
	if age < 20 {
		category = 0
	} else if age < 25 {
		category = 1
	} else if age < 30 {
		category = 2
	} else {
		category = 3
	}

	// Additional validation for Age1st=0 category
	// BCRA has issues with Age1st=0 when combined with certain AgeMen values
	// that can cause calculation to result in NaN
	if category == 0 && ageMen != 99 {
		// If menarche age is >= 12 with first birth before 20, it can cause issues
		// Return 99 (unknown) as a conservative approach
		if ageMen >= 12 {
			return 99
		}
	}

	// Validate the age is not absurdly high
	if age > 60 {
		return 99
	}

	return category
}

// mapFirstDegreeBreastCancerRelatives counts first-degree relatives with breast cancer
func mapFirstDegreeBreastCancerRelatives(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 99 // Unknown
	}

	count := 0

	// Check mother
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
		*familyInfo.MotherCancerType == enum.CancerTypeBreast {
		count++
	}

	// Check father
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
		*familyInfo.FatherCancerType == enum.CancerTypeBreast {
		count++
	}

	// Check siblings
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
		*familyInfo.SiblingCancerType == enum.CancerTypeBreast {
		count++
	}

	// Return categorized count
	if count >= 2 {
		return 2
	} else if count == 1 {
		return 1
	} else if count == 0 {
		return 0
	}

	return 99 // Unknown
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
func mapGenderToPremm5(gender enum.Gender) int {
	switch gender {
	case enum.GenderMale:
		return 1
	case enum.GenderFemale:
		return 0
	default:
		return 0 // Default to female
	}
}

// Map personal CRC count (0, 1, or 2+)
func mapPersonalCrcCount(cancerInfo *entity.CancerInfo) int {
	if cancerInfo != nil && cancerInfo.Cancer &&
		(cancerInfo.CancerType == nil || *cancerInfo.CancerType == enum.CancerTypeColon) {
		// TODO: You may need additional field to track multiple CRCs
		// For now, assume single CRC = 1, multiple = 2
		return 1
	}
	return 0
}

// Map count of first-degree relatives with CRC (0, 1, or 2+)
func mapNumFdrCrc(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	count := 0
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
		*familyInfo.MotherCancerType == enum.CancerTypeColon {
		count++
	}
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
		*familyInfo.FatherCancerType == enum.CancerTypeColon {
		count++
	}
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
		*familyInfo.SiblingCancerType == enum.CancerTypeColon {
		count++
	}
	// Cap at 2 for >=2
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among FDR with CRC
func mapYoungestFdrCrcAge(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	youngestAge := 999
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
		*familyInfo.MotherCancerType == enum.CancerTypeColon &&
		familyInfo.MotherCancerAge != nil {
		age := int(*familyInfo.MotherCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
		*familyInfo.FatherCancerType == enum.CancerTypeColon &&
		familyInfo.FatherCancerAge != nil {
		age := int(*familyInfo.FatherCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
		*familyInfo.SiblingCancerType == enum.CancerTypeColon &&
		familyInfo.SiblingCancerAge != nil {
		age := int(*familyInfo.SiblingCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of FDR with endometrial cancer
func mapNumFdrEc(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	count := 0
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
		*familyInfo.MotherCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
		*familyInfo.FatherCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
		*familyInfo.SiblingCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among FDR with endometrial cancer
func mapYoungestFdrEcAge(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	youngestAge := 999
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil &&
		*familyInfo.MotherCancerType == enum.CancerTypeEndometrial &&
		familyInfo.MotherCancerAge != nil {
		age := int(*familyInfo.MotherCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil &&
		*familyInfo.FatherCancerType == enum.CancerTypeEndometrial &&
		familyInfo.FatherCancerAge != nil {
		age := int(*familyInfo.FatherCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil &&
		*familyInfo.SiblingCancerType == enum.CancerTypeEndometrial &&
		familyInfo.SiblingCancerAge != nil {
		age := int(*familyInfo.SiblingCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of second-degree relatives with CRC (0, 1, or 2+)
func mapNumSdrCrc(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	count := 0
	if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil &&
		*familyInfo.AmeAmoCancerType == enum.CancerTypeColon {
		count++
	}
	if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil &&
		*familyInfo.KhaleDaeiCancerType == enum.CancerTypeColon {
		count++
	}
	if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
		familyInfo.OtherRelativeCancerType != nil &&
		*familyInfo.OtherRelativeCancerType == enum.CancerTypeColon {
		count++
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among SDR with CRC
func mapYoungestSdrCrcAge(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	youngestAge := 999
	if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil &&
		*familyInfo.AmeAmoCancerType == enum.CancerTypeColon &&
		familyInfo.AmeAmoCancerAge != nil {
		age := int(*familyInfo.AmeAmoCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil &&
		*familyInfo.KhaleDaeiCancerType == enum.CancerTypeColon &&
		familyInfo.KhaleDaeiCancerAge != nil {
		age := int(*familyInfo.KhaleDaeiCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
		familyInfo.OtherRelativeCancerType != nil &&
		*familyInfo.OtherRelativeCancerType == enum.CancerTypeColon &&
		familyInfo.OtherRelativeCancerAge != nil {
		age := int(*familyInfo.OtherRelativeCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of SDR with endometrial cancer
func mapNumSdrEc(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	count := 0
	if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil &&
		*familyInfo.AmeAmoCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil &&
		*familyInfo.KhaleDaeiCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
		familyInfo.OtherRelativeCancerType != nil &&
		*familyInfo.OtherRelativeCancerType == enum.CancerTypeEndometrial {
		count++
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among SDR with endometrial cancer
func mapYoungestSdrEcAge(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	youngestAge := 999
	if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil &&
		*familyInfo.AmeAmoCancerType == enum.CancerTypeEndometrial &&
		familyInfo.AmeAmoCancerAge != nil {
		age := int(*familyInfo.AmeAmoCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil &&
		*familyInfo.KhaleDaeiCancerType == enum.CancerTypeEndometrial &&
		familyInfo.KhaleDaeiCancerAge != nil {
		age := int(*familyInfo.KhaleDaeiCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
		familyInfo.OtherRelativeCancerType != nil &&
		*familyInfo.OtherRelativeCancerType == enum.CancerTypeEndometrial &&
		familyInfo.OtherRelativeCancerAge != nil {
		age := int(*familyInfo.OtherRelativeCancerAge)
		if age < youngestAge {
			youngestAge = age
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map if FDR has other Lynch syndrome cancers
func mapHasFdrOtherLs(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	// Check mother, father, siblings for other LS cancers
	lscancers := []enum.CancerType{
		enum.CancerTypeOvarian, enum.CancerTypeStomach,
		enum.CancerTypePancreatic, enum.CancerTypeBrain,
		enum.CancerTypeLiver,
	}

	// Check mother
	if familyInfo.MotherCancer && familyInfo.MotherCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.MotherCancerType == lsc {
				return 1
			}
		}
	}

	// Check father
	if familyInfo.FatherCancer && familyInfo.FatherCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.FatherCancerType == lsc {
				return 1
			}
		}
	}

	// Check siblings
	if familyInfo.SiblingCancer && familyInfo.SiblingCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.SiblingCancerType == lsc {
				return 1
			}
		}
	}

	return 0
}

// Map if SDR has other Lynch syndrome cancers
func mapHasSdrOtherLs(familyInfo *entity.FamilyCancerInfo) int {
	if familyInfo == nil {
		return 0
	}
	// Check aunt/uncle for other LS cancers
	lscancers := []enum.CancerType{
		enum.CancerTypeOvarian, enum.CancerTypeStomach,
		enum.CancerTypePancreatic, enum.CancerTypeBrain,
		enum.CancerTypeLiver,
	}

	// Check AmeAmo
	if familyInfo.AmeAmoCancer && familyInfo.AmeAmoCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.AmeAmoCancerType == lsc {
				return 1
			}
		}
	}

	// Check KhaleDaei
	if familyInfo.KhaleDaeiCancer && familyInfo.KhaleDaeiCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.KhaleDaeiCancerType == lsc {
				return 1
			}
		}
	}

	// Check other relatives
	if familyInfo.OtherRelativeCancer != nil && *familyInfo.OtherRelativeCancer &&
		familyInfo.OtherRelativeCancerType != nil {
		for _, lsc := range lscancers {
			if *familyInfo.OtherRelativeCancerType == lsc {
				return 1
			}
		}
	}

	return 0
}

// Map youngest age at CRC diagnosis
func mapAgeCrcDx(cancerInfo *entity.CancerInfo) int {
	if cancerInfo != nil && cancerInfo.CancerAge != nil &&
		(cancerInfo.CancerType == nil || *cancerInfo.CancerType == enum.CancerTypeColon) {
		return int(*cancerInfo.CancerAge)
	}
	return 0
}

// Map endometrial cancer
func mapPersonalEndometrial(cancerInfo *entity.CancerInfo) int {
	if cancerInfo != nil && cancerInfo.Cancer &&
		cancerInfo.CancerType != nil && *cancerInfo.CancerType == enum.CancerTypeEndometrial {
		return 1
	}
	return 0
}

// Map age at endometrial cancer diagnosis
func mapAgeEcDx(cancerInfo *entity.CancerInfo) int {
	if cancerInfo != nil && cancerInfo.CancerAge != nil &&
		cancerInfo.CancerType != nil && *cancerInfo.CancerType == enum.CancerTypeEndometrial {
		return int(*cancerInfo.CancerAge)
	}
	return 0
}

// Map other LS-associated cancers
func mapPersonalLsOther(cancerInfo *entity.CancerInfo) int {
	if cancerInfo != nil && cancerInfo.Cancer && cancerInfo.CancerType != nil {
		switch *cancerInfo.CancerType {
		case enum.CancerTypeOvarian, enum.CancerTypeStomach, enum.CancerTypePancreatic:
			return 1
		}
	}
	return 0
}

// callPremm5API makes the HTTP request to the PREMM5 calculator API
func (calcService *CalcService) callPremm5API(request calcdto.SendFormToPremm5Request) (calcdto.Premm5Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return calcdto.Premm5Response{}, fmt.Errorf("failed to marshal PREMM5 request: %w", err)
	}

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.Premm5)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return calcdto.Premm5Response{}, fmt.Errorf("failed to create PREMM5 request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return calcdto.Premm5Response{}, fmt.Errorf("failed to call PREMM5 API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return calcdto.Premm5Response{}, fmt.Errorf("failed to read PREMM5 response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return calcdto.Premm5Response{}, fmt.Errorf("PREMM5 API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var premm5Response calcdto.Premm5Response
	err = json.Unmarshal(body, &premm5Response)
	if err != nil {
		return calcdto.Premm5Response{}, fmt.Errorf("failed to parse PREMM5 response: %w", err)
	}

	return premm5Response, nil
}

// savePremm5Result creates or updates the PREMM5 result in the database
func (calcService *CalcService) savePremm5Result(formID uint, premm5Response calcdto.Premm5Response) error {
	existingResult, err := calcService.formRepository.FindPremm5ResultByFormID(calcService.db, formID)
	if err != nil {
		return err
	}

	if existingResult == nil {
		premm5Result := &entity.Premm5Result{
			FormID:          formID,
			MLH1Probability: premm5Response.GeneProbs["MLH1"],
			MSH2Probability: premm5Response.GeneProbs["MSH2"],
			MSH6Probability: premm5Response.GeneProbs["MSH6"],
			PMS2Probability: premm5Response.GeneProbs["PMS2"],
			PAny:            premm5Response.PAny,
			PNone:           premm5Response.PNone,
		}
		return calcService.formRepository.CreatePremm5Result(calcService.db, premm5Result)
	}

	existingResult.MLH1Probability = premm5Response.GeneProbs["MLH1"]
	existingResult.MSH2Probability = premm5Response.GeneProbs["MSH2"]
	existingResult.MSH6Probability = premm5Response.GeneProbs["MSH6"]
	existingResult.PMS2Probability = premm5Response.GeneProbs["PMS2"]
	existingResult.PAny = premm5Response.PAny
	existingResult.PNone = premm5Response.PNone

	return calcService.formRepository.UpdatePremm5Result(calcService.db, existingResult)
}
