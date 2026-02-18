package formdto

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/jalaali/go-jalaali"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
)

type BirthDate jalaali.Jalaali

func (d *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = BirthDate(jalaali.From(t))
	return nil
}

func (d BirthDate) String() string {
	return d.Time.Format("2006-01-02")
}

type CreateBasicFormRequest struct {
	UserID             uint
	FilledByOperatorID *uint

	// page 1
	BirthDate            BirthDate
	SocialSecurityNumber string
	Gender               uint
	IsAtba               bool
	Height               float64
	Weight               float64
	FormType             *enum.FormType
}

type UpdateBasicFormRequest struct {
	UserID uint
	FormID uint

	// page 1
	BirthDate            *BirthDate
	SocialSecurityNumber *string
	Gender               *uint
	IsAtba               *bool
	Height               *float64
	Weight               *float64
}

type GetUserFormsRequest struct {
	UserID uint
	Offset int
	Limit  int
}

type GetPartialFormRequest struct {
	UserID uint
	FormID uint
}

type DeleteFormRequest struct {
	UserID uint
	FormID uint
}

type UpsertGeneralHealthRequest struct {
	UserID uint
	FormID uint

	DrinksAlcohol             *enum.Answer
	CupsPerWeek               *string
	LastMonthSabzijatMeal     string
	LastMonthSabzijatWeight   string
	MediumActivityMonthInYear uint
	MediumActivityHourInWeek  string
	HardActivityMonthInYear   uint
	HardActivityHourInWeek    string
	SmokeAtLeast100           *enum.Answer
	SmokingAge                *uint
	SmokingNow                *enum.Answer
	YearSmoke                 *uint
	LeaveSmokingAge           *uint
	CountSmokingDaily         *string
	CountGheliandaily         *string
	CountSmokingDailyPast     *string
	CountGheliandailyPast     *string

	// Attention question answer
	AttentionCorrect *bool
}

type UpsertMamographyRequest struct {
	UserID uint
	FormID uint

	GhaedeAge                    uint
	HasChildren                  *enum.Answer
	NumberOfChildren             *uint
	DaughterCount                *uint
	SonCount                     *uint
	AgeOfFirstBirth              *uint
	MenopausalStatus             uint
	MenopauseAge                 *string
	HRT                          *enum.Answer
	HRTUseLength                 *uint
	LastFiveYearsHRTUse          *enum.Answer
	CurrentHRTUse                *enum.Answer
	IntendedHRTUse               *uint
	HRTType                      *string
	Oral                         *enum.Answer
	OralDuration                 *string
	OralTwoLastYears             *enum.Answer
	MamoGraphy                   *enum.Answer
	MamoGraphyPictures           []*multipart.FileHeader
	BreastDensity                *string
	Falop                        *enum.Answer
	Andometrioz                  *enum.Answer
	LeavePestan                  bool
	LeaveTokhmdan                bool
	LaDeColon                    *enum.Answer
	LaDePol                      *enum.Answer
	AspLaMo                      *enum.Answer
	NsaiDLaMo                    *enum.Answer
	LastFiveYearBloodTestInStool *enum.Answer
	NumberOfBreastBiopsies       *uint
	HyperplasiaInBiopsy          *uint

	// Attention question answer
	AttentionCorrect *bool
}

type CancerRequest struct {
	CancerType uint
	CancerAge  uint
	Pictures   []*multipart.FileHeader
}

type CreateCancerRequest struct {
	UserID uint
	FormID uint

	CancerType uint
	CancerAge  uint
	Pictures   []*multipart.FileHeader
}

type UpdateCancerRequest struct {
	UserID   uint
	FormID   uint
	CancerID uint

	CancerType uint
	CancerAge  uint
	Pictures   []*multipart.FileHeader
}

type DeleteCancerRequest struct {
	UserID   uint
	FormID   uint
	CancerID uint
}

type CreateFamilyCancerRequest struct {
	UserID uint
	FormID uint

	Relative         enum.Relative
	RelativeRelation *string
	Name             *string
	NumberRelative   *uint
	LifeStatus       *enum.LifeStatus
	CancerType       uint
	CancerAge        uint
	Pictures         []*multipart.FileHeader
}

type UpdateFamilyCancerRequest struct {
	UserID           uint
	FormID           uint
	FamilyCancerID   uint
	Relative         enum.Relative
	RelativeRelation *string
	Name             *string
	NumberRelative   *uint
	LifeStatus       *enum.LifeStatus
	CancerType       uint
	CancerAge        uint
	Pictures         []*multipart.FileHeader
}

type DeleteFamilyCancerRequest struct {
	UserID         uint
	FormID         uint
	FamilyCancerID uint
}

