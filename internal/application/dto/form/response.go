package formdto

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
)

type BasicFormResponse struct {
	FormID               uint                 `json:"id"`
	Status               string               `json:"status"`
	FormType             enum.FormType        `json:"formType"`
	UserID               uint                 `json:"user_id"`
	OperatorID           *uint                `json:"operatorId,omitempty"`
	FilledByOperatorID   *uint                `json:"filledByOperatorId,omitempty"`
	SocialSecurityNumber string               `json:"socialSecurityNumber"`
	Name                 *string              `json:"name,omitempty"`
	FilledForms          *FilledFormsResponse `json:"filledForms,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AttentionQuestionsCorrect *int `json:"attentionQuestionsCorrect"`
}

type FilledFormsResponse struct {
	Basic         *bool `json:"basic"`
	GeneralHealth *bool `json:"generalhealth"`
	Mamography    *bool `json:"mamography"`
	Cancer        *bool `json:"cancer"`
	FamilyCancer  *bool `json:"familycancer"`
	Contact       *bool `json:"contact"`
	LungCancer    *bool `json:"lungcancer"`
	NavidForm     *bool `json:"navid"`
}

type CreateFormResponse struct {
	Form BasicFormResponse `json:"form"`
}

type UpdateFormResponse struct {
	Form BasicFormResponse `json:"form"`
}

type GetUserFormsResponse struct {
	Forms []BasicFormResponse `json:"forms"`
	Total int                 `json:"total"`
}

type DeleteFormResponse struct {
	Message string `json:"message"`
}

type UpsertGeneralHealthResponse struct {
	Form BasicFormResponse `json:"form"`
}

type UpsertMamographyResponse struct {
	Form BasicFormResponse `json:"form"`
}

type CreateCancerResponse struct {
	Cancer CancerResponse `json:"cancer"`
}

type UpdateCancerResponse struct {
	Cancer CancerResponse `json:"cancer"`
}

type DeleteCancerResponse struct {
	Message string `json:"message"`
}

type FamilyCancerItemResponse struct {
	ID               uint             `json:"id"`
	Relative         enum.Relative    `json:"relative"`
	RelativeRelation *string          `json:"relation,omitempty"`
	NumberRelative   *uint            `json:"numberRelative,omitempty"`
	Name             *string          `json:"name,omitempty"`
	LifeStatus       *enum.LifeStatus `json:"lifeStatus,omitempty"`
	CancerType       enum.CancerType  `json:"cancerType"`
	CancerAge        uint             `json:"cancerAge"`
	Pictures         []string         `json:"pictures,omitempty"`
}

type CreateFamilyCancerResponse struct {
	FamilyCancer FamilyCancerItemResponse `json:"familyCancer"`
}

type UpdateFamilyCancerResponse struct {
	FamilyCancer FamilyCancerItemResponse `json:"familyCancer"`
}

type DeleteFamilyCancerResponse struct {
	Message string `json:"message"`
}

type UpsertFamilyCancerResponse struct {
	Form BasicFormResponse `json:"form"`
}

type UpsertContactResponse struct {
	Form BasicFormResponse `json:"form"`
}

type UpsertLungCancerResponse struct {
	Form BasicFormResponse `json:"form"`
}

type UpsertNavidFormResponse struct {
	Form BasicFormResponse `json:"form"`
}

type GetBasicFormResponse struct {
	ID                   uint          `json:"id"`
	FormType             enum.FormType `json:"formType"`
	Gender               enum.Gender   `json:"gender"`
	BirthDate            BirthDate     `json:"birthDate"`
	IsAtba               bool          `json:"isAtba"`
	SocialSecurityNumber string        `json:"socialSecurityNumber"`
	Height               float64       `json:"height"`
	Weight               float64       `json:"weight"`
}

type GetGeneralHealthResponse struct {
	ID                        uint         `json:"id"`
	DrinksAlcohol             *enum.Answer `json:"drinksAlcohol"`
	CupsPerWeek               *string      `json:"cupsPerWeek"`
	LastMonthSabzijatMeal     string       `json:"lastMonthSabzijatMeal"`
	LastMonthSabzijatWeight   string       `json:"lastMonthSabzijatWeight"`
	MediumActivityMonthInYear uint         `json:"mediumActivityMonthInYear"`
	MediumActivityHourInWeek  string       `json:"mediumActivityHourInWeek"`
	HardActivityMonthInYear   uint         `json:"hardActivityMonthInYear"`
	HardActivityHourInWeek    string       `json:"hardActivityHourInWeek"`
	SmokeAtLeast100           *enum.Answer `json:"smokeAtLeast100"`
	SmokingAge                *uint        `json:"smokingAge"`
	SmokingNow                *enum.Answer `json:"smokingNow"`
	YearSmoke                 *uint        `json:"yearSmoke"`
	LeaveSmokingAge           *uint        `json:"leaveSmokingAge"`
	CountSmokingDaily         *string      `json:"countSmokingDaily"`
	CountGheliandaily         *string      `json:"countGheliandaily"`
	CountSmokingDailyPast     *string      `json:"countSmokingDailyPast"`
	CountGheliandailyPast     *string      `json:"countGheliandailyPast"`
}

type CancerResponse struct {
	ID         uint            `json:"id"`
	CancerType enum.CancerType `json:"cancerType"`
	CancerAge  uint            `json:"cancerAge"`
	Pictures   []string        `json:"pictures,omitempty"`
}

type GetCancersResponse struct {
	Cancer  bool             `json:"cancer"`
	Cancers []CancerResponse `json:"cancers"`
}

type FamilyCancerResponse struct {
	Relative         enum.Relative    `json:"relative"`
	RelativeRelation *string          `json:"relation,omitempty"`
	Name             *string          `json:"name,omitempty"`
	NumberRelative   *uint            `json:"numberRelative,omitempty"`
	LifeStatus       *enum.LifeStatus `json:"lifeStatus,omitempty"`
	Cancers          []CancerResponse `json:"cancers,omitempty"`
}

type GetFamilyCancerResponse struct {
	FamilyCancers []FamilyCancerResponse `json:"familyCancers"`
}
type GetFamilyCancerListResponse struct {
	AmeAmoCancer    bool `json:"ameAmoCancer"`
	FatherCancer    bool `json:"fatherCancer"`
	KhaleDaeiCancer bool `json:"khaleDaeiCancer"`
	MotherCancer    bool `json:"motherCancer"`
	SiblingCancer   bool `json:"siblingCancer"`
}

type GetMamographyResponse struct {
	ID                           uint                  `json:"id"`
	GhaedeAge                    uint                  `json:"ghaedeAge"`
	HasChildren                  *enum.Answer          `json:"hasChildren"`
	NumberOfChildren             *uint                 `json:"numberOfChildren,omitempty"`
	DaughterCount                *uint                 `json:"daughterCount,omitempty"`
	SonCount                     *uint                 `json:"sonCount,omitempty"`
	AgeOfFirstBirth              *uint                 `json:"ageOfFirstBirth,omitempty"`
	MenopausalStatus             enum.MenopausalStatus `json:"menopausalStatus"`
	MenopauseAge                 *string               `json:"menopauseAge,omitempty"`
	HRT                          *enum.Answer          `json:"hrt"`
	HRTUseLength                 *uint                 `json:"hrtUseLength,omitempty"`
	LastFiveYearsHRTUse          *enum.Answer          `json:"lastFiveYearsHrtUse"`
	CurrentHRTUse                *enum.Answer          `json:"currentHrtUse"`
	IntendedHRTUse               *uint                 `json:"intendedHrtUse,omitempty"`
	HRTType                      *string               `json:"hrtType,omitempty"`
	Oral                         *enum.Answer          `json:"oral"`
	OralDuration                 *string               `json:"oralDuration,omitempty"`
	OralTwoLastYears             *enum.Answer          `json:"oralTwoLastYears"`
	MamoGraphy                   *enum.Answer          `json:"mamoGraphy"`
	MamoGraphyPictures           []string              `json:"mamoGraphyPictures,omitempty"`
	BreastDensity                string                `json:"breastDensity"`
	Falop                        *enum.Answer          `json:"falop"`
	Andometrioz                  *enum.Answer          `json:"andometrioz"`
	LeavePestan                  bool                  `json:"leavePestan"`
	LeaveTokhmdan                bool                  `json:"leaveTokhmdan"`
	LaDeColon                    *enum.Answer          `json:"laDeColon"`
	LaDePol                      *enum.Answer          `json:"laDePol"`
	AspLaMo                      *enum.Answer          `json:"aspLaMo"`
	NsaiDLaMo                    *enum.Answer          `json:"nsaiDLaMo"`
	LastFiveYearBloodTestInStool *enum.Answer          `json:"lastFiveYearBloodTestInStool"`
	NumberOfBreastBiopsies       *uint                 `json:"numberOfBreastBiopsies,omitempty"`
	HyperplasiaInBiopsy          *uint                 `json:"hyperplasiaInBiopsy,omitempty"`
}

type ChangeFormStatusResponse struct {
	Form BasicFormResponse `json:"form"`
}

type GetContactResponse struct {
	ID                    uint         `json:"id"`
	TestGen               *enum.Answer `json:"testGen"`
	TestGenPictures       []string     `json:"testGenPictures,omitempty"`
	Name                  string       `json:"name"`
	FmTestGen             *enum.Answer `json:"fmTestGen"`
	FatherTestGenPictures []string     `json:"fatherTestGenPictures,omitempty"`
	MotherTestGenPictures []string     `json:"motherTestGenPictures,omitempty"`
	CallExpert            bool         `json:"callExpert"`
	BirthCountry          *string      `json:"birthCountry,omitempty"`
	Province              *string      `json:"province,omitempty"`
	City                  *string      `json:"city,omitempty"`
	Country               *string      `json:"country,omitempty"`
	Address               string       `json:"address"`
	PostalCode            string       `json:"postalCode"`
	Education             string       `json:"education"`
	Phone2                *string      `json:"phone2,omitempty"`
	Phone3                *string      `json:"phone3,omitempty"`
	BrotherNumber         *uint        `json:"brotherNumber"`
	SisterNumber          *uint        `json:"sisterNumber"`
	PaternalAuntNumber    *uint        `json:"paternalAuntNumber"`
	PaternalUncleNumber   *uint        `json:"paternalUncleNumber"`
	MaternalAuntNumber    *uint        `json:"maternalAuntNumber"`
	MaternalUncleNumber   *uint        `json:"maternalUncleNumber"`
}

type GetLungCancerResponse struct {
	ID                           uint             `json:"id"`
	SupplementaryInsuranceStatus *enum.Answer     `json:"takmilBime"`
	InsuranceStatus              *string          `json:"insuranceStatus,omitempty"`
	SupplementaryInsurances      *string          `json:"supplementaryInsurances,omitempty"`
	ChronicLungDisease           *enum.Answer     `json:"chronicLungDisease"`
	ChronicLungDiseaseType       *string          `json:"chronicLungDiseaseType,omitempty"`
	LungCancerHistory            *enum.Answer     `json:"lungCancerHistory"`
	OtherCancerHistory           *enum.Answer     `json:"otherCancerHistory"`
	OtherCancerType              *enum.CancerType `json:"otherCancerType,omitempty"`
	LungCancerFamily             *enum.Answer     `json:"lungCancerFamily"`
	LungCancerFamilyRelation     *string          `json:"lungCancerFamilyRelation,omitempty"`
	OtherCancerFamily            *enum.Answer     `json:"otherCancerFamily"`
	OtherCancerFamilyType        *enum.CancerType `json:"otherCancerFamilyType,omitempty"`
	OtherCancerFamilyRelation    *string          `json:"otherCancerFamilyRelation,omitempty"`
	OccupationalExposure         *string          `json:"occupationalExposure,omitempty"`
	CurrentSmoking               *enum.Answer     `json:"currentSmoking"`
	SmokingStartAgeCurrent       *uint            `json:"smokingStartAgeCurrent,omitempty"`
	SmokingTypesCurrent          *string          `json:"smokingTypesCurrent,omitempty"`
	CigarettesPerDayCurrent      *uint            `json:"cigarettesPerDayCurrent,omitempty"`
	CigarPerDayCurrent           *uint            `json:"cigarPerDayCurrent,omitempty"`
	ECigPerDayCurrent            *uint            `json:"eCigPerDayCurrent,omitempty"`
	PipePerDayCurrent            *uint            `json:"pipePerDayCurrent,omitempty"`
	ChapoghPerDayCurrent         *uint            `json:"chapoghPerDayCurrent,omitempty"`
	SmokedOpiumPerDayCurrent     *uint            `json:"smokedOpiumPerDayCurrent,omitempty"`
	ChewedOpiumPerDayCurrent     *uint            `json:"chewedOpiumPerDayCurrent,omitempty"`
	HookahPerWeekCurrent         *uint            `json:"hookahPerWeekCurrent,omitempty"`
	PastSmoking                  *enum.Answer     `json:"pastSmoking,omitempty"`
	SmokePastAvg                 *uint            `json:"smokePastAvg"`
	SmokeCurrentAvg              *uint            `json:"smokeCurrentAvg"`
	Bronchitis                   *enum.Answer     `json:"bronshit"`
	LungIll                      *enum.Answer     `json:"fibroz"`
	Fibrosis                     *enum.Answer     `json:"lungill"`
	LeaveSmoke                   *uint            `json:"leaveSmoke,omitempty"`
	SmokingStartAgePast          *uint            `json:"smokingStartAgePast,omitempty"`
	SmokingTypesPast             *string          `json:"smokingTypesPast,omitempty"`
	CigarettesPerDayPast         *uint            `json:"cigarettesPerDayPast,omitempty"`
	CigarPerDayPast              *uint            `json:"cigarPerDayPast,omitempty"`
	ECigPerDayPast               *uint            `json:"eCigPerDayPast,omitempty"`
	PipePerDayPast               *uint            `json:"pipePerDayPast,omitempty"`
	ChapoghPerDayPast            *uint            `json:"chapoghPerDayPast,omitempty"`
	SmokedOpiumPerDayPast        *uint            `json:"smokedOpiumPerDayPast,omitempty"`
	ChewedOpiumPerDayPast        *uint            `json:"chewedOpiumPerDayPast,omitempty"`
	HookahPerWeekPast            *uint            `json:"hookahPerWeekPast,omitempty"`
	SecondhandSmoke              *enum.Answer     `json:"secondhandSmoke"`
	SecondhandSmokeLocation      *string          `json:"secondhandSmokeLocation,omitempty"`
	LungDiseaseHistory           *string          `json:"lungDiseaseHistory,omitempty"`
}

type PostalCodeInfoResponse struct {
	Province     *string `json:"province,omitempty"`
	City         *string `json:"city,omitempty"`
	Town         *string `json:"town,omitempty"`
	District     *string `json:"district,omitempty"`
	Street       *string `json:"street,omitempty"`
	Street2      *string `json:"street2,omitempty"`
	Number       *string `json:"number,omitempty"`
	Floor        *string `json:"floor,omitempty"`
	SideFloor    *string `json:"sideFloor,omitempty"`
	BuildingName *string `json:"buildingName,omitempty"`
	Description  *string `json:"description,omitempty"`
}

type GetNavidFormResponse struct {
	ID                           uint             `json:"id"`
	SupplementaryInsuranceStatus *enum.Answer     `json:"takmilBime"`
	InsuranceStatus              *string          `json:"insuranceStatus,omitempty"`
	SupplementaryInsurances      *string          `json:"supplementaryInsurances,omitempty"`
	Hypertension                 string           `json:"hypertension"`
	HypertensionTreatment        *enum.Answer     `json:"hypertensionTreatment"`
	HeartDisease                 string           `json:"heartDisease"`
	HeartDiseaseTreatment        *enum.Answer     `json:"heartDiseaseTreatment"`
	Diabetes                     string           `json:"diabetes"`
	DiabetesTreatment            *enum.Answer     `json:"diabetesTreatment"`
	ChronicLungDisease           *enum.Answer     `json:"chronicLungDisease"`
	ChronicLungDiseaseType       *string          `json:"chronicLungDiseaseType,omitempty"`
	LungCancerHistory            *enum.Answer     `json:"lungCancerHistory"`
	OtherCancerHistory           *enum.Answer     `json:"otherCancerHistory"`
	OtherCancerType              *enum.CancerType `json:"otherCancerType,omitempty"`
	LungCancerFamily             *enum.Answer     `json:"lungCancerFamily"`
	LungCancerFamilyRelation     *string          `json:"lungCancerFamilyRelation,omitempty"`
	OtherCancerFamily            *enum.Answer     `json:"otherCancerFamily"`
	OtherCancerFamilyType        *enum.CancerType `json:"otherCancerFamilyType,omitempty"`
	OtherCancerFamilyRelation    *string          `json:"otherCancerFamilyRelation,omitempty"`
	OccupationalExposure         *string          `json:"occupationalExposure,omitempty"`
	CurrentSmoking               *enum.Answer     `json:"currentSmoking"`
	SmokingStartAgeCurrent       *uint            `json:"smokingStartAgeCurrent,omitempty"`
	SmokingTypesCurrent          *string          `json:"smokingTypesCurrent,omitempty"`
	CigarettesPerDayCurrent      *uint            `json:"cigarettesPerDayCurrent,omitempty"`
	CigarPerDayCurrent           *uint            `json:"cigarPerDayCurrent,omitempty"`
	ECigPerDayCurrent            *uint            `json:"eCigPerDayCurrent,omitempty"`
	PipePerDayCurrent            *uint            `json:"pipePerDayCurrent,omitempty"`
	ChapoghPerDayCurrent         *uint            `json:"chapoghPerDayCurrent,omitempty"`
	SmokedOpiumPerDayCurrent     *uint            `json:"smokedOpiumPerDayCurrent,omitempty"`
	ChewedOpiumPerDayCurrent     *uint            `json:"chewedOpiumPerDayCurrent,omitempty"`
	HookahPerWeekCurrent         *uint            `json:"hookahPerWeekCurrent,omitempty"`
	PastSmoking                  *enum.Answer     `json:"pastSmoking,omitempty"`
	LeaveSmoke                   *uint            `json:"leaveSmoke,omitempty"`
	SmokingStartAgePast          *uint            `json:"smokingStartAgePast,omitempty"`
	SmokingTypesPast             *string          `json:"smokingTypesPast,omitempty"`
	CigarettesPerDayPast         *uint            `json:"cigarettesPerDayPast,omitempty"`
	CigarPerDayPast              *uint            `json:"cigarPerDayPast,omitempty"`
	ECigPerDayPast               *uint            `json:"eCigPerDayPast,omitempty"`
	PipePerDayPast               *uint            `json:"pipePerDayPast,omitempty"`
	ChapoghPerDayPast            *uint            `json:"chapoghPerDayPast,omitempty"`
	SmokedOpiumPerDayPast        *uint            `json:"smokedOpiumPerDayPast,omitempty"`
	ChewedOpiumPerDayPast        *uint            `json:"chewedOpiumPerDayPast,omitempty"`
	HookahPerWeekPast            *uint            `json:"hookahPerWeekPast,omitempty"`
	SecondhandSmoke              *enum.Answer     `json:"secondhandSmoke"`
	SecondhandSmokeLocation      *string          `json:"secondhandSmokeLocation,omitempty"`
	LungDiseaseHistory           *string          `json:"lungDiseaseHistory,omitempty"`
	CurrentCigaretteSmoking      *enum.Answer     `json:"Csig"`
	CurrentRolledTobacco         *enum.Answer     `json:"CsigBarg"`
	CurrentPipeSmoking           *enum.Answer     `json:"Cpip"`
	CurrentHookahUse             *enum.Answer     `json:"Cghel"`
	CurrentChiboukSmoking        *enum.Answer     `json:"Cchop"`
	CurrentOpiumUse              *enum.Answer     `json:"Cteryak"`
	FormerCigaretteSmoking       *enum.Answer     `json:"Psig"`
	FormerRolledTobacco          *enum.Answer     `json:"PsigBarg"`
	FormerPipeSmoking            *enum.Answer     `json:"Ppip"`
	FormerHookahUse              *enum.Answer     `json:"Pghel"`
	FormerChiboukSmoking         *enum.Answer     `json:"Pchop"`
	FormerOpiumUse               *enum.Answer     `json:"Pteryak"`
	PelecSig                     *enum.Answer     `json:"PelecSig"`
	CelecSig                     *enum.Answer     `json:"CelecSig"`
}
