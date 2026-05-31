package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	calcdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	ptime "github.com/yaa110/go-persian-calendar"
)

type CalcService struct {
	constants        *bootstrap.Constants
	formRepository   postgres.FormRepository
	formService      usecase.FormService
	actionLogService usecase.ActionLogService
	db               database.Database
	calcURL          *bootstrap.CalcURL
}

func NewCalcService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	formService usecase.FormService,
	actionLogService usecase.ActionLogService,
	db database.Database,
	calcURL *bootstrap.CalcURL,
) *CalcService {
	return &CalcService{
		constants:        constants,
		formRepository:   formRepository,
		formService:      formService,
		actionLogService: actionLogService,
		db:               db,
		calcURL:          calcURL,
	}
}

func calculateAge(birthDate time.Time) int {
	// Convert to Iranian calendar (Jalali)
	today := ptime.Now()
	birthDateJalali := ptime.New(birthDate)

	age := today.Year() - birthDateJalali.Year()
	if today.Month() < birthDateJalali.Month() ||
		(today.Month() == birthDateJalali.Month() && today.Day() < birthDateJalali.Day()) {
		age--
	}
	return age
}

func (calcService *CalcService) SendFormToCalc(request calcdto.SendFormToCalcRequest) (calcdto.ModelResponse, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		log.Printf("[CALC] Database error finding form ID %d: %v", request.FormID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Message: "failed to fetch form from database",
			OrigErr: err,
		}
	}
	if form == nil {
		log.Printf("[CALC] Form not found - FormID: %d, UserID: %d", request.FormID, request.UserID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: calcService.constants.Field.Form}
	}

	// Validate form is in ReadyForCalculation status
	if form.Status != enum.FormStatusReadyForCalculation && form.Status != enum.FormStatusCalculated {
		log.Printf("[CALC] Form not ready for calculation - FormID: %d, CurrentStatus: %s, UserID: %d", request.FormID, form.Status.String(), request.UserID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Message: fmt.Sprintf("form must be in ReadyForCalculation status, current status: %s", form.Status.String()),
		}
	}

	// Validate calculation ID
	validCalcID := false
	var modelName string
	switch enum.Calc(request.CalcID) {
	case enum.CalcPremm5:
		validCalcID = true
		modelName = "PREMM5"
	case enum.CalcBCRA:
		validCalcID = true
		modelName = "BCRA"
	case enum.CalcGail:
		validCalcID = true
		modelName = "Gail"
	case enum.CalcPLCO:
		validCalcID = true
		modelName = "PLCO"
	}

	if !validCalcID {
		log.Printf("[CALC] Invalid calculation model - CalcID: %d, UserID: %d", request.CalcID, request.UserID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Message: fmt.Sprintf("invalid calculation model ID: %d", request.CalcID),
		}
	}

	var response calcdto.ModelResponse
	switch enum.Calc(request.CalcID) {
	case enum.CalcPremm5:
		response, err = calcService.sendFormToPremm5(form, request.UserID)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	case enum.CalcBCRA:
		response, err = calcService.sendFormToBCRA(form, request.UserID)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	case enum.CalcGail:
		response, err = calcService.sendFormToGail(form, request.UserID)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	case enum.CalcPLCO:
		response, err = calcService.sendFormToPLCO(form, request.UserID)
		if err != nil {
			return calcdto.ModelResponse{}, err
		}
	}

	// Update form status to calculated
	if form.Status != enum.FormStatusCalculated {
		form.Status = enum.FormStatusCalculated
		err = calcService.formRepository.UpdateForm(calcService.db, form)
		if err != nil {
			log.Printf("[CALC] Database error updating form status - FormID: %d, UserID: %d, Error: %v", request.FormID, request.UserID, err)
			return calcdto.ModelResponse{}, &exception.CalcError{
				Type:    exception.ErrorTypeDatabaseError,
				Model:   modelName,
				Message: "failed to update form status in database",
				OrigErr: err,
			}
		}
	}

	log.Printf("[CALC] Form successfully sent to %s - FormID: %d, UserID: %d", modelName, request.FormID, request.UserID)
	return response, nil
}

func (calcService *CalcService) sendFormToPremm5(form *entity.Form, userID uint) (calcdto.ModelResponse, error) {
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PREMM5] Database error fetching basic info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PREMM5",
			Message: "failed to fetch basic information",
			OrigErr: err,
		}
	}
	if basicInfo == nil {
		log.Printf("[CALC:PREMM5] Missing basic info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeMissingData,
			Model:   "PREMM5",
			Message: "basic information not found in form",
		}
	}

	cancerInfo, err := calcService.formRepository.FindCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PREMM5] Database error fetching cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PREMM5",
			Message: "failed to fetch cancer information",
			OrigErr: err,
		}
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PREMM5] Database error fetching family cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PREMM5",
			Message: "failed to fetch family cancer information",
			OrigErr: err,
		}
	}

	currentAge := calculateAge(basicInfo.BirthDate)

	// Validate age is within acceptable range
	if currentAge < 0 || currentAge > 120 {
		log.Printf("[CALC:PREMM5] Invalid age calculated - FormID: %d, Age: %d, BirthDate: %v", form.ID, currentAge, basicInfo.BirthDate)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "PREMM5",
			Message: fmt.Sprintf("calculated age %d is outside valid range (0-120)", currentAge),
		}
	}

	request := calcdto.SendFormToPremm5Request{
		Sex:        mapGenderToPremm5(basicInfo.Gender),
		CurrentAge: currentAge,

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
		log.Printf("[CALC:PREMM5] Database error saving result - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PREMM5",
			Message: "failed to save calculation result",
			OrigErr: err,
		}
	}

	actionLog := actionlogdto.LogAction{
		ActorID:    userID,
		Action:     enum.ActionTypeFormSentToPremm5,
		ResourceID: &form.ID,
	}
	calcService.actionLogService.LogAction(actionLog)

	log.Printf("[CALC:PREMM5] Calculation successful - FormID: %d, Probability: %.4f", form.ID, premm5Response.PAny)
	return calcdto.ModelResponse{
		Name:        "PREMM5",
		Probability: premm5Response.PAny,
	}, nil
}