type FamilyCancerRequest struct {
	UserID uint
	FormID uint

	Relative         enum.Relative
	RelativeRelation *string
	Name             *string
	LifeStatus       *enum.LifeStatus
	Cancer           bool
	Cancers          []CancerRequest
}

type UpsertContactRequest struct {
	UserID uint
	FormID uint

	Name                  string
	TestGen               *enum.Answer
	TestGenPictures       []*multipart.FileHeader
	FmTestGen             *enum.Answer
	FatherTestGenPictures []*multipart.FileHeader
	MotherTestGenPictures []*multipart.FileHeader
	CallExpert            bool
	BirthCountry          *string
	Province              *string
	City                  *string
	Country               *string
	Address               string
	PostalCode            string
	Education             string
	Phone2                *string
	Phone3                *string
}

type UpsertNavidFormRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus              *string
	SupplementaryInsurances      *string
	SupplementaryInsuranceStatus *enum.Answer
	Hypertension                 string
	HypertensionTreatment        *enum.Answer
	HeartDisease                 string
	HeartDiseaseTreatment        *enum.Answer
	Diabetes                     string
	DiabetesTreatment            *enum.Answer
	ChronicLungDisease           *enum.Answer
	ChronicLungDiseaseType       *string
	LungCancerHistory            *enum.Answer
	OtherCancerHistory           *enum.Answer
	OtherCancerType              *uint
	LungCancerFamily             *enum.Answer
	LungCancerFamilyRelation     *string
	OtherCancerFamily            *enum.Answer
	OtherCancerFamilyType        *uint
	OtherCancerFamilyRelation    *string
	OccupationalExposure         *string
	CurrentSmoking               *enum.Answer
	SmokingStartAgeCurrent       *uint
	SmokingTypesCurrent          *string
	CigarettesPerDayCurrent      *uint
	CigarPerDayCurrent           *uint
	ECigPerDayCurrent            *uint
	PipePerDayCurrent            *uint
	ChapoghPerDayCurrent         *uint
	SmokedOpiumPerDayCurrent     *uint
	ChewedOpiumPerDayCurrent     *uint
	HookahPerWeekCurrent         *uint
	PastSmoking                  *string
	LeaveSmoke                   *uint
	SmokingStartAgePast          *uint
	SmokingTypesPast             *string
	CigarettesPerDayPast         *uint
	CigarPerDayPast              *uint
	ECigPerDayPast               *uint
	PipePerDayPast               *uint
	ChapoghPerDayPast            *uint
	SmokedOpiumPerDayPast        *uint
	ChewedOpiumPerDayPast        *uint
	HookahPerWeekPast            *uint
	SecondhandSmoke              *enum.Answer
	SecondhandSmokeLocation      *string
	LungDiseaseHistory           *string
	CurrentCigaretteSmoking      *enum.Answer
	CurrentRolledTobacco         *enum.Answer
	CurrentPipeSmoking           *enum.Answer
	CurrentHookahUse             *enum.Answer
	CurrentChiboukSmoking        *enum.Answer
	CurrentOpiumUse              *enum.Answer
	FormerCigaretteSmoking       *enum.Answer
	FormerRolledTobacco          *enum.Answer
	FormerPipeSmoking            *enum.Answer
	FormerHookahUse              *enum.Answer
	FormerChiboukSmoking         *enum.Answer
	FormerOpiumUse               *enum.Answer
	PelecSig                     *enum.Answer
	CelecSig                     *enum.Answer

	// Attention question answer
	AttentionCorrect *bool
}

type UpsertLungCancerRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus              *string
	SupplementaryInsurances      *string
	SupplementaryInsuranceStatus *enum.Answer
	Hypertension                 *enum.Answer
	HypertensionTreatment        *enum.Answer
	HeartDisease                 *enum.Answer
	HeartDiseaseTreatment        *enum.Answer
	Diabetes                     *enum.Answer
	DiabetesTreatment            *enum.Answer
	ChronicLungDisease           *enum.Answer
	ChronicLungDiseaseType       *string
	LungCancerHistory            *enum.Answer
	OtherCancerHistory           *enum.Answer
	OtherCancerType              *uint
	LungCancerFamily             *enum.Answer
	LungCancerFamilyRelation     *string
	OtherCancerFamily            *enum.Answer
	OtherCancerFamilyType        *uint
	OtherCancerFamilyRelation    *string
	OccupationalExposure         *string
	CurrentSmoking               *enum.Answer
	SmokingStartAgeCurrent       *uint
	SmokingTypesCurrent          *string
	CigarettesPerDayCurrent      *uint
	CigarPerDayCurrent           *uint
	ECigPerDayCurrent            *uint
	PipePerDayCurrent            *uint
	ChapoghPerDayCurrent         *uint
	SmokedOpiumPerDayCurrent     *uint
	ChewedOpiumPerDayCurrent     *uint
	HookahPerWeekCurrent         *uint
	PastSmoking                  *string
	SmokePastAvg                 *uint
	SmokeCurrentAvg              *uint
	Bronchitis                   *enum.Answer
	LungIll                      *enum.Answer
	Fibrosis                     *enum.Answer
	LeaveSmoke                   *uint
	SmokingStartAgePast          *uint
	SmokingTypesPast             *string
	CigarettesPerDayPast         *uint
	CigarPerDayPast              *uint
	ECigPerDayPast               *uint
	PipePerDayPast               *uint
	ChapoghPerDayPast            *uint
	SmokedOpiumPerDayPast        *uint
	ChewedOpiumPerDayPast        *uint
	HookahPerWeekPast            *uint
	SecondhandSmoke              *enum.Answer
	SecondhandSmokeLocation      *string
	LungDiseaseHistory           *string

	// Attention question answer
	AttentionCorrect *bool
}

type ChangeFormStatusRequest struct {
	UserID uint
	FormID uint
}

type AssignOperatorRequest struct {
	UserID     uint
	FormID     uint
	OperatorID uint
}

type UnassignOperatorRequest struct {
	UserID uint
	FormID uint
}

type UpdateGeneralHealthRequest struct {
	UserID uint
	FormID uint

	DrinksAlcohol             *enum.Answer
	CupsPerWeek               *string
	LastMonthSabzijatMeal     *string
	LastMonthSabzijatWeight   *string
	MediumActivityMonthInYear *uint
	MediumActivityHourInWeek  *string
	HardActivityMonthInYear   *uint
	HardActivityHourInWeek    *string
	SmokeAtLeast100           *enum.Answer
	SmokingAge                *uint
	SmokingNow                *enum.Answer
	YearSmoke                 *uint
	LeaveSmokingAge           *uint
	CountSmokingDaily         *string
	CountGheliandaily         *string
	CountSmokingDailyPast     *string
	CountGheliandailyPast     *string

	// Attention question answer
	AttentionCorrect *bool
}
type UpdateMamographyRequest struct {
	UserID uint
	FormID uint

	GhaedeAge                    *uint
	HasChildren                  *enum.Answer
	NumberOfChildren             *uint
	DaughterCount                *uint
	SonCount                     *uint
	AgeOfFirstBirth              *uint
	MenopausalStatus             *uint
	MenopauseAge                 *string
	HRT                          *enum.Answer
	HRTUseLength                 *uint
	LastFiveYearsHRTUse          *enum.Answer
	CurrentHRTUse                *enum.Answer
	IntendedHRTUse               *uint
	HRTType                      *string
	Oral                         *enum.Answer
	OralDuration                 *string
	OralTwoLastYears             *enum.Answer
	MamoGraphy                   *enum.Answer
	MamoGraphyPictures           []*multipart.FileHeader
	BreastDensity                *string
	Falop                        *enum.Answer
	Andometrioz                  *enum.Answer
	LeavePestan                  *bool
	LeaveTokhmdan                *bool
	LaDeColon                    *enum.Answer
	LaDePol                      *enum.Answer
	AspLaMo                      *enum.Answer
	NsaiDLaMo                    *enum.Answer
	LastFiveYearBloodTestInStool *enum.Answer
	NumberOfBreastBiopsies       *uint
	HyperplasiaInBiopsy          *uint

	// Attention question answer
	AttentionCorrect *bool
}
type UpdateContactRequest struct {
	UserID uint
	FormID uint

	Name                  *string
	TestGen               *enum.Answer
	TestGenPictures       []*multipart.FileHeader
	FmTestGen             *enum.Answer
	FatherTestGenPictures []*multipart.FileHeader
	MotherTestGenPictures []*multipart.FileHeader
	CallExpert            *bool
	BirthCountry          *string
	Province              *string
	City                  *string
	Country               *string
	Address               *string
	PostalCode            *string
	Education             *string
	Phone2                *string
	Phone3                *string
}
type UpdateLungCancerRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus              *string
	SupplementaryInsuranceStatus *enum.Answer
	SupplementaryInsurances      *string
	Hypertension                 *enum.Answer
	HypertensionTreatment        *enum.Answer
	HeartDisease                 *enum.Answer
	HeartDiseaseTreatment        *enum.Answer
	Diabetes                     *enum.Answer
	DiabetesTreatment            *enum.Answer
	ChronicLungDisease           *enum.Answer
	ChronicLungDiseaseType       *string
	LungCancerHistory            *enum.Answer
	OtherCancerHistory           *enum.Answer
	OtherCancerType              *uint
	LungCancerFamily             *enum.Answer
	LungCancerFamilyRelation     *string
	OtherCancerFamily            *enum.Answer
	OtherCancerFamilyType        *uint
	OtherCancerFamilyRelation    *string
	OccupationalExposure         *string
	CurrentSmoking               *enum.Answer
	SmokingStartAgeCurrent       *uint
	SmokingTypesCurrent          *string
	CigarettesPerDayCurrent      *uint
	CigarPerDayCurrent           *uint
	ECigPerDayCurrent            *uint
	PipePerDayCurrent            *uint
	ChapoghPerDayCurrent         *uint
	SmokedOpiumPerDayCurrent     *uint
	ChewedOpiumPerDayCurrent     *uint
	HookahPerWeekCurrent         *uint
	PastSmoking                  *string
	SmokePastAvg                 *uint
	SmokeCurrentAvg              *uint
	Bronchitis                   *enum.Answer
	LungIll                      *enum.Answer
	Fibrosis                     *enum.Answer
	LeaveSmoke                   *uint
	SmokingStartAgePast          *uint
	SmokingTypesPast             *string
	CigarettesPerDayPast         *uint
	CigarPerDayPast              *uint
	ECigPerDayPast               *uint
	PipePerDayPast               *uint
	ChapoghPerDayPast            *uint
	SmokedOpiumPerDayPast        *uint
	ChewedOpiumPerDayPast        *uint
	HookahPerWeekPast            *uint
	SecondhandSmoke              *enum.Answer
	SecondhandSmokeLocation      *string
	LungDiseaseHistory           *string

	// Attention question answer
	AttentionCorrect *bool
}

