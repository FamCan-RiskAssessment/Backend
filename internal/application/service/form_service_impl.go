package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormService struct {
	constants      *bootstrap.Constants
	formRepository postgres.FormRepository
	userService    usecase.UserService
	db             database.Database
}

func NewFormService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	userService usecase.UserService,
	db database.Database,
) *FormService {
	return &FormService{
		constants:      constants,
		formRepository: formRepository,
		userService:    userService,
		db:             db,
	}
}

func (formService *FormService) CreateBasicInfoForm(request formdto.CreateBasicFormRequest) (formdto.BasicFormResponse, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return formdto.BasicFormResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return formdto.BasicFormResponse{}, notFoundError
	}

	form := &entity.Form{
		UserID: request.UserID,
		Status: enum.FormStatusInComplete,
	}

	if err = formService.formRepository.CreateForm(formService.db, form); err != nil {
		return formdto.BasicFormResponse{}, err
	}

	basic := &entity.BasicInfo{
		FormID:               form.ID,
		Gender:               request.Gender,
		BirthYear:            request.BirthYear,
		BirthMonth:           request.BirthMonth,
		BirthDay:             request.BirthDay,
		IsAtba:               request.IsAtba,
		SocialSecurityNumber: request.SocialSecurityNumber,
		Height:               request.Height,
		Weight:               request.Weight,
	}

	if err = formService.formRepository.CreateBasicInfo(formService.db, basic); err != nil {
		return formdto.BasicFormResponse{}, err
	}

	response := formdto.BasicFormResponse{
		FormID:    form.ID,
		Status:    form.Status.String(),
		UserID:    form.UserID,
		CreatedAt: form.CreatedAt,
		UpdatedAt: form.UpdatedAt,
	}

	return response, nil
}

func (formService *FormService) UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}

	if info == nil {
		info = &entity.GeneralHealthInfo{FormID: request.FormID}
	}

	info.DrinksAlcohol = request.DrinksAlcohol
	info.CupsPerWeek = request.CupsPerWeek
	info.LastMonthSabzijatMeal = request.LastMonthSabzijatMeal
	info.LastMonthSabzijatWeight = request.LastMonthSabzijatWeight
	info.MediumActivityMonthInYear = request.MediumActivityMonthInYear
	info.MediumActivityHourInWeek = request.MediumActivityHourInWeek
	info.HardActivityMonthInYear = request.HardActivityMonthInYear
	info.HardActivityHourInWeek = request.HardActivityHourInWeek
	info.SmokeAtLeast100 = request.SmokeAtLeast100
	info.SmokingAge = request.SmokingAge
	info.SmokingNow = request.SmokingNow
	info.LeaveSmokingAge = request.LeaveSmokingAge
	info.CountSmokingDaily = request.CountSmokingDaily
	info.CountGheliandaily = request.CountGheliandaily
	info.CountSmokingDailyPast = request.CountSmokingDailyPast
	info.CountGheliandailyPast = request.CountGheliandailyPast

	if info.ID == 0 {
		return formService.formRepository.CreateGeneralHealth(formService.db, info)
	}
	return formService.formRepository.UpdateGeneralHealth(formService.db, info)
}

