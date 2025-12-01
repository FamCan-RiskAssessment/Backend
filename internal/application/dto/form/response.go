package formdto

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
)

type BasicFormResponse struct {
	FormID               uint    `json:"id"`
	Status               string  `json:"status"`
	UserID               uint    `json:"user_id"`
	OperatorID           *uint   `json:"operatorId,omitempty"`
	FilledByOperatorID   *uint   `json:"filledByOperatorId,omitempty"`
	SocialSecurityNumber string  `json:"socialSecurityNumber"`
	Name                 *string `json:"name,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AttentionQuestionsCorrect *int `json:"attentionQuestionsCorrect"`
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

type GetBasicFormResponse struct {
	ID                   uint        `json:"id"`
	Gender               enum.Gender `json:"gender"`
	BirthDate            time.Time   `json:"birthDate"`
	IsAtba               bool        `json:"isAtba"`
	SocialSecurityNumber string      `json:"socialSecurityNumber"`
	Height               float64     `json:"height"`
	Weight               float64     `json:"weight"`
}

type GetGeneralHealthResponse struct {
	ID                        uint    `json:"id"`
	DrinksAlcohol             *bool   `json:"drinksAlcohol,omitempty"`
	CupsPerWeek               *string `json:"cupsPerWeek,omitempty"`
	LastMonthSabzijatMeal     string  `json:"lastMonthSabzijatMeal"`
	LastMonthSabzijatWeight   string  `json:"lastMonthSabzijatWeight"`
	MediumActivityMonthInYear uint    `json:"mediumActivityMonthInYear"`
	MediumActivityHourInWeek  string  `json:"mediumActivityHourInWeek"`
	HardActivityMonthInYear   uint    `json:"hardActivityMonthInYear"`
	HardActivityHourInWeek    string  `json:"hardActivityHourInWeek"`
	SmokeAtLeast100           *bool   `json:"smokeAtLeast100,omitempty"`
	SmokingAge                *uint   `json:"smokingAge,omitempty"`
	SmokingNow                bool    `json:"smokingNow"`
	LeaveSmokingAge           *uint   `json:"leaveSmokingAge,omitempty"`
	CountSmokingDaily         *string `json:"countSmokingDaily,omitempty"`
	CountGheliandaily         *string `json:"countGheliandaily,omitempty"`
	CountSmokingDailyPast     *string `json:"countSmokingDailyPast,omitempty"`
	CountGheliandailyPast     *string `json:"countGheliandailyPast,omitempty"`
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
	LifeStatus       *enum.LifeStatus `json:"lifeStatus,omitempty"`
	Cancers          []CancerResponse `json:"cancers,omitempty"`
}

type GetFamilyCancerResponse struct {
	FamilyCancers []FamilyCancerResponse `json:"familyCancers"`
}

type GetMamographyResponse struct {
	ID                           uint                  `json:"id"`
	GhaedeAge                    uint                  `json:"ghaedeAge"`
	HasChildren                  bool                  `json:"hasChildren"`
	NumberOfChildren             *uint                 `json:"numberOfChildren,omitempty"`
	AgeOfFirstBirth              *uint                 `json:"ageOfFirstBirth,omitempty"`
	MenopausalStatus             enum.MenopausalStatus `json:"menopausalStatus"`
	MenopauseAge                 *string               `json:"menopauseAge,omitempty"`
	HRT                          *bool                 `json:"hrt,omitempty"`
	HRTUseLength                 *uint                 `json:"hrtUseLength,omitempty"`
	LastFiveYearsHRTUse          bool                  `json:"lastFiveYearsHrtUse"`
	CurrentHRTUse                *bool                 `json:"currentHrtUse,omitempty"`
	IntendedHRTUse               *uint                 `json:"intendedHrtUse,omitempty"`
	HRTType                      *string               `json:"hrtType,omitempty"`
	Oral                         *bool                 `json:"oral,omitempty"`
	OralDuration                 *string               `json:"oralDuration,omitempty"`
	OralTwoLastYears             *bool                 `json:"oralTwoLastYears,omitempty"`
	MamoGraphy                   *bool                 `json:"mamoGraphy,omitempty"`
	MamoGraphyPictures           []string              `json:"mamoGraphyPictures,omitempty"`
	Falop                        *bool                 `json:"falop,omitempty"`
	Andometrioz                  *bool                 `json:"andometrioz,omitempty"`
	LeavePestan                  bool                  `json:"leavePestan"`
	LeaveTokhmdan                bool                  `json:"leaveTokhmdan"`
	LaDeColon                    *bool                 `json:"laDeColon,omitempty"`
	LaDePol                      *bool                 `json:"laDePol,omitempty"`
	AspLaMo                      *bool                 `json:"aspLaMo,omitempty"`
	NsaiDLaMo                    *bool                 `json:"nsaiDLaMo,omitempty"`
	LastFiveYearBloodTestInStool *bool                 `json:"lastFiveYearBloodTestInStool,omitempty"`
	NumberOfBreastBiopsies       *uint                 `json:"numberOfBreastBiopsies,omitempty"`
	HyperplasiaInBiopsy          *uint                 `json:"hyperplasiaInBiopsy,omitempty"`
}

type ChangeFormStatusResponse struct {
	Form BasicFormResponse `json:"form"`
}

type GetContactResponse struct {
	ID                    uint     `json:"id"`
	TestGen               *bool    `json:"testGen,omitempty"`
	TestGenPictures       []string `json:"testGenPictures,omitempty"`
	Name                  string   `json:"name"`
	FmTestGen             *bool    `json:"fmTestGen,omitempty"`
	FatherTestGenPictures []string `json:"fatherTestGenPictures,omitempty"`
	MotherTestGenPictures []string `json:"motherTestGenPictures,omitempty"`
	CallExpert            bool     `json:"callExpert"`
	BirthCountry          *string  `json:"birthCountry,omitempty"`
	Province              *string  `json:"province,omitempty"`
	City                  *string  `json:"city,omitempty"`
	Country               *string  `json:"country,omitempty"`
	Address               string   `json:"address"`
	PostalCode            string   `json:"postalCode"`
	Education             string   `json:"education"`
	Phone2                *string  `json:"phone2,omitempty"`
	Phone3                *string  `json:"phone3,omitempty"`
}

type GetLungCancerResponse struct {
	ID                           uint             `json:"id"`
	SupplementaryInsuranceStatus *bool            `json:"takmilBime,omitempty"`
	InsuranceStatus              *string          `json:"insuranceStatus,omitempty"`
	SupplementaryInsurances      *string          `json:"supplementaryInsurances,omitempty"`
	Hypertension                 bool             `json:"hypertension"`
	HypertensionTreatment        *bool            `json:"hypertensionTreatment,omitempty"`
	HeartDisease                 bool             `json:"heartDisease"`
	HeartDiseaseTreatment        *bool            `json:"heartDiseaseTreatment,omitempty"`
	Diabetes                     bool             `json:"diabetes"`
	DiabetesTreatment            *bool            `json:"diabetesTreatment,omitempty"`
	ChronicLungDisease           *bool            `json:"chronicLungDisease,omitempty"`
	ChronicLungDiseaseType       *string          `json:"chronicLungDiseaseType,omitempty"`
	LungCancerHistory            bool             `json:"lungCancerHistory"`
	OtherCancerHistory           bool             `json:"otherCancerHistory"`
	OtherCancerType              *enum.CancerType `json:"otherCancerType,omitempty"`
	LungCancerFamily             *bool            `json:"lungCancerFamily,omitempty"`
	LungCancerFamilyRelation     *string          `json:"lungCancerFamilyRelation,omitempty"`
	OtherCancerFamily            *bool            `json:"otherCancerFamily,omitempty"`
	OtherCancerFamilyType        *enum.CancerType `json:"otherCancerFamilyType,omitempty"`
	OtherCancerFamilyRelation    *string          `json:"otherCancerFamilyRelation,omitempty"`
	OccupationalExposure         *string          `json:"occupationalExposure,omitempty"`
	CurrentSmoking               bool             `json:"currentSmoking"`
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
	PastSmoking                  *string          `json:"pastSmoking,omitempty"`
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
	SecondhandSmoke              bool             `json:"secondhandSmoke"`
	SecondhandSmokeLocation      *string          `json:"secondhandSmokeLocation,omitempty"`
	LungDiseaseHistory           *string          `json:"lungDiseaseHistory,omitempty"`
}
