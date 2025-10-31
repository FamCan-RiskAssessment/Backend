package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormService struct {
	constants        *bootstrap.Constants
	formRepository   postgres.FormRepository
	userService      usecase.UserService
	actionLogService usecase.ActionLogService
	db               database.Database
}

func NewFormService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	userService usecase.UserService,
	actionLogService usecase.ActionLogService,
	db database.Database,
) *FormService {
	return &FormService{
		constants:        constants,
		formRepository:   formRepository,
		userService:      userService,
		actionLogService: actionLogService,
		db:               db,
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
		Gender:               enum.Gender(uint(request.Gender)),
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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

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
	info.MenopausalStatus = enum.MenopausalStatus(uint(request.MenopausalStatus))
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
	info.NumberOfBreastBiopsies = request.NumberOfBreastBiopsies
	info.HyperplasiaInBiopsy = (*enum.HyperplasiaInBiopsyStatus)(request.HyperplasiaInBiopsy)

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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	if err := formService.formRepository.DeleteAllCancersByFormID(formService.db, request.FormID); err != nil {
		return err
	}

	for _, v := range request.Cancers {
		info := &entity.CancerInfo{FormID: request.FormID}
		info.Cancer = &entity.CancerSpec{
			CancerAge:  v.CancerAge,
			CancerType: enum.CancerType(v.CancerType),
		}

		if err := formService.formRepository.CreateCancer(formService.db, info); err != nil {
			return err
		}
	}
	return nil
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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindFamilyCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.FamilyCancerInfo{FormID: request.FormID}
	}

	info.ChildCancer = request.ChildCancer
	info.ChildName = request.ChildName
	if request.ChildCancerType != nil {
		info.ChildCancerType = (*enum.CancerType)(request.ChildCancerType)
	}
	info.ChildCancerAge = request.ChildCancerAge
	info.ChildLifeStatus = request.ChildLifeStatus

	info.MotherCancer = request.MotherCancer
	info.MotherName = request.MotherName
	info.MotherLifeStatus = request.MotherLifeStatus
	if request.MotherCancerType != nil {
		info.MotherCancerType = (*enum.CancerType)(request.MotherCancerType)
	}
	info.MotherCancerAge = request.MotherCancerAge

	info.FatherCancer = request.FatherCancer
	info.FatherName = request.FatherName
	info.FatherLifeStatus = request.FatherLifeStatus
	if request.FatherCancerType != nil {
		info.FatherCancerType = (*enum.CancerType)(request.FatherCancerType)
	}
	info.FatherCancerAge = request.FatherCancerAge

	info.SiblingCancer = request.SiblingCancer
	info.SiblingName = request.SiblingName
	info.SiblingLifeStatus = request.SiblingLifeStatus
	if request.SiblingCancerType != nil {
		info.SiblingCancerType = (*enum.CancerType)(request.SiblingCancerType)
	}
	info.SiblingCancerAge = request.SiblingCancerAge

	info.AmeAmoCancer = request.AmeAmoCancer
	info.AmeAmoName = request.AmeAmoName
	info.AmeAmoLifeStatus = request.AmeAmoLifeStatus
	if request.AmeAmoCancerType != nil {
		info.AmeAmoCancerType = (*enum.CancerType)(request.AmeAmoCancerType)
	}
	info.AmeAmoCancerAge = request.AmeAmoCancerAge

	info.KhaleDaeiCancer = request.KhaleDaeiCancer
	info.KhaleDaeiName = request.KhaleDaeiName
	info.KhaleDaeiLifeStatus = request.KhaleDaeiLifeStatus
	if request.KhaleDaeiCancerType != nil {
		info.KhaleDaeiCancerType = (*enum.CancerType)(request.KhaleDaeiCancerType)
	}
	info.KhaleDaeiCancerAge = request.KhaleDaeiCancerAge

	info.OtherRelativeCancer = request.OtherRelativeCancer
	info.OtherRelativeName = request.OtherRelativeName
	info.OtherRelativeRelation = request.OtherRelativeRelation
	info.OtherRelativeLifeStatus = request.OtherRelativeLifeStatus
	if request.OtherRelativeCancerType != nil {
		info.OtherRelativeCancerType = (*enum.CancerType)(request.OtherRelativeCancerType)
	}
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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

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
	if request.OtherCancerType != nil {
		info.OtherCancerType = (*enum.CancerType)(request.OtherCancerType)
	}
	info.LungCancerFamily = request.LungCancerFamily
	info.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	info.OtherCancerFamily = request.OtherCancerFamily
	if request.OtherCancerFamilyType != nil {
		info.OtherCancerFamilyType = (*enum.CancerType)(request.OtherCancerFamilyType)
	}
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

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.ChangeFormStatusResponse{}, ForbiddenError
	// }

	form.Status = enum.FormStatusReady
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return formdto.ChangeFormStatusResponse{}, err
	}

	response := formdto.ChangeFormStatusResponse{
		Form: formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			OperatorID: form.OperatorID,
			UserID:     form.UserID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		},
	}

	return response, nil
}