func (calcService *CalcService) sendFormToBCRA(form *entity.Form, userID uint) (calcdto.ModelResponse, error) {
	// Load all required data
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:BCRA] Database error fetching basic info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "BCRA",
			Message: "failed to fetch basic information",
			OrigErr: err,
		}
	}
	if basicInfo == nil {
		log.Printf("[CALC:BCRA] Missing basic info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeMissingData,
			Model:   "BCRA",
			Message: "basic information not found in form",
		}
	}

	mamographyInfo, err := calcService.formRepository.FindMamographyByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:BCRA] Database error fetching mamography info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "BCRA",
			Message: "failed to fetch mamography information",
			OrigErr: err,
		}
	}

	if mamographyInfo == nil {
		log.Printf("[CALC:BCRA] Missing mamography info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: calcService.constants.Field.MamoGraphyInfo}
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:BCRA] Database error fetching family cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "BCRA",
			Message: "failed to fetch family cancer information",
			OrigErr: err,
		}
	}

	if len(familyCancerInfo) == 0 {
		log.Printf("[CALC:BCRA] Missing family cancer info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: calcService.constants.Field.FamilyCancerInfo}
	}

	// Calculate current age
	currentAge := float64(calculateAge(basicInfo.BirthDate))

	// Validate age
	if currentAge < 0 || currentAge > 120 {
		log.Printf("[CALC:BCRA] Invalid age - FormID: %d, Age: %.1f", form.ID, currentAge)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "BCRA",
			Message: fmt.Sprintf("calculated age %.1f is outside valid range", currentAge),
		}
	}

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
		log.Printf("[CALC:BCRA] Database error saving result - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "BCRA",
			Message: "failed to save calculation result",
			OrigErr: err,
		}
	}

	actionLog := actionlogdto.LogAction{
		ActorID:    userID,
		Action:     enum.ActionTypeFormSentToBCRA,
		ResourceID: &form.ID,
	}
	calcService.actionLogService.LogAction(actionLog)

	log.Printf("[CALC:BCRA] Calculation successful - FormID: %d, AbsRisk: %.4f", form.ID, bcraResponse.AbsRisk)
	return calcdto.ModelResponse{
		Name:        "BCRA",
		Probability: bcraResponse.AbsRisk,
	}, nil
}

func (calcService *CalcService) sendFormToGail(form *entity.Form, userID uint) (calcdto.ModelResponse, error) {
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:Gail] Database error fetching basic info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "Gail",
			Message: "failed to fetch basic information",
			OrigErr: err,
		}
	}
	if basicInfo == nil {
		log.Printf("[CALC:Gail] Missing basic info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeMissingData,
			Model:   "Gail",
			Message: "basic information not found in form",
		}
	}

	mamographyInfo, err := calcService.formRepository.FindMamographyByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:Gail] Database error fetching mamography info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "Gail",
			Message: "failed to fetch mamography information",
			OrigErr: err,
		}
	}

	if mamographyInfo == nil {
		log.Printf("[CALC:Gail] Missing mamography info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: calcService.constants.Field.MamoGraphyInfo}
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:Gail] Database error fetching family cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "Gail",
			Message: "failed to fetch family cancer information",
			OrigErr: err,
		}
	}

	if len(familyCancerInfo) == 0 {
		log.Printf("[CALC:Gail] Missing family cancer info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: calcService.constants.Field.FamilyCancerInfo}
	}

	// Calculate current age
	currentAge := calculateAge(basicInfo.BirthDate)

	// Validate age
	if currentAge < 0 || currentAge > 120 {
		log.Printf("[CALC:Gail] Invalid age - FormID: %d, Age: %d", form.ID, currentAge)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "Gail",
			Message: fmt.Sprintf("calculated age %d is outside valid range (0-120)", currentAge),
		}
	}

	// Build Gail request
	request := calcdto.SendFormToGailRequest{
		Age:          currentAge,
		HorizonYears: 5,
		MenarcheAge:  mapAgeAtMenarche(mamographyInfo),
		NumBiopsies:  mapBiopsyCount(mamographyInfo),
		FLBAge:       mapAgeAtFirstBirth(mamographyInfo),
		NumRelatives: mapFirstDegreeBreastCancerRelatives(familyCancerInfo),
		Race:         "white", // Default to White
		ShowRR:       false,
	}

	// Make API call
	gailResponse, err := calcService.callGailAPI(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	// Save result to database
	err = calcService.saveGailResult(form.ID, gailResponse)
	if err != nil {
		log.Printf("[CALC:Gail] Database error saving result - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "Gail",
			Message: "failed to save calculation result",
			OrigErr: err,
		}
	}

	actionLog := actionlogdto.LogAction{
		ActorID:    userID,
		Action:     enum.ActionTypeFormSentToGail,
		ResourceID: &form.ID,
	}
	calcService.actionLogService.LogAction(actionLog)

	log.Printf("[CALC:Gail] Calculation successful - FormID: %d, AbsoluteRisk: %.4f", form.ID, gailResponse.AbsoluteRisk)
	return calcdto.ModelResponse{
		Name:        "Gail",
		Probability: gailResponse.AbsoluteRisk,
	}, nil
}

