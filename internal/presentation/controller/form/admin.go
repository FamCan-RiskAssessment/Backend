package form

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminFormController struct {
	constants   *bootstrap.Constants
	formService usecase.FormService
	userService usecase.UserService
	pagination  *bootstrap.Pagination
}

func NewAdminFormController(
	constants *bootstrap.Constants,
	formService usecase.FormService,
	userService usecase.UserService,
	pagination *bootstrap.Pagination,
) *AdminFormController {
	return &AdminFormController{
		constants:   constants,
		formService: formService,
		userService: userService,
		pagination:  pagination,
	}
}

func (formController *AdminFormController) GetAllForms(ctx *gin.Context) {
	type GetAllFormsParams struct {
		Page               int     `form:"page"`
		PageSize           int     `form:"pageSize"`
		Status             *uint   `form:"status"`
		Gender             *string `form:"gender"`
		BirthYear          *uint   `form:"birthYear"`
		DrinksAlcohol      *bool   `form:"drinksAlcohol"`
		SmokingNow         *bool   `form:"smokingNow"`
		Cancer             *bool   `form:"cancer"`
		FilledByOperatorID *uint   `form:"filledByOperatorID"`
	}

	params := controller.Validate[GetAllFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	filters := &postgres.FormFilters{
		Status:             params.Status,
		Gender:             params.Gender,
		BirthYear:          params.BirthYear,
		DrinksAlcohol:      params.DrinksAlcohol,
		SmokingNow:         params.SmokingNow,
		Cancer:             params.Cancer,
		FilledByOperatorID: params.FilledByOperatorID,
	}

	forms, count, err := formController.formService.GetAllForms(offset, limit, filters)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(forms, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (formController *AdminFormController) GetAllOperatorForms(ctx *gin.Context) {
	type GetAllOperatorFormsParams struct {
		Page          int     `form:"page"`
		PageSize      int     `form:"pageSize"`
		Status        *uint   `form:"status"`
		Gender        *string `form:"gender"`
		BirthYear     *uint   `form:"birthYear"`
		DrinksAlcohol *bool   `form:"drinksAlcohol"`
		SmokingNow    *bool   `form:"smokingNow"`
		Cancer        *bool   `form:"cancer"`
	}

	params := controller.Validate[GetAllOperatorFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	filters := &postgres.OperatorFormFilters{
		OperatorID:    userID.(uint),
		Status:        params.Status,
		Gender:        params.Gender,
		BirthYear:     params.BirthYear,
		DrinksAlcohol: params.DrinksAlcohol,
		SmokingNow:    params.SmokingNow,
		Cancer:        params.Cancer,
	}

	forms, count, err := formController.formService.GetAllOperatorForms(offset, limit, filters)
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
		FormID               uint       `uri:"formID" validate:"required"`
		BirthDate            *time.Time `json:"birthDate" validate:"required" time_format:"2006-01-02"`
		SocialSecurityNumber *string    `json:"socialSecurityNumber"`
		Gender               *uint      `json:"gender"`
		IsAtba               *bool      `json:"isAtba"`
		Height               *float64   `json:"height"`
		Weight               *float64   `json:"weight"`
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

func (formController *AdminFormController) UpdateGeneralHealth(ctx *gin.Context) {
	type UpdateGeneralHealthParams struct {
		FormID uint `uri:"formID" validate:"required"`

		DrinksAlcohol             *bool   `json:"drinksAlcohol"`
		CupsPerWeek               *string `json:"cupsPerWeek"`
		LastMonthSabzijatMeal     *string `json:"lastMonthSabzijatMeal"`
		LastMonthSabzijatWeight   *string `json:"lastMonthSabzijatWeight"`
		MediumActivityMonthInYear *uint   `json:"mediumActivityMonthInYear"`
		MediumActivityHourInWeek  *string `json:"mediumActivityHourInWeek"`
		HardActivityMonthInYear   *uint   `json:"hardActivityMonthInYear"`
		HardActivityHourInWeek    *string `json:"hardActivityHourInWeek"`
		SmokeAtLeast100           *bool   `json:"smokeAtLeast100"`
		SmokingAge                *uint   `json:"smokingAge"`
		SmokingNow                *bool   `json:"smokingNow"`
		LeaveSmokingAge           *uint   `json:"leaveSmokingAge"`
		CountSmokingDaily         *string `json:"countSmokingDaily"`
		CountGheliandaily         *string `json:"countGheliandaily"`
		CountSmokingDailyPast     *string `json:"countSmokingDailyPast"`
		CountGheliandailyPast     *string `json:"countGheliandailyPast"`
	}

	params := controller.Validate[UpdateGeneralHealthParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateGeneralHealthRequest{
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

	if err := formController.formService.UpdateGeneralHealth(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertGeneralHealthResponse{})
}

func (formController *AdminFormController) UpdateMamography(ctx *gin.Context) {
	type UpdateMamographyParams struct {
		FormID uint `uri:"formID" validate:"required"`

		GhaedeAge                    *uint   `json:"ghaedeAge"`
		HasChildren                  *bool   `json:"hasChildren"`
		NumberOfChildren             *uint   `json:"numberOfChildren"`
		AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth"`
		MenopausalStatus             *uint   `json:"menopausalStatus"`
		MenopauseAge                 *string `json:"menopauseAge"`
		HRT                          *bool   `json:"hrt"`
		HRTUseLength                 *uint   `json:"hrtUseLength"`
		LastFiveYearsHRTUse          *bool   `json:"lastFiveYearsHrtUse"`
		CurrentHRTUse                *bool   `json:"currentHrtUse"`
		IntendedHRTUse               *uint   `json:"intendedHrtUse"`
		HRTType                      *string `json:"hrtType"`
		Oral                         *bool   `json:"oral"`
		OralDuration                 *string `json:"oralDuration"`
		OralTwoLastYears             *bool   `json:"oralTwoLastYears"`
		MamoGraphy                   *bool   `json:"mamoGraphy"`
		Falop                        *bool   `json:"falop"`
		Andometrioz                  *bool   `json:"andometrioz"`
		LeavePestan                  *bool   `json:"leavePestan"`
		LeaveTokhmdan                *bool   `json:"leaveTokhmdan"`
		LaDeColon                    *bool   `json:"laDeColon"`
		LaDePol                      *bool   `json:"laDePol"`
		AspLaMo                      *bool   `json:"aspLaMo"`
		NsaiDLaMo                    *bool   `json:"nsaiDLaMo"`
		LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool"`
	}

	params := controller.Validate[UpdateMamographyParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateMamographyRequest{
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

	if err := formController.formService.UpdateMamography(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertMamographyResponse{})
}

func (formController *AdminFormController) UpdateCancer(ctx *gin.Context) {
	type UpdateCancerParams struct {
		FormID     uint  `uri:"formID" validate:"required"`
		Cancer     *bool `json:"cancer"`
		CancerType *uint `json:"cancerType"`
		CancerAge  *uint `json:"cancerAge"`
	}

	params := controller.Validate[UpdateCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateCancerRequest{
		UserID:     userID.(uint),
		FormID:     params.FormID,
		Cancer:     params.Cancer,
		CancerType: params.CancerType,
		CancerAge:  params.CancerAge,
	}

	if err := formController.formService.UpdateCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertCancerResponse{})
}

func (formController *AdminFormController) UpdateFamilyCancer(ctx *gin.Context) {
	type UpdateFamilyCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`

		ChildCancer     *bool   `json:"childCancer"`
		ChildName       *string `json:"childName"`
		ChildCancerType *uint   `json:"childCancerType"`
		ChildCancerAge  *uint   `json:"childCancerAge"`
		ChildLifeStatus *string `json:"childLifeStatus"`

		MotherCancer     *bool   `json:"motherCancer"`
		MotherName       *string `json:"motherName"`
		MotherLifeStatus *string `json:"motherLifeStatus"`
		MotherCancerType *uint   `json:"motherCancerType"`
		MotherCancerAge  *uint   `json:"motherCancerAge"`

		FatherCancer     *bool   `json:"fatherCancer"`
		FatherName       *string `json:"fatherName"`
		FatherLifeStatus *string `json:"fatherLifeStatus"`
		FatherCancerType *uint   `json:"fatherCancerType"`
		FatherCancerAge  *uint   `json:"fatherCancerAge"`

		SiblingCancer     *bool   `json:"siblingCancer"`
		SiblingName       *string `json:"siblingName"`
		SiblingLifeStatus *string `json:"siblingLifeStatus"`
		SiblingCancerType *uint   `json:"siblingCancerType"`
		SiblingCancerAge  *uint   `json:"siblingCancerAge"`

		AmeAmoCancer     *bool   `json:"ameAmoCancer"`
		AmeAmoName       *string `json:"ameAmoName"`
		AmeAmoLifeStatus *string `json:"ameAmoLifeStatus"`
		AmeAmoCancerType *uint   `json:"ameAmoCancerType"`
		AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge"`

		KhaleDaeiCancer     *bool   `json:"khaleDaeiCancer"`
		KhaleDaeiName       *string `json:"khaleDaeiName"`
		KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus"`
		KhaleDaeiCancerType *uint   `json:"khaleDaeiCancerType"`
		KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge"`

		OtherRelativeCancer     *bool   `json:"otherRelativeCancer"`
		OtherRelativeName       *string `json:"otherRelativeName"`
		OtherRelativeRelation   *string `json:"otherRelativeRelation"`
		OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus"`
		OtherRelativeCancerType *uint   `json:"otherRelativeCancerType"`
		OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge"`
	}

	params := controller.Validate[UpdateFamilyCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateFamilyCancerRequest{
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

	if err := formController.formService.UpdateFamilyCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertFamilyCancerResponse{})
}

func (formController *AdminFormController) UpdateContact(ctx *gin.Context) {
	type UpdateContactParams struct {
		FormID uint `uri:"formID" validate:"required"`

		Name         *string `json:"name"`
		TestGen      *bool   `json:"testGen"`
		FmTestGen    *bool   `json:"fmTestGen"`
		CallExpert   *bool   `json:"callExpert"`
		BirthCountry *string `json:"birthCountry"`
		Province     *string `json:"province"`
		City         *string `json:"city"`
		Country      *string `json:"country"`
		Address      *string `json:"address"`
		PostalCode   *string `json:"postalCode"`
	}

	params := controller.Validate[UpdateContactParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateContactRequest{
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

	if err := formController.formService.UpdateContact(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertContactResponse{})
}

func (formController *AdminFormController) UpdateLungCancer(ctx *gin.Context) {
	type UpdateLungCancerParams struct {
		FormID uint `uri:"formID" validate:"required"`

		InsuranceStatus           *string `json:"insuranceStatus"`
		SupplementaryInsurances   *string `json:"supplementaryInsurances"`
		Hypertension              *bool   `json:"hypertension"`
		HypertensionTreatment     *bool   `json:"hypertensionTreatment"`
		HeartDisease              *bool   `json:"heartDisease"`
		HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment"`
		Diabetes                  *bool   `json:"diabetes"`
		DiabetesTreatment         *bool   `json:"diabetesTreatment"`
		ChronicLungDisease        *bool   `json:"chronicLungDisease"`
		ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType"`
		LungCancerHistory         *bool   `json:"lungCancerHistory"`
		OtherCancerHistory        *bool   `json:"otherCancerHistory"`
		OtherCancerType           *uint   `json:"otherCancerType"`
		LungCancerFamily          *bool   `json:"lungCancerFamily"`
		LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation"`
		OtherCancerFamily         *bool   `json:"otherCancerFamily"`
		OtherCancerFamilyType     *uint   `json:"otherCancerFamilyType"`
		OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation"`
		OccupationalExposure      *string `json:"occupationalExposure"`
		CurrentSmoking            *bool   `json:"currentSmoking"`
		SmokingStartAgeCurrent    *uint   `json:"smokingStartAgeCurrent"`
		SmokingTypesCurrent       *string `json:"smokingTypesCurrent"`
		CigarettesPerDayCurrent   *uint   `json:"cigarettesPerDayCurrent"`
		CigarPerDayCurrent        *uint   `json:"cigarPerDayCurrent"`
		ECigPerDayCurrent         *uint   `json:"eCigPerDayCurrent"`
		PipePerDayCurrent         *uint   `json:"pipePerDayCurrent"`
		ChapoghPerDayCurrent      *uint   `json:"chapoghPerDayCurrent"`
		SmokedOpiumPerDayCurrent  *uint   `json:"smokedOpiumPerDayCurrent"`
		ChewedOpiumPerDayCurrent  *uint   `json:"chewedOpiumPerDayCurrent"`
		HookahPerWeekCurrent      *uint   `json:"hookahPerWeekCurrent"`
		PastSmoking               *string `json:"pastSmoking"`
		SmokingStartAgePast       *uint   `json:"smokingStartAgePast"`
		SmokingTypesPast          *string `json:"smokingTypesPast"`
		CigarettesPerDayPast      *uint   `json:"cigarettesPerDayPast"`
		CigarPerDayPast           *uint   `json:"cigarPerDayPast"`
		ECigPerDayPast            *uint   `json:"eCigPerDayPast"`
		PipePerDayPast            *uint   `json:"pipePerDayPast"`
		ChapoghPerDayPast         *uint   `json:"chapoghPerDayPast"`
		SmokedOpiumPerDayPast     *uint   `json:"smokedOpiumPerDayPast"`
		ChewedOpiumPerDayPast     *uint   `json:"chewedOpiumPerDayPast"`
		HookahPerWeekPast         *uint   `json:"hookahPerWeekPast"`
		SecondhandSmoke           *bool   `json:"secondhandSmoke"`
		SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation"`
	}

	params := controller.Validate[UpdateLungCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateLungCancerRequest{
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

	if err := formController.formService.UpdateLungCancer(req); err != nil {
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

func (formController *AdminFormController) AssignOperator(ctx *gin.Context) {
	type AssignOperatorParams struct {
		FormID     uint `uri:"formID" validate:"required"`
		OperatorID uint `json:"operatorId" validate:"required"`
	}

	params := controller.Validate[AssignOperatorParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.AssignOperatorRequest{
		UserID:     userID.(uint),
		FormID:     params.FormID,
		OperatorID: params.OperatorID,
	}

	err := formController.formService.AssignOperator(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, nil)
}
func (formController *AdminFormController) UnassignOperator(ctx *gin.Context) {
	type UnassignOperatorParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}

	params := controller.Validate[UnassignOperatorParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.UnassignOperatorRequest{
		UserID: userID.(uint),
		FormID: params.FormID,
	}

	err := formController.formService.UnassignOperator(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) CreateFormForUser(ctx *gin.Context) {
	type CreateFormForUserParams struct {
		UserID               uint      `json:"userId" validate:"required"`
		BirthDate            time.Time `json:"birthDate" validate:"required" time_format:"2006-01-02"`
		SocialSecurityNumber string    `json:"socialSecurityNumber" validate:"required"`
		Gender               uint      `json:"gender" validate:"required"`
		IsAtba               bool      `json:"isAtba"`
		Height               float64   `json:"height" validate:"required"`
		Weight               float64   `json:"weight" validate:"required"`
	}

	params := controller.Validate[CreateFormForUserParams](ctx)

	operatorID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.CreateBasicFormRequest{
		UserID:               params.UserID,
		FilledByOperatorID:   &[]uint{operatorID.(uint)}[0],
		BirthDate:            params.BirthDate,
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

func (formController *AdminFormController) RequestUserValidationOTP(ctx *gin.Context) {
	type RequestUserValidationOTPParams struct {
		Phone string `json:"phone" validate:"required,min=11,max=11"`
	}

	params := controller.Validate[RequestUserValidationOTPParams](ctx)
	operatorID, _ := ctx.Get(formController.constants.Context.ID)

	request := userdto.RequestUserValidationOTPRequest{
		Phone: params.Phone,
	}

	err := formController.userService.RequestUserValidationOTP(operatorID.(uint), request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.sendOTP")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) VerifyUserValidationOTP(ctx *gin.Context) {
	type VerifyUserValidationOTPParams struct {
		Phone string `json:"phone" validate:"required,min=11,max=11"`
		OTP   string `json:"otp" validate:"required,min=6,max=6"`
	}

	params := controller.Validate[VerifyUserValidationOTPParams](ctx)
	operatorID, _ := ctx.Get(formController.constants.Context.ID)

	request := userdto.VerifyUserValidationOTPRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}

	response, err := formController.userService.VerifyUserValidationOTP(operatorID.(uint), request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.userVerified")
	controller.Response(ctx, 200, message, response)
}

func (formController *AdminFormController) GetOperatorFormsFilled(ctx *gin.Context) {

}