func (formService *FormService) GetBasicForm(request formdto.GetPartialFormRequest) (formdto.GetBasicFormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetBasicFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetBasicFormResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetBasicFormResponse{}, ForbiddenError
	// }

	basic, err := formService.formRepository.FindBasicInfoByFormID(formService.db, request.FormID)
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

func (formService *FormService) GetGeneralHealth(request formdto.GetPartialFormRequest) (formdto.GetGeneralHealthResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetGeneralHealthResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetGeneralHealthResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetGeneralHealthResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
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

func (formService *FormService) GetMamography(request formdto.GetPartialFormRequest) (formdto.GetMamographyResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetMamographyResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetMamographyResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetMamographyResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
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

func (formService *FormService) GetCancers(request formdto.GetPartialFormRequest) (formdto.GetCancersResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetCancersResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetCancersResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindAllCancersByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetCancersResponse{}, err
	}

	cancersResponse := formdto.GetCancersResponse{}
	if len(info) == 0 {
		cancersResponse.Cancer = false
		cancersResponse.Cancers = make([]*formdto.CancerResponse, len(info))
		return cancersResponse, nil
	}

	cancersResponse.Cancer = true
	for _, v := range info {
		cancersResponse.Cancers = append(cancersResponse.Cancers, &formdto.CancerResponse{ID: v.ID, CancerType: v.Cancer.CancerType, CancerAge: v.Cancer.CancerAge})
	}
	return cancersResponse, nil
}

func (formService *FormService) GetFamilyCancer(request formdto.GetPartialFormRequest) (formdto.GetFamilyCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetFamilyCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetFamilyCancerResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetFamilyCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindFamilyCancerByFormID(formService.db, request.FormID)
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

func (formService *FormService) GetContact(request formdto.GetPartialFormRequest) (formdto.GetContactResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetContactResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	return formdto.GetContactResponse{
		ID:           info.ID,
		Name:         info.Name,
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

func (formService *FormService) GetLungCancer(request formdto.GetPartialFormRequest) (formdto.GetLungCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetLungCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetLungCancerResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetLungCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
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
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
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

	info, err := formService.formRepository.FindBasicInfoByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	if request.BirthDay != nil {
		info.BirthDay = *request.BirthDay
	}
	if request.BirthMonth != nil {
		info.BirthMonth = *request.BirthMonth
	}
	if request.SocialSecurityNumber != nil {
		info.SocialSecurityNumber = *request.SocialSecurityNumber
	}
	if request.Gender != nil {
		info.Gender = enum.Gender(uint(*request.Gender))
	}
	if request.IsAtba != nil {
		info.IsAtba = *request.IsAtba
	}
	if request.Height != nil {
		info.Height = *request.Height
	}
	if request.Weight != nil {
		info.Weight = *request.Weight
	}

	err = formService.formRepository.UpdateBasicInfo(formService.db, info)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) DeleteForm(request formdto.DeleteFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	err = formService.formRepository.DeleteForm(formService.db, request.FormID)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.BasicFormResponse, int64, error) {
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
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) GetAllOperatorForms(offset, limit int, filters *postgres.OperatorFormFilters) ([]formdto.BasicFormResponse, int64, error) {
	operator, err := formService.userService.GetUserByID(filters.OperatorID)
	if err != nil {
		return nil, 0, err
	}
	if operator == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return nil, 0, notFoundError
	}
	userRoles, err := formService.userService.GetUserRoles(filters.OperatorID)
	if err != nil {
		return nil, 0, err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Role}
		return nil, 0, forbiddenError
	}

	forms, err := formService.formRepository.FindAllOperatorForms(formService.db, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllOperatorForms(formService.db, filters)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
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

func (formService *FormService) UpdateGeneralHealth(request formdto.UpdateGeneralHealthRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}

	if info == nil {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info.DrinksAlcohol = request.DrinksAlcohol
	info.CupsPerWeek = request.CupsPerWeek
	if request.LastMonthSabzijatMeal != nil {
		info.LastMonthSabzijatMeal = *request.LastMonthSabzijatMeal
	}
	if request.LastMonthSabzijatWeight != nil {
		info.LastMonthSabzijatWeight = *request.LastMonthSabzijatWeight
	}
	if request.MediumActivityMonthInYear != nil {
		info.MediumActivityMonthInYear = *request.MediumActivityMonthInYear
	}
	if request.MediumActivityHourInWeek != nil {
		info.MediumActivityHourInWeek = *request.MediumActivityHourInWeek
	}
	if request.HardActivityMonthInYear != nil {
		info.HardActivityMonthInYear = *request.HardActivityMonthInYear
	}
	if request.HardActivityHourInWeek != nil {
		info.HardActivityHourInWeek = *request.HardActivityHourInWeek
	}
	info.SmokeAtLeast100 = request.SmokeAtLeast100
	info.SmokingAge = request.SmokingAge
	if request.SmokingNow != nil {
		info.SmokingNow = *request.SmokingNow
	}
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

func (formService *FormService) UpdateMamography(request formdto.UpdateMamographyRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	if request.GhaedeAge != nil {
		info.GhaedeAge = *request.GhaedeAge
	}
	if request.HasChildren != nil {
		info.HasChildren = *request.HasChildren
	}
	info.NumberOfChildren = request.NumberOfChildren
	info.AgeOfFirstBirth = request.AgeOfFirstBirth
	if request.MenopausalStatus != nil {
		info.MenopausalStatus = enum.MenopausalStatus(uint(*request.MenopausalStatus))
	}
	info.MenopauseAge = request.MenopauseAge
	info.HRT = request.HRT
	info.HRTUseLength = request.HRTUseLength
	if request.LastFiveYearsHRTUse != nil {
		info.LastFiveYearsHRTUse = *request.LastFiveYearsHRTUse
	}
	info.CurrentHRTUse = request.CurrentHRTUse
	info.IntendedHRTUse = request.IntendedHRTUse
	info.HRTType = request.HRTType
	info.Oral = request.Oral
	info.OralDuration = request.OralDuration
	info.OralTwoLastYears = request.OralTwoLastYears
	info.MamoGraphy = request.MamoGraphy
	info.Falop = request.Falop
	info.Andometrioz = request.Andometrioz
	if request.LeavePestan != nil {
		info.LeavePestan = *request.LeavePestan
	}
	if request.LeaveTokhmdan != nil {
		info.LeaveTokhmdan = *request.LeaveTokhmdan
	}
	info.LaDeColon = request.LaDeColon
	info.LaDePol = request.LaDePol
	info.AspLaMo = request.AspLaMo
	info.NsaiDLaMo = request.NsaiDLaMo
	info.LastFiveYearBloodTestInStool = request.LastFiveYearBloodTestInStool
	if request.NumberOfBreastBiopsies != nil {
		info.NumberOfBreastBiopsies = request.NumberOfBreastBiopsies
	}
	if request.HyperplasiaInBiopsy != nil {
		info.HyperplasiaInBiopsy = (*enum.HyperplasiaInBiopsyStatus)(request.HyperplasiaInBiopsy)
	}

	if info.ID == 0 {
		return formService.formRepository.CreateMamography(formService.db, info)
	}
	return formService.formRepository.UpdateMamography(formService.db, info)
}

func (formService *FormService) UpdateCancer(request formdto.UpdateCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	for _, v := range request.Cancers {
		info := &entity.CancerInfo{FormID: request.FormID}
		info.Cancer = &entity.CancerSpec{
			CancerAge:  v.CancerAge,
			CancerType: enum.CancerType(v.CancerType),
		}

		if err := formService.formRepository.CreateCancer(formService.db, info); err != nil {
			return err
		}
	}
	return nil
}

func (formService *FormService) UpdateFamilyCancer(request formdto.UpdateFamilyCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindFamilyCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	if request.ChildCancer != nil {
		info.ChildCancer = *request.ChildCancer
	}
	info.ChildName = request.ChildName
	if request.ChildCancerType != nil {
		info.ChildCancerType = (*enum.CancerType)(request.ChildCancerType)
	}
	info.ChildCancerAge = request.ChildCancerAge
	info.ChildLifeStatus = request.ChildLifeStatus

	if request.MotherCancer != nil {
		info.MotherCancer = *request.MotherCancer
	}
	info.MotherName = request.MotherName
	info.MotherLifeStatus = request.MotherLifeStatus
	if request.MotherCancerType != nil {
		info.MotherCancerType = (*enum.CancerType)(request.MotherCancerType)
	}
	info.MotherCancerAge = request.MotherCancerAge

	if request.FatherCancer != nil {
		info.FatherCancer = *request.FatherCancer
	}
	info.FatherName = request.FatherName
	info.FatherLifeStatus = request.FatherLifeStatus
	if request.FatherCancerType != nil {
		info.FatherCancerType = (*enum.CancerType)(request.FatherCancerType)
	}
	info.FatherCancerAge = request.FatherCancerAge

	if request.SiblingCancer != nil {
		info.SiblingCancer = *request.SiblingCancer
	}
	info.SiblingName = request.SiblingName
	info.SiblingLifeStatus = request.SiblingLifeStatus
	if request.SiblingCancerType != nil {
		info.SiblingCancerType = (*enum.CancerType)(request.SiblingCancerType)
	}
	info.SiblingCancerAge = request.SiblingCancerAge

	if request.AmeAmoCancer != nil {
		info.AmeAmoCancer = *request.AmeAmoCancer
	}
	info.AmeAmoName = request.AmeAmoName
	info.AmeAmoLifeStatus = request.AmeAmoLifeStatus
	if request.AmeAmoCancerType != nil {
		info.AmeAmoCancerType = (*enum.CancerType)(request.AmeAmoCancerType)
	}
	info.AmeAmoCancerAge = request.AmeAmoCancerAge

	if request.KhaleDaeiCancer != nil {
		info.KhaleDaeiCancer = *request.KhaleDaeiCancer
	}
	info.KhaleDaeiName = request.KhaleDaeiName
	info.KhaleDaeiLifeStatus = request.KhaleDaeiLifeStatus
	if request.KhaleDaeiCancerType != nil {
		info.KhaleDaeiCancerType = (*enum.CancerType)(request.KhaleDaeiCancerType)
	}
	info.KhaleDaeiCancerAge = request.KhaleDaeiCancerAge

	info.OtherRelativeCancer = request.OtherRelativeCancer
	info.OtherRelativeName = request.OtherRelativeName
	info.OtherRelativeRelation = request.OtherRelativeRelation
	info.OtherRelativeLifeStatus = request.OtherRelativeLifeStatus
	if request.OtherRelativeCancerType != nil {
		info.OtherRelativeCancerType = (*enum.CancerType)(request.OtherRelativeCancerType)
	}
	info.OtherRelativeCancerAge = request.OtherRelativeCancerAge

	if info.ID == 0 {
		return formService.formRepository.CreateFamilyCancer(formService.db, info)
	}
	return formService.formRepository.UpdateFamilyCancer(formService.db, info)
}

func (formService *FormService) UpdateContact(request formdto.UpdateContactRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.ContactInfo{FormID: request.FormID}
	}

	if request.Name != nil {
		info.Name = *request.Name
	}
	info.TestGen = request.TestGen
	info.FmTestGen = request.FmTestGen
	if request.CallExpert != nil {
		info.CallExpert = *request.CallExpert
	}
	info.BirthCountry = request.BirthCountry
	info.Province = request.Province
	info.City = request.City
	info.Country = request.Country
	if request.Address != nil {
		info.Address = *request.Address
	}
	if request.PostalCode != nil {
		info.PostalCode = *request.PostalCode
	}

	if info.ID == 0 {
		return formService.formRepository.CreateContact(formService.db, info)
	}
	return formService.formRepository.UpdateContact(formService.db, info)
}

func (formService *FormService) UpdateLungCancer(request formdto.UpdateLungCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.LungCancerInfo{FormID: request.FormID}
	}

	info.InsuranceStatus = request.InsuranceStatus
	info.SupplementaryInsurances = request.SupplementaryInsurances
	if request.Hypertension != nil {
		info.Hypertension = *request.Hypertension
	}
	info.HypertensionTreatment = request.HypertensionTreatment
	if request.HeartDisease != nil {
		info.HeartDisease = *request.HeartDisease
	}
	info.HeartDiseaseTreatment = request.HeartDiseaseTreatment
	if request.Diabetes != nil {
		info.Diabetes = *request.Diabetes
	}
	info.DiabetesTreatment = request.DiabetesTreatment
	info.ChronicLungDisease = request.ChronicLungDisease
	info.ChronicLungDiseaseType = request.ChronicLungDiseaseType
	if request.LungCancerHistory != nil {
		info.LungCancerHistory = *request.LungCancerHistory
	}
	if request.OtherCancerHistory != nil {
		info.OtherCancerHistory = *request.OtherCancerHistory
	}
	if request.OtherCancerType != nil {
		info.OtherCancerType = (*enum.CancerType)(request.OtherCancerType)
	}
	info.LungCancerFamily = request.LungCancerFamily
	info.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	info.OtherCancerFamily = request.OtherCancerFamily
	if request.OtherCancerFamilyType != nil {
		info.OtherCancerFamilyType = (*enum.CancerType)(request.OtherCancerFamilyType)
	}
	info.OtherCancerFamilyRelation = request.OtherCancerFamilyRelation
	info.OccupationalExposure = request.OccupationalExposure
	if request.CurrentSmoking != nil {
		info.CurrentSmoking = *request.CurrentSmoking
	}
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
	if request.SecondhandSmoke != nil {
		info.SecondhandSmoke = *request.SecondhandSmoke
	}
	info.SecondhandSmokeLocation = request.SecondhandSmokeLocation

	if info.ID == 0 {
		return formService.formRepository.CreateLungCancer(formService.db, info)
	}
	return formService.formRepository.UpdateLungCancer(formService.db, info)
}

func (formService *FormService) AssignOperator(request formdto.AssignOperatorRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	operator, err := formService.userService.GetUserByID(request.OperatorID)
	if err != nil {
		return err
	}
	if operator == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return notFoundError
	}
	userRoles, err := formService.userService.GetUserRoles(request.OperatorID)
	if err != nil {
		return err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Role}
		return forbiddenError
	}

	form.OperatorID = &request.OperatorID
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	log := actionlogdto.LogAction{
		ActorID:    request.UserID,
		TargetID:   &request.OperatorID,
		Action:     enum.ActionTypeFormAssigned,
		ResourceID: &request.FormID,
		Details:    "فرم به اپراتور اساین شد",
	}
	formService.actionLogService.LogAction(log)

	return nil
}
func (formService *FormService) UnassignOperator(request formdto.UnassignOperatorRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	form.OperatorID = nil
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	log := actionlogdto.LogAction{
		ActorID:    request.UserID,
		ResourceID: &request.FormID,
		Action:     enum.ActionTypeFormUnAssigned,
		Details:    "فرم از اپراتور گرفته شد",
	}
	formService.actionLogService.LogAction(log)

	return nil
}

func (formService *FormService) GetAllCancerTypes() ([]generaldto.EnumResponse, error) {
	cancerTypes := enum.GetAllCancerTypes()
	response := make([]generaldto.EnumResponse, len(cancerTypes))
	for i, cancerType := range cancerTypes {
		response[i] = generaldto.EnumResponse{
			ID:   uint(cancerType),
			Name: cancerType.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllGenders() ([]generaldto.EnumResponse, error) {
	genders := enum.GetAllGenders()
	response := make([]generaldto.EnumResponse, len(genders))
	for i, gender := range genders {
		response[i] = generaldto.EnumResponse{
			ID:   uint(gender),
			Name: gender.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllMenopausalStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllMenopausalStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllFormStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllFormStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllHyperplasiaInBiopsyStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllHyperplasiaInBiopsyStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}