func (calcService *CalcService) GetCalcBrowse(request calcdto.GetCalcBrowseRequest) ([]calcdto.CalcBrowseItemResponse, int64, error) {
	forms, count, err := calcService.formService.GetAllForms(
		request.Offset,
		request.Limit,
		request.Filters,
		request.SortBy,
		request.SortOrder,
		request.Search,
	)
	if err != nil {
		return nil, 0, err
	}

	if len(forms) == 0 {
		return []calcdto.CalcBrowseItemResponse{}, count, nil
	}

	formIDs := make([]uint, len(forms))
	for i, form := range forms {
		formIDs[i] = form.FormID
	}

	premm5ByFormID, err := calcService.loadPremm5ResultsByFormID(formIDs)
	if err != nil {
		return nil, 0, err
	}
	bcraByFormID, err := calcService.loadBCRAResultsByFormID(formIDs)
	if err != nil {
		return nil, 0, err
	}
	gailByFormID, err := calcService.loadGailResultsByFormID(formIDs)
	if err != nil {
		return nil, 0, err
	}
	plcoByFormID, err := calcService.loadPLCOResultsByFormID(formIDs)
	if err != nil {
		return nil, 0, err
	}

	items := make([]calcdto.CalcBrowseItemResponse, len(forms))
	for i, form := range forms {
		items[i] = calcdto.CalcBrowseItemResponse{
			Form: form,
			Results: calcdto.CalcResultsBundle{
				Premm5: mapPremm5Result(premm5ByFormID[form.FormID]),
				BCRA:   mapBCRAResult(bcraByFormID[form.FormID]),
				Gail:   mapGailResult(gailByFormID[form.FormID]),
				PLCO:   mapPLCOResult(plcoByFormID[form.FormID]),
			},
		}
	}

	return items, count, nil
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
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.Premm5Response{}, notFoundError
	}

	response := mapPremm5Result(result)
	if response == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.Premm5Response{}, notFoundError
	}

	return *response, nil

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
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.BCRAResponse{}, notFoundError
	}

	response := mapBCRAResult(result)
	if response == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.BCRAResponse{}, notFoundError
	}

	return *response, nil
}

func (calcService *CalcService) GetGailResults(request calcdto.SendFormToCalcRequest) (calcdto.GailResponse, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.GailResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return calcdto.GailResponse{}, notFoundError
	}

	result, err := calcService.formRepository.FindGailResultByFormID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.GailResponse{}, err
	}

	if result == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.GailResponse{}, notFoundError
	}

	response := mapGailResult(result)
	if response == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.GailResponse{}, notFoundError
	}

	return *response, nil
}

func (calcService *CalcService) GetPLCOResults(request calcdto.SendFormToCalcRequest) (calcdto.PLCOResponse, error) {
	form, err := calcService.formRepository.FindFormByID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.PLCOResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Form}
		return calcdto.PLCOResponse{}, notFoundError
	}

	result, err := calcService.formRepository.FindPLCOResultByFormID(calcService.db, request.FormID)
	if err != nil {
		return calcdto.PLCOResponse{}, err
	}

	if result == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.PLCOResponse{}, notFoundError
	}

	response := mapPLCOResult(result)
	if response == nil {
		notFoundError := exception.NotFoundError{Item: calcService.constants.Field.Record}
		return calcdto.PLCOResponse{}, notFoundError
	}

	return *response, nil
}

func (calcService *CalcService) loadPremm5ResultsByFormID(formIDs []uint) (map[uint]*entity.Premm5Result, error) {
	results, err := calcService.formRepository.FindPremm5ResultsByFormIDs(calcService.db, formIDs)
	if err != nil {
		return nil, err
	}
	byFormID := make(map[uint]*entity.Premm5Result, len(results))
	for i := range results {
		byFormID[results[i].FormID] = &results[i]
	}
	return byFormID, nil
}

func (calcService *CalcService) loadBCRAResultsByFormID(formIDs []uint) (map[uint]*entity.BCRAResult, error) {
	results, err := calcService.formRepository.FindBCRAResultsByFormIDs(calcService.db, formIDs)
	if err != nil {
		return nil, err
	}
	byFormID := make(map[uint]*entity.BCRAResult, len(results))
	for i := range results {
		byFormID[results[i].FormID] = &results[i]
	}
	return byFormID, nil
}

func (calcService *CalcService) loadGailResultsByFormID(formIDs []uint) (map[uint]*entity.GailResult, error) {
	results, err := calcService.formRepository.FindGailResultsByFormIDs(calcService.db, formIDs)
	if err != nil {
		return nil, err
	}
	byFormID := make(map[uint]*entity.GailResult, len(results))
	for i := range results {
		byFormID[results[i].FormID] = &results[i]
	}
	return byFormID, nil
}

func (calcService *CalcService) loadPLCOResultsByFormID(formIDs []uint) (map[uint]*entity.PLCOResult, error) {
	results, err := calcService.formRepository.FindPLCOResultsByFormIDs(calcService.db, formIDs)
	if err != nil {
		return nil, err
	}
	byFormID := make(map[uint]*entity.PLCOResult, len(results))
	for i := range results {
		byFormID[results[i].FormID] = &results[i]
	}
	return byFormID, nil
}