func (formService *FormService) UpsertMamography(request formdto.UpsertMamographyRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.MamoGraphyInfo{FormID: request.FormID}
	}

	info.GhaedeAge = request.GhaedeAge
	info.HasChildren = request.HasChildren
	info.NumberOfChildren = request.NumberOfChildren
	info.AgeOfFirstBirth = request.AgeOfFirstBirth
	info.MenopausalStatus = request.MenopausalStatus
	info.MenopauseAge = request.MenopauseAge
	info.HRT = request.HRT
	info.HRTUseLength = request.HRTUseLength
	info.LastFiveYearsHRTUse = request.LastFiveYearsHRTUse
	info.CurrentHRTUse = request.CurrentHRTUse
	info.IntendedHRTUse = request.IntendedHRTUse
	info.HRTType = request.HRTType
	info.Oral = request.Oral
	info.OralDuration = request.OralDuration
	info.OralTwoLastYears = request.OralTwoLastYears
	info.MamoGraphy = request.MamoGraphy
	info.Falop = request.Falop
	info.Andometrioz = request.Andometrioz
	info.LeavePestan = request.LeavePestan
	info.LeaveTokhmdan = request.LeaveTokhmdan
	info.LaDeColon = request.LaDeColon
	info.LaDePol = request.LaDePol
	info.AspLaMo = request.AspLaMo
	info.NsaiDLaMo = request.NsaiDLaMo
	info.LastFiveYearBloodTestInStool = request.LastFiveYearBloodTestInStool

	if info.ID == 0 {
		return formService.formRepository.CreateMamography(formService.db, info)
	}
	return formService.formRepository.UpdateMamography(formService.db, info)
}

func (formService *FormService) UpsertCancer(request formdto.UpsertCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.CancerInfo{FormID: request.FormID}
	}

	info.Cancer = request.Cancer
	info.CancerType = request.CancerType
	info.CancerAge = request.CancerAge

	if info.ID == 0 {
		return formService.formRepository.CreateCancer(formService.db, info)
	}
	return formService.formRepository.UpdateCancer(formService.db, info)
}

func (formService *FormService) UpsertFamilyCancer(request formdto.UpsertFamilyCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindFamilyCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.FamilyCancerInfo{FormID: request.FormID}
	}

	info.ChildCancer = request.ChildCancer
	info.ChildName = request.ChildName
	info.ChildCancerType = request.ChildCancerType
	info.ChildCancerAge = request.ChildCancerAge
	info.ChildLifeStatus = request.ChildLifeStatus

	info.MotherCancer = request.MotherCancer
	info.MotherName = request.MotherName
	info.MotherLifeStatus = request.MotherLifeStatus
	info.MotherCancerType = request.MotherCancerType
	info.MotherCancerAge = request.MotherCancerAge

	info.FatherCancer = request.FatherCancer
	info.FatherName = request.FatherName
	info.FatherLifeStatus = request.FatherLifeStatus
	info.FatherCancerType = request.FatherCancerType
	info.FatherCancerAge = request.FatherCancerAge

	info.SiblingCancer = request.SiblingCancer
	info.SiblingName = request.SiblingName
	info.SiblingLifeStatus = request.SiblingLifeStatus
	info.SiblingCancerType = request.SiblingCancerType
	info.SiblingCancerAge = request.SiblingCancerAge

	info.AmeAmoCancer = request.AmeAmoCancer
	info.AmeAmoName = request.AmeAmoName
	info.AmeAmoLifeStatus = request.AmeAmoLifeStatus
	info.AmeAmoCancerType = request.AmeAmoCancerType
	info.AmeAmoCancerAge = request.AmeAmoCancerAge

	info.KhaleDaeiCancer = request.KhaleDaeiCancer
	info.KhaleDaeiName = request.KhaleDaeiName
	info.KhaleDaeiLifeStatus = request.KhaleDaeiLifeStatus
	info.KhaleDaeiCancerType = request.KhaleDaeiCancerType
	info.KhaleDaeiCancerAge = request.KhaleDaeiCancerAge

	info.OtherRelativeCancer = request.OtherRelativeCancer
	info.OtherRelativeName = request.OtherRelativeName
	info.OtherRelativeRelation = request.OtherRelativeRelation
	info.OtherRelativeLifeStatus = request.OtherRelativeLifeStatus
	info.OtherRelativeCancerType = request.OtherRelativeCancerType
	info.OtherRelativeCancerAge = request.OtherRelativeCancerAge

	if info.ID == 0 {
		return formService.formRepository.CreateFamilyCancer(formService.db, info)
	}
	return formService.formRepository.UpdateFamilyCancer(formService.db, info)
}