type PostalCodeInfoRequest struct {
	PostalCode string
}

type UpdateNavidFormRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus              *string
	SupplementaryInsuranceStatus *enum.Answer
	SupplementaryInsurances      *string
	Hypertension                 *string
	HypertensionTreatment        *enum.Answer
	HeartDisease                 *string
	HeartDiseaseTreatment        *enum.Answer
	Diabetes                     *string
	DiabetesTreatment            *enum.Answer
	ChronicLungDisease           *enum.Answer
	ChronicLungDiseaseType       *string
	LungCancerHistory            *enum.Answer
	OtherCancerHistory           *enum.Answer
	OtherCancerType              *uint
	LungCancerFamily             *enum.Answer
	LungCancerFamilyRelation     *string
	OtherCancerFamily            *enum.Answer
	OtherCancerFamilyType        *uint
	OtherCancerFamilyRelation    *string
	OccupationalExposure         *string
	CurrentSmoking               *enum.Answer
	SmokingStartAgeCurrent       *uint
	SmokingTypesCurrent          *string
	CigarettesPerDayCurrent      *uint
	CigarPerDayCurrent           *uint
	ECigPerDayCurrent            *uint
	PipePerDayCurrent            *uint
	ChapoghPerDayCurrent         *uint
	SmokedOpiumPerDayCurrent     *uint
	ChewedOpiumPerDayCurrent     *uint
	HookahPerWeekCurrent         *uint
	PastSmoking                  *string
	LeaveSmoke                   *uint
	SmokingStartAgePast          *uint
	SmokingTypesPast             *string
	CigarettesPerDayPast         *uint
	CigarPerDayPast              *uint
	ECigPerDayPast               *uint
	PipePerDayPast               *uint
	ChapoghPerDayPast            *uint
	SmokedOpiumPerDayPast        *uint
	ChewedOpiumPerDayPast        *uint
	HookahPerWeekPast            *uint
	SecondhandSmoke              *enum.Answer
	SecondhandSmokeLocation      *string
	LungDiseaseHistory           *string
	CurrentCigaretteSmoking      *enum.Answer
	CurrentRolledTobacco         *enum.Answer
	CurrentPipeSmoking           *enum.Answer
	CurrentHookahUse             *enum.Answer
	CurrentChiboukSmoking        *enum.Answer
	CurrentOpiumUse              *enum.Answer
	FormerCigaretteSmoking       *enum.Answer
	FormerRolledTobacco          *enum.Answer
	FormerPipeSmoking            *enum.Answer
	FormerHookahUse              *enum.Answer
	FormerChiboukSmoking         *enum.Answer
	FormerOpiumUse               *enum.Answer
	PelecSig                     *enum.Answer
	CelecSig                     *enum.Answer

	// Attention question answer
	AttentionCorrect *bool
}