func mapPremm5Result(result *entity.Premm5Result) *calcdto.Premm5Response {
	if result == nil {
		return nil
	}
	return &calcdto.Premm5Response{
		GeneProbs: map[string]float64{
			"MLH1": result.MLH1Probability,
			"MSH2": result.MSH2Probability,
			"MSH6": result.MSH6Probability,
			"PMS2": result.PMS2Probability,
		},
		PAny:  result.PAny,
		PNone: result.PNone,
	}
}

func mapBCRAResult(result *entity.BCRAResult) *calcdto.BCRAResponse {
	if result == nil {
		return nil
	}
	return &calcdto.BCRAResponse{
		AbsRisk:    result.AbsRisk,
		AbsRiskAvg: result.AbsRiskAvg,
		RRStar1:    result.RRStar1,
		RRStar2:    result.RRStar2,
		ProjIntvl:  result.ProjIntvl,
	}
}

func mapGailResult(result *entity.GailResult) *calcdto.GailResponse {
	if result == nil {
		return nil
	}
	return &calcdto.GailResponse{
		AbsoluteRisk: result.AbsoluteRisk,
		RelativeRisk: result.RelativeRisk,
	}
}

func mapPLCOResult(result *entity.PLCOResult) *calcdto.PLCOResponse {
	if result == nil {
		return nil
	}
	response := &calcdto.PLCOResponse{
		PLCOM20126YrRisk:     result.PLCOM20126YrRisk,
		PLCOM20123YrRisk:     result.PLCOM20123YrRisk,
		PLCOM2012RiskPercent: result.PLCOM2012RiskPercent,
	}
	if result.PLCO2012Results3Yr != nil && result.PLCO2012Results6Yr != nil {
		response.PLCO2012Results = map[string]float64{
			"risk_3yr": *result.PLCO2012Results3Yr,
			"risk_6yr": *result.PLCO2012Results6Yr,
		}
	}
	return response
}

func (calcService *CalcService) callPremm5API(request calcdto.SendFormToPremm5Request) (calcdto.Premm5Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[CALC:PREMM5:API] Error marshaling request: %v", err)
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "PREMM5",
			Message: "failed to marshal request data",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:PREMM5:API] Calling PREMM5 API - CurrentAge=%d, PersonalCrcCount=%d, NumFdrCrc=%d",
		request.CurrentAge, request.PersonalCrcCount, request.NumFdrCrc)

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.Premm5)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[CALC:PREMM5:API] Error creating request: %v", err)
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PREMM5",
			Message: "failed to create API request",
			OrigErr: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CALC:PREMM5:API] API call failed: %v", err)
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PREMM5",
			Message: "failed to call PREMM5 API",
			OrigErr: err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CALC:PREMM5:API] Error reading response: %v", err)
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PREMM5",
			Message: "failed to read API response",
			OrigErr: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CALC:PREMM5:API] API returned error status %d: %s", resp.StatusCode, string(body))
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PREMM5",
			Message: fmt.Sprintf("API returned error status %d", resp.StatusCode),
			OrigErr: fmt.Errorf("response: %s", string(body)),
		}
	}

	var premm5Response calcdto.Premm5Response
	err = json.Unmarshal(body, &premm5Response)
	if err != nil {
		log.Printf("[CALC:PREMM5:API] Error parsing response: %v", err)
		return calcdto.Premm5Response{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PREMM5",
			Message: "failed to parse API response",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:PREMM5:API] API call successful - PAny: %.4f", premm5Response.PAny)
	return premm5Response, nil
}

func (calcService *CalcService) callBCRAAPI(request calcdto.SendFormToBCRARequest) (calcdto.BCRAResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[CALC:BCRA:API] Error marshaling request: %v", err)
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "BCRA",
			Message: "failed to marshal request data",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:BCRA:API] Calling BCRA API - T1=%.1f, T2=%.1f, N_Biop=%d, HypPlas=%d, AgeMen=%d, Age1st=%d, N_Rels=%d",
		request.T1, request.T2, request.N_Biop, request.HypPlas, request.AgeMen, request.Age1st, request.N_Rels)

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.BCRA)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[CALC:BCRA:API] Error creating request: %v", err)
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "BCRA",
			Message: "failed to create API request",
			OrigErr: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CALC:BCRA:API] API call failed: %v", err)
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "BCRA",
			Message: "failed to call BCRA API",
			OrigErr: err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CALC:BCRA:API] Error reading response: %v", err)
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "BCRA",
			Message: "failed to read API response",
			OrigErr: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CALC:BCRA:API] API returned error status %d: %s", resp.StatusCode, string(body))
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "BCRA",
			Message: fmt.Sprintf("API returned error status %d", resp.StatusCode),
			OrigErr: fmt.Errorf("response: %s", string(body)),
		}
	}

	var bcraResponse calcdto.BCRAResponse
	err = json.Unmarshal(body, &bcraResponse)
	if err != nil {
		log.Printf("[CALC:BCRA:API] Error parsing response: %v", err)
		return calcdto.BCRAResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "BCRA",
			Message: "failed to parse API response",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:BCRA:API] API call successful - AbsRisk: %.4f", bcraResponse.AbsRisk)
	return bcraResponse, nil
}

