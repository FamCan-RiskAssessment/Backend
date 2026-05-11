package form

import (
	"mime/multipart"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
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
		Page               int            `form:"page"`
		PageSize           int            `form:"pageSize"`
		FormType           *enum.FormType `form:"formType"`
		Status             *uint          `form:"status"`
		Gender             *string        `form:"gender"`
		BirthYear          *uint          `form:"birthYear"`
		DrinksAlcohol      *bool          `form:"drinksAlcohol"`
		SmokingNow         *bool          `form:"smokingNow"`
		Cancer             *bool          `form:"cancer"`
		FilledByOperatorID *uint          `form:"filledByOperatorID"`
		SortBy             *string        `form:"sortBy"`
		SortOrder          *string        `form:"sortOrder"`
		Search             *string        `form:"search"`
	}

	params := controller.Validate[GetAllFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	filters := &postgres.FormFilters{
		Status:             params.Status,
		FormType:           params.FormType,
		Gender:             params.Gender,
		BirthYear:          params.BirthYear,
		DrinksAlcohol:      params.DrinksAlcohol,
		SmokingNow:         params.SmokingNow,
		Cancer:             params.Cancer,
		FilledByOperatorID: params.FilledByOperatorID,
	}

	forms, count, err := formController.formService.GetAllForms(offset, limit, filters, params.SortBy, params.SortOrder, params.Search)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(forms, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (formController *AdminFormController) GetAllOperatorForms(ctx *gin.Context) {
	type GetAllOperatorFormsParams struct {
		Page          int            `form:"page"`
		PageSize      int            `form:"pageSize"`
		FormType      *enum.FormType `form:"formType"`
		Status        *uint          `form:"status"`
		Gender        *string        `form:"gender"`
		BirthYear     *uint          `form:"birthYear"`
		DrinksAlcohol *bool          `form:"drinksAlcohol"`
		SmokingNow    *bool          `form:"smokingNow"`
		Cancer        *bool          `form:"cancer"`
		SortBy        *string        `form:"sortBy"`
		SortOrder     *string        `form:"sortOrder"`
		Search        *string        `form:"search"`
	}

	params := controller.Validate[GetAllOperatorFormsParams](ctx)
	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, formController.pagination.DefaultPage, formController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	filters := &postgres.OperatorFormFilters{
		OperatorID:    userID.(uint),
		Status:        params.Status,
		FormType:      params.FormType,
		Gender:        params.Gender,
		BirthYear:     params.BirthYear,
		DrinksAlcohol: params.DrinksAlcohol,
		SmokingNow:    params.SmokingNow,
		Cancer:        params.Cancer,
	}

	forms, count, err := formController.formService.GetAllOperatorForms(offset, limit, filters, params.SortBy, params.SortOrder, params.Search)
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

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.AcceptForm(params.FormID, userID.(uint))
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

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.RejectForm(params.FormID, userID.(uint))
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectForm")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) RequestPatientResponse(ctx *gin.Context) {
	type RequestPatientResponseParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[RequestPatientResponseParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.RequestPatientResponse(params.FormID, userID.(uint))
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.requestPatientResponse")
	controller.Response(ctx, 200, message, nil)
}

func (formController *AdminFormController) RequestDocuments(ctx *gin.Context) {
	type RequestDocumentsParams struct {
		FormID uint `uri:"formID" validate:"required"`
	}
	params := controller.Validate[RequestDocumentsParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	err := formController.formService.RequestDocuments(params.FormID, userID.(uint))
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.requestDocuments")
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

	operatorID, _ := ctx.Get(formController.constants.Context.ID)

	request := formdto.UpdateBasicFormRequest{
		FormID:               params.FormID,
		BirthDate:            params.BirthDate,
		SocialSecurityNumber: params.SocialSecurityNumber,
		Gender:               params.Gender,
		IsAtba:               params.IsAtba,
		Height:               params.Height,
		Weight:               params.Weight,
		UserID:               userID.(uint),
		FilledByOperatorID:   &[]uint{operatorID.(uint)}[0],
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

		DrinksAlcohol             *enum.Answer `json:"drinksAlcohol"`
		CupsPerWeek               *string      `json:"cupsPerWeek"`
		LastMonthSabzijatMeal     *string      `json:"lastMonthSabzijatMeal"`
		LastMonthSabzijatWeight   *string      `json:"lastMonthSabzijatWeight"`
		MediumActivityMonthInYear *uint        `json:"mediumActivityMonthInYear"`
		MediumActivityHourInWeek  *string      `json:"mediumActivityHourInWeek"`
		HardActivityMonthInYear   *uint        `json:"hardActivityMonthInYear"`
		HardActivityHourInWeek    *string      `json:"hardActivityHourInWeek"`
		SmokeAtLeast100           *enum.Answer `json:"smokeAtLeast100"`
		SmokingAge                *uint        `json:"smokingAge"`
		SmokingNow                *enum.Answer `json:"smokingNow"`
		YearSmoke                 *uint        `json:"yearSmoke"`
		LeaveSmokingAge           *uint        `json:"leaveSmokingAge"`
		CountSmokingDaily         *string      `json:"countSmokingDaily"`
		CountGheliandaily         *string      `json:"countGheliandaily"`
		CountSmokingDailyPast     *string      `json:"countSmokingDailyPast"`
		CountGheliandailyPast     *string      `json:"countGheliandailyPast"`
		AttentionCorrect          *bool        `json:"attentionCorrect"`
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
		YearSmoke:                 params.YearSmoke,
		LeaveSmokingAge:           params.LeaveSmokingAge,
		CountSmokingDaily:         params.CountSmokingDaily,
		CountGheliandaily:         params.CountGheliandaily,
		CountSmokingDailyPast:     params.CountSmokingDailyPast,
		CountGheliandailyPast:     params.CountGheliandailyPast,
		AttentionCorrect:          params.AttentionCorrect,
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
		Oral                         *enum.Answer            `form:"oral"`
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
		AttentionCorrect             *bool                   `form:"attentionCorrect"`
	}

	params := controller.Validate[UpdateMamographyParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateMamographyRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		GhaedeAge:                    &params.GhaedeAge,
		HasChildren:                  params.HasChildren,
		NumberOfChildren:             params.NumberOfChildren,
		SonCount:                     params.SonCount,
		DaughterCount:                params.DaughterCount,
		AgeOfFirstBirth:              params.AgeOfFirstBirth,
		MenopausalStatus:             &params.MenopausalStatus,
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
		LeavePestan:                  &params.LeavePestan,
		LeaveTokhmdan:                &params.LeaveTokhmdan,
		LaDeColon:                    params.LaDeColon,
		LaDePol:                      params.LaDePol,
		AspLaMo:                      params.AspLaMo,
		NsaiDLaMo:                    params.NsaiDLaMo,
		LastFiveYearBloodTestInStool: params.LastFiveYearBloodTestInStool,
		AttentionCorrect:             params.AttentionCorrect,
	}

	if err := formController.formService.UpdateMamography(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertMamographyResponse{})
}

func (formController *AdminFormController) UpdateContact(ctx *gin.Context) {
	type UpdateContactParams struct {
		FormID uint `uri:"formID" validate:"required"`

		Name                  *string                 `form:"name"`
		TestGen               *enum.Answer            `form:"testGen"`
		TestGenPictures       []*multipart.FileHeader `form:"testGenPictures,omitempty"`
		FmTestGen             *enum.Answer            `form:"fmTestGen"`
		FatherTestGenPictures []*multipart.FileHeader `form:"fatherTestGenPictures,omitempty"`
		MotherTestGenPictures []*multipart.FileHeader `form:"motherTestGenPictures,omitempty"`
		CallExpert            *bool                   `form:"callExpert"`
		BirthCountry          *string                 `form:"birthCountry"`
		Province              *string                 `form:"province"`
		City                  *string                 `form:"city"`
		Country               *string                 `form:"country"`
		Address               *string                 `form:"address"`
		PostalCode            *string                 `form:"postalCode"`
		Education             *string                 `form:"education"`
		Phone2                *string                 `form:"phone2"`
		Phone3                *string                 `form:"phone3"`
		BrotherNumber         *uint                   `form:"brotherNumber,omitempty"`
		SisterNumber          *uint                   `form:"sisterNumber,omitempty"`
		PaternalAuntNumber    *uint                   `form:"paternalAuntNumber,omitempty"`
		PaternalUncleNumber   *uint                   `form:"paternalUncleNumber,omitempty"`
		MaternalAuntNumber    *uint                   `form:"maternalAuntNumber,omitempty"`
		MaternalUncleNumber   *uint                   `form:"maternalUncleNumber,omitempty"`
	}

	params := controller.Validate[UpdateContactParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateContactRequest{
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

		InsuranceStatus              *string      `json:"insuranceStatus"`
		SupplementaryInsuranceStatus *enum.Answer `json:"takmilBime"`
		SupplementaryInsurances      *string      `json:"supplementaryInsurances"`
		Hypertension                 *enum.Answer `json:"hypertension"`
		HypertensionTreatment        *enum.Answer `json:"hypertensionTreatment"`
		HeartDisease                 *enum.Answer `json:"heartDisease"`
		HeartDiseaseTreatment        *enum.Answer `json:"heartDiseaseTreatment"`
		Diabetes                     *enum.Answer `json:"diabetes"`
		DiabetesTreatment            *enum.Answer `json:"diabetesTreatment"`
		ChronicLungDisease           *enum.Answer `json:"chronicLungDisease"`
		ChronicLungDiseaseType       *string      `json:"chronicLungDiseaseType"`
		LungCancerHistory            *enum.Answer `json:"lungCancerHistory"`
		OtherCancerHistory           *enum.Answer `json:"otherCancerHistory"`
		OtherCancerType              *uint        `json:"otherCancerType"`
		LungCancerFamily             *enum.Answer `json:"lungCancerFamily"`
		LungCancerFamilyRelation     *string      `json:"lungCancerFamilyRelation"`
		OtherCancerFamily            *enum.Answer `json:"otherCancerFamily"`
		OtherCancerFamilyType        *uint        `json:"otherCancerFamilyType"`
		OtherCancerFamilyRelation    *string      `json:"otherCancerFamilyRelation"`
		OccupationalExposure         *string      `json:"occupationalExposure"`
		CurrentSmoking               *enum.Answer `json:"currentSmoking"`
		SmokingStartAgeCurrent       *uint        `json:"smokingStartAgeCurrent"`
		SmokingTypesCurrent          *string      `json:"smokingTypesCurrent"`
		CigarettesPerDayCurrent      *uint        `json:"cigarettesPerDayCurrent"`
		CigarPerDayCurrent           *uint        `json:"cigarPerDayCurrent"`
		ECigPerDayCurrent            *uint        `json:"eCigPerDayCurrent"`
		PipePerDayCurrent            *uint        `json:"pipePerDayCurrent"`
		ChapoghPerDayCurrent         *uint        `json:"chapoghPerDayCurrent"`
		SmokedOpiumPerDayCurrent     *uint        `json:"smokedOpiumPerDayCurrent"`
		ChewedOpiumPerDayCurrent     *uint        `json:"chewedOpiumPerDayCurrent"`
		HookahPerWeekCurrent         *uint        `json:"hookahPerWeekCurrent"`
		PastSmoking                  *enum.Answer `json:"pastSmoking"`
		SmokePastAvg                 *uint        `json:"smokePastAvg"`
		SmokeCurrentdAvg             *uint        `json:"smokeCurrentAvg"`
		Bronchitis                   *enum.Answer `json:"bronshit"`
		LungIll                      *enum.Answer `json:"lungill"`
		Fibrosis                     *enum.Answer `json:"fibroz"`
		LeaveSmoke                   *uint        `json:"leaveSmoke"`
		SmokingStartAgePast          *uint        `json:"smokingStartAgePast"`
		SmokingTypesPast             *string      `json:"smokingTypesPast"`
		CigarettesPerDayPast         *uint        `json:"cigarettesPerDayPast"`
		CigarPerDayPast              *uint        `json:"cigarPerDayPast"`
		ECigPerDayPast               *uint        `json:"eCigPerDayPast"`
		PipePerDayPast               *uint        `json:"pipePerDayPast"`
		ChapoghPerDayPast            *uint        `json:"chapoghPerDayPast"`
		SmokedOpiumPerDayPast        *uint        `json:"smokedOpiumPerDayPast"`
		ChewedOpiumPerDayPast        *uint        `json:"chewedOpiumPerDayPast"`
		HookahPerWeekPast            *uint        `json:"hookahPerWeekPast"`
		SecondhandSmoke              *enum.Answer `json:"secondhandSmoke"`
		SecondhandSmokeLocation      *string      `json:"secondhandSmokeLocation"`
		AttentionCorrect             *bool        `json:"attentionCorrect"`
		LungDiseaseHistory           *string      `json:"lungDiseaseHistory"`
	}

	params := controller.Validate[UpdateLungCancerParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateLungCancerRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		InsuranceStatus:              params.InsuranceStatus,
		SupplementaryInsuranceStatus: params.SupplementaryInsuranceStatus,
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
		SmokeCurrentAvg:              params.SmokeCurrentdAvg,
		Bronchitis:                   params.Bronchitis,
		LungIll:                      params.LungIll,
		Fibrosis:                     params.Fibrosis,
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
	}

	if err := formController.formService.UpdateLungCancer(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertLungCancerResponse{})
}

func (formController *AdminFormController) UpdateNavidForm(ctx *gin.Context) {
	type UpdateNavidFormParams struct {
		FormID uint `uri:"formID" validate:"required"`

		InsuranceStatus              *string      `json:"insuranceStatus"`
		SupplementaryInsuranceStatus *enum.Answer `json:"takmilBime"`
		SupplementaryInsurances      *string      `json:"supplementaryInsurances"`
		Hypertension                 *string      `json:"hypertension"`
		HypertensionTreatment        *enum.Answer `json:"hypertensionTreatment"`
		HeartDisease                 *string      `json:"heartDisease"`
		HeartDiseaseTreatment        *enum.Answer `json:"heartDiseaseTreatment"`
		Diabetes                     *string      `json:"diabetes"`
		DiabetesTreatment            *enum.Answer `json:"diabetesTreatment"`
		ChronicLungDisease           *enum.Answer `json:"chronicLungDisease"`
		ChronicLungDiseaseType       *string      `json:"chronicLungDiseaseType"`
		LungCancerHistory            *enum.Answer `json:"lungCancerHistory"`
		OtherCancerHistory           *enum.Answer `json:"otherCancerHistory"`
		OtherCancerType              *uint        `json:"otherCancerType"`
		LungCancerFamily             *enum.Answer `json:"lungCancerFamily"`
		LungCancerFamilyRelation     *string      `json:"lungCancerFamilyRelation"`
		OtherCancerFamily            *enum.Answer `json:"otherCancerFamily"`
		OtherCancerFamilyType        *uint        `json:"otherCancerFamilyType"`
		OtherCancerFamilyRelation    *string      `json:"otherCancerFamilyRelation"`
		OccupationalExposure         *string      `json:"occupationalExposure"`
		CurrentSmoking               *enum.Answer `json:"currentSmoking"`
		SmokingStartAgeCurrent       *uint        `json:"smokingStartAgeCurrent"`
		SmokingTypesCurrent          *string      `json:"smokingTypesCurrent"`
		CigarettesPerDayCurrent      *uint        `json:"cigarettesPerDayCurrent"`
		CigarPerDayCurrent           *uint        `json:"cigarPerDayCurrent"`
		ECigPerDayCurrent            *uint        `json:"eCigPerDayCurrent"`
		PipePerDayCurrent            *uint        `json:"pipePerDayCurrent"`
		ChapoghPerDayCurrent         *uint        `json:"chapoghPerDayCurrent"`
		SmokedOpiumPerDayCurrent     *uint        `json:"smokedOpiumPerDayCurrent"`
		ChewedOpiumPerDayCurrent     *uint        `json:"chewedOpiumPerDayCurrent"`
		HookahPerWeekCurrent         *uint        `json:"hookahPerWeekCurrent"`
		PastSmoking                  *enum.Answer `json:"pastSmoking"`
		LeaveSmoke                   *uint        `json:"leaveSmoke"`
		SmokingStartAgePast          *uint        `json:"smokingStartAgePast"`
		SmokingTypesPast             *string      `json:"smokingTypesPast"`
		CigarettesPerDayPast         *uint        `json:"cigarettesPerDayPast"`
		CigarPerDayPast              *uint        `json:"cigarPerDayPast"`
		ECigPerDayPast               *uint        `json:"eCigPerDayPast"`
		PipePerDayPast               *uint        `json:"pipePerDayPast"`
		ChapoghPerDayPast            *uint        `json:"chapoghPerDayPast"`
		SmokedOpiumPerDayPast        *uint        `json:"smokedOpiumPerDayPast"`
		ChewedOpiumPerDayPast        *uint        `json:"chewedOpiumPerDayPast"`
		HookahPerWeekPast            *uint        `json:"hookahPerWeekPast"`
		SecondhandSmoke              *enum.Answer `json:"secondhandSmoke"`
		SecondhandSmokeLocation      *string      `json:"secondhandSmokeLocation"`
		AttentionCorrect             *bool        `json:"attentionCorrect"`
		LungDiseaseHistory           *string      `json:"lungDiseaseHistory"`
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

	params := controller.Validate[UpdateNavidFormParams](ctx)

	userID, _ := ctx.Get(formController.constants.Context.ID)

	req := formdto.UpdateNavidFormRequest{
		UserID:                       userID.(uint),
		FormID:                       params.FormID,
		InsuranceStatus:              params.InsuranceStatus,
		SupplementaryInsuranceStatus: params.SupplementaryInsuranceStatus,
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

	if err := formController.formService.UpdateNavidForm(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateForm")
	controller.Response(ctx, 201, message, formdto.UpsertNavidFormResponse{})
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
func (formController *AdminFormController) GetAllCancers(ctx *gin.Context) {
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

func (formController *AdminFormController) CreateCancer(ctx *gin.Context) {
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

func (formController *AdminFormController) UpdateCancer(ctx *gin.Context) {
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

func (formController *AdminFormController) DeleteCancer(ctx *gin.Context) {
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

func (formController *AdminFormController) CreateFamilyCancer(ctx *gin.Context) {
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

func (formController *AdminFormController) UpdateFamilyCancer(ctx *gin.Context) {
	type UpdateFamilyCancerParams struct {
		FormID           uint                    `uri:"formID" validate:"required"`
		FamilyCancerID   uint                    `uri:"familyCancerID" validate:"required"`
		Relative         uint                    `form:"relative" validate:"required,gt=0"`
		RelativeRelation *string                 `form:"relativeRelation,omitempty"`
		Name             *string                 `form:"name,omitempty"`
		NumberRelative   *uint                   `form:"numberRelative,omitempty"`
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
		Name:             params.Name,
		NumberRelative:   params.NumberRelative,
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

func (formController *AdminFormController) DeleteFamilyCancer(ctx *gin.Context) {
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

	trans := controller.GetTranslator(ctx, formController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteForm")
	controller.Response(ctx, 200, message, nil)
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

func (formController *AdminFormController) GetNavidForm(ctx *gin.Context) {
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
		UserID               uint              `json:"userId" validate:"required"`
		BirthDate            formdto.BirthDate `json:"birthDate" validate:"required" time_format:"2006-01-02"`
		SocialSecurityNumber string            `json:"socialSecurityNumber" validate:"required"`
		Gender               uint              `json:"gender" validate:"required"`
		IsAtba               bool              `json:"isAtba"`
		Height               float64           `json:"height" validate:"required"`
		Weight               float64           `json:"weight" validate:"required"`
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
