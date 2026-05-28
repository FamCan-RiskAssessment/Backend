package form

import (
	"mime/multipart"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
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
		BirthDate            formdto.BirthDate `json:"birthDate" validate:"required" time_format:"2006-01-02"`
		SocialSecurityNumber string            `json:"socialSecurityNumber" validate:"required"`
		Gender               uint              `json:"gender" validate:"required"`
		FormType             *enum.FormType    `json:"formType"`
		IsAtba               bool              `json:"isAtba"`
		Height               float64           `json:"height" validate:"required"`
		Weight               float64           `json:"weight" validate:"required"`
	}

	params := controller.Validate[CreateFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.CreateBasicFormRequest{
		UserID:               userID.(uint),
		BirthDate:            params.BirthDate,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,
		FormType:             params.FormType,
	}

	form, err := formController.formService.CreateBasicInfoForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, formdto.CreateFormResponse{Form: form})
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

func (formController *CustomerFormController) UpdateBasicInfo(ctx *gin.Context) {
	type UpdateBasicInfoParams struct {
		FormID               uint               `uri:"formID" validate:"required"`
		BirthDate            *formdto.BirthDate `json:"birthDate" validate:"required" time_format:"2006-01-02"`
		SocialSecurityNumber *string            `json:"socialSecurityNumber" validate:"required"`
		Gender               *uint              `json:"gender" validate:"required"`
		IsAtba               *bool              `json:"isAtba"`
		Height               *float64           `json:"height" validate:"required"`
		Weight               *float64           `json:"weight" validate:"required"`
	}

	params := controller.Validate[UpdateBasicInfoParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.UpdateBasicFormRequest{
		FormID:               params.FormID,
		BirthDate:            params.BirthDate,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,
		UserID:               userID.(uint),
	}

	err := formController.formService.UpdateBasicInfo(request)
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

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.DeleteFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	err := formController.formService.DeleteForm(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *CustomerFormController) UpsertGeneralHealth(ctx *gin.Context) {
	type UpsertGeneralHealthParams struct {
		FormID uint `uri:"formID" validate:"required"`

		DrinksAlcohol             *enum.Answer `json:"drinksAlcohol,omitempty"`
		CupsPerWeek               *string      `json:"cupsPerWeek,omitempty"`
		LastMonthSabzijatMeal     string       `json:"lastMonthSabzijatMeal" validate:"required"`
		LastMonthSabzijatWeight   string       `json:"lastMonthSabzijatWeight" validate:"required"`
		MediumActivityMonthInYear *uint        `json:"mediumActivityMonthInYear" validate:"required"`
		MediumActivityHourInWeek  string       `json:"mediumActivityHourInWeek" validate:"required"`
		HardActivityMonthInYear   *uint        `json:"hardActivityMonthInYear" validate:"required"`
		HardActivityHourInWeek    string       `json:"hardActivityHourInWeek" validate:"required"`
		SmokeAtLeast100           *enum.Answer `json:"smokeAtLeast100,omitempty"`
		SmokingAge                *uint        `json:"smokingAge"`
		SmokingNow                *enum.Answer `json:"smokingNow"`
		YearSmoke                 *uint        `json:"yearSmoke"`
		LeaveSmokingAge           *uint        `json:"leaveSmokingAge"`
		CountSmokingDaily         *string      `json:"countSmokingDaily,omitempty"`
		CountGheliandaily         *string      `json:"countGheliandaily,omitempty"`
		CountSmokingDailyPast     *string      `json:"countSmokingDailyPast,omitempty"`
		CountGheliandailyPast     *string      `json:"countGheliandailyPast,omitempty"`
		AttentionCorrect          *bool        `json:"attentionCorrect,omitempty"`
	}

	params := controller.Validate[UpsertGeneralHealthParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertGeneralHealthRequest{
		UserID:                    userID.(uint),
		FormID:                    params.FormID,
		DrinksAlcohol:             params.DrinksAlcohol,
		CupsPerWeek:               params.CupsPerWeek,
		LastMonthSabzijatMeal:     params.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   params.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: *params.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  params.MediumActivityHourInWeek,
		HardActivityMonthInYear:   *params.HardActivityMonthInYear,
		HardActivityHourInWeek:    params.HardActivityHourInWeek,
		SmokeAtLeast100:           params.SmokeAtLeast100,
		SmokingAge:                params.SmokingAge,
		SmokingNow:                params.SmokingNow,
		YearSmoke:                 params.YearSmoke,
		LeaveSmokingAge:           params.LeaveSmokingAge,
		CountSmokingDaily:         params.CountSmokingDaily,
		CountGheliandaily:         params.CountGheliandaily,
		CountSmokingDailyPast:     params.CountSmokingDailyPast,
		CountGheliandailyPast:     params.CountGheliandailyPast,
		AttentionCorrect:          params.AttentionCorrect,
	}

	if err := formController.formService.UpsertGeneralHealth(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertGeneralHealthResponse{})
}
func (formController *CustomerFormController) UpsertMamography(ctx *gin.Context) {
	type UpsertMamographyParams struct {
		FormID uint `uri:"formID" validate:"required"`

		GhaedeAge                    uint                    `form:"ghaedeAge"`
		HasChildren                  *enum.Answer            `form:"hasChildren"`
		NumberOfChildren             *uint                   `form:"numberOfChildren"`
		SonCount                     *uint                   `form:"sonCount"`
		DaughterCount                *uint                   `form:"daughterCount"`
		AgeOfFirstBirth              *uint                   `form:"ageOfFirstBirth"`
		MenopausalStatus             uint                    `form:"menopausalStatus"`
		MenopauseAge                 *string                 `form:"menopauseAge,omitempty"`
		HRT                          *enum.Answer            `form:"hrt,omitempty"`
		HRTUseLength                 *uint                   `form:"hrtUseLength"`
		LastFiveYearsHRTUse          *enum.Answer            `form:"lastFiveYearsHrtUse"`
		CurrentHRTUse                *enum.Answer            `form:"currentHrtUse,omitempty"`
		IntendedHRTUse               *uint                   `form:"intendedHrtUse"`
		HRTType                      *string                 `form:"hrtType,omitempty"`
		Oral                         *enum.Answer            `form:"oral,omitempty"`
		OralDuration                 *string                 `form:"oralDuration,omitempty"`
		OralTwoLastYears             *enum.Answer            `form:"oralTwoLastYears,omitempty"`
		MamoGraphy                   *enum.Answer            `form:"mamoGraphy,omitempty"`
		MamoGraphyPictures           []*multipart.FileHeader `form:"mamoGraphyPictures"`
		BreastDensity                *string                 `form:"breastDensity,omitempty"`
		Falop                        *enum.Answer            `form:"falop,omitempty"`
		Andometrioz                  *enum.Answer            `form:"andometrioz,omitempty"`
		LeavePestan                  bool                    `form:"leavePestan"`
		LeaveTokhmdan                bool                    `form:"leaveTokhmdan"`
		LaDeColon                    *enum.Answer            `form:"laDeColon,omitempty"`
		LaDePol                      *enum.Answer            `form:"laDePol,omitempty"`
		AspLaMo                      *enum.Answer            `form:"aspLaMo,omitempty"`
		NsaiDLaMo                    *enum.Answer            `form:"nsaiDLaMo,omitempty"`
		LastFiveYearBloodTestInStool *enum.Answer            `form:"lastFiveYearBloodTestInStool,omitempty"`
		NumberOfBreastBiopsies       *uint                   `form:"numberOfBreastBiopsies"`
		HyperplasiaInBiopsy          *uint                   `form:"hyperplasiaInBiopsy"`
		AttentionCorrect             *bool                   `form:"attentionCorrect,omitempty"`
	}

	params := controller.Validate[UpsertMamographyParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertMamographyRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		GhaedeAge:                    params.GhaedeAge,
		HasChildren:                  params.HasChildren,
		NumberOfChildren:             params.NumberOfChildren,
		SonCount:                     params.SonCount,
		DaughterCount:                params.DaughterCount,
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
		MamoGraphyPictures:           params.MamoGraphyPictures,
		BreastDensity:                params.BreastDensity,
		Falop:                        params.Falop,
		Andometrioz:                  params.Andometrioz,
		LeavePestan:                  params.LeavePestan,
		LeaveTokhmdan:                params.LeaveTokhmdan,
		LaDeColon:                    params.LaDeColon,
		LaDePol:                      params.LaDePol,
		AspLaMo:                      params.AspLaMo,
		NsaiDLaMo:                    params.NsaiDLaMo,
		LastFiveYearBloodTestInStool: params.LastFiveYearBloodTestInStool,
		NumberOfBreastBiopsies:       params.NumberOfBreastBiopsies,
		HyperplasiaInBiopsy:          params.HyperplasiaInBiopsy,
		AttentionCorrect:             params.AttentionCorrect,
	}

	if err := formController.formService.UpsertMamography(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertMamographyResponse{})
}

func (formController *CustomerFormController) UpsertContact(ctx *gin.Context) {
	type UpsertContactParams struct {
		FormID uint `uri:"formID" validate:"required"`

		Name                  string                  `form:"name" validate:"required"`
		TestGen               *enum.Answer            `form:"testGen,omitempty"`
		TestGenPictures       []*multipart.FileHeader `form:"testGenPictures,omitempty"`
		FmTestGen             *enum.Answer            `form:"fmTestGen,omitempty"`
		FatherTestGenPictures []*multipart.FileHeader `form:"fatherTestGenPictures,omitempty"`
		MotherTestGenPictures []*multipart.FileHeader `form:"motherTestGenPictures,omitempty"`
		CallExpert            bool                    `form:"callExpert"`
		BirthCountry          *string                 `form:"birthCountry,omitempty"`
		Province              *string                 `form:"province,omitempty"`
		City                  *string                 `form:"city,omitempty"`
		Country               *string                 `form:"country,omitempty"`
		Address               string                  `form:"address" validate:"required"`
		PostalCode            string                  `form:"postalCode" validate:"required"`
		Education             string                  `form:"education" validate:"required"`
		Phone2                *string                 `form:"phone2"`
		Phone3                *string                 `form:"phone3"`
		BrotherNumber         *uint                   `form:"brotherNumber,omitempty"`
		SisterNumber          *uint                   `form:"sisterNumber,omitempty"`
		PaternalAuntNumber    *uint                   `form:"paternalAuntNumber,omitempty"`
		PaternalUncleNumber   *uint                   `form:"paternalUncleNumber,omitempty"`
		MaternalAuntNumber    *uint                   `form:"maternalAuntNumber,omitempty"`
		MaternalUncleNumber   *uint                   `form:"maternalUncleNumber,omitempty"`
	}

	params := controller.Validate[UpsertContactParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertContactRequest{
		UserID:                userID.(uint),
		FormID:                params.FormID,
		Name:                  params.Name,
		TestGen:               params.TestGen,
		TestGenPictures:       params.TestGenPictures,
		FmTestGen:             params.FmTestGen,
		FatherTestGenPictures: params.FatherTestGenPictures,
		MotherTestGenPictures: params.MotherTestGenPictures,
		CallExpert:            params.CallExpert,
		BirthCountry:          params.BirthCountry,
		Province:              params.Province,
		City:                  params.City,
		Country:               params.Country,
		Address:               params.Address,
		PostalCode:            params.PostalCode,
		Education:             params.Education,
		Phone2:                params.Phone2,
		Phone3:                params.Phone3,
		BrotherNumber:         params.BrotherNumber,
		SisterNumber:          params.SisterNumber,
		PaternalAuntNumber:    params.PaternalAuntNumber,
		PaternalUncleNumber:   params.PaternalUncleNumber,
		MaternalAuntNumber:    params.MaternalAuntNumber,
		MaternalUncleNumber:   params.MaternalUncleNumber,
	}

	if err := formController.formService.UpsertContact(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertContactResponse{})
}
func (formController *CustomerFormController) UpsertLungCancer(ctx *gin.Context) {
	type UpsertLungCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`

		InsuranceStatus              *string      `json:"insuranceStatus,omitempty"`
		SupplementaryInsuranceStatus *enum.Answer `json:"takmilBime,omitempty"`
		SupplementaryInsurances      *string      `json:"supplementaryInsurances,omitempty"`
		Hypertension                 *enum.Answer `json:"hypertension"`
		HypertensionTreatment        *enum.Answer `json:"hypertensionTreatment,omitempty"`
		HeartDisease                 *enum.Answer `json:"heartDisease"`
		HeartDiseaseTreatment        *enum.Answer `json:"heartDiseaseTreatment,omitempty"`
		Diabetes                     *enum.Answer `json:"diabetes"`
		DiabetesTreatment            *enum.Answer `json:"diabetesTreatment,omitempty"`
		ChronicLungDisease           *enum.Answer `json:"chronicLungDisease,omitempty"`
		ChronicLungDiseaseType       *string      `json:"chronicLungDiseaseType,omitempty"`
		LungCancerHistory            *enum.Answer `json:"lungCancerHistory"`
		OtherCancerHistory           *enum.Answer `json:"otherCancerHistory"`
		OtherCancerType              *uint        `json:"otherCancerType"`
		LungCancerFamily             *enum.Answer `json:"lungCancerFamily,omitempty"`
		LungCancerFamilyRelation     *string      `json:"lungCancerFamilyRelation,omitempty"`
		OtherCancerFamily            *enum.Answer `json:"otherCancerFamily,omitempty"`
		OtherCancerFamilyType        *uint        `json:"otherCancerFamilyType"`
		OtherCancerFamilyRelation    *string      `json:"otherCancerFamilyRelation,omitempty"`
		OccupationalExposure         *string      `json:"occupationalExposure,omitempty"`
		CurrentSmoking               *enum.Answer `json:"currentSmoking"`
		SmokingStartAgeCurrent       *uint        `json:"smokingStartAgeCurrent"`
		SmokingTypesCurrent          *string      `json:"smokingTypesCurrent,omitempty"`
		CigarettesPerDayCurrent      *uint        `json:"cigarettesPerDayCurrent"`
		CigarPerDayCurrent           *uint        `json:"cigarPerDayCurrent"`
		ECigPerDayCurrent            *uint        `json:"eCigPerDayCurrent"`
		PipePerDayCurrent            *uint        `json:"pipePerDayCurrent"`
		ChapoghPerDayCurrent         *uint        `json:"chapoghPerDayCurrent"`
		SmokedOpiumPerDayCurrent     *uint        `json:"smokedOpiumPerDayCurrent"`
		ChewedOpiumPerDayCurrent     *uint        `json:"chewedOpiumPerDayCurrent"`
		HookahPerWeekCurrent         *uint        `json:"hookahPerWeekCurrent"`
		PastSmoking                  *enum.Answer `json:"pastSmoking,omitempty"`
		SmokePastAvg                 *uint        `json:"smokePastAvg"`
		SmokeCurrentAvg              *uint        `json:"smokeCurrentAvg"`
		Bronchitis                   *enum.Answer `json:"bronshit"`
		LungIll                      *enum.Answer `json:"lungill"`
		Fibrosis                     *enum.Answer `json:"fibroz"`
		SmokingStartAgePast          *uint        `json:"smokingStartAgePast"`
		LeaveSmoke                   *uint        `json:"leaveSmoke"`
		SmokeTypePast                *string      `json:"smokeTypePast,omitempty"`
		SmokingTypesPast             *string      `json:"smokingTypesPast,omitempty"`
		CigarettesPerDayPast         *uint        `json:"cigarettesPerDayPast"`
		CigarPerDayPast              *uint        `json:"cigarPerDayPast"`
		ECigPerDayPast               *uint        `json:"eCigPerDayPast"`
		PipePerDayPast               *uint        `json:"pipePerDayPast"`
		ChapoghPerDayPast            *uint        `json:"chapoghPerDayPast"`
		SmokedOpiumPerDayPast        *uint        `json:"smokedOpiumPerDayPast"`
		ChewedOpiumPerDayPast        *uint        `json:"chewedOpiumPerDayPast"`
		HookahPerWeekPast            *uint        `json:"hookahPerWeekPast"`
		SecondhandSmoke              *enum.Answer `json:"secondhandSmoke"`
		SecondhandSmokeLocation      *string      `json:"secondhandSmokeLocation,omitempty"`
		AttentionCorrect             *bool        `json:"attentionCorrect"`
		LungDiseaseHistory           *string      `json:"lungDiseaseHistory,omitempty"`
	}

	params := controller.Validate[UpsertLungCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertLungCancerRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		SupplementaryInsuranceStatus: params.SupplementaryInsuranceStatus,
		InsuranceStatus:              params.InsuranceStatus,
		SupplementaryInsurances:      params.SupplementaryInsurances,
		Hypertension:                 params.Hypertension,
		HypertensionTreatment:        params.HypertensionTreatment,
		HeartDisease:                 params.HeartDisease,
		HeartDiseaseTreatment:        params.HeartDiseaseTreatment,
		Diabetes:                     params.Diabetes,
		DiabetesTreatment:            params.DiabetesTreatment,
		ChronicLungDisease:           params.ChronicLungDisease,
		ChronicLungDiseaseType:       params.ChronicLungDiseaseType,
		LungCancerHistory:            params.LungCancerHistory,
		OtherCancerHistory:           params.OtherCancerHistory,
		OtherCancerType:              params.OtherCancerType,
		LungCancerFamily:             params.LungCancerFamily,
		LungCancerFamilyRelation:     params.LungCancerFamilyRelation,
		OtherCancerFamily:            params.OtherCancerFamily,
		OtherCancerFamilyType:        params.OtherCancerFamilyType,
		OtherCancerFamilyRelation:    params.OtherCancerFamilyRelation,
		OccupationalExposure:         params.OccupationalExposure,
		CurrentSmoking:               params.CurrentSmoking,
		SmokingStartAgeCurrent:       params.SmokingStartAgeCurrent,
		SmokingTypesCurrent:          params.SmokingTypesCurrent,
		CigarettesPerDayCurrent:      params.CigarettesPerDayCurrent,
		CigarPerDayCurrent:           params.CigarPerDayCurrent,
		ECigPerDayCurrent:            params.ECigPerDayCurrent,
		PipePerDayCurrent:            params.PipePerDayCurrent,
		ChapoghPerDayCurrent:         params.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:     params.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:     params.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:         params.HookahPerWeekCurrent,
		PastSmoking:                  params.PastSmoking,
		SmokePastAvg:                 params.SmokePastAvg,
		SmokeCurrentAvg:              params.SmokeCurrentAvg,
		Bronchitis:                   params.Bronchitis,
		LungIll:                      params.LungIll,
		Fibrosis:                     params.Fibrosis,
		LeaveSmoke:                   params.LeaveSmoke,
		SmokingStartAgePast:          params.SmokingStartAgePast,
		SmokeTypePast:                params.SmokeTypePast,
		SmokingTypesPast:             params.SmokingTypesPast,
		CigarettesPerDayPast:         params.CigarettesPerDayPast,
		CigarPerDayPast:              params.CigarPerDayPast,
		ECigPerDayPast:               params.ECigPerDayPast,
		PipePerDayPast:               params.PipePerDayPast,
		ChapoghPerDayPast:            params.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:        params.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:        params.ChewedOpiumPerDayPast,
		HookahPerWeekPast:            params.HookahPerWeekPast,
		SecondhandSmoke:              params.SecondhandSmoke,
		SecondhandSmokeLocation:      params.SecondhandSmokeLocation,
		AttentionCorrect:             params.AttentionCorrect,
		LungDiseaseHistory:           params.LungDiseaseHistory,
	}

	if err := formController.formService.UpsertLungCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertLungCancerResponse{})
}

func (formController *CustomerFormController) UpsertNavidForm(ctx *gin.Context) {
	type UpsertNavidFormParams struct {
		FormID uint `uri:"formID" validate:"required"`

		InsuranceStatus              *string      `json:"insuranceStatus,omitempty"`
		SupplementaryInsuranceStatus *enum.Answer `json:"takmilBime,omitempty"`
		SupplementaryInsurances      *string      `json:"supplementaryInsurances,omitempty"`
		Hypertension                 string       `json:"hypertension"`
		HypertensionTreatment        *enum.Answer `json:"hypertensionTreatment,omitempty"`
		HeartDisease                 string       `json:"heartDisease"`
		HeartDiseaseTreatment        *enum.Answer `json:"heartDiseaseTreatment,omitempty"`
		Diabetes                     string       `json:"diabetes"`
		DiabetesTreatment            *enum.Answer `json:"diabetesTreatment,omitempty"`
		ChronicLungDisease           *enum.Answer `json:"chronicLungDisease,omitempty"`
		ChronicLungDiseaseType       *string      `json:"chronicLungDiseaseType,omitempty"`
		LungCancerHistory            *enum.Answer `json:"lungCancerHistory"`
		OtherCancerHistory           *enum.Answer `json:"otherCancerHistory"`
		OtherCancerType              *uint        `json:"otherCancerType"`
		LungCancerFamily             *enum.Answer `json:"lungCancerFamily,omitempty"`
		LungCancerFamilyRelation     *string      `json:"lungCancerFamilyRelation,omitempty"`
		OtherCancerFamily            *enum.Answer `json:"otherCancerFamily,omitempty"`
		OtherCancerFamilyType        *uint        `json:"otherCancerFamilyType"`
		OtherCancerFamilyRelation    *string      `json:"otherCancerFamilyRelation,omitempty"`
		OccupationalExposure         *string      `json:"occupationalExposure,omitempty"`
		CurrentSmoking               *enum.Answer `json:"currentSmoking"`
		SmokingStartAgeCurrent       *uint        `json:"smokingStartAgeCurrent"`
		SmokingTypesCurrent          *string      `json:"smokingTypesCurrent,omitempty"`
		CigarettesPerDayCurrent      *uint        `json:"cigarettesPerDayCurrent"`
		CigarPerDayCurrent           *uint        `json:"cigarPerDayCurrent"`
		ECigPerDayCurrent            *uint        `json:"eCigPerDayCurrent"`
		PipePerDayCurrent            *uint        `json:"pipePerDayCurrent"`
		ChapoghPerDayCurrent         *uint        `json:"chapoghPerDayCurrent"`
		SmokedOpiumPerDayCurrent     *uint        `json:"smokedOpiumPerDayCurrent"`
		ChewedOpiumPerDayCurrent     *uint        `json:"chewedOpiumPerDayCurrent"`
		HookahPerWeekCurrent         *uint        `json:"hookahPerWeekCurrent"`
		PastSmoking                  *enum.Answer `json:"pastSmoking,omitempty"`
		SmokingStartAgePast          *uint        `json:"smokingStartAgePast"`
		LeaveSmoke                   *uint        `json:"leaveSmoke"`
		SmokingTypesPast             *string      `json:"smokingTypesPast,omitempty"`
		CigarettesPerDayPast         *uint        `json:"cigarettesPerDayPast"`
		CigarPerDayPast              *uint        `json:"cigarPerDayPast"`
		ECigPerDayPast               *uint        `json:"eCigPerDayPast"`
		PipePerDayPast               *uint        `json:"pipePerDayPast"`
		ChapoghPerDayPast            *uint        `json:"chapoghPerDayPast"`
		SmokedOpiumPerDayPast        *uint        `json:"smokedOpiumPerDayPast"`
		ChewedOpiumPerDayPast        *uint        `json:"chewedOpiumPerDayPast"`
		HookahPerWeekPast            *uint        `json:"hookahPerWeekPast"`
		SecondhandSmoke              *enum.Answer `json:"secondhandSmoke"`
		SecondhandSmokeLocation      *string      `json:"secondhandSmokeLocation,omitempty"`
		AttentionCorrect             *bool        `json:"attentionCorrect"`
		LungDiseaseHistory           *string      `json:"lungDiseaseHistory,omitempty"`
		CurrentCigaretteSmoking      *enum.Answer `json:"Csig,omitempty"`
		CurrentRolledTobacco         *enum.Answer `json:"CsigBarg,omitempty"`
		CurrentPipeSmoking           *enum.Answer `json:"Cpip,omitempty"`
		CurrentHookahUse             *enum.Answer `json:"Cghel,omitempty"`
		CurrentChiboukSmoking        *enum.Answer `json:"Cchop,omitempty"`
		CurrentOpiumUse              *enum.Answer `json:"Cteryak,omitempty"`
		FormerCigaretteSmoking       *enum.Answer `json:"Psig,omitempty"`
		FormerRolledTobacco          *enum.Answer `json:"PsigBarg,omitempty"`
		FormerPipeSmoking            *enum.Answer `json:"Ppip,omitempty"`
		FormerHookahUse              *enum.Answer `json:"Pghel,omitempty"`
		FormerChiboukSmoking         *enum.Answer `json:"Pchop,omitempty"`
		FormerOpiumUse               *enum.Answer `json:"Pteryak,omitempty"`
		PelecSig                     *enum.Answer `json:"PelecSig,omitempty"`
		CelecSig                     *enum.Answer `json:"CelecSig,omitempty"`
	}

	params := controller.Validate[UpsertNavidFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertNavidFormRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		SupplementaryInsuranceStatus: params.SupplementaryInsuranceStatus,
		InsuranceStatus:              params.InsuranceStatus,
		SupplementaryInsurances:      params.SupplementaryInsurances,
		Hypertension:                 params.Hypertension,
		HypertensionTreatment:        params.HypertensionTreatment,
		HeartDisease:                 params.HeartDisease,
		HeartDiseaseTreatment:        params.HeartDiseaseTreatment,
		Diabetes:                     params.Diabetes,
		DiabetesTreatment:            params.DiabetesTreatment,
		ChronicLungDisease:           params.ChronicLungDisease,
		ChronicLungDiseaseType:       params.ChronicLungDiseaseType,
		LungCancerHistory:            params.LungCancerHistory,
		OtherCancerHistory:           params.OtherCancerHistory,
		OtherCancerType:              params.OtherCancerType,
		LungCancerFamily:             params.LungCancerFamily,
		LungCancerFamilyRelation:     params.LungCancerFamilyRelation,
		OtherCancerFamily:            params.OtherCancerFamily,
		OtherCancerFamilyType:        params.OtherCancerFamilyType,
		OtherCancerFamilyRelation:    params.OtherCancerFamilyRelation,
		OccupationalExposure:         params.OccupationalExposure,
		CurrentSmoking:               params.CurrentSmoking,
		SmokingStartAgeCurrent:       params.SmokingStartAgeCurrent,
		SmokingTypesCurrent:          params.SmokingTypesCurrent,
		CigarettesPerDayCurrent:      params.CigarettesPerDayCurrent,
		CigarPerDayCurrent:           params.CigarPerDayCurrent,
		ECigPerDayCurrent:            params.ECigPerDayCurrent,
		PipePerDayCurrent:            params.PipePerDayCurrent,
		ChapoghPerDayCurrent:         params.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:     params.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:     params.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:         params.HookahPerWeekCurrent,
		PastSmoking:                  params.PastSmoking,
		LeaveSmoke:                   params.LeaveSmoke,
		SmokingStartAgePast:          params.SmokingStartAgePast,
		SmokingTypesPast:             params.SmokingTypesPast,
		CigarettesPerDayPast:         params.CigarettesPerDayPast,
		CigarPerDayPast:              params.CigarPerDayPast,
		ECigPerDayPast:               params.ECigPerDayPast,
		PipePerDayPast:               params.PipePerDayPast,
		ChapoghPerDayPast:            params.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:        params.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:        params.ChewedOpiumPerDayPast,
		HookahPerWeekPast:            params.HookahPerWeekPast,
		SecondhandSmoke:              params.SecondhandSmoke,
		SecondhandSmokeLocation:      params.SecondhandSmokeLocation,
		AttentionCorrect:             params.AttentionCorrect,
		LungDiseaseHistory:           params.LungDiseaseHistory,
		CurrentCigaretteSmoking:      params.CurrentCigaretteSmoking,
		CurrentRolledTobacco:         params.CurrentRolledTobacco,
		CurrentPipeSmoking:           params.CurrentPipeSmoking,
		CurrentHookahUse:             params.CurrentHookahUse,
		CurrentChiboukSmoking:        params.CurrentChiboukSmoking,
		CurrentOpiumUse:              params.CurrentOpiumUse,
		FormerCigaretteSmoking:       params.FormerCigaretteSmoking,
		FormerRolledTobacco:          params.FormerRolledTobacco,
		FormerPipeSmoking:            params.FormerPipeSmoking,
		FormerHookahUse:              params.FormerHookahUse,
		FormerChiboukSmoking:         params.FormerChiboukSmoking,
		FormerOpiumUse:               params.FormerOpiumUse,
		PelecSig:                     params.PelecSig,
		CelecSig:                     params.CelecSig,
	}

	if err := formController.formService.UpsertNavidForm(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertNavidFormResponse{})
}

func (formController *CustomerFormController) ChangeFormStatus(ctx *gin.Context) {
	type ChangeFormStatusParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[ChangeFormStatusParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.ChangeFormStatusRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.ChangeFormStatus(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.changeFormStatus")
	controller.Response(ctx, 200, message, response)
}

func (formController *CustomerFormController) ResubmitRejectedForm(ctx *gin.Context) {
	type ResubmitRejectedFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[ResubmitRejectedFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.ResubmitRejectedForm(params.FormID, userID.(uint))
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.formResubmitted")
	controller.Response(ctx, 200, message, nil)
}

func (formController *CustomerFormController) SubmitDocuments(ctx *gin.Context) {
	type SubmitDocumentsParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[SubmitDocumentsParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.SubmitDocuments(params.FormID, userID.(uint))
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.documentsSubmitted")
	controller.Response(ctx, 200, message, nil)
}

func (formController *CustomerFormController) GetBasicForm(ctx *gin.Context) {
	type GetBasicFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetBasicFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetBasicForm(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *CustomerFormController) GetGeneralHealth(ctx *gin.Context) {
	type GetGeneralHealthParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetGeneralHealthParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetGeneralHealth(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *CustomerFormController) GetMamography(ctx *gin.Context) {
	type GetMamographyParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetMamographyParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetMamography(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *CustomerFormController) GetAllCancers(ctx *gin.Context) {
	type GetCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetCancers(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) VisitCancer(ctx *gin.Context) {
	type VisitCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[VisitCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.VisitCancer(userID.(uint), params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, nil)
}

func (formController *CustomerFormController) CreateCancer(ctx *gin.Context) {
	type CreateCancerParams struct {
		FormID     uint                    `uri:"formID" validate:"required"`
		CancerType uint                    `form:"cancerType" validate:"required,gt=0"`
		CancerAge  uint                    `form:"cancerAge" validate:"required,gte=0"`
		Pictures   []*multipart.FileHeader `form:"pictures,omitempty"`
	}

	params := controller.Validate[CreateCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.CreateCancerRequest{
		UserID:     userID.(uint),
		FormID:     params.FormID,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,
		Pictures:   params.Pictures,
	}

	err := formController.formService.CreateCancer(req)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, nil)
}

func (formController *CustomerFormController) UpdateCancer(ctx *gin.Context) {
	type UpdateCancerParams struct {
		FormID     uint                    `uri:"formID" validate:"required"`
		CancerID   uint                    `uri:"cancerID" validate:"required"`
		CancerType uint                    `form:"cancerType" validate:"required,gt=0"`
		CancerAge  uint                    `form:"cancerAge" validate:"required,gte=0"`
		Pictures   []*multipart.FileHeader `form:"pictures,omitempty"`
	}

	params := controller.Validate[UpdateCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateCancerRequest{
		UserID:     userID.(uint),
		FormID:     params.FormID,
		CancerID:   params.CancerID,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,
		Pictures:   params.Pictures,
	}

	response, err := formController.formService.UpdateCancer(req)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, response)
}

func (formController *CustomerFormController) DeleteCancer(ctx *gin.Context) {
	type DeleteCancerParams struct {
		FormID   uint `uri:"formID" validate:"required"`
		CancerID uint `uri:"cancerID" validate:"required"`
	}

	params := controller.Validate[DeleteCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.DeleteCancerRequest{
		UserID:   userID.(uint),
		FormID:   params.FormID,
		CancerID: params.CancerID,
	}

	err := formController.formService.DeleteCancer(req)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *CustomerFormController) VisitFamilyCancer(ctx *gin.Context) {
	type VisitFamilyCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[VisitFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.VisitFamilyCancer(userID.(uint), params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, nil)
}

func (formController *CustomerFormController) CreateFamilyCancer(ctx *gin.Context) {
	type CreateFamilyCancerParams struct {
		FormID           uint                    `uri:"formID" validate:"required"`
		Relative         uint                    `form:"relative" validate:"required,gt=0"`
		RelativeRelation *string                 `form:"relativeRelation,omitempty"`
		NumberRelative   *uint                   `form:"numberRelative,omitempty"`
		Name             *string                 `form:"name,omitempty"`
		LifeStatus       *uint                   `form:"lifeStatus"`
		CancerType       uint                    `form:"cancerType" validate:"required,gt=0"`
		CancerAge        uint                    `form:"cancerAge" validate:"required,gte=0"`
		Pictures         []*multipart.FileHeader `form:"pictures,omitempty"`
	}

	params := controller.Validate[CreateFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	var lifeStatus *enum.LifeStatus
	if params.LifeStatus != nil {
		ls := enum.LifeStatus(*params.LifeStatus)
		lifeStatus = &ls
	}

	req := formdto.CreateFamilyCancerRequest{
		UserID:           userID.(uint),
		FormID:           params.FormID,
		Relative:         enum.Relative(params.Relative),
		RelativeRelation: params.RelativeRelation,
		NumberRelative:   params.NumberRelative,
		Name:             params.Name,
		LifeStatus:       lifeStatus,
		CancerType:       params.CancerType,
		CancerAge:        params.CancerAge,
		Pictures:         params.Pictures,
	}

	response, err := formController.formService.CreateFamilyCancer(req)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createForm")
	controller.Response(ctx, 201, message, response)
}

func (formController *CustomerFormController) UpdateFamilyCancer(ctx *gin.Context) {
	type UpdateFamilyCancerParams struct {
		FormID           uint                    `uri:"formID" validate:"required"`
		FamilyCancerID   uint                    `uri:"familyCancerID" validate:"required"`
		Relative         uint                    `form:"relative" validate:"required,gt=0"`
		RelativeRelation *string                 `form:"relativeRelation,omitempty"`
		NumberRelative   *uint                   `form:"numberRelative,omitempty"`
		Name             *string                 `form:"name,omitempty"`
		LifeStatus       *uint                   `form:"lifeStatus"`
		CancerType       uint                    `form:"cancerType" validate:"required,gt=0"`
		CancerAge        uint                    `form:"cancerAge" validate:"required,gte=0"`
		Pictures         []*multipart.FileHeader `form:"pictures,omitempty"`
	}

	params := controller.Validate[UpdateFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	var lifeStatus *enum.LifeStatus
	if params.LifeStatus != nil {
		ls := enum.LifeStatus(*params.LifeStatus)
		lifeStatus = &ls
	}

	req := formdto.UpdateFamilyCancerRequest{
		UserID:           userID.(uint),
		FormID:           params.FormID,
		FamilyCancerID:   params.FamilyCancerID,
		Relative:         enum.Relative(params.Relative),
		RelativeRelation: params.RelativeRelation,
		NumberRelative:   params.NumberRelative,
		Name:             params.Name,
		LifeStatus:       lifeStatus,
		CancerType:       params.CancerType,
		CancerAge:        params.CancerAge,
		Pictures:         params.Pictures,
	}

	response, err := formController.formService.UpdateFamilyCancer(req)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, response)
}

func (formController *CustomerFormController) DeleteFamilyCancer(ctx *gin.Context) {
	type DeleteFamilyCancerParams struct {
		FormID         uint `uri:"formID" validate:"required"`
		FamilyCancerID uint `uri:"familyCancerID" validate:"required"`
	}

	params := controller.Validate[DeleteFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.DeleteFamilyCancerRequest{
		UserID:         userID.(uint),
		FormID:         params.FormID,
		FamilyCancerID: params.FamilyCancerID,
	}

	err := formController.formService.DeleteFamilyCancer(req)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", nil)
}

func (formController *CustomerFormController) GetFamilyCancer(ctx *gin.Context) {
	type GetFamilyCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetFamilyCancer(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *CustomerFormController) GetFamilyCancerList(ctx *gin.Context) {
	type GetFamilyCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetFamilyCancerList(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *CustomerFormController) GetContact(ctx *gin.Context) {
	type GetContactParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetContactParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetContact(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) GetLungCancer(ctx *gin.Context) {
	type GetLungCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetLungCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetLungCancer(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) GetAddressByPostalCode(ctx *gin.Context) {
	type PostalCodeParams struct {
		PostalCode string `json:"postalCode" validate:"required,len=10"`
	}

	params := controller.Validate[PostalCodeParams](ctx)

	response, err := formController.formService.GetAddressByPostalCode(params.PostalCode)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}

func (formController *CustomerFormController) GetNavidForm(ctx *gin.Context) {
	type GetNavidFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetNavidFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetNavidForm(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