func (calcService *CalcService) callGailAPI(request calcdto.SendFormToGailRequest) (calcdto.GailResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[CALC:Gail:API] Error marshaling request: %v", err)
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "Gail",
			Message: "failed to marshal request data",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:Gail:API] Calling Gail API - Age=%d, MenarcheAge=%d, NumBiopsies=%d, FLBAge=%d, NumRelatives=%d, Race=%s",
		request.Age, request.MenarcheAge, request.NumBiopsies, request.FLBAge, request.NumRelatives, request.Race)

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.Gail)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[CALC:Gail:API] Error creating request: %v", err)
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "Gail",
			Message: "failed to create API request",
			OrigErr: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CALC:Gail:API] API call failed: %v", err)
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "Gail",
			Message: "failed to call Gail API",
			OrigErr: err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CALC:Gail:API] Error reading response: %v", err)
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "Gail",
			Message: "failed to read API response",
			OrigErr: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CALC:Gail:API] API returned error status %d: %s", resp.StatusCode, string(body))
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "Gail",
			Message: fmt.Sprintf("API returned error status %d", resp.StatusCode),
			OrigErr: fmt.Errorf("response: %s", string(body)),
		}
	}

	var gailResponse calcdto.GailResponse
	err = json.Unmarshal(body, &gailResponse)
	if err != nil {
		log.Printf("[CALC:Gail:API] Error parsing response: %v", err)
		return calcdto.GailResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "Gail",
			Message: "failed to parse API response",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:Gail:API] API call successful - AbsoluteRisk: %.4f", gailResponse.AbsoluteRisk)
	return gailResponse, nil
}

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

func (calcService *CalcService) saveGailResult(formID uint, gailResponse calcdto.GailResponse) error {
	existingResult, err := calcService.formRepository.FindGailResultByFormID(calcService.db, formID)
	if err != nil {
		return err
	}

	if existingResult == nil {
		gailResult := &entity.GailResult{
			FormID:       formID,
			AbsoluteRisk: gailResponse.AbsoluteRisk,
			RelativeRisk: gailResponse.RelativeRisk,
		}
		return calcService.formRepository.CreateGailResult(calcService.db, gailResult)
	}

	existingResult.AbsoluteRisk = gailResponse.AbsoluteRisk
	if gailResponse.RelativeRisk != nil {
		existingResult.RelativeRisk = gailResponse.RelativeRisk
	}

	return calcService.formRepository.UpdateGailResult(calcService.db, existingResult)
}

