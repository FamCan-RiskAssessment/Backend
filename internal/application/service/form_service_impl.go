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

func (formService *FormService) entityToResponse(form *entity.Form) formdto.FormResponse {
	return formdto.FormResponse{
		ID:                   form.ID,
		Status:               form.Status.String(),
		UserID:               form.UserID,
		Name:                 form.Name,
		BirthDay:             form.BirthDay,
		BirthMonth:           form.BirthMonth,
		BirthYear:            form.BirthYear,
		Address:              form.Address,
		PostalCode:           form.PostalCode,
		SocialSecurityNumber: form.SocialSecurityNumber,
		Gender:               form.Gender,
		IsAtba:               form.IsAtba,
		Height:               form.Height,
		Weight:               form.Weight,

		DrinksAlcohol:             form.DrinksAlcohol,
		CupsPerWeek:               form.CupsPerWeek,
		LastMonthSabzijatMeal:     form.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   form.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: form.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  form.MediumActivityHourInWeek,
		HardActivityMonthInYear:   form.HardActivityMonthInYear,
		HardActivityHourInWeek:    form.HardActivityHourInWeek,
		SmokeAtLeast100:           form.SmokeAtLeast100,
		SmokingAge:                form.SmokingAge,
		SmokingNow:                form.SmokingNow,
		LeaveSmokingAge:           form.LeaveSmokingAge,
		CountSmokingDaily:         form.CountSmokingDaily,
		CountGheliandaily:         form.CountGheliandaily,
		CountSmokingDailyPast:     form.CountSmokingDailyPast,
		CountGheliandailyPast:     form.CountGheliandailyPast,

		GhaedeAge:                    form.GhaedeAge,
		HasChildren:                  form.HasChildren,
		NumberOfChildren:             form.NumberOfChildren,
		AgeOfFirstBirth:              form.AgeOfFirstBirth,
		MenopausalStatus:             form.MenopausalStatus,
		MenopauseAge:                 form.MenopauseAge,
		HRT:                          form.HRT,
		HRTUseLength:                 form.HRTUseLength,
		LastFiveYearsHRTUse:          form.LastFiveYearsHRTUse,
		CurrentHRTUse:                form.CurrentHRTUse,
		IntendedHRTUse:               form.IntendedHRTUse,
		HRTType:                      form.HRTType,
		Oral:                         form.Oral,
		OralDuration:                 form.OralDuration,
		OralTwoLastYears:             form.OralTwoLastYears,
		MamoGraphy:                   form.MamoGraphy,
		Falop:                        form.Falop,
		Andometrioz:                  form.Andometrioz,
		LeavePestan:                  form.LeavePestan,
		LeaveTokhmdan:                form.LeaveTokhmdan,
		LaDeColon:                    form.LaDeColon,
		LaDePol:                      form.LaDePol,
		AspLaMo:                      form.AspLaMo,
		NsaiDLaMo:                    form.NsaiDLaMo,
		LastFiveYearBloodTestInStool: form.LastFiveYearBloodTestInStool,

		Cancer:     form.Cancer,
		CancerType: form.CancerType,
		CancerAge:  form.CancerAge,

		ChildCancer:     form.ChildCancer,
		ChildName:       form.ChildName,
		ChildCancerType: form.ChildCancerType,
		ChildCancerAge:  form.ChildCancerAge,
		ChildLifeStatus: form.ChildLifeStatus,

		MotherCancer:     form.MotherCancer,
		MotherName:       form.MotherName,
		MotherLifeStatus: form.MotherLifeStatus,
		MotherCancerType: form.MotherCancerType,
		MotherCancerAge:  form.MotherCancerAge,

		FatherCancer:     form.FatherCancer,
		FatherName:       form.FatherName,
		FatherLifeStatus: form.FatherLifeStatus,
		FatherCancerType: form.FatherCancerType,
		FatherCancerAge:  form.FatherCancerAge,

		SiblingCancer:     form.SiblingCancer,
		SiblingName:       form.SiblingName,
		SiblingLifeStatus: form.SiblingLifeStatus,
		SiblingCancerType: form.SiblingCancerType,
		SiblingCancerAge:  form.SiblingCancerAge,

		AmeAmoCancer:     form.AmeAmoCancer,
		AmeAmoName:       form.AmeAmoName,
		AmeAmoLifeStatus: form.AmeAmoLifeStatus,
		AmeAmoCancerType: form.AmeAmoCancerType,
		AmeAmoCancerAge:  form.AmeAmoCancerAge,

		KhaleDaeiCancer:     form.KhaleDaeiCancer,
		KhaleDaeiName:       form.KhaleDaeiName,
		KhaleDaeiLifeStatus: form.KhaleDaeiLifeStatus,
		KhaleDaeiCancerType: form.KhaleDaeiCancerType,
		KhaleDaeiCancerAge:  form.KhaleDaeiCancerAge,

		OtherRelativeCancer:     form.OtherRelativeCancer,
		OtherRelativeName:       form.OtherRelativeName,
		OtherRelativeRelation:   form.OtherRelativeRelation,
		OtherRelativeLifeStatus: form.OtherRelativeLifeStatus,
		OtherRelativeCancerType: form.OtherRelativeCancerType,
		OtherRelativeCancerAge:  form.OtherRelativeCancerAge,

		TestGen:   form.TestGen,
		FmTestGen: form.FmTestGen,

		CallExpert:   form.CallExpert,
		BirthCountry: form.BirthCountry,
		Province:     form.Province,
		City:         form.City,
		Country:      form.Country,

		InsuranceStatus:           form.InsuranceStatus,
		SupplementaryInsurances:   form.SupplementaryInsurances,
		Hypertension:              form.Hypertension,
		HypertensionTreatment:     form.HypertensionTreatment,
		HeartDisease:              form.HeartDisease,
		HeartDiseaseTreatment:     form.HeartDiseaseTreatment,
		Diabetes:                  form.Diabetes,
		DiabetesTreatment:         form.DiabetesTreatment,
		ChronicLungDisease:        form.ChronicLungDisease,
		ChronicLungDiseaseType:    form.ChronicLungDiseaseType,
		LungCancerHistory:         form.LungCancerHistory,
		OtherCancerHistory:        form.OtherCancerHistory,
		OtherCancerType:           form.OtherCancerType,
		LungCancerFamily:          form.LungCancerFamily,
		LungCancerFamilyRelation:  form.LungCancerFamilyRelation,
		OtherCancerFamily:         form.OtherCancerFamily,
		OtherCancerFamilyType:     form.OtherCancerFamilyType,
		OtherCancerFamilyRelation: form.OtherCancerFamilyRelation,
		OccupationalExposure:      form.OccupationalExposure,
		CurrentSmoking:            form.CurrentSmoking,
		SmokingStartAgeCurrent:    form.SmokingStartAgeCurrent,
		SmokingTypesCurrent:       form.SmokingTypesCurrent,
		CigarettesPerDayCurrent:   form.CigarettesPerDayCurrent,
		CigarPerDayCurrent:        form.CigarPerDayCurrent,
		ECigPerDayCurrent:         form.ECigPerDayCurrent,
		PipePerDayCurrent:         form.PipePerDayCurrent,
		ChapoghPerDayCurrent:      form.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:  form.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:  form.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:      form.HookahPerWeekCurrent,
		PastSmoking:               form.PastSmoking,
		SmokingStartAgePast:       form.SmokingStartAgePast,
		SmokingTypesPast:          form.SmokingTypesPast,
		CigarettesPerDayPast:      form.CigarettesPerDayPast,
		CigarPerDayPast:           form.CigarPerDayPast,
		ECigPerDayPast:            form.ECigPerDayPast,
		PipePerDayPast:            form.PipePerDayPast,
		ChapoghPerDayPast:         form.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:     form.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:     form.ChewedOpiumPerDayPast,
		HookahPerWeekPast:         form.HookahPerWeekPast,
		SecondhandSmoke:           form.SecondhandSmoke,
		SecondhandSmokeLocation:   form.SecondhandSmokeLocation,

		CreatedAt: form.CreatedAt,
		UpdatedAt: form.UpdatedAt,
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
		Name:                 request.Name,
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

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.ContactInfo{FormID: request.FormID}
	}

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

func (formService *FormService) GetForm(formID uint) (formdto.FormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return formdto.FormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.FormResponse{}, notFoundError
	}

	return formService.entityToResponse(form), nil
}

