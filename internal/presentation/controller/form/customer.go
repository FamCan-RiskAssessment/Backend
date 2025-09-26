package form

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
	pagination  *bootstrap.Pagination
}

func NewCustomerFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
	pagination *bootstrap.Pagination,
) *CustomerFormController {
	return &CustomerFormController{
		constants:   constants,
		formService: formService,
		pagination:  pagination,
	}
}

func (formController *CustomerFormController) CreateForm(ctx *gin.Context) {
	type CreateFormParams struct {
		Name                 string  `json:"name" validate:"required"`
		BirthDay             uint    `json:"birthDay" validate:"required"`
		BirthMonth           string  `json:"birthMonth" validate:"required"`
		BirthYear            uint    `json:"birthYear" validate:"required"`
		Address              string  `json:"address" validate:"required"`
		PostalCode           string  `json:"postalCode" validate:"required"`
		SocialSecurityNumber string  `json:"socialSecurityNumber" validate:"required"`
		Gender               string  `json:"gender" validate:"required"`
		IsAtba               bool    `json:"isAtba"`
		Height               float64 `json:"height" validate:"required"`
		Weight               float64 `json:"weight" validate:"required"`

		DrinksAlcohol             *bool   `json:"drinksAlcohol,omitempty"`
		CupsPerWeek               *string `json:"cupsPerWeek,omitempty"`
		LastMonthSabzijatMeal     string  `json:"lastMonthSabzijatMeal" validate:"required"`
		LastMonthSabzijatWeight   string  `json:"lastMonthSabzijatWeight" validate:"required"`
		MediumActivityMonthInYear uint    `json:"mediumActivityMonthInYear" validate:"required"`
		MediumActivityHourInWeek  string  `json:"mediumActivityHourInWeek" validate:"required"`
		HardActivityMonthInYear   uint    `json:"hardActivityMonthInYear" validate:"required"`
		HardActivityHourInWeek    string  `json:"hardActivityHourInWeek" validate:"required"`
		SmokeAtLeast100           *bool   `json:"smokeAtLeast100,omitempty"`
		SmokingAge                *uint   `json:"smokingAge,omitempty"`
		SmokingNow                bool    `json:"smokingNow"`
		LeaveSmokingAge           *uint   `json:"leaveSmokingAge,omitempty"`
		CountSmokingDaily         *string `json:"countSmokingDaily,omitempty"`
		CountGheliandaily         *string `json:"countGheliandaily,omitempty"`
		CountSmokingDailyPast     *string `json:"countSmokingDailyPast,omitempty"`
		CountGheliandailyPast     *string `json:"countGheliandailyPast,omitempty"`

		GhaedeAge                    uint    `json:"ghaedeAge" validate:"required"`
		HasChildren                  bool    `json:"hasChildren"`
		NumberOfChildren             *uint   `json:"numberOfChildren,omitempty"`
		AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth,omitempty"`
		MenopausalStatus             string  `json:"menopausalStatus" validate:"required"`
		MenopauseAge                 *string `json:"menopauseAge,omitempty"`
		HRT                          *bool   `json:"hrt,omitempty"`
		HRTUseLength                 *uint   `json:"hrtUseLength,omitempty"`
		LastFiveYearsHRTUse          bool    `json:"lastFiveYearsHrtUse"`
		CurrentHRTUse                *bool   `json:"currentHrtUse,omitempty"`
		IntendedHRTUse               *uint   `json:"intendedHrtUse,omitempty"`
		HRTType                      *string `json:"hrtType,omitempty"`
		Oral                         *bool   `json:"oral,omitempty"`
		OralDuration                 *string `json:"oralDuration,omitempty"`
		OralTwoLastYears             *bool   `json:"oralTwoLastYears,omitempty"`
		MamoGraphy                   *bool   `json:"mamoGraphy,omitempty"`
		Falop                        *bool   `json:"falop,omitempty"`
		Andometrioz                  *bool   `json:"andometrioz,omitempty"`
		LeavePestan                  bool    `json:"leavePestan"`
		LeaveTokhmdan                bool    `json:"leaveTokhmdan"`
		LaDeColon                    *bool   `json:"laDeColon,omitempty"`
		LaDePol                      *bool   `json:"laDePol,omitempty"`
		AspLaMo                      *bool   `json:"aspLaMo,omitempty"`
		NsaiDLaMo                    *bool   `json:"nsaiDLaMo,omitempty"`
		LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool,omitempty"`

		Cancer     bool    `json:"cancer"`
		CancerType *string `json:"cancerType,omitempty"`
		CancerAge  *uint   `json:"cancerAge,omitempty"`

		ChildCancer     bool    `json:"childCancer"`
		ChildName       *string `json:"childName,omitempty"`
		ChildCancerType *string `json:"childCancerType,omitempty"`
		ChildCancerAge  *uint   `json:"childCancerAge,omitempty"`
		ChildLifeStatus *string `json:"childLifeStatus,omitempty"`

		MotherCancer     bool    `json:"motherCancer"`
		MotherName       *string `json:"motherName,omitempty"`
		MotherLifeStatus *string `json:"motherLifeStatus,omitempty"`
		MotherCancerType *string `json:"motherCancerType,omitempty"`
		MotherCancerAge  *uint   `json:"motherCancerAge,omitempty"`

		FatherCancer     bool    `json:"fatherCancer"`
		FatherName       *string `json:"fatherName,omitempty"`
		FatherLifeStatus *string `json:"fatherLifeStatus,omitempty"`
		FatherCancerType *string `json:"fatherCancerType,omitempty"`
		FatherCancerAge  *uint   `json:"fatherCancerAge,omitempty"`

		SiblingCancer     bool    `json:"siblingCancer"`
		SiblingName       *string `json:"siblingName,omitempty"`
		SiblingLifeStatus *string `json:"siblingLifeStatus,omitempty"`
		SiblingCancerType *string `json:"siblingCancerType,omitempty"`
		SiblingCancerAge  *uint   `json:"siblingCancerAge,omitempty"`

		AmeAmoCancer     bool    `json:"ameAmoCancer"`
		AmeAmoName       *string `json:"ameAmoName,omitempty"`
		AmeAmoLifeStatus *string `json:"ameAmoLifeStatus,omitempty"`
		AmeAmoCancerType *string `json:"ameAmoCancerType,omitempty"`
		AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge,omitempty"`

		KhaleDaeiCancer     bool    `json:"khaleDaeiCancer"`
		KhaleDaeiName       *string `json:"khaleDaeiName,omitempty"`
		KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus,omitempty"`
		KhaleDaeiCancerType *string `json:"khaleDaeiCancerType,omitempty"`
		KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge,omitempty"`

		OtherRelativeCancer     *bool   `json:"otherRelativeCancer,omitempty"`
		OtherRelativeName       *string `json:"otherRelativeName,omitempty"`
		OtherRelativeRelation   *string `json:"otherRelativeRelation,omitempty"`
		OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus,omitempty"`
		OtherRelativeCancerType *string `json:"otherRelativeCancerType,omitempty"`
		OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge,omitempty"`

		TestGen   *bool `json:"testGen,omitempty"`
		FmTestGen *bool `json:"fmTestGen,omitempty"`

		CallExpert   bool    `json:"callExpert"`
		BirthCountry *string `json:"birthCountry,omitempty"`
		Province     *string `json:"province,omitempty"`
		City         *string `json:"city,omitempty"`
		Country      *string `json:"country,omitempty"`

		InsuranceStatus           *string `json:"insuranceStatus,omitempty"`
		SupplementaryInsurances   *string `json:"supplementaryInsurances,omitempty"`
		Hypertension              bool    `json:"hypertension"`
		HypertensionTreatment     *bool   `json:"hypertensionTreatment,omitempty"`
		HeartDisease              bool    `json:"heartDisease"`
		HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment,omitempty"`
		Diabetes                  bool    `json:"diabetes"`
		DiabetesTreatment         *bool   `json:"diabetesTreatment,omitempty"`
		ChronicLungDisease        *bool   `json:"chronicLungDisease,omitempty"`
		ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType,omitempty"`
		LungCancerHistory         bool    `json:"lungCancerHistory"`
		OtherCancerHistory        bool    `json:"otherCancerHistory"`
		OtherCancerType           *string `json:"otherCancerType,omitempty"`
		LungCancerFamily          *bool   `json:"lungCancerFamily,omitempty"`
		LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation,omitempty"`
		OtherCancerFamily         *bool   `json:"otherCancerFamily,omitempty"`
		OtherCancerFamilyType     *string `json:"otherCancerFamilyType,omitempty"`
		OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation,omitempty"`
		OccupationalExposure      *string `json:"occupationalExposure,omitempty"`
		CurrentSmoking            bool    `json:"currentSmoking"`
		SmokingStartAgeCurrent    *uint   `json:"smokingStartAgeCurrent,omitempty"`
		SmokingTypesCurrent       *string `json:"smokingTypesCurrent,omitempty"`
		CigarettesPerDayCurrent   *uint   `json:"cigarettesPerDayCurrent,omitempty"`
		CigarPerDayCurrent        *uint   `json:"cigarPerDayCurrent,omitempty"`
		ECigPerDayCurrent         *uint   `json:"eCigPerDayCurrent,omitempty"`
		PipePerDayCurrent         *uint   `json:"pipePerDayCurrent,omitempty"`
		ChapoghPerDayCurrent      *uint   `json:"chapoghPerDayCurrent,omitempty"`
		SmokedOpiumPerDayCurrent  *uint   `json:"smokedOpiumPerDayCurrent,omitempty"`
		ChewedOpiumPerDayCurrent  *uint   `json:"chewedOpiumPerDayCurrent,omitempty"`
		HookahPerWeekCurrent      *uint   `json:"hookahPerWeekCurrent,omitempty"`
		PastSmoking               *string `json:"pastSmoking,omitempty"`
		SmokingStartAgePast       *uint   `json:"smokingStartAgePast,omitempty"`
		SmokingTypesPast          *string `json:"smokingTypesPast,omitempty"`
		CigarettesPerDayPast      *uint   `json:"cigarettesPerDayPast,omitempty"`
		CigarPerDayPast           *uint   `json:"cigarPerDayPast,omitempty"`
		ECigPerDayPast            *uint   `json:"eCigPerDayPast,omitempty"`
		PipePerDayPast            *uint   `json:"pipePerDayPast,omitempty"`
		ChapoghPerDayPast         *uint   `json:"chapoghPerDayPast,omitempty"`
		SmokedOpiumPerDayPast     *uint   `json:"smokedOpiumPerDayPast,omitempty"`
		ChewedOpiumPerDayPast     *uint   `json:"chewedOpiumPerDayPast,omitempty"`
		HookahPerWeekPast         *uint   `json:"hookahPerWeekPast,omitempty"`
		SecondhandSmoke           bool    `json:"secondhandSmoke"`
		SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation,omitempty"`
	}

	params := controller.Validate[CreateFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.CreateFormRequest{
		UserID:               userID.(uint),
		Name:                 params.Name,
		BirthDay:             params.BirthDay,
		BirthMonth:           params.BirthMonth,
		BirthYear:            params.BirthYear,
		Address:              params.Address,
		PostalCode:           params.PostalCode,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,

		DrinksAlcohol:             params.DrinksAlcohol,
		CupsPerWeek:               params.CupsPerWeek,
		LastMonthSabzijatMeal:     params.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   params.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: params.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  params.MediumActivityHourInWeek,
		HardActivityMonthInYear:   params.HardActivityMonthInYear,
		HardActivityHourInWeek:    params.HardActivityHourInWeek,
		SmokeAtLeast100:           params.SmokeAtLeast100,
		SmokingAge:                params.SmokingAge,
		SmokingNow:                params.SmokingNow,
		LeaveSmokingAge:           params.LeaveSmokingAge,
		CountSmokingDaily:         params.CountSmokingDaily,
		CountGheliandaily:         params.CountGheliandaily,
		CountSmokingDailyPast:     params.CountSmokingDailyPast,
		CountGheliandailyPast:     params.CountGheliandailyPast,

		GhaedeAge:                    params.GhaedeAge,
		HasChildren:                  params.HasChildren,
		NumberOfChildren:             params.NumberOfChildren,
		AgeOfFirstBirth:              params.AgeOfFirstBirth,
		MenopausalStatus:             params.MenopausalStatus,
		MenopauseAge:                 params.MenopauseAge,
		HRT:                          params.HRT,
		HRTUseLength:                 params.HRTUseLength,
		LastFiveYearsHRTUse:          params.LastFiveYearsHRTUse,
		CurrentHRTUse:                params.CurrentHRTUse,
		IntendedHRTUse:               params.IntendedHRTUse,
		HRTType:                      params.HRTType,
		Oral:                         params.Oral,
		OralDuration:                 params.OralDuration,
		OralTwoLastYears:             params.OralTwoLastYears,
		MamoGraphy:                   params.MamoGraphy,
		Falop:                        params.Falop,
		Andometrioz:                  params.Andometrioz,
		LeavePestan:                  params.LeavePestan,
		LeaveTokhmdan:                params.LeaveTokhmdan,
		LaDeColon:                    params.LaDeColon,
		LaDePol:                      params.LaDePol,
		AspLaMo:                      params.AspLaMo,
		NsaiDLaMo:                    params.NsaiDLaMo,
		LastFiveYearBloodTestInStool: params.LastFiveYearBloodTestInStool,

		Cancer:     params.Cancer,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,

		ChildCancer:     params.ChildCancer,
		ChildName:       params.ChildName,
		ChildCancerType: params.ChildCancerType,
		ChildCancerAge:  params.ChildCancerAge,
		ChildLifeStatus: params.ChildLifeStatus,

		MotherCancer:     params.MotherCancer,
		MotherName:       params.MotherName,
		MotherLifeStatus: params.MotherLifeStatus,
		MotherCancerType: params.MotherCancerType,
		MotherCancerAge:  params.MotherCancerAge,

		FatherCancer:     params.FatherCancer,
		FatherName:       params.FatherName,
		FatherLifeStatus: params.FatherLifeStatus,
		FatherCancerType: params.FatherCancerType,
		FatherCancerAge:  params.FatherCancerAge,

		SiblingCancer:     params.SiblingCancer,
		SiblingName:       params.SiblingName,
		SiblingLifeStatus: params.SiblingLifeStatus,
		SiblingCancerType: params.SiblingCancerType,
		SiblingCancerAge:  params.SiblingCancerAge,

		AmeAmoCancer:     params.AmeAmoCancer,
		AmeAmoName:       params.AmeAmoName,
		AmeAmoLifeStatus: params.AmeAmoLifeStatus,
		AmeAmoCancerType: params.AmeAmoCancerType,
		AmeAmoCancerAge:  params.AmeAmoCancerAge,

		KhaleDaeiCancer:     params.KhaleDaeiCancer,
		KhaleDaeiName:       params.KhaleDaeiName,
		KhaleDaeiLifeStatus: params.KhaleDaeiLifeStatus,
		KhaleDaeiCancerType: params.KhaleDaeiCancerType,
		KhaleDaeiCancerAge:  params.KhaleDaeiCancerAge,

		OtherRelativeCancer:     params.OtherRelativeCancer,
		OtherRelativeName:       params.OtherRelativeName,
		OtherRelativeRelation:   params.OtherRelativeRelation,
		OtherRelativeLifeStatus: params.OtherRelativeLifeStatus,
		OtherRelativeCancerType: params.OtherRelativeCancerType,
		OtherRelativeCancerAge:  params.OtherRelativeCancerAge,

		TestGen:   params.TestGen,
		FmTestGen: params.FmTestGen,

		CallExpert:   params.CallExpert,
		BirthCountry: params.BirthCountry,
		Province:     params.Province,
		City:         params.City,
		Country:      params.Country,

		InsuranceStatus:           params.InsuranceStatus,
		SupplementaryInsurances:   params.SupplementaryInsurances,
		Hypertension:              params.Hypertension,
		HypertensionTreatment:     params.HypertensionTreatment,
		HeartDisease:              params.HeartDisease,
		HeartDiseaseTreatment:     params.HeartDiseaseTreatment,
		Diabetes:                  params.Diabetes,
		DiabetesTreatment:         params.DiabetesTreatment,
		ChronicLungDisease:        params.ChronicLungDisease,
		ChronicLungDiseaseType:    params.ChronicLungDiseaseType,
		LungCancerHistory:         params.LungCancerHistory,
		OtherCancerHistory:        params.OtherCancerHistory,
		OtherCancerType:           params.OtherCancerType,
		LungCancerFamily:          params.LungCancerFamily,
		LungCancerFamilyRelation:  params.LungCancerFamilyRelation,
		OtherCancerFamily:         params.OtherCancerFamily,
		OtherCancerFamilyType:     params.OtherCancerFamilyType,
		OtherCancerFamilyRelation: params.OtherCancerFamilyRelation,
		OccupationalExposure:      params.OccupationalExposure,
		CurrentSmoking:            params.CurrentSmoking,
		SmokingStartAgeCurrent:    params.SmokingStartAgeCurrent,
		SmokingTypesCurrent:       params.SmokingTypesCurrent,
		CigarettesPerDayCurrent:   params.CigarettesPerDayCurrent,
		CigarPerDayCurrent:        params.CigarPerDayCurrent,
		ECigPerDayCurrent:         params.ECigPerDayCurrent,
		PipePerDayCurrent:         params.PipePerDayCurrent,
		ChapoghPerDayCurrent:      params.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:  params.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:  params.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:      params.HookahPerWeekCurrent,
		PastSmoking:               params.PastSmoking,
		SmokingStartAgePast:       params.SmokingStartAgePast,
		SmokingTypesPast:          params.SmokingTypesPast,
		CigarettesPerDayPast:      params.CigarettesPerDayPast,
		CigarPerDayPast:           params.CigarPerDayPast,
		ECigPerDayPast:            params.ECigPerDayPast,
		PipePerDayPast:            params.PipePerDayPast,
		ChapoghPerDayPast:         params.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:     params.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:     params.ChewedOpiumPerDayPast,
		HookahPerWeekPast:         params.HookahPerWeekPast,
		SecondhandSmoke:           params.SecondhandSmoke,
		SecondhandSmokeLocation:   params.SecondhandSmokeLocation,
	}

	err := formController.formService.CreateForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, nil)
}

