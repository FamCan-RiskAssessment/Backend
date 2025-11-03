package form

import (
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
		BirthDay             uint    `json:"birthDay" validate:"required"`
		BirthMonth           string  `json:"birthMonth" validate:"required"`
		BirthYear            uint    `json:"birthYear" validate:"required"`
		SocialSecurityNumber string  `json:"socialSecurityNumber" validate:"required"`
		Gender               uint    `json:"gender" validate:"required"`
		IsAtba               bool    `json:"isAtba"`
		Height               float64 `json:"height" validate:"required"`
		Weight               float64 `json:"weight" validate:"required"`
	}

	params := controller.Validate[CreateFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.CreateBasicFormRequest{
		UserID:               userID.(uint),
		BirthDay:             params.BirthDay,
		BirthMonth:           params.BirthMonth,
		BirthYear:            params.BirthYear,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,
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
		FormID               uint     `uri:"formID" validate:"required"`
		BirthDay             *uint    `json:"birthDay" validate:"required"`
		BirthMonth           *string  `json:"birthMonth" validate:"required"`
		BirthYear            *uint    `json:"birthYear" validate:"required"`
		SocialSecurityNumber *string  `json:"socialSecurityNumber" validate:"required"`
		Gender               *uint    `json:"gender" validate:"required"`
		IsAtba               *bool    `json:"isAtba"`
		Height               *float64 `json:"height" validate:"required"`
		Weight               *float64 `json:"weight" validate:"required"`
	}

	params := controller.Validate[UpdateBasicInfoParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.UpdateBasicFormRequest{
		FormID:               params.FormID,
		BirthDay:             params.BirthDay,
		BirthMonth:           params.BirthMonth,
		BirthYear:            params.BirthYear,
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

		GhaedeAge                    uint    `json:"ghaedeAge" validate:"required"`
		HasChildren                  bool    `json:"hasChildren"`
		NumberOfChildren             *uint   `json:"numberOfChildren,omitempty"`
		AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth,omitempty"`
		MenopausalStatus             uint    `json:"menopausalStatus" validate:"required"`
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
		NumberOfBreastBiopsies       *uint   `json:"numberOfBreastBiopsies,omitempty"`
		HyperplasiaInBiopsy          *uint   `json:"hyperplasiaInBiopsy,omitempty"`
	}

	params := controller.Validate[UpsertMamographyParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertMamographyRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
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
		NumberOfBreastBiopsies:       params.NumberOfBreastBiopsies,
		HyperplasiaInBiopsy:          params.HyperplasiaInBiopsy,
	}

	if err := formController.formService.UpsertMamography(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertMamographyResponse{})
}
func (formController *CustomerFormController) UpsertCancer(ctx *gin.Context) {
	type CancerParams struct {
		CancerType uint `json:"cancerType" validate:"required,gt=0"`
		CancerAge  uint `json:"cancerAge" validate:"required,gte=0"`
	}
	type UpsertCancerParams struct {
		FormID  uint           `uri:"formID" validate:"required"`
		Cancer  bool           `json:"cancer"`
		Cancers []CancerParams `json:"cancers" validate:"required_if=Cancer true,dive"`
	}

	params := controller.Validate[UpsertCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	var cancers []formdto.CancerRequest
	for _, v := range params.Cancers {
		cancers = append(cancers, formdto.CancerRequest{CancerType: v.CancerType, CancerAge: v.CancerAge})
	}
	req := formdto.UpsertCancerRequest{
		UserID:  userID.(uint),
		FormID:  params.FormID,
		Cancer:  params.Cancer,
		Cancers: cancers,
	}

	if err := formController.formService.UpsertCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertCancerResponse{})
}
func (formController *CustomerFormController) UpsertFamilyCancer(ctx *gin.Context) {
	type CancerParams struct {
		CancerType uint `json:"cancerType" validate:"required,gt=0"`
		CancerAge  uint `json:"cancerAge" validate:"required,gte=0"`
	}
	type FamilyCancerParams struct {
		FormID           uint           `uri:"formID" validate:"required"`
		Relative         uint           `json:"relative" validate:"required,gt=0"`
		RelativeRelation *string        `json:"relativeRelation,omitempty"`
		Name             *string        `json:"name,omitempty"`
		LifeStatus       *uint          `json:"lifeStatus,omitempty"`
		Cancer           bool           `json:"cancer" validate:"required"`
		Cancers          []CancerParams `json:"cancers" validate:"required_if=Cancer true,dive"`
	}

	params := controller.Validate[FamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	cancers := []formdto.CancerRequest{}
	for _, c := range params.Cancers {
		Cancer := formdto.CancerRequest{CancerType: c.CancerType, CancerAge: c.CancerAge}
		cancers = append(cancers, Cancer)
	}

	req := formdto.FamilyCancerRequest{
		UserID:           userID.(uint),
		FormID:           params.FormID,
		Relative:         enum.Relative(params.Relative),
		RelativeRelation: params.RelativeRelation,
		Name:             params.Name,
		LifeStatus:       (*enum.LifeStatus)(params.LifeStatus),
		Cancer:           params.Cancer,
		Cancers:          cancers,
	}

	if err := formController.formService.UpsertFamilyCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertFamilyCancerResponse{})
}
func (formController *CustomerFormController) UpsertContact(ctx *gin.Context) {
	type UpsertContactParams struct {
		FormID uint `uri:"formID" validate:"required"`

		Name         string  `json:"name" validate:"required"`
		TestGen      *bool   `json:"testGen,omitempty"`
		FmTestGen    *bool   `json:"fmTestGen,omitempty"`
		CallExpert   bool    `json:"callExpert"`
		BirthCountry *string `json:"birthCountry,omitempty"`
		Province     *string `json:"province,omitempty"`
		City         *string `json:"city,omitempty"`
		Country      *string `json:"country,omitempty"`
		Address      string  `json:"address" validate:"required"`
		PostalCode   string  `json:"postalCode" validate:"required"`
	}

	params := controller.Validate[UpsertContactParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertContactRequest{
		UserID:       userID.(uint),
		FormID:       params.FormID,
		Name:         params.Name,
		TestGen:      params.TestGen,
		FmTestGen:    params.FmTestGen,
		CallExpert:   params.CallExpert,
		BirthCountry: params.BirthCountry,
		Province:     params.Province,
		City:         params.City,
		Country:      params.Country,
		Address:      params.Address,
		PostalCode:   params.PostalCode,
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
		OtherCancerType           *uint   `json:"otherCancerType,omitempty"`
		LungCancerFamily          *bool   `json:"lungCancerFamily,omitempty"`
		LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation,omitempty"`
		OtherCancerFamily         *bool   `json:"otherCancerFamily,omitempty"`
		OtherCancerFamilyType     *uint   `json:"otherCancerFamilyType,omitempty"`
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

	params := controller.Validate[UpsertLungCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertLungCancerRequest{
		UserID:                    userID.(uint),
		FormID:                    params.FormID,
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

	if err := formController.formService.UpsertLungCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertLungCancerResponse{})
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