func (formService *FormService) GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.FormResponse, int64, error) {
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

	formResponses := make([]formdto.FormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formService.entityToResponse(form)
	}

	return formResponses, count, nil
}

func (formService *FormService) UpdateForm(request formdto.UpdateFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	if request.Name != nil {
		form.Name = *request.Name
	}
	if request.BirthDay != nil {
		form.BirthDay = *request.BirthDay
	}
	if request.BirthMonth != nil {
		form.BirthMonth = *request.BirthMonth
	}
	if request.BirthYear != nil {
		form.BirthYear = *request.BirthYear
	}
	if request.Address != nil {
		form.Address = *request.Address
	}
	if request.PostalCode != nil {
		form.PostalCode = *request.PostalCode
	}
	if request.SocialSecurityNumber != nil {
		form.SocialSecurityNumber = *request.SocialSecurityNumber
	}
	if request.Gender != nil {
		form.Gender = *request.Gender
	}
	if request.IsAtba != nil {
		form.IsAtba = *request.IsAtba
	}
	if request.Height != nil {
		form.Height = *request.Height
	}
	if request.Weight != nil {
		form.Weight = *request.Weight
	}
	if request.DrinksAlcohol != nil {
		form.DrinksAlcohol = request.DrinksAlcohol
	}
	if request.CupsPerWeek != nil {
		form.CupsPerWeek = request.CupsPerWeek
	}
	if request.LastMonthSabzijatMeal != nil {
		form.LastMonthSabzijatMeal = *request.LastMonthSabzijatMeal
	}
	if request.LastMonthSabzijatWeight != nil {
		form.LastMonthSabzijatWeight = *request.LastMonthSabzijatWeight
	}
	if request.MediumActivityMonthInYear != nil {
		form.MediumActivityMonthInYear = *request.MediumActivityMonthInYear
	}
	if request.MediumActivityHourInWeek != nil {
		form.MediumActivityHourInWeek = *request.MediumActivityHourInWeek
	}
	if request.HardActivityMonthInYear != nil {
		form.HardActivityMonthInYear = *request.HardActivityMonthInYear
	}
	if request.HardActivityHourInWeek != nil {
		form.HardActivityHourInWeek = *request.HardActivityHourInWeek
	}
	if request.SmokeAtLeast100 != nil {
		form.SmokeAtLeast100 = request.SmokeAtLeast100
	}
	if request.SmokingNow != nil {
		form.SmokingNow = *request.SmokingNow
	}
	if request.LeaveSmokingAge != nil {
		form.LeaveSmokingAge = request.LeaveSmokingAge
	}
	if request.CountSmokingDaily != nil {
		form.CountSmokingDaily = request.CountSmokingDaily
	}
	if request.CountGheliandaily != nil {
		form.CountGheliandaily = request.CountGheliandaily
	}
	if request.CountSmokingDailyPast != nil {
		form.CountSmokingDailyPast = request.CountSmokingDailyPast
	}
	if request.CountGheliandailyPast != nil {
		form.CountGheliandailyPast = request.CountGheliandailyPast
	}
	if request.GhaedeAge != nil {
		form.GhaedeAge = *request.GhaedeAge
	}
	if request.HasChildren != nil {
		form.HasChildren = *request.HasChildren
	}
	if request.NumberOfChildren != nil {
		form.NumberOfChildren = request.NumberOfChildren
	}
	if request.AgeOfFirstBirth != nil {
		form.AgeOfFirstBirth = request.AgeOfFirstBirth
	}
	if request.MenopausalStatus != nil {
		form.MenopausalStatus = *request.MenopausalStatus
	}
	if request.HRT != nil {
		form.HRT = request.HRT
	}
	if request.HRTUseLength != nil {
		form.HRTUseLength = request.HRTUseLength
	}
	if request.LastFiveYearsHRTUse != nil {
		form.LastFiveYearsHRTUse = *request.LastFiveYearsHRTUse
	}
	if request.CurrentHRTUse != nil {
		form.CurrentHRTUse = request.CurrentHRTUse
	}
	if request.IntendedHRTUse != nil {
		form.IntendedHRTUse = request.IntendedHRTUse
	}
	if request.HRTType != nil {
		form.HRTType = request.HRTType
	}
	if request.Oral != nil {
		form.Oral = request.Oral
	}
	if request.OralDuration != nil {
		form.OralDuration = request.OralDuration
	}
	if request.OralTwoLastYears != nil {
		form.OralTwoLastYears = request.OralTwoLastYears
	}
	if request.MamoGraphy != nil {
		form.MamoGraphy = request.MamoGraphy
	}
	if request.Falop != nil {
		form.Falop = request.Falop
	}
	if request.Andometrioz != nil {
		form.Andometrioz = request.Andometrioz
	}
	if request.LeavePestan != nil {
		form.LeavePestan = *request.LeavePestan
	}
	if request.LeaveTokhmdan != nil {
		form.LeaveTokhmdan = *request.LeaveTokhmdan
	}
	if request.LaDeColon != nil {
		form.LaDeColon = request.LaDeColon
	}
	if request.LaDePol != nil {
		form.LaDePol = request.LaDePol
	}
	if request.AspLaMo != nil {
		form.AspLaMo = request.AspLaMo
	}
	if request.NsaiDLaMo != nil {
		form.NsaiDLaMo = request.NsaiDLaMo
	}
	if request.LastFiveYearBloodTestInStool != nil {
		form.LastFiveYearBloodTestInStool = request.LastFiveYearBloodTestInStool
	}
	if request.Cancer != nil {
		form.Cancer = *request.Cancer
	}
	if request.CancerType != nil {
		form.CancerType = request.CancerType
	}
	if request.CancerAge != nil {
		form.CancerAge = request.CancerAge
	}
	if request.ChildCancer != nil {
		form.ChildCancer = *request.ChildCancer
	}
	if request.ChildName != nil {
		form.ChildName = request.ChildName
	}
	if request.ChildCancerType != nil {
		form.ChildCancerType = request.ChildCancerType
	}
	if request.ChildCancerAge != nil {
		form.ChildCancerAge = request.ChildCancerAge
	}
	if request.ChildLifeStatus != nil {
		form.ChildLifeStatus = request.ChildLifeStatus
	}

	if request.MotherCancer != nil {
		form.MotherCancer = *request.MotherCancer
	}
	if request.MotherName != nil {
		form.MotherName = request.MotherName
	}
	if request.MotherLifeStatus != nil {
		form.MotherLifeStatus = request.MotherLifeStatus
	}
	if request.MotherCancerType != nil {
		form.MotherCancerType = request.MotherCancerType
	}
	if request.MotherCancerAge != nil {
		form.MotherCancerAge = request.MotherCancerAge
	}

	if request.FatherCancer != nil {
		form.FatherCancer = *request.FatherCancer
	}
	if request.FatherName != nil {
		form.FatherName = request.FatherName
	}
	if request.FatherLifeStatus != nil {
		form.FatherLifeStatus = request.FatherLifeStatus
	}
	if request.FatherCancerType != nil {
		form.FatherCancerType = request.FatherCancerType
	}
	if request.FatherCancerAge != nil {
		form.FatherCancerAge = request.FatherCancerAge
	}

	if request.SiblingCancer != nil {
		form.SiblingCancer = *request.SiblingCancer
	}
	if request.SiblingName != nil {
		form.SiblingName = request.SiblingName
	}
	if request.SiblingLifeStatus != nil {
		form.SiblingLifeStatus = request.SiblingLifeStatus
	}
	if request.SiblingCancerType != nil {
		form.SiblingCancerType = request.SiblingCancerType
	}
	if request.SiblingCancerAge != nil {
		form.SiblingCancerAge = request.SiblingCancerAge
	}

	if request.AmeAmoCancer != nil {
		form.AmeAmoCancer = *request.AmeAmoCancer
	}
	if request.AmeAmoName != nil {
		form.AmeAmoName = request.AmeAmoName
	}
	if request.AmeAmoLifeStatus != nil {
		form.AmeAmoLifeStatus = request.AmeAmoLifeStatus
	}
	if request.AmeAmoCancerType != nil {
		form.AmeAmoCancerType = request.AmeAmoCancerType
	}
	if request.AmeAmoCancerAge != nil {
		form.AmeAmoCancerAge = request.AmeAmoCancerAge
	}

	if request.KhaleDaeiCancer != nil {
		form.KhaleDaeiCancer = *request.KhaleDaeiCancer
	}
	if request.KhaleDaeiName != nil {
		form.KhaleDaeiName = request.KhaleDaeiName
	}
	if request.KhaleDaeiLifeStatus != nil {
		form.KhaleDaeiLifeStatus = request.KhaleDaeiLifeStatus
	}
	if request.KhaleDaeiCancerType != nil {
		form.KhaleDaeiCancerType = request.KhaleDaeiCancerType
	}
	if request.KhaleDaeiCancerAge != nil {
		form.KhaleDaeiCancerAge = request.KhaleDaeiCancerAge
	}

	if request.OtherRelativeCancer != nil {
		form.OtherRelativeCancer = request.OtherRelativeCancer
	}
	if request.OtherRelativeName != nil {
		form.OtherRelativeName = request.OtherRelativeName
	}
	if request.OtherRelativeRelation != nil {
		form.OtherRelativeRelation = request.OtherRelativeRelation
	}
	if request.OtherRelativeLifeStatus != nil {
		form.OtherRelativeLifeStatus = request.OtherRelativeLifeStatus
	}
	if request.OtherRelativeCancerType != nil {
		form.OtherRelativeCancerType = request.OtherRelativeCancerType
	}
	if request.OtherRelativeCancerAge != nil {
		form.OtherRelativeCancerAge = request.OtherRelativeCancerAge
	}
	if request.TestGen != nil {
		form.TestGen = request.TestGen
	}
	if request.FmTestGen != nil {
		form.FmTestGen = request.FmTestGen
	}
	if request.CallExpert != nil {
		form.CallExpert = *request.CallExpert
	}
	if request.BirthCountry != nil {
		form.BirthCountry = request.BirthCountry
	}
	if request.Province != nil {
		form.Province = request.Province
	}
	if request.City != nil {
		form.City = request.City
	}
	if request.Country != nil {
		form.Country = request.Country
	}
	if request.SupplementaryInsurances != nil {
		form.SupplementaryInsurances = request.SupplementaryInsurances
	}
	if request.Hypertension != nil {
		form.Hypertension = *request.Hypertension
	}
	if request.HypertensionTreatment != nil {
		form.HypertensionTreatment = request.HypertensionTreatment
	}
	if request.HeartDisease != nil {
		form.HeartDisease = *request.HeartDisease
	}
	if request.HeartDiseaseTreatment != nil {
		form.HeartDiseaseTreatment = request.HeartDiseaseTreatment
	}
	if request.Diabetes != nil {
		form.Diabetes = *request.Diabetes
	}
	if request.DiabetesTreatment != nil {
		form.DiabetesTreatment = request.DiabetesTreatment
	}
	if request.ChronicLungDisease != nil {
		form.ChronicLungDisease = request.ChronicLungDisease
	}
	if request.ChronicLungDiseaseType != nil {
		form.ChronicLungDiseaseType = request.ChronicLungDiseaseType
	}
	if request.LungCancerHistory != nil {
		form.LungCancerHistory = *request.LungCancerHistory
	}
	if request.OtherCancerHistory != nil {
		form.OtherCancerHistory = *request.OtherCancerHistory
	}
	if request.OtherCancerType != nil {
		form.OtherCancerType = request.OtherCancerType
	}
	if request.LungCancerFamily != nil {
		form.LungCancerFamily = request.LungCancerFamily
	}
	if request.LungCancerFamilyRelation != nil {
		form.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	}
	if request.OtherCancerFamily != nil {
		form.OtherCancerFamily = request.OtherCancerFamily
	}
	if request.OtherCancerFamilyType != nil {
		form.OtherCancerFamilyType = request.OtherCancerFamilyType
	}
	if request.OtherCancerFamilyRelation != nil {
		form.OtherCancerFamilyRelation = request.OtherCancerFamilyRelation
	}
	if request.OccupationalExposure != nil {
		form.OccupationalExposure = request.OccupationalExposure
	}
	if request.CurrentSmoking != nil {
		form.CurrentSmoking = *request.CurrentSmoking
	}
	if request.SmokingStartAgeCurrent != nil {
		form.SmokingStartAgeCurrent = request.SmokingStartAgeCurrent
	}
	if request.SmokingTypesCurrent != nil {
		form.SmokingTypesCurrent = request.SmokingTypesCurrent
	}
	if request.CigarettesPerDayCurrent != nil {
		form.CigarettesPerDayCurrent = request.CigarettesPerDayCurrent
	}
	if request.CigarPerDayCurrent != nil {
		form.CigarPerDayCurrent = request.CigarPerDayCurrent
	}
	if request.ECigPerDayCurrent != nil {
		form.ECigPerDayCurrent = request.ECigPerDayCurrent
	}
	if request.PipePerDayCurrent != nil {
		form.PipePerDayCurrent = request.PipePerDayCurrent
	}
	if request.ChapoghPerDayCurrent != nil {
		form.ChapoghPerDayCurrent = request.ChapoghPerDayCurrent
	}
	if request.SmokedOpiumPerDayCurrent != nil {
		form.SmokedOpiumPerDayCurrent = request.SmokedOpiumPerDayCurrent
	}
	if request.ChewedOpiumPerDayCurrent != nil {
		form.ChewedOpiumPerDayCurrent = request.ChewedOpiumPerDayCurrent
	}
	if request.HookahPerWeekCurrent != nil {
		form.HookahPerWeekCurrent = request.HookahPerWeekCurrent
	}
	if request.PastSmoking != nil {
		form.PastSmoking = request.PastSmoking
	}
	if request.SmokingStartAgePast != nil {
		form.SmokingStartAgePast = request.SmokingStartAgePast
	}
	if request.SmokingTypesPast != nil {
		form.SmokingTypesPast = request.SmokingTypesPast
	}
	if request.CigarettesPerDayPast != nil {
		form.CigarettesPerDayPast = request.CigarettesPerDayPast
	}
	if request.CigarPerDayPast != nil {
		form.CigarPerDayPast = request.CigarPerDayPast
	}
	if request.ECigPerDayPast != nil {
		form.ECigPerDayPast = request.ECigPerDayPast
	}
	if request.PipePerDayPast != nil {
		form.PipePerDayPast = request.PipePerDayPast
	}
	if request.ChapoghPerDayPast != nil {
		form.ChapoghPerDayPast = request.ChapoghPerDayPast
	}
	if request.SmokedOpiumPerDayPast != nil {
		form.SmokedOpiumPerDayPast = request.SmokedOpiumPerDayPast
	}
	if request.ChewedOpiumPerDayPast != nil {
		form.ChewedOpiumPerDayPast = request.ChewedOpiumPerDayPast
	}
	if request.HookahPerWeekPast != nil {
		form.HookahPerWeekPast = request.HookahPerWeekPast
	}
	if request.SecondhandSmoke != nil {
		form.SecondhandSmoke = *request.SecondhandSmoke
	}
	if request.SecondhandSmokeLocation != nil {
		form.SecondhandSmokeLocation = request.SecondhandSmokeLocation
	}

	err = formService.formRepository.UpdateForm(formService.db, form)
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

func (formService *FormService) GetAllForms(offset, limit int) ([]formdto.FormResponse, int64, error) {
	forms, err := formService.formRepository.FindAllForms(formService.db, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllForms(formService.db)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.FormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formService.entityToResponse(form)
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