func (calcService *CalcService) sendFormToPLCO(form *entity.Form, userID uint) (calcdto.ModelResponse, error) {
	basicInfo, err := calcService.formRepository.FindBasicInfoByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PLCO] Database error fetching basic info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PLCO",
			Message: "failed to fetch basic information",
			OrigErr: err,
		}
	}
	if basicInfo == nil {
		log.Printf("[CALC:PLCO] Missing basic info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeMissingData,
			Model:   "PLCO",
			Message: "basic information not found in form",
		}
	}

	lungCancerInfo, err := calcService.formRepository.FindLungCancerByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PLCO] Database error fetching lung cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PLCO",
			Message: "failed to fetch lung cancer information",
			OrigErr: err,
		}
	}

	if lungCancerInfo == nil {
		log.Printf("[CALC:PLCO] Missing lung cancer info - FormID: %d", form.ID)
		return calcdto.ModelResponse{}, exception.NotFoundError{Item: "lung cancer information"}
	}

	cancerInfo, err := calcService.formRepository.FindCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PLCO] Database error fetching cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PLCO",
			Message: "failed to fetch cancer information",
			OrigErr: err,
		}
	}

	familyCancerInfo, err := calcService.formRepository.FindFamilyCancersByFormID(calcService.db, form.ID)
	if err != nil {
		log.Printf("[CALC:PLCO] Database error fetching family cancer info - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PLCO",
			Message: "failed to fetch family cancer information",
			OrigErr: err,
		}
	}

	// Calculate current age
	currentAge := calculateAge(basicInfo.BirthDate)

	// Validate age is within acceptable range (0-120)
	if currentAge < 0 || currentAge > 120 {
		log.Printf("[CALC:PLCO] Invalid age - FormID: %d, Age: %d, BirthDate: %v", form.ID, currentAge, basicInfo.BirthDate)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "PLCO",
			Message: fmt.Sprintf("calculated age %d is outside valid range (0-120)", currentAge),
		}
	}

	// Validate BMI parameters
	if basicInfo.Height <= 0 || basicInfo.Weight < 0 {
		log.Printf("[CALC:PLCO] Invalid height/weight - FormID: %d, Height: %.2f, Weight: %.2f", form.ID, basicInfo.Height, basicInfo.Weight)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "PLCO",
			Message: "invalid height or weight values",
		}
	}

	// Calculate BMI from height and weight
	bmi := basicInfo.Weight / ((basicInfo.Height / 100) * (basicInfo.Height / 100))

	// Map COPD from ChronicLungDisease
	copd := 0
	if answeredYes(lungCancerInfo.ChronicLungDisease) {
		copd = 1
	}

	// Map personal cancer history
	personalCancerHistory := 0
	if len(cancerInfo) > 0 {
		personalCancerHistory = 1
	}

	// Map family lung cancer history
	familyLungCancer := 0
	for _, info := range familyCancerInfo {
		if info.CancerType == enum.CancerTypeLung {
			familyLungCancer = 1
			break
		}
	}

	// Map smoking status and history
	smokingStatus := 0 // 0 = Former smoker
	cigarettesPerDay := 0.0
	smokingDuration := 0
	yearsQuit := 0

	if answeredYes(lungCancerInfo.CurrentSmoking) {
		smokingStatus = 1 // 1 = Current smoker
		if lungCancerInfo.CigarettesPerDayCurrent != nil {
			cigarettesPerDay = float64(*lungCancerInfo.CigarettesPerDayCurrent)
		}
		if lungCancerInfo.SmokingStartAgeCurrent != nil {
			smokingDuration = currentAge - int(*lungCancerInfo.SmokingStartAgeCurrent)
			if smokingDuration < 0 {
				smokingDuration = 0
			}
		}
	} else if answerinList(lungCancerInfo.PastSmoking, []enum.Answer{enum.AnswerYes, enum.AnswerAgo, enum.AnswerLongAgo}) {
		// Former smoker
		if lungCancerInfo.CigarettesPerDayPast != nil {
			cigarettesPerDay = float64(*lungCancerInfo.CigarettesPerDayPast)
		}
		if lungCancerInfo.SmokingStartAgePast != nil {
			// Estimate smoking duration - assume they smoked until recently or use a reasonable estimate
			// Since we don't have quit age, we'll estimate years quit as 0 (recently quit)
			// This is conservative and the model can handle it
			estimatedQuitAge := currentAge - 1 // Assume they quit 1 year ago
			if estimatedQuitAge > int(*lungCancerInfo.SmokingStartAgePast) {
				smokingDuration = estimatedQuitAge - int(*lungCancerInfo.SmokingStartAgePast)
				yearsQuit = 1 // Assume quit 1 year ago
			} else {
				smokingDuration = currentAge - int(*lungCancerInfo.SmokingStartAgePast)
				yearsQuit = 0
			}
			if smokingDuration < 0 {
				smokingDuration = 0
			}
		}
	}

	// Build PLCO request
	request := calcdto.SendFormToPLCORequest{
		Age:                   currentAge,
		Education:             4, // Default to "Some college" (can be updated when education field is added)
		BMI:                   bmi,
		COPD:                  copd,
		PersonalCancerHistory: personalCancerHistory,
		FamilyLungCancer:      familyLungCancer,
		RaceWhite:             1, // Default to White (can be updated when race field is added)
		RaceBlack:             0,
		RaceHispanic:          0,
		RaceAsian:             0,
		RaceNHPI:              0,
		RaceAmericanIndian:    0,
		SmokingStatus:         smokingStatus,
		CigarettesPerDay:      cigarettesPerDay,
		SmokingDuration:       smokingDuration,
		YearsQuit:             yearsQuit,
		ScreeningResult:       nil, // Optional - not currently collected
	}

	// Make API call
	plcoResponse, err := calcService.callPLCOAPI(request)
	if err != nil {
		return calcdto.ModelResponse{}, err
	}

	// Save result to database
	err = calcService.savePLCOResult(form.ID, plcoResponse)
	if err != nil {
		log.Printf("[CALC:PLCO] Database error saving result - FormID: %d, Error: %v", form.ID, err)
		return calcdto.ModelResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeDatabaseError,
			Model:   "PLCO",
			Message: "failed to save calculation result",
			OrigErr: err,
		}
	}

	actionLog := actionlogdto.LogAction{
		ActorID:    userID,
		Action:     enum.ActionTypeFormSentToPLCO,
		ResourceID: &form.ID,
	}
	calcService.actionLogService.LogAction(actionLog)

	log.Printf("[CALC:PLCO] Calculation successful - FormID: %d, RiskPercent: %.2f", form.ID, plcoResponse.PLCOM2012RiskPercent)
	return calcdto.ModelResponse{
		Name:        "PLCO",
		Probability: plcoResponse.PLCOM2012RiskPercent / 100.0,
	}, nil
}

func (calcService *CalcService) callPLCOAPI(request calcdto.SendFormToPLCORequest) (calcdto.PLCOResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[CALC:PLCO:API] Error marshaling request: %v", err)
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeInvalidData,
			Model:   "PLCO",
			Message: "failed to marshal request data",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:PLCO:API] Calling PLCO API - Age=%d, BMI=%.2f, SmokingStatus=%d, CigarettesPerDay=%.1f, SmokingDuration=%d, YearsQuit=%d",
		request.Age, request.BMI, request.SmokingStatus, request.CigarettesPerDay, request.SmokingDuration, request.YearsQuit)

	url := fmt.Sprintf("%s/calculate", calcService.calcURL.PLCO)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[CALC:PLCO:API] Error creating request: %v", err)
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PLCO",
			Message: "failed to create API request",
			OrigErr: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CALC:PLCO:API] API call failed: %v", err)
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PLCO",
			Message: "failed to call PLCO API",
			OrigErr: err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CALC:PLCO:API] Error reading response: %v", err)
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PLCO",
			Message: "failed to read API response",
			OrigErr: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CALC:PLCO:API] API returned error status %d: %s", resp.StatusCode, string(body))
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PLCO",
			Message: fmt.Sprintf("API returned error status %d", resp.StatusCode),
			OrigErr: fmt.Errorf("response: %s", string(body)),
		}
	}

	var plcoResponse calcdto.PLCOResponse
	err = json.Unmarshal(body, &plcoResponse)
	if err != nil {
		log.Printf("[CALC:PLCO:API] Error parsing response: %v", err)
		return calcdto.PLCOResponse{}, &exception.CalcError{
			Type:    exception.ErrorTypeAPIFailure,
			Model:   "PLCO",
			Message: "failed to parse API response",
			OrigErr: err,
		}
	}

	log.Printf("[CALC:PLCO:API] API call successful - RiskPercent: %.2f", plcoResponse.PLCOM2012RiskPercent)
	return plcoResponse, nil
}

