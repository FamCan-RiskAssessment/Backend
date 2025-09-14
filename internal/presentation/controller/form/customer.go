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
		BirthDay             string  `json:"birth_day" validate:"required"`
		BirthMonth           string  `json:"birth_month" validate:"required"`
		BirthYear            string  `json:"birth_year" validate:"required"`
		Address              string  `json:"address" validate:"required"`
		PostalCode           string  `json:"postal_code" validate:"required"`
		SocialSecurityNumber string  `json:"social_security_number" validate:"required"`
		Gender               string  `json:"gender" validate:"required"`
		IsAtba               bool    `json:"is_atba"`
		Height               float64 `json:"height" validate:"required"`
		Weight               float64 `json:"weight" validate:"required"`

		DrinksAlcohol             *bool   `json:"drinks_alcohol,omitempty"`
		CupsPerWeek               *string `json:"cups_per_week,omitempty"`
		LastMonthSabzijatMeal     string  `json:"last_month_sabzijat_meal" validate:"required"`
		LastMonthSabzijatWeight   string  `json:"last_month_sabzijat_weight" validate:"required"`
		MediumActivityMonthInYear uint    `json:"medium_activity_month_in_year" validate:"required"`
		MediumActivityHourInWeek  string  `json:"medium_activity_hour_in_week" validate:"required"`
		HardActivityMonthInYear   uint    `json:"hard_activity_month_in_year" validate:"required"`
		HardActivityHourInWeek    string  `json:"hard_activity_hour_in_week" validate:"required"`
		SmokeAtLeast100           *bool   `json:"smoke_at_least_100,omitempty"`
		SmokingAge                uint    `json:"smoking_age" validate:"required"`
		SmokingNow                bool    `json:"smoking_now"`
		LeaveSmokingAge           *uint   `json:"leave_smoking_age,omitempty"`
		CountSmokingDaily         *uint   `json:"count_smoking_daily,omitempty"`
		CountGheliandaily         *uint   `json:"count_gheliandaily,omitempty"`
		CountSmokingDailyPast     *uint   `json:"count_smoking_daily_past,omitempty"`
		CountGheliandailyPast     *uint   `json:"count_gheliandaily_past,omitempty"`

		GhaedeAge                    uint    `json:"ghaede_age" validate:"required"`
		HasChildren                  bool    `json:"has_children"`
		NumberOfChildren             *uint   `json:"number_of_children,omitempty"`
		AgeOfFirstBirth              *uint   `json:"age_of_first_birth,omitempty"`
		MenopausalStatus             string  `json:"menopausal_status" validate:"required"`
		MenopauseAge                 string  `json:"menopause_age" validate:"required"`
		HRT                          *bool   `json:"hrt,omitempty"`
		HRTUseLength                 *uint   `json:"hrt_use_length,omitempty"`
		LastFiveYearsHRTUse          bool    `json:"last_five_years_hrt_use"`
		CurrentHRTUse                *bool   `json:"current_hrt_use,omitempty"`
		IntendedHRTUse               *uint   `json:"intended_hrt_use,omitempty"`
		HRTType                      *string `json:"hrt_type,omitempty"`
		Oral                         *bool   `json:"oral,omitempty"`
		OralDuration                 *string `json:"oral_duration,omitempty"`
		OralTwoLastYears             *bool   `json:"oral_two_last_years,omitempty"`
		MamoGraphy                   *bool   `json:"mamo_graphy,omitempty"`
		Falop                        *bool   `json:"falop,omitempty"`
		Andometrioz                  *bool   `json:"andometrioz,omitempty"`
		LeavePestan                  bool    `json:"leave_pestan"`
		LeaveTokhmdan                bool    `json:"leave_tokhmdan"`
		LaDeColon                    *bool   `json:"la_de_colon,omitempty"`
		LaDePol                      *bool   `json:"la_de_pol,omitempty"`
		AspLaMo                      *bool   `json:"asp_la_mo,omitempty"`
		NsaiDLaMo                    *bool   `json:"nsai_d_la_mo,omitempty"`
		LastFiveYearBloodTestInStool *bool   `json:"last_five_year_blood_test_in_stool,omitempty"`

		Cancer     bool    `json:"cancer"`
		CancerType *string `json:"cancer_type,omitempty"`
		CancerAge  *uint   `json:"cancer_age,omitempty"`

		ChildCancer     bool    `json:"child_cancer"`
		ChildName       *string `json:"child_name,omitempty"`
		ChildCancerType *string `json:"child_cancer_type,omitempty"`
		ChildCancerAge  *string `json:"child_cancer_age,omitempty"`
		ChildLifeStatus *string `json:"child_life_status,omitempty"`

		MotherCancer     bool    `json:"mother_cancer"`
		MotherName       *string `json:"mother_name,omitempty"`
		MotherLifeStatus *string `json:"mother_life_status,omitempty"`
		MotherCancerType *string `json:"mother_cancer_type,omitempty"`
		MotherCancerAge  *string `json:"mother_cancer_age,omitempty"`

		FatherCancer     bool    `json:"father_cancer"`
		FatherName       *string `json:"father_name,omitempty"`
		FatherLifeStatus *string `json:"father_life_status,omitempty"`
		FatherCancerType *string `json:"father_cancer_type,omitempty"`
		FatherCancerAge  *string `json:"father_cancer_age,omitempty"`

		SiblingCancer     bool    `json:"sibling_cancer"`
		SiblingName       *string `json:"sibling_name,omitempty"`
		SiblingLifeStatus *string `json:"sibling_life_status,omitempty"`
		SiblingCancerType *string `json:"sibling_cancer_type,omitempty"`
		SiblingCancerAge  *string `json:"sibling_cancer_age,omitempty"`

		AmeAmoCancer     bool    `json:"ame_amo_cancer"`
		AmeAmoName       *string `json:"ame_amo_name,omitempty"`
		AmeAmoLifeStatus *string `json:"ame_amo_life_status,omitempty"`
		AmeAmoCancerType *string `json:"ame_amo_cancer_type,omitempty"`
		AmeAmoCancerAge  *string `json:"ame_amo_cancer_age,omitempty"`

		KhaleDaeiCancer     bool    `json:"khale_daei_cancer"`
		KhaleDaeiName       *string `json:"khale_daei_name,omitempty"`
		KhaleDaeiLifeStatus *string `json:"khale_daei_life_status,omitempty"`
		KhaleDaeiCancerType *string `json:"khale_daei_cancer_type,omitempty"`
		KhaleDaeiCancerAge  *string `json:"khale_daei_cancer_age,omitempty"`

		OtherRelativeCancer     *bool   `json:"other_relative_cancer,omitempty"`
		OtherRelativeName       *string `json:"other_relative_name,omitempty"`
		OtherRelativeRelation   *string `json:"other_relative_relation,omitempty"`
		OtherRelativeLifeStatus *string `json:"other_relative_life_status,omitempty"`
		OtherRelativeCancerType *string `json:"other_relative_cancer_type,omitempty"`
		OtherRelativeCancerAge  *string `json:"other_relative_cancer_age,omitempty"`

		TestGen   *bool `json:"test_gen,omitempty"`
		FmTestGen *bool `json:"fm_test_gen,omitempty"`

		CallExpert   bool    `json:"call_expert"`
		BirthCountry *string `json:"birth_country,omitempty"`
		Province     *string `json:"province,omitempty"`
		City         *string `json:"city,omitempty"`
		Country      *string `json:"country,omitempty"`

		InsuranceStatus           string  `json:"insurance_status" validate:"required"`
		SupplementaryInsurances   *string `json:"supplementary_insurances,omitempty"`
		Hypertension              bool    `json:"hypertension"`
		HypertensionTreatment     *bool   `json:"hypertension_treatment,omitempty"`
		HeartDisease              bool    `json:"heart_disease"`
		HeartDiseaseTreatment     *bool   `json:"heart_disease_treatment,omitempty"`
		Diabetes                  bool    `json:"diabetes"`
		DiabetesTreatment         *bool   `json:"diabetes_treatment,omitempty"`
		ChronicLungDisease        *bool   `json:"chronic_lung_disease,omitempty"`
		ChronicLungDiseaseType    *string `json:"chronic_lung_disease_type,omitempty"`
		LungCancerHistory         bool    `json:"lung_cancer_history"`
		OtherCancerHistory        bool    `json:"other_cancer_history"`
		OtherCancerType           *string `json:"other_cancer_type,omitempty"`
		LungCancerFamily          *bool   `json:"lung_cancer_family,omitempty"`
		LungCancerFamilyRelation  *string `json:"lung_cancer_family_relation,omitempty"`
		OtherCancerFamily         *bool   `json:"other_cancer_family,omitempty"`
		OtherCancerFamilyType     *string `json:"other_cancer_family_type,omitempty"`
		OtherCancerFamilyRelation *string `json:"other_cancer_family_relation,omitempty"`
		OccupationalExposure      *string `json:"occupational_exposure,omitempty"`
		CurrentSmoking            bool    `json:"current_smoking"`
		SmokingStartAgeCurrent    *uint   `json:"smoking_start_age_current,omitempty"`
		SmokingTypesCurrent       *string `json:"smoking_types_current,omitempty"`
		CigarettesPerDayCurrent   *uint   `json:"cigarettes_per_day_current,omitempty"`
		CigarPerDayCurrent        *uint   `json:"cigar_per_day_current,omitempty"`
		ECigPerDayCurrent         *uint   `json:"e_cig_per_day_current,omitempty"`
		PipePerDayCurrent         *uint   `json:"pipe_per_day_current,omitempty"`
		ChapoghPerDayCurrent      *uint   `json:"chapogh_per_day_current,omitempty"`
		SmokedOpiumPerDayCurrent  *uint   `json:"smoked_opium_per_day_current,omitempty"`
		ChewedOpiumPerDayCurrent  *uint   `json:"chewed_opium_per_day_current,omitempty"`
		HookahPerWeekCurrent      *uint   `json:"hookah_per_week_current,omitempty"`
		PastSmoking               *string `json:"past_smoking,omitempty"`
		SmokingStartAgePast       *uint   `json:"smoking_start_age_past,omitempty"`
		SmokingTypesPast          *string `json:"smoking_types_past,omitempty"`
		CigarettesPerDayPast      *uint   `json:"cigarettes_per_day_past,omitempty"`
		CigarPerDayPast           *uint   `json:"cigar_per_day_past,omitempty"`
		ECigPerDayPast            *uint   `json:"e_cig_per_day_past,omitempty"`
		PipePerDayPast            *uint   `json:"pipe_per_day_past,omitempty"`
		ChapoghPerDayPast         *uint   `json:"chapogh_per_day_past,omitempty"`
		SmokedOpiumPerDayPast     *uint   `json:"smoked_opium_per_day_past,omitempty"`
		ChewedOpiumPerDayPast     *uint   `json:"chewed_opium_per_day_past,omitempty"`
		HookahPerWeekPast         *uint   `json:"hookah_per_week_past,omitempty"`
		SecondhandSmoke           bool    `json:"secondhand_smoke"`
		SecondhandSmokeLocation   *string `json:"secondhand_smoke_location,omitempty"`
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

	response, err := formController.formService.CreateForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, response)
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
		BirthDay             *string  `json:"birth_day,omitempty"`
		BirthMonth           *string  `json:"birth_month,omitempty"`
		BirthYear            *string  `json:"birth_year,omitempty"`
		Address              *string  `json:"address,omitempty"`
		PostalCode           *string  `json:"postal_code,omitempty"`
		SocialSecurityNumber *string  `json:"social_security_number,omitempty"`
		Gender               *string  `json:"gender,omitempty"`
		IsAtba               *bool    `json:"is_atba,omitempty"`
		Height               *float64 `json:"height,omitempty"`
		Weight               *float64 `json:"weight,omitempty"`

		DrinksAlcohol             *bool   `json:"drinks_alcohol,omitempty"`
		CupsPerWeek               *string `json:"cups_per_week,omitempty"`
		LastMonthSabzijatMeal     *string `json:"last_month_sabzijat_meal,omitempty"`
		LastMonthSabzijatWeight   *string `json:"last_month_sabzijat_weight,omitempty"`
		MediumActivityMonthInYear *uint   `json:"medium_activity_month_in_year,omitempty"`
		MediumActivityHourInWeek  *string `json:"medium_activity_hour_in_week,omitempty"`
		HardActivityMonthInYear   *uint   `json:"hard_activity_month_in_year,omitempty"`
		HardActivityHourInWeek    *string `json:"hard_activity_hour_in_week,omitempty"`
		SmokeAtLeast100           *bool   `json:"smoke_at_least_100,omitempty"`
		SmokingAge                *uint   `json:"smoking_age,omitempty"`
		SmokingNow                *bool   `json:"smoking_now,omitempty"`
		LeaveSmokingAge           *uint   `json:"leave_smoking_age,omitempty"`
		CountSmokingDaily         *uint   `json:"count_smoking_daily,omitempty"`
		CountGheliandaily         *uint   `json:"count_gheliandaily,omitempty"`
		CountSmokingDailyPast     *uint   `json:"count_smoking_daily_past,omitempty"`
		CountGheliandailyPast     *uint   `json:"count_gheliandaily_past,omitempty"`

		GhaedeAge                    *uint   `json:"ghaede_age,omitempty"`
		HasChildren                  *bool   `json:"has_children,omitempty"`
		NumberOfChildren             *uint   `json:"number_of_children,omitempty"`
		AgeOfFirstBirth              *uint   `json:"age_of_first_birth,omitempty"`
		MenopausalStatus             *string `json:"menopausal_status,omitempty"`
		MenopauseAge                 *string `json:"menopause_age,omitempty"`
		HRT                          *bool   `json:"hrt,omitempty"`
		HRTUseLength                 *uint   `json:"hrt_use_length,omitempty"`
		LastFiveYearsHRTUse          *bool   `json:"last_five_years_hrt_use,omitempty"`
		CurrentHRTUse                *bool   `json:"current_hrt_use,omitempty"`
		IntendedHRTUse               *uint   `json:"intended_hrt_use,omitempty"`
		HRTType                      *string `json:"hrt_type,omitempty"`
		Oral                         *bool   `json:"oral,omitempty"`
		OralDuration                 *string `json:"oral_duration,omitempty"`
		OralTwoLastYears             *bool   `json:"oral_two_last_years,omitempty"`
		MamoGraphy                   *bool   `json:"mamo_graphy,omitempty"`
		Falop                        *bool   `json:"falop,omitempty"`
		Andometrioz                  *bool   `json:"andometrioz,omitempty"`
		LeavePestan                  *bool   `json:"leave_pestan,omitempty"`
		LeaveTokhmdan                *bool   `json:"leave_tokhmdan,omitempty"`
		LaDeColon                    *bool   `json:"la_de_colon,omitempty"`
		LaDePol                      *bool   `json:"la_de_pol,omitempty"`
		AspLaMo                      *bool   `json:"asp_la_mo,omitempty"`
		NsaiDLaMo                    *bool   `json:"nsai_d_la_mo,omitempty"`
		LastFiveYearBloodTestInStool *bool   `json:"last_five_year_blood_test_in_stool,omitempty"`

		Cancer     *bool   `json:"cancer,omitempty"`
		CancerType *string `json:"cancer_type,omitempty"`
		CancerAge  *uint   `json:"cancer_age,omitempty"`

		ChildCancer     *bool   `json:"child_cancer,omitempty"`
		ChildName       *string `json:"child_name,omitempty"`
		ChildCancerType *string `json:"child_cancer_type,omitempty"`
		ChildCancerAge  *string `json:"child_cancer_age,omitempty"`
		ChildLifeStatus *string `json:"child_life_status,omitempty"`

		MotherCancer     *bool   `json:"mother_cancer,omitempty"`
		MotherName       *string `json:"mother_name,omitempty"`
		MotherLifeStatus *string `json:"mother_life_status,omitempty"`
		MotherCancerType *string `json:"mother_cancer_type,omitempty"`
		MotherCancerAge  *string `json:"mother_cancer_age,omitempty"`

		FatherCancer     *bool   `json:"father_cancer,omitempty"`
		FatherName       *string `json:"father_name,omitempty"`
		FatherLifeStatus *string `json:"father_life_status,omitempty"`
		FatherCancerType *string `json:"father_cancer_type,omitempty"`
		FatherCancerAge  *string `json:"father_cancer_age,omitempty"`

		SiblingCancer     *bool   `json:"sibling_cancer,omitempty"`
		SiblingName       *string `json:"sibling_name,omitempty"`
		SiblingLifeStatus *string `json:"sibling_life_status,omitempty"`
		SiblingCancerType *string `json:"sibling_cancer_type,omitempty"`
		SiblingCancerAge  *string `json:"sibling_cancer_age,omitempty"`

		AmeAmoCancer     *bool   `json:"ame_amo_cancer,omitempty"`
		AmeAmoName       *string `json:"ame_amo_name,omitempty"`
		AmeAmoLifeStatus *string `json:"ame_amo_life_status,omitempty"`
		AmeAmoCancerType *string `json:"ame_amo_cancer_type,omitempty"`
		AmeAmoCancerAge  *string `json:"ame_amo_cancer_age,omitempty"`

		KhaleDaeiCancer     *bool   `json:"khale_daei_cancer,omitempty"`
		KhaleDaeiName       *string `json:"khale_daei_name,omitempty"`
		KhaleDaeiLifeStatus *string `json:"khale_daei_life_status,omitempty"`
		KhaleDaeiCancerType *string `json:"khale_daei_cancer_type,omitempty"`
		KhaleDaeiCancerAge  *string `json:"khale_daei_cancer_age,omitempty"`

		OtherRelativeCancer     *bool   `json:"other_relative_cancer,omitempty"`
		OtherRelativeName       *string `json:"other_relative_name,omitempty"`
		OtherRelativeRelation   *string `json:"other_relative_relation,omitempty"`
		OtherRelativeLifeStatus *string `json:"other_relative_life_status,omitempty"`
		OtherRelativeCancerType *string `json:"other_relative_cancer_type,omitempty"`
		OtherRelativeCancerAge  *string `json:"other_relative_cancer_age,omitempty"`

		TestGen   *bool `json:"test_gen,omitempty"`
		FmTestGen *bool `json:"fm_test_gen,omitempty"`

		CallExpert   *bool   `json:"call_expert,omitempty"`
		BirthCountry *string `json:"birth_country,omitempty"`
		Province     *string `json:"province,omitempty"`
		City         *string `json:"city,omitempty"`
		Country      *string `json:"country,omitempty"`

		InsuranceStatus           *string `json:"insurance_status,omitempty"`
		SupplementaryInsurances   *string `json:"supplementary_insurances,omitempty"`
		Hypertension              *bool   `json:"hypertension,omitempty"`
		HypertensionTreatment     *bool   `json:"hypertension_treatment,omitempty"`
		HeartDisease              *bool   `json:"heart_disease,omitempty"`
		HeartDiseaseTreatment     *bool   `json:"heart_disease_treatment,omitempty"`
		Diabetes                  *bool   `json:"diabetes,omitempty"`
		DiabetesTreatment         *bool   `json:"diabetes_treatment,omitempty"`
		ChronicLungDisease        *bool   `json:"chronic_lung_disease,omitempty"`
		ChronicLungDiseaseType    *string `json:"chronic_lung_disease_type,omitempty"`
		LungCancerHistory         *bool   `json:"lung_cancer_history,omitempty"`
		OtherCancerHistory        *bool   `json:"other_cancer_history,omitempty"`
		OtherCancerType           *string `json:"other_cancer_type,omitempty"`
		LungCancerFamily          *bool   `json:"lung_cancer_family,omitempty"`
		LungCancerFamilyRelation  *string `json:"lung_cancer_family_relation,omitempty"`
		OtherCancerFamily         *bool   `json:"other_cancer_family,omitempty"`
		OtherCancerFamilyType     *string `json:"other_cancer_family_type,omitempty"`
		OtherCancerFamilyRelation *string `json:"other_cancer_family_relation,omitempty"`
		OccupationalExposure      *string `json:"occupational_exposure,omitempty"`
		CurrentSmoking            *bool   `json:"current_smoking,omitempty"`
		SmokingStartAgeCurrent    *uint   `json:"smoking_start_age_current,omitempty"`
		SmokingTypesCurrent       *string `json:"smoking_types_current,omitempty"`
		CigarettesPerDayCurrent   *uint   `json:"cigarettes_per_day_current,omitempty"`
		CigarPerDayCurrent        *uint   `json:"cigar_per_day_current,omitempty"`
		ECigPerDayCurrent         *uint   `json:"e_cig_per_day_current,omitempty"`
		PipePerDayCurrent         *uint   `json:"pipe_per_day_current,omitempty"`
		ChapoghPerDayCurrent      *uint   `json:"chapogh_per_day_current,omitempty"`
		SmokedOpiumPerDayCurrent  *uint   `json:"smoked_opium_per_day_current,omitempty"`
		ChewedOpiumPerDayCurrent  *uint   `json:"chewed_opium_per_day_current,omitempty"`
		HookahPerWeekCurrent      *uint   `json:"hookah_per_week_current,omitempty"`
		PastSmoking               *string `json:"past_smoking,omitempty"`
		SmokingStartAgePast       *uint   `json:"smoking_start_age_past,omitempty"`
		SmokingTypesPast          *string `json:"smoking_types_past,omitempty"`
		CigarettesPerDayPast      *uint   `json:"cigarettes_per_day_past,omitempty"`
		CigarPerDayPast           *uint   `json:"cigar_per_day_past,omitempty"`
		ECigPerDayPast            *uint   `json:"e_cig_per_day_past,omitempty"`
		PipePerDayPast            *uint   `json:"pipe_per_day_past,omitempty"`
		ChapoghPerDayPast         *uint   `json:"chapogh_per_day_past,omitempty"`
		SmokedOpiumPerDayPast     *uint   `json:"smoked_opium_per_day_past,omitempty"`
		ChewedOpiumPerDayPast     *uint   `json:"chewed_opium_per_day_past,omitempty"`
		HookahPerWeekPast         *uint   `json:"hookah_per_week_past,omitempty"`
		SecondhandSmoke           *bool   `json:"secondhand_smoke,omitempty"`
		SecondhandSmokeLocation   *string `json:"secondhand_smoke_location,omitempty"`
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

	response, err := formController.formService.UpdateForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, response)
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