func (formController *CustomerFormController) GetUserForms(ctx *gin.Context) {
	type GetUserFormsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	params := controller.Validate[GetUserFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetUserFormsRequest{
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}

	forms, count, err := formController.formService.GetUserForms(request)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(forms, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (formController *CustomerFormController) GetForm(ctx *gin.Context) {
	type GetFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[GetFormParams](ctx)

	response, err := formController.formService.GetForm(params.FormID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) UpdateForm(ctx *gin.Context) {
	type UpdateFormParams struct {
		FormID               uint     `json:"form_id" validate:"required"`
		Name                 *string  `json:"name,omitempty"`
		BirthDay             *uint    `json:"birthDay,omitempty"`
		BirthMonth           *string  `json:"birthMonth,omitempty"`
		BirthYear            *uint    `json:"birthYear,omitempty"`
		Address              *string  `json:"address,omitempty"`
		PostalCode           *string  `json:"postalCode,omitempty"`
		SocialSecurityNumber *string  `json:"socialSecurityNumber,omitempty"`
		Gender               *string  `json:"gender,omitempty"`
		IsAtba               *bool    `json:"isAtba,omitempty"`
		Height               *float64 `json:"height,omitempty"`
		Weight               *float64 `json:"weight,omitempty"`

		DrinksAlcohol             *bool   `json:"drinksAlcohol,omitempty"`
		CupsPerWeek               *string `json:"cupsPerWeek,omitempty"`
		LastMonthSabzijatMeal     *string `json:"lastMonthSabzijatMeal,omitempty"`
		LastMonthSabzijatWeight   *string `json:"lastMonthSabzijatWeight,omitempty"`
		MediumActivityMonthInYear *uint   `json:"mediumActivityMonthInYear,omitempty"`
		MediumActivityHourInWeek  *string `json:"mediumActivityHourInWeek,omitempty"`
		HardActivityMonthInYear   *uint   `json:"hardActivityMonthInYear,omitempty"`
		HardActivityHourInWeek    *string `json:"hardActivityHourInWeek,omitempty"`
		SmokeAtLeast100           *bool   `json:"smokeAtLeast100,omitempty"`
		SmokingAge                *uint   `json:"smokingAge,omitempty"`
		SmokingNow                *bool   `json:"smokingNow,omitempty"`
		LeaveSmokingAge           *uint   `json:"leaveSmokingAge,omitempty"`
		CountSmokingDaily         *string `json:"countSmokingDaily,omitempty"`
		CountGheliandaily         *string `json:"countGheliandaily,omitempty"`
		CountSmokingDailyPast     *string `json:"countSmokingDailyPast,omitempty"`
		CountGheliandailyPast     *string `json:"countGheliandailyPast,omitempty"`

		GhaedeAge                    *uint   `json:"ghaedeAge,omitempty"`
		HasChildren                  *bool   `json:"hasChildren,omitempty"`
		NumberOfChildren             *uint   `json:"numberOfChildren,omitempty"`
		AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth,omitempty"`
		MenopausalStatus             *string `json:"menopausalStatus,omitempty"`
		MenopauseAge                 *string `json:"menopauseAge,omitempty"`
		HRT                          *bool   `json:"hrt,omitempty"`
		HRTUseLength                 *uint   `json:"hrtUseLength,omitempty"`
		LastFiveYearsHRTUse          *bool   `json:"lastFiveYearsHrtUse,omitempty"`
		CurrentHRTUse                *bool   `json:"currentHrtUse,omitempty"`
		IntendedHRTUse               *uint   `json:"intendedHrtUse,omitempty"`
		HRTType                      *string `json:"hrtType,omitempty"`
		Oral                         *bool   `json:"oral,omitempty"`
		OralDuration                 *string `json:"oralDuration,omitempty"`
		OralTwoLastYears             *bool   `json:"oralTwoLastYears,omitempty"`
		MamoGraphy                   *bool   `json:"mamoGraphy,omitempty"`
		Falop                        *bool   `json:"falop,omitempty"`
		Andometrioz                  *bool   `json:"andometrioz,omitempty"`
		LeavePestan                  *bool   `json:"leavePestan,omitempty"`
		LeaveTokhmdan                *bool   `json:"leaveTokhmdan,omitempty"`
		LaDeColon                    *bool   `json:"laDeColon,omitempty"`
		LaDePol                      *bool   `json:"laDePol,omitempty"`
		AspLaMo                      *bool   `json:"aspLaMo,omitempty"`
		NsaiDLaMo                    *bool   `json:"nsaiDLaMo,omitempty"`
		LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool,omitempty"`

		Cancer     *bool   `json:"cancer,omitempty"`
		CancerType *string `json:"cancerType,omitempty"`
		CancerAge  *uint   `json:"cancerAge,omitempty"`

		ChildCancer     *bool   `json:"childCancer,omitempty"`
		ChildName       *string `json:"childName,omitempty"`
		ChildCancerType *string `json:"childCancerType,omitempty"`
		ChildCancerAge  *uint   `json:"childCancerAge,omitempty"`
		ChildLifeStatus *string `json:"childLifeStatus,omitempty"`

		MotherCancer     *bool   `json:"motherCancer,omitempty"`
		MotherName       *string `json:"motherName,omitempty"`
		MotherLifeStatus *string `json:"motherLifeStatus,omitempty"`
		MotherCancerType *string `json:"motherCancerType,omitempty"`
		MotherCancerAge  *uint   `json:"motherCancerAge,omitempty"`

		FatherCancer     *bool   `json:"fatherCancer,omitempty"`
		FatherName       *string `json:"fatherName,omitempty"`
		FatherLifeStatus *string `json:"fatherLifeStatus,omitempty"`
		FatherCancerType *string `json:"fatherCancerType,omitempty"`
		FatherCancerAge  *uint   `json:"fatherCancerAge,omitempty"`

		SiblingCancer     *bool   `json:"siblingCancer,omitempty"`
		SiblingName       *string `json:"siblingName,omitempty"`
		SiblingLifeStatus *string `json:"siblingLifeStatus,omitempty"`
		SiblingCancerType *string `json:"siblingCancerType,omitempty"`
		SiblingCancerAge  *uint   `json:"siblingCancerAge,omitempty"`

		AmeAmoCancer     *bool   `json:"ameAmoCancer,omitempty"`
		AmeAmoName       *string `json:"ameAmoName,omitempty"`
		AmeAmoLifeStatus *string `json:"ameAmoLifeStatus,omitempty"`
		AmeAmoCancerType *string `json:"ameAmoCancerType,omitempty"`
		AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge,omitempty"`

		KhaleDaeiCancer     *bool   `json:"khaleDaeiCancer,omitempty"`
		KhaleDaeiName       *string `json:"khaleDaeiName,omitempty"`
		KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus,omitempty"`
		KhaleDaeiCancerType *string `json:"khaleDaeiCancerType,omitempty"`
		KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge,omitempty"`

		OtherRelativeCancer     *bool   `json:"otherRelativeCancer,omitempty"`
		OtherRelativeName       *string `json:"otherRelativeName,omitempty"`
		OtherRelativeRelation   *string `json:"otherRelativeRelation,omitempty"`
		OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus,omitempty"`
		OtherRelativeCancerType *string `json:"otherRelativeCancerType,omitempty"`
		OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge,omitempty"`

		TestGen   *bool `json:"testGen,omitempty"`
		FmTestGen *bool `json:"fmTestGen,omitempty"`

		CallExpert   *bool   `json:"callExpert,omitempty"`
		BirthCountry *string `json:"birthCountry,omitempty"`
		Province     *string `json:"province,omitempty"`
		City         *string `json:"city,omitempty"`
		Country      *string `json:"country,omitempty"`

		InsuranceStatus           *string `json:"insuranceStatus,omitempty"`
		SupplementaryInsurances   *string `json:"supplementaryInsurances,omitempty"`
		Hypertension              *bool   `json:"hypertension,omitempty"`
		HypertensionTreatment     *bool   `json:"hypertensionTreatment,omitempty"`
		HeartDisease              *bool   `json:"heartDisease,omitempty"`
		HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment,omitempty"`
		Diabetes                  *bool   `json:"diabetes,omitempty"`
		DiabetesTreatment         *bool   `json:"diabetesTreatment,omitempty"`
		ChronicLungDisease        *bool   `json:"chronicLungDisease,omitempty"`
		ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType,omitempty"`
		LungCancerHistory         *bool   `json:"lungCancerHistory,omitempty"`
		OtherCancerHistory        *bool   `json:"otherCancerHistory,omitempty"`
		OtherCancerType           *string `json:"otherCancerType,omitempty"`
		LungCancerFamily          *bool   `json:"lungCancerFamily,omitempty"`
		LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation,omitempty"`
		OtherCancerFamily         *bool   `json:"otherCancerFamily,omitempty"`
		OtherCancerFamilyType     *string `json:"otherCancerFamilyType,omitempty"`
		OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation,omitempty"`
		OccupationalExposure      *string `json:"occupationalExposure,omitempty"`
		CurrentSmoking            *bool   `json:"currentSmoking,omitempty"`
		SmokingStartAgeCurrent    *uint   `json:"smokingStartAgeCurrent,omitempty"`
		SmokingTypesCurrent       *string `json:"smokingTypesCurrent,omitempty"`
		CigarettesPerDayCurrent   *uint   `json:"cigarettesPerDayCurrent,omitempty"`
		CigarPerDayCurrent        *uint   `json:"cigarPerDayCurrent,omitempty"`
		ECigPerDayCurrent         *uint   `json:"eCigPerDayCurrent,omitempty"`
		PipePerDayCurrent         *uint   `json:"pipePerDayCurrent,omitempty"`
		ChapoghPerDayCurrent      *uint   `json:"chapoghPerDayCurrent,omitempty"`
		SmokedOpiumPerDayCurrent  *uint   `json:"smokedOpiumPerDayCurrent,omitempty"`
		ChewedOpiumPerDayCurrent  *uint   `json:"chewedOpiumPerDayCurrent,omitempty"`
		HookahPerWeekCurrent      *uint   `json:"hookahPerWeekCurrent,omitempty"`
		PastSmoking               *string `json:"pastSmoking,omitempty"`
		SmokingStartAgePast       *uint   `json:"smokingStartAgePast,omitempty"`
		SmokingTypesPast          *string `json:"smokingTypesPast,omitempty"`
		CigarettesPerDayPast      *uint   `json:"cigarettesPerDayPast,omitempty"`
		CigarPerDayPast           *uint   `json:"cigarPerDayPast,omitempty"`
		ECigPerDayPast            *uint   `json:"eCigPerDayPast,omitempty"`
		PipePerDayPast            *uint   `json:"pipePerDayPast,omitempty"`
		ChapoghPerDayPast         *uint   `json:"chapoghPerDayPast,omitempty"`
		SmokedOpiumPerDayPast     *uint   `json:"smokedOpiumPerDayPast,omitempty"`
		ChewedOpiumPerDayPast     *uint   `json:"chewedOpiumPerDayPast,omitempty"`
		HookahPerWeekPast         *uint   `json:"hookahPerWeekPast,omitempty"`
		SecondhandSmoke           *bool   `json:"secondhandSmoke,omitempty"`
		SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation,omitempty"`
	}

	params := controller.Validate[UpdateFormParams](ctx)

	request := formdto.UpdateFormRequest{
		FormID:               params.FormID,
		Name:                 params.Name,
		BirthDay:             params.BirthDay,
		BirthMonth:           params.BirthMonth,
		BirthYear:            params.BirthYear,
		Address:              params.Address,
		PostalCode:           params.PostalCode,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,

		DrinksAlcohol:             params.DrinksAlcohol,
		CupsPerWeek:               params.CupsPerWeek,
		LastMonthSabzijatMeal:     params.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   params.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: params.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  params.MediumActivityHourInWeek,
		HardActivityMonthInYear:   params.HardActivityMonthInYear,
		HardActivityHourInWeek:    params.HardActivityHourInWeek,
		SmokeAtLeast100:           params.SmokeAtLeast100,
		SmokingAge:                params.SmokingAge,
		SmokingNow:                params.SmokingNow,
		LeaveSmokingAge:           params.LeaveSmokingAge,
		CountSmokingDaily:         params.CountSmokingDaily,
		CountGheliandaily:         params.CountGheliandaily,
		CountSmokingDailyPast:     params.CountSmokingDailyPast,
		CountGheliandailyPast:     params.CountGheliandailyPast,

		GhaedeAge:                    params.GhaedeAge,
		HasChildren:                  params.HasChildren,
		NumberOfChildren:             params.NumberOfChildren,
		AgeOfFirstBirth:              params.AgeOfFirstBirth,
		MenopausalStatus:             params.MenopausalStatus,
		MenopauseAge:                 params.MenopauseAge,
		HRT:                          params.HRT,
		HRTUseLength:                 params.HRTUseLength,
		LastFiveYearsHRTUse:          params.LastFiveYearsHRTUse,
		CurrentHRTUse:                params.CurrentHRTUse,
		IntendedHRTUse:               params.IntendedHRTUse,
		HRTType:                      params.HRTType,
		Oral:                         params.Oral,
		OralDuration:                 params.OralDuration,
		OralTwoLastYears:             params.OralTwoLastYears,
		MamoGraphy:                   params.MamoGraphy,
		Falop:                        params.Falop,
		Andometrioz:                  params.Andometrioz,
		LeavePestan:                  params.LeavePestan,
		LeaveTokhmdan:                params.LeaveTokhmdan,
		LaDeColon:                    params.LaDeColon,
		LaDePol:                      params.LaDePol,
		AspLaMo:                      params.AspLaMo,
		NsaiDLaMo:                    params.NsaiDLaMo,
		LastFiveYearBloodTestInStool: params.LastFiveYearBloodTestInStool,

		Cancer:     params.Cancer,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,

		ChildCancer:     params.ChildCancer,
		ChildName:       params.ChildName,
		ChildCancerType: params.ChildCancerType,
		ChildCancerAge:  params.ChildCancerAge,
		ChildLifeStatus: params.ChildLifeStatus,

		MotherCancer:     params.MotherCancer,
		MotherName:       params.MotherName,
		MotherLifeStatus: params.MotherLifeStatus,
		MotherCancerType: params.MotherCancerType,
		MotherCancerAge:  params.MotherCancerAge,

		FatherCancer:     params.FatherCancer,
		FatherName:       params.FatherName,
		FatherLifeStatus: params.FatherLifeStatus,
		FatherCancerType: params.FatherCancerType,
		FatherCancerAge:  params.FatherCancerAge,

		SiblingCancer:     params.SiblingCancer,
		SiblingName:       params.SiblingName,
		SiblingLifeStatus: params.SiblingLifeStatus,
		SiblingCancerType: params.SiblingCancerType,
		SiblingCancerAge:  params.SiblingCancerAge,

		AmeAmoCancer:     params.AmeAmoCancer,
		AmeAmoName:       params.AmeAmoName,
		AmeAmoLifeStatus: params.AmeAmoLifeStatus,
		AmeAmoCancerType: params.AmeAmoCancerType,
		AmeAmoCancerAge:  params.AmeAmoCancerAge,

		KhaleDaeiCancer:     params.KhaleDaeiCancer,
		KhaleDaeiName:       params.KhaleDaeiName,
		KhaleDaeiLifeStatus: params.KhaleDaeiLifeStatus,
		KhaleDaeiCancerType: params.KhaleDaeiCancerType,
		KhaleDaeiCancerAge:  params.KhaleDaeiCancerAge,

		OtherRelativeCancer:     params.OtherRelativeCancer,
		OtherRelativeName:       params.OtherRelativeName,
		OtherRelativeRelation:   params.OtherRelativeRelation,
		OtherRelativeLifeStatus: params.OtherRelativeLifeStatus,
		OtherRelativeCancerType: params.OtherRelativeCancerType,
		OtherRelativeCancerAge:  params.OtherRelativeCancerAge,

		TestGen:   params.TestGen,
		FmTestGen: params.FmTestGen,

		CallExpert:   params.CallExpert,
		BirthCountry: params.BirthCountry,
		Province:     params.Province,
		City:         params.City,
		Country:      params.Country,

		InsuranceStatus:           params.InsuranceStatus,
		SupplementaryInsurances:   params.SupplementaryInsurances,
		Hypertension:              params.Hypertension,
		HypertensionTreatment:     params.HypertensionTreatment,
		HeartDisease:              params.HeartDisease,
		HeartDiseaseTreatment:     params.HeartDiseaseTreatment,
		Diabetes:                  params.Diabetes,
		DiabetesTreatment:         params.DiabetesTreatment,
		ChronicLungDisease:        params.ChronicLungDisease,
		ChronicLungDiseaseType:    params.ChronicLungDiseaseType,
		LungCancerHistory:         params.LungCancerHistory,
		OtherCancerHistory:        params.OtherCancerHistory,
		OtherCancerType:           params.OtherCancerType,
		LungCancerFamily:          params.LungCancerFamily,
		LungCancerFamilyRelation:  params.LungCancerFamilyRelation,
		OtherCancerFamily:         params.OtherCancerFamily,
		OtherCancerFamilyType:     params.OtherCancerFamilyType,
		OtherCancerFamilyRelation: params.OtherCancerFamilyRelation,
		OccupationalExposure:      params.OccupationalExposure,
		CurrentSmoking:            params.CurrentSmoking,
		SmokingStartAgeCurrent:    params.SmokingStartAgeCurrent,
		SmokingTypesCurrent:       params.SmokingTypesCurrent,
		CigarettesPerDayCurrent:   params.CigarettesPerDayCurrent,
		CigarPerDayCurrent:        params.CigarPerDayCurrent,
		ECigPerDayCurrent:         params.ECigPerDayCurrent,
		PipePerDayCurrent:         params.PipePerDayCurrent,
		ChapoghPerDayCurrent:      params.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:  params.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:  params.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:      params.HookahPerWeekCurrent,
		PastSmoking:               params.PastSmoking,
		SmokingStartAgePast:       params.SmokingStartAgePast,
		SmokingTypesPast:          params.SmokingTypesPast,
		CigarettesPerDayPast:      params.CigarettesPerDayPast,
		CigarPerDayPast:           params.CigarPerDayPast,
		ECigPerDayPast:            params.ECigPerDayPast,
		PipePerDayPast:            params.PipePerDayPast,
		ChapoghPerDayPast:         params.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:     params.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:     params.ChewedOpiumPerDayPast,
		HookahPerWeekPast:         params.HookahPerWeekPast,
		SecondhandSmoke:           params.SecondhandSmoke,
		SecondhandSmokeLocation:   params.SecondhandSmokeLocation,
	}

	err := formController.formService.UpdateForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *CustomerFormController) DeleteForm(ctx *gin.Context) {
	type DeleteFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[DeleteFormParams](ctx)

	err := formController.formService.DeleteForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, nil)
}