func (calcService *CalcService) savePLCOResult(formID uint, plcoResponse calcdto.PLCOResponse) error {
	existingResult, err := calcService.formRepository.FindPLCOResultByFormID(calcService.db, formID)
	if err != nil {
		return err
	}

	if existingResult == nil {
		plcoResult := &entity.PLCOResult{
			FormID:               formID,
			PLCOM20126YrRisk:     plcoResponse.PLCOM20126YrRisk,
			PLCOM20123YrRisk:     plcoResponse.PLCOM20123YrRisk,
			PLCOM2012RiskPercent: plcoResponse.PLCOM2012RiskPercent,
		}

		if plcoResponse.PLCO2012Results != nil {
			if risk3yr, ok := plcoResponse.PLCO2012Results["risk_3yr"]; ok {
				plcoResult.PLCO2012Results3Yr = &risk3yr
			}
			if risk6yr, ok := plcoResponse.PLCO2012Results["risk_6yr"]; ok {
				plcoResult.PLCO2012Results6Yr = &risk6yr
			}
		}

		return calcService.formRepository.CreatePLCOResult(calcService.db, plcoResult)
	}

	existingResult.PLCOM20126YrRisk = plcoResponse.PLCOM20126YrRisk
	existingResult.PLCOM20123YrRisk = plcoResponse.PLCOM20123YrRisk
	existingResult.PLCOM2012RiskPercent = plcoResponse.PLCOM2012RiskPercent

	if plcoResponse.PLCO2012Results != nil {
		if risk3yr, ok := plcoResponse.PLCO2012Results["risk_3yr"]; ok {
			existingResult.PLCO2012Results3Yr = &risk3yr
		}
		if risk6yr, ok := plcoResponse.PLCO2012Results["risk_6yr"]; ok {
			existingResult.PLCO2012Results6Yr = &risk6yr
		}
	}

	return calcService.formRepository.UpdatePLCOResult(calcService.db, existingResult)
}

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

func mapAgeAtMenarche(mamographyInfo *entity.MamoGraphyInfo) int {
	if mamographyInfo.GhaedeAge > 0 {
		return int(mamographyInfo.GhaedeAge)
	}
	return 99 // Unknown
}