func (formService *FormService) UpsertContact(request formdto.UpsertContactRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.ContactInfo{FormID: request.FormID}
	}

	info.Name = request.Name
	info.TestGen = request.TestGen
	info.FmTestGen = request.FmTestGen
	info.CallExpert = request.CallExpert
	info.BirthCountry = request.BirthCountry
	info.Province = request.Province
	info.City = request.City
	info.Country = request.Country
	info.Address = request.Address
	info.PostalCode = request.PostalCode

	if info.ID == 0 {
		return formService.formRepository.CreateContact(formService.db, info)
	}
	return formService.formRepository.UpdateContact(formService.db, info)
}

func (formService *FormService) UpsertLungCancer(request formdto.UpsertLungCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.LungCancerInfo{FormID: request.FormID}
	}

	info.InsuranceStatus = request.InsuranceStatus
	info.SupplementaryInsurances = request.SupplementaryInsurances
	info.Hypertension = request.Hypertension
	info.HypertensionTreatment = request.HypertensionTreatment
	info.HeartDisease = request.HeartDisease
	info.HeartDiseaseTreatment = request.HeartDiseaseTreatment
	info.Diabetes = request.Diabetes
	info.DiabetesTreatment = request.DiabetesTreatment
	info.ChronicLungDisease = request.ChronicLungDisease
	info.ChronicLungDiseaseType = request.ChronicLungDiseaseType
	info.LungCancerHistory = request.LungCancerHistory
	info.OtherCancerHistory = request.OtherCancerHistory
	info.OtherCancerType = request.OtherCancerType
	info.LungCancerFamily = request.LungCancerFamily
	info.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	info.OtherCancerFamily = request.OtherCancerFamily
	info.OtherCancerFamilyType = request.OtherCancerFamilyType
	info.OtherCancerFamilyRelation = request.OtherCancerFamilyRelation
	info.OccupationalExposure = request.OccupationalExposure
	info.CurrentSmoking = request.CurrentSmoking
	info.SmokingStartAgeCurrent = request.SmokingStartAgeCurrent
	info.SmokingTypesCurrent = request.SmokingTypesCurrent
	info.CigarettesPerDayCurrent = request.CigarettesPerDayCurrent
	info.CigarPerDayCurrent = request.CigarPerDayCurrent
	info.ECigPerDayCurrent = request.ECigPerDayCurrent
	info.PipePerDayCurrent = request.PipePerDayCurrent
	info.ChapoghPerDayCurrent = request.ChapoghPerDayCurrent
	info.SmokedOpiumPerDayCurrent = request.SmokedOpiumPerDayCurrent
	info.ChewedOpiumPerDayCurrent = request.ChewedOpiumPerDayCurrent
	info.HookahPerWeekCurrent = request.HookahPerWeekCurrent
	info.PastSmoking = request.PastSmoking
	info.SmokingStartAgePast = request.SmokingStartAgePast
	info.SmokingTypesPast = request.SmokingTypesPast
	info.CigarettesPerDayPast = request.CigarettesPerDayPast
	info.CigarPerDayPast = request.CigarPerDayPast
	info.ECigPerDayPast = request.ECigPerDayPast
	info.PipePerDayPast = request.PipePerDayPast
	info.ChapoghPerDayPast = request.ChapoghPerDayPast
	info.SmokedOpiumPerDayPast = request.SmokedOpiumPerDayPast
	info.ChewedOpiumPerDayPast = request.ChewedOpiumPerDayPast
	info.HookahPerWeekPast = request.HookahPerWeekPast
	info.SecondhandSmoke = request.SecondhandSmoke
	info.SecondhandSmokeLocation = request.SecondhandSmokeLocation

	if info.ID == 0 {
		return formService.formRepository.CreateLungCancer(formService.db, info)
	}
	return formService.formRepository.UpdateLungCancer(formService.db, info)
}

