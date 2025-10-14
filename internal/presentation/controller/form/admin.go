package form

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
	pagination  *bootstrap.Pagination
}

func NewAdminFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
	pagination *bootstrap.Pagination,
) *AdminFormController {
	return &AdminFormController{
		constants:   constants,
		formService: formService,
		pagination:  pagination,
	}
}

func (formController *AdminFormController) GetAllForms(ctx *gin.Context) {
	type GetAllFormsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	params := controller.Validate[GetAllFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	forms, count, err := formController.formService.GetAllForms(offset, limit)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(forms, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (formController *AdminFormController) DeleteForm(ctx *gin.Context) {
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

func (formController *AdminFormController) AcceptForm(ctx *gin.Context) {
	type AcceptFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[AcceptFormParams](ctx)

	err := formController.formService.AcceptForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.acceptForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) RejectForm(ctx *gin.Context) {
	type RejectFormParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[RejectFormParams](ctx)

	err := formController.formService.RejectForm(params.FormID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) GetBasicForm(ctx *gin.Context) {
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
func (formController *AdminFormController) UpdateBasicInfo(ctx *gin.Context) {
	type UpdateBasicInfoParams struct {
		FormID               uint    `uri:"formID" validate:"required"`
		BirthDay             uint    `json:"birthDay" validate:"required"`
		BirthMonth           string  `json:"birthMonth" validate:"required"`
		BirthYear            uint    `json:"birthYear" validate:"required"`
		SocialSecurityNumber string  `json:"socialSecurityNumber" validate:"required"`
		Gender               string  `json:"gender" validate:"required"`
		IsAtba               bool    `json:"isAtba"`
		Height               float64 `json:"height" validate:"required"`
		Weight               float64 `json:"weight" validate:"required"`
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
func (formController *AdminFormController) UpsertGeneralHealth(ctx *gin.Context) {
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
func (formController *AdminFormController) UpsertMamography(ctx *gin.Context) {
	type UpsertMamographyParams struct {
		FormID uint `uri:"formID" validate:"required"`

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
	}

	if err := formController.formService.UpsertMamography(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertMamographyResponse{})
}
func (formController *AdminFormController) UpsertCancer(ctx *gin.Context) {
	type UpsertCancerParams struct {
		FormID     uint    `uri:"formID" validate:"required"`
		Cancer     bool    `json:"cancer"`
		CancerType *string `json:"cancerType,omitempty"`
		CancerAge  *uint   `json:"cancerAge,omitempty"`
	}

	params := controller.Validate[UpsertCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertCancerRequest{
		UserID:     userID.(uint),
		FormID:     params.FormID,
		Cancer:     params.Cancer,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,
	}

	if err := formController.formService.UpsertCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertCancerResponse{})
}
func (formController *AdminFormController) UpsertFamilyCancer(ctx *gin.Context) {
	type UpsertFamilyCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`

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
	}

	params := controller.Validate[UpsertFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpsertFamilyCancerRequest{
		UserID:                  userID.(uint),
		FormID:                  params.FormID,
		ChildCancer:             params.ChildCancer,
		ChildName:               params.ChildName,
		ChildCancerType:         params.ChildCancerType,
		ChildCancerAge:          params.ChildCancerAge,
		ChildLifeStatus:         params.ChildLifeStatus,
		MotherCancer:            params.MotherCancer,
		MotherName:              params.MotherName,
		MotherLifeStatus:        params.MotherLifeStatus,
		MotherCancerType:        params.MotherCancerType,
		MotherCancerAge:         params.MotherCancerAge,
		FatherCancer:            params.FatherCancer,
		FatherName:              params.FatherName,
		FatherLifeStatus:        params.FatherLifeStatus,
		FatherCancerType:        params.FatherCancerType,
		FatherCancerAge:         params.FatherCancerAge,
		SiblingCancer:           params.SiblingCancer,
		SiblingName:             params.SiblingName,
		SiblingLifeStatus:       params.SiblingLifeStatus,
		SiblingCancerType:       params.SiblingCancerType,
		SiblingCancerAge:        params.SiblingCancerAge,
		AmeAmoCancer:            params.AmeAmoCancer,
		AmeAmoName:              params.AmeAmoName,
		AmeAmoLifeStatus:        params.AmeAmoLifeStatus,
		AmeAmoCancerType:        params.AmeAmoCancerType,
		AmeAmoCancerAge:         params.AmeAmoCancerAge,
		KhaleDaeiCancer:         params.KhaleDaeiCancer,
		KhaleDaeiName:           params.KhaleDaeiName,
		KhaleDaeiLifeStatus:     params.KhaleDaeiLifeStatus,
		KhaleDaeiCancerType:     params.KhaleDaeiCancerType,
		KhaleDaeiCancerAge:      params.KhaleDaeiCancerAge,
		OtherRelativeCancer:     params.OtherRelativeCancer,
		OtherRelativeName:       params.OtherRelativeName,
		OtherRelativeRelation:   params.OtherRelativeRelation,
		OtherRelativeLifeStatus: params.OtherRelativeLifeStatus,
		OtherRelativeCancerType: params.OtherRelativeCancerType,
		OtherRelativeCancerAge:  params.OtherRelativeCancerAge,
	}

	if err := formController.formService.UpsertFamilyCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertFamilyCancerResponse{})
}
func (formController *AdminFormController) UpsertContact(ctx *gin.Context) {
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
func (formController *AdminFormController) UpsertLungCancer(ctx *gin.Context) {
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
func (formController *AdminFormController) GetGeneralHealth(ctx *gin.Context) {
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
func (formController *AdminFormController) GetMamography(ctx *gin.Context) {
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
func (formController *AdminFormController) GetCancer(ctx *gin.Context) {
	type GetCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[GetCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.GetPartialFormRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	response, err := formController.formService.GetCancer(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", response)
}
func (formController *AdminFormController) GetFamilyCancer(ctx *gin.Context) {
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
func (formController *AdminFormController) GetContact(ctx *gin.Context) {
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
func (formController *AdminFormController) GetLungCancer(ctx *gin.Context) {
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
func (formController *AdminFormController) GetUserForms(ctx *gin.Context) {
	type GetUserFormsParams struct {
		UserID   uint `form:"userId"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}

	params := controller.Validate[GetUserFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	request := formdto.GetUserFormsRequest{
		UserID: params.UserID,
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