func mapAgeAtFirstBirth(mamographyInfo *entity.MamoGraphyInfo) int {
	// First check if patient has children
	if !answeredYes(mamographyInfo.HasChildren) {
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

func mapFirstDegreeBreastCancerRelatives(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 99 // Unknown
	}

	count := 0

	for _, info := range familyInfo {
		// Check mother, father, siblings for breast cancer
		if (info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister) &&
			info.CancerType == enum.CancerTypeBreast {
			count++
		}
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
func mapPersonalCrcCount(cancerInfo []*entity.CancerInfo) int {
	// TODO: You may need additional field to track multiple CRCs
	// For now, assume single CRC = 1, multiple = 2
	var crcCount int = 0
	for _, v := range cancerInfo {
		if v.CancerType == enum.CancerTypeColon {
			crcCount++
		}
	}
	return crcCount
}

// Map count of first-degree relatives with CRC (0, 1, or 2+)
func mapNumFdrCrc(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	count := 0
	for _, info := range familyInfo {
		// First-degree relatives: Mother, Father, Brother, Sister
		if (info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister) &&
			info.CancerType == enum.CancerTypeColon {
			count++
		}
	}
	// Cap at 2 for >=2
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among FDR with CRC
func mapYoungestFdrCrcAge(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	youngestAge := 999
	for _, info := range familyInfo {
		// First-degree relatives: Mother, Father, Brother, Sister
		if (info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister) &&
			info.CancerType == enum.CancerTypeColon {
			age := int(info.CancerAge)
			if age < youngestAge {
				youngestAge = age
			}
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of FDR with endometrial cancer
func mapNumFdrEc(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	count := 0
	for _, info := range familyInfo {
		// First-degree relatives: Mother, Father, Brother, Sister
		if (info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister) &&
			info.CancerType == enum.CancerTypeEndometrial {
			count++
		}
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among FDR with endometrial cancer
func mapYoungestFdrEcAge(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	youngestAge := 999
	for _, info := range familyInfo {
		// First-degree relatives: Mother, Father, Brother, Sister
		if (info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister) &&
			info.CancerType == enum.CancerTypeEndometrial {
			age := int(info.CancerAge)
			if age < youngestAge {
				youngestAge = age
			}
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of second-degree relatives with CRC (0, 1, or 2+)
func mapNumSdrCrc(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	count := 0
	for _, info := range familyInfo {
		// Second-degree relatives: Uncles, Aunts, Grandparents
		if (info.Relative == enum.PaternalUncle || info.Relative == enum.PaternalAunt ||
			info.Relative == enum.MaternalUncle || info.Relative == enum.MaternalAunt ||
			info.Relative == enum.PaternalGrandFather || info.Relative == enum.PaternalGrandMother ||
			info.Relative == enum.MaternalGrandFather || info.Relative == enum.MaternalGrandMother ||
			info.Relative == enum.DistantRelative) &&
			info.CancerType == enum.CancerTypeColon {
			count++
		}
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among SDR with CRC
func mapYoungestSdrCrcAge(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	youngestAge := 999
	for _, info := range familyInfo {
		// Second-degree relatives: Uncles, Aunts, Grandparents
		if (info.Relative == enum.PaternalUncle || info.Relative == enum.PaternalAunt ||
			info.Relative == enum.MaternalUncle || info.Relative == enum.MaternalAunt ||
			info.Relative == enum.PaternalGrandFather || info.Relative == enum.PaternalGrandMother ||
			info.Relative == enum.MaternalGrandFather || info.Relative == enum.MaternalGrandMother ||
			info.Relative == enum.DistantRelative) &&
			info.CancerType == enum.CancerTypeColon {
			age := int(info.CancerAge)
			if age < youngestAge {
				youngestAge = age
			}
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map count of SDR with endometrial cancer
func mapNumSdrEc(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	count := 0
	for _, info := range familyInfo {
		// Second-degree relatives: Uncles, Aunts, Grandparents
		if (info.Relative == enum.PaternalUncle || info.Relative == enum.PaternalAunt ||
			info.Relative == enum.MaternalUncle || info.Relative == enum.MaternalAunt ||
			info.Relative == enum.PaternalGrandFather || info.Relative == enum.PaternalGrandMother ||
			info.Relative == enum.MaternalGrandFather || info.Relative == enum.MaternalGrandMother ||
			info.Relative == enum.DistantRelative) &&
			info.CancerType == enum.CancerTypeEndometrial {
			count++
		}
	}
	if count >= 2 {
		return 2
	}
	return count
}

// Map youngest age among SDR with endometrial cancer
func mapYoungestSdrEcAge(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	youngestAge := 999
	for _, info := range familyInfo {
		// Second-degree relatives: Uncles, Aunts, Grandparents
		if (info.Relative == enum.PaternalUncle || info.Relative == enum.PaternalAunt ||
			info.Relative == enum.MaternalUncle || info.Relative == enum.MaternalAunt ||
			info.Relative == enum.PaternalGrandFather || info.Relative == enum.PaternalGrandMother ||
			info.Relative == enum.MaternalGrandFather || info.Relative == enum.MaternalGrandMother ||
			info.Relative == enum.DistantRelative) &&
			info.CancerType == enum.CancerTypeEndometrial {
			age := int(info.CancerAge)
			if age < youngestAge {
				youngestAge = age
			}
		}
	}
	if youngestAge == 999 {
		return 0
	}
	return youngestAge
}

// Map if FDR has other Lynch syndrome cancers
func mapHasFdrOtherLs(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	// Check mother, father, siblings for other LS cancers
	lscancers := []enum.CancerType{
		enum.CancerTypeOvarian, enum.CancerTypeStomach,
		enum.CancerTypePancreatic, enum.CancerTypeBrain,
		enum.CancerTypeLiver,
	}

	for _, info := range familyInfo {
		// First-degree relatives: Mother, Father, Brother, Sister
		if info.Relative == enum.Mother || info.Relative == enum.Father ||
			info.Relative == enum.Brother || info.Relative == enum.Sister {
			for _, lsc := range lscancers {
				if info.CancerType == lsc {
					return 1
				}
			}
		}
	}

	return 0
}

// Map if SDR has other Lynch syndrome cancers
func mapHasSdrOtherLs(familyInfo []*entity.FamilyCancerInfo) int {
	if len(familyInfo) == 0 {
		return 0
	}
	// Check aunt/uncle for other LS cancers
	lscancers := []enum.CancerType{
		enum.CancerTypeOvarian, enum.CancerTypeStomach,
		enum.CancerTypePancreatic, enum.CancerTypeBrain,
		enum.CancerTypeLiver,
	}

	for _, info := range familyInfo {
		// Second-degree relatives: Uncles, Aunts, Grandparents
		if info.Relative == enum.PaternalUncle || info.Relative == enum.PaternalAunt ||
			info.Relative == enum.MaternalUncle || info.Relative == enum.MaternalAunt ||
			info.Relative == enum.PaternalGrandFather || info.Relative == enum.PaternalGrandMother ||
			info.Relative == enum.MaternalGrandFather || info.Relative == enum.MaternalGrandMother ||
			info.Relative == enum.DistantRelative {
			for _, lsc := range lscancers {
				if info.CancerType == lsc {
					return 1
				}
			}
		}
	}

	return 0
}

// Map youngest age at CRC diagnosis
func mapAgeCrcDx(cancerInfo []*entity.CancerInfo) int {
	for _, v := range cancerInfo {
		if v.CancerType == enum.CancerTypeColon {
			return int(v.CancerAge)
		}
	}
	return 0
}

// Map endometrial cancer
func mapPersonalEndometrial(cancerInfo []*entity.CancerInfo) int {
	for _, v := range cancerInfo {
		if v.CancerType == enum.CancerTypeEndometrial {
			return 1
		}
	}
	return 0
}

// Map age at endometrial cancer diagnosis
func mapAgeEcDx(cancerInfo []*entity.CancerInfo) int {
	for _, v := range cancerInfo {
		if v.CancerType == enum.CancerTypeEndometrial {
			return int(v.CancerAge)
		}
	}
	return 0
}

// Map other LS-associated cancers
func mapPersonalLsOther(cancerInfo []*entity.CancerInfo) int {
	for _, v := range cancerInfo {
		if v.CancerType == enum.CancerTypeOvarian || v.CancerType == enum.CancerTypeStomach || v.CancerType == enum.CancerTypePancreatic {
			return 1
		}
	}
	return 0
}

// Return true if answered Yes, false otherwise
func answeredYes(answer *enum.Answer) bool {
	if answer != nil && *answer == enum.AnswerYes {
		return true
	}
	return false
}

func answerinList(answer *enum.Answer, validAnswers []enum.Answer) bool {
	if answer == nil {
		return false
	}
	for _, valid := range validAnswers {
		if *answer == valid {
			return true
		}
	}
	return false
}