func (formService *FormService) ChangeFormStatus(request formdto.ChangeFormStatusRequest) (formdto.ChangeFormStatusResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.ChangeFormStatusResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.ChangeFormStatusResponse{}, notFoundError
	}

	form.Status = enum.FormStatusReady
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return formdto.ChangeFormStatusResponse{}, err
	}

	response := formdto.ChangeFormStatusResponse{
		Form: formdto.BasicFormResponse{
			FormID:    form.ID,
			Status:    form.Status.String(),
			UserID:    form.UserID,
			CreatedAt: form.CreatedAt,
			UpdatedAt: form.UpdatedAt,
		},
	}

	return response, nil
}

func (formService *FormService) GetBasicForm(formID uint) (formdto.GetBasicFormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetBasicFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetBasicFormResponse{}, notFoundError
	}

	basic, err := formService.formRepository.FindBasicInfoByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetBasicFormResponse{}, err
	}
	if basic == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetBasicFormResponse{}, notFoundError
	}

	return formdto.GetBasicFormResponse{
		ID:                   basic.ID,
		Gender:               basic.Gender,
		BirthYear:            basic.BirthYear,
		BirthMonth:           basic.BirthMonth,
		BirthDay:             basic.BirthDay,
		IsAtba:               basic.IsAtba,
		SocialSecurityNumber: basic.SocialSecurityNumber,
		Height:               basic.Height,
		Weight:               basic.Weight,
	}, nil
}

func (formService *FormService) GetGeneralHealth(formID uint) (formdto.GetGeneralHealthResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetGeneralHealthResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetGeneralHealthResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetGeneralHealthResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetGeneralHealthResponse{}, notFoundError
	}

	return formdto.GetGeneralHealthResponse{
		ID:                        info.ID,
		DrinksAlcohol:             info.DrinksAlcohol,
		CupsPerWeek:               info.CupsPerWeek,
		LastMonthSabzijatMeal:     info.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   info.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: info.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  info.MediumActivityHourInWeek,
		HardActivityMonthInYear:   info.HardActivityMonthInYear,
		HardActivityHourInWeek:    info.HardActivityHourInWeek,
		SmokeAtLeast100:           info.SmokeAtLeast100,
		SmokingAge:                info.SmokingAge,
		SmokingNow:                info.SmokingNow,
		LeaveSmokingAge:           info.LeaveSmokingAge,
		CountSmokingDaily:         info.CountSmokingDaily,
		CountGheliandaily:         info.CountGheliandaily,
		CountSmokingDailyPast:     info.CountSmokingDailyPast,
		CountGheliandailyPast:     info.CountGheliandailyPast,
	}, nil
}

func (formService *FormService) GetMamography(formID uint) (formdto.GetMamographyResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetMamographyResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetMamographyResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetMamographyResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetMamographyResponse{}, notFoundError
	}

	return formdto.GetMamographyResponse{
		ID:                           info.ID,
		GhaedeAge:                    info.GhaedeAge,
		HasChildren:                  info.HasChildren,
		NumberOfChildren:             info.NumberOfChildren,
		AgeOfFirstBirth:              info.AgeOfFirstBirth,
		MenopausalStatus:             info.MenopausalStatus,
		MenopauseAge:                 info.MenopauseAge,
		HRT:                          info.HRT,
		HRTUseLength:                 info.HRTUseLength,
		LastFiveYearsHRTUse:          info.LastFiveYearsHRTUse,
		CurrentHRTUse:                info.CurrentHRTUse,
		IntendedHRTUse:               info.IntendedHRTUse,
		HRTType:                      info.HRTType,
		Oral:                         info.Oral,
		OralDuration:                 info.OralDuration,
		OralTwoLastYears:             info.OralTwoLastYears,
		MamoGraphy:                   info.MamoGraphy,
		Falop:                        info.Falop,
		Andometrioz:                  info.Andometrioz,
		LeavePestan:                  info.LeavePestan,
		LeaveTokhmdan:                info.LeaveTokhmdan,
		LaDeColon:                    info.LaDeColon,
		LaDePol:                      info.LaDePol,
		AspLaMo:                      info.AspLaMo,
		NsaiDLaMo:                    info.NsaiDLaMo,
		LastFiveYearBloodTestInStool: info.LastFiveYearBloodTestInStool,
	}, nil
}

func (formService *FormService) GetCancer(formID uint) (formdto.GetCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetCancerResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindCancerByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetCancerResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetCancerResponse{}, notFoundError
	}

	return formdto.GetCancerResponse{
		ID:         info.ID,
		Cancer:     info.Cancer,
		CancerType: info.CancerType,
		CancerAge:  info.CancerAge,
	}, nil
}

func (formService *FormService) GetFamilyCancer(formID uint) (formdto.GetFamilyCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetFamilyCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetFamilyCancerResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindFamilyCancerByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetFamilyCancerResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetFamilyCancerResponse{}, notFoundError
	}

	return formdto.GetFamilyCancerResponse{
		ID:                      info.ID,
		ChildCancer:             info.ChildCancer,
		ChildName:               info.ChildName,
		ChildCancerType:         info.ChildCancerType,
		ChildCancerAge:          info.ChildCancerAge,
		ChildLifeStatus:         info.ChildLifeStatus,
		MotherCancer:            info.MotherCancer,
		MotherName:              info.MotherName,
		MotherLifeStatus:        info.MotherLifeStatus,
		MotherCancerType:        info.MotherCancerType,
		MotherCancerAge:         info.MotherCancerAge,
		FatherCancer:            info.FatherCancer,
		FatherName:              info.FatherName,
		FatherLifeStatus:        info.FatherLifeStatus,
		FatherCancerType:        info.FatherCancerType,
		FatherCancerAge:         info.FatherCancerAge,
		SiblingCancer:           info.SiblingCancer,
		SiblingName:             info.SiblingName,
		SiblingLifeStatus:       info.SiblingLifeStatus,
		SiblingCancerType:       info.SiblingCancerType,
		SiblingCancerAge:        info.SiblingCancerAge,
		AmeAmoCancer:            info.AmeAmoCancer,
		AmeAmoName:              info.AmeAmoName,
		AmeAmoLifeStatus:        info.AmeAmoLifeStatus,
		AmeAmoCancerType:        info.AmeAmoCancerType,
		AmeAmoCancerAge:         info.AmeAmoCancerAge,
		KhaleDaeiCancer:         info.KhaleDaeiCancer,
		KhaleDaeiName:           info.KhaleDaeiName,
		KhaleDaeiLifeStatus:     info.KhaleDaeiLifeStatus,
		KhaleDaeiCancerType:     info.KhaleDaeiCancerType,
		KhaleDaeiCancerAge:      info.KhaleDaeiCancerAge,
		OtherRelativeCancer:     info.OtherRelativeCancer,
		OtherRelativeName:       info.OtherRelativeName,
		OtherRelativeRelation:   info.OtherRelativeRelation,
		OtherRelativeLifeStatus: info.OtherRelativeLifeStatus,
		OtherRelativeCancerType: info.OtherRelativeCancerType,
		OtherRelativeCancerAge:  info.OtherRelativeCancerAge,
	}, nil
}

func (formService *FormService) GetContact(formID uint) (formdto.GetContactResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindContactByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	return formdto.GetContactResponse{
		ID:           info.ID,
		TestGen:      info.TestGen,
		FmTestGen:    info.FmTestGen,
		CallExpert:   info.CallExpert,
		BirthCountry: info.BirthCountry,
		Province:     info.Province,
		City:         info.City,
		Country:      info.Country,
		Address:      info.Address,
		PostalCode:   info.PostalCode,
	}, nil
}

func (formService *FormService) GetLungCancer(formID uint) (formdto.GetLungCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.GetLungCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetLungCancerResponse{}, notFoundError
	}

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, formID)
	if err != nil {
		return formdto.GetLungCancerResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetLungCancerResponse{}, notFoundError
	}

	return formdto.GetLungCancerResponse{
		ID:                        info.ID,
		InsuranceStatus:           info.InsuranceStatus,
		SupplementaryInsurances:   info.SupplementaryInsurances,
		Hypertension:              info.Hypertension,
		HypertensionTreatment:     info.HypertensionTreatment,
		HeartDisease:              info.HeartDisease,
		HeartDiseaseTreatment:     info.HeartDiseaseTreatment,
		Diabetes:                  info.Diabetes,
		DiabetesTreatment:         info.DiabetesTreatment,
		ChronicLungDisease:        info.ChronicLungDisease,
		ChronicLungDiseaseType:    info.ChronicLungDiseaseType,
		LungCancerHistory:         info.LungCancerHistory,
		OtherCancerHistory:        info.OtherCancerHistory,
		OtherCancerType:           info.OtherCancerType,
		LungCancerFamily:          info.LungCancerFamily,
		LungCancerFamilyRelation:  info.LungCancerFamilyRelation,
		OtherCancerFamily:         info.OtherCancerFamily,
		OtherCancerFamilyType:     info.OtherCancerFamilyType,
		OtherCancerFamilyRelation: info.OtherCancerFamilyRelation,
		OccupationalExposure:      info.OccupationalExposure,
		CurrentSmoking:            info.CurrentSmoking,
		SmokingStartAgeCurrent:    info.SmokingStartAgeCurrent,
		SmokingTypesCurrent:       info.SmokingTypesCurrent,
		CigarettesPerDayCurrent:   info.CigarettesPerDayCurrent,
		CigarPerDayCurrent:        info.CigarPerDayCurrent,
		ECigPerDayCurrent:         info.ECigPerDayCurrent,
		PipePerDayCurrent:         info.PipePerDayCurrent,
		ChapoghPerDayCurrent:      info.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:  info.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:  info.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:      info.HookahPerWeekCurrent,
		PastSmoking:               info.PastSmoking,
		SmokingStartAgePast:       info.SmokingStartAgePast,
		SmokingTypesPast:          info.SmokingTypesPast,
		CigarettesPerDayPast:      info.CigarettesPerDayPast,
		CigarPerDayPast:           info.CigarPerDayPast,
		ECigPerDayPast:            info.ECigPerDayPast,
		PipePerDayPast:            info.PipePerDayPast,
		ChapoghPerDayPast:         info.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:     info.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:     info.ChewedOpiumPerDayPast,
		HookahPerWeekPast:         info.HookahPerWeekPast,
		SecondhandSmoke:           info.SecondhandSmoke,
		SecondhandSmokeLocation:   info.SecondhandSmokeLocation,
	}, nil
}

func (formService *FormService) GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return nil, 0, notFoundError
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	forms, err := formService.formRepository.FindFormsByUserID(formService.db, request.UserID, options)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountFormsByUserID(formService.db, request.UserID)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:    form.ID,
			Status:    form.Status.String(),
			UserID:    form.UserID,
			CreatedAt: form.CreatedAt,
			UpdatedAt: form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) UpdateForm(request formdto.UpdateBasicFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) UpdateBasicInfo(request formdto.UpdateBasicFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}
	if form.Status != enum.FormStatusInComplete {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info, err := formService.formRepository.FindBasicInfoByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	info.BirthDay = request.BirthDay
	info.BirthMonth = request.BirthMonth
	info.BirthYear = request.BirthYear
	info.SocialSecurityNumber = request.SocialSecurityNumber
	info.Gender = request.Gender
	info.IsAtba = request.IsAtba
	info.Height = request.Height
	info.Weight = request.Weight

	err = formService.formRepository.UpdateBasicInfo(formService.db, info)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) DeleteForm(formID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	err = formService.formRepository.DeleteForm(formService.db, formID)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.FormResponse, int64, error) {
	forms, err := formService.formRepository.FindAllForms(formService.db, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllForms(formService.db, filters)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:    form.ID,
			Status:    form.Status.String(),
			UserID:    form.UserID,
			CreatedAt: form.CreatedAt,
			UpdatedAt: form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) AcceptForm(formID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}

	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	form.Status = enum.FormStatusApproved
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) RejectForm(formID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}

	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	form.Status = enum.FormStatusRejected
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}
