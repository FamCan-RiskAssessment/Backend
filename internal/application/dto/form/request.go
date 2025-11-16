package formdto

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
)

type BirthDate time.Time

func (d *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = BirthDate(t)
	return nil
}

func (d BirthDate) Time() time.Time {
	return time.Time(d)
}

func (d BirthDate) String() string {
	return time.Time(d).Format("2006-01-02")
}

type CreateBasicFormRequest struct {
	UserID             uint
	FilledByOperatorID *uint

	// page 1
	BirthDate            time.Time
	SocialSecurityNumber string
	Gender               uint
	IsAtba               bool
	Height               float64
	Weight               float64
}

type UpdateBasicFormRequest struct {
	UserID uint
	FormID uint

	// page 1
	BirthDate            *time.Time
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

	DrinksAlcohol             *bool
	CupsPerWeek               *string
	LastMonthSabzijatMeal     string
	LastMonthSabzijatWeight   string
	MediumActivityMonthInYear uint
	MediumActivityHourInWeek  string
	HardActivityMonthInYear   uint
	HardActivityHourInWeek    string
	SmokeAtLeast100           *bool
	SmokingAge                *uint
	SmokingNow                bool
	LeaveSmokingAge           *uint
	CountSmokingDaily         *string
	CountGheliandaily         *string
	CountSmokingDailyPast     *string
	CountGheliandailyPast     *string
}

type UpsertMamographyRequest struct {
	UserID uint
	FormID uint

	GhaedeAge                    uint
	HasChildren                  bool
	NumberOfChildren             *uint
	AgeOfFirstBirth              *uint
	MenopausalStatus             uint
	MenopauseAge                 *string
	HRT                          *bool
	HRTUseLength                 *uint
	LastFiveYearsHRTUse          bool
	CurrentHRTUse                *bool
	IntendedHRTUse               *uint
	HRTType                      *string
	Oral                         *bool
	OralDuration                 *string
	OralTwoLastYears             *bool
	MamoGraphy                   *bool
	MamoGraphyPicture            *multipart.FileHeader
	Falop                        *bool
	Andometrioz                  *bool
	LeavePestan                  bool
	LeaveTokhmdan                bool
	LaDeColon                    *bool
	LaDePol                      *bool
	AspLaMo                      *bool
	NsaiDLaMo                    *bool
	LastFiveYearBloodTestInStool *bool
	NumberOfBreastBiopsies       *uint
	HyperplasiaInBiopsy          *uint
}

type CancerRequest struct {
	CancerType uint
	CancerAge  uint
	Picture    *multipart.FileHeader
}

type CreateCancerRequest struct {
	UserID uint
	FormID uint

	CancerType uint
	CancerAge  uint
	Picture    *multipart.FileHeader
}

type UpdateCancerRequest struct {
	UserID   uint
	FormID   uint
	CancerID uint

	CancerType uint
	CancerAge  uint
	Picture    *multipart.FileHeader
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
	LifeStatus       *enum.LifeStatus
	CancerType       uint
	CancerAge        uint
	Picture          *multipart.FileHeader
}

type UpdateFamilyCancerRequest struct {
	UserID           uint
	FormID           uint
	FamilyCancerID   uint
	Relative         enum.Relative
	RelativeRelation *string
	Name             *string
	LifeStatus       *enum.LifeStatus
	CancerType       uint
	CancerAge        uint
	Picture          *multipart.FileHeader
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

	Name                 string
	TestGen              *bool
	TestGenPicture       *multipart.FileHeader
	FmTestGen            *bool
	FatherTestGenPicture *multipart.FileHeader
	MotherTestGenPicture *multipart.FileHeader
	CallExpert           bool
	BirthCountry         *string
	Province             *string
	City                 *string
	Country              *string
	Address              string
	PostalCode           string
}

type UpsertLungCancerRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus           *string
	SupplementaryInsurances   *string
	Hypertension              bool
	HypertensionTreatment     *bool
	HeartDisease              bool
	HeartDiseaseTreatment     *bool
	Diabetes                  bool
	DiabetesTreatment         *bool
	ChronicLungDisease        *bool
	ChronicLungDiseaseType    *string
	LungCancerHistory         bool
	OtherCancerHistory        bool
	OtherCancerType           *uint
	LungCancerFamily          *bool
	LungCancerFamilyRelation  *string
	OtherCancerFamily         *bool
	OtherCancerFamilyType     *uint
	OtherCancerFamilyRelation *string
	OccupationalExposure      *string
	CurrentSmoking            bool
	SmokingStartAgeCurrent    *uint
	SmokingTypesCurrent       *string
	CigarettesPerDayCurrent   *uint
	CigarPerDayCurrent        *uint
	ECigPerDayCurrent         *uint
	PipePerDayCurrent         *uint
	ChapoghPerDayCurrent      *uint
	SmokedOpiumPerDayCurrent  *uint
	ChewedOpiumPerDayCurrent  *uint
	HookahPerWeekCurrent      *uint
	PastSmoking               *string
	SmokingStartAgePast       *uint
	SmokingTypesPast          *string
	CigarettesPerDayPast      *uint
	CigarPerDayPast           *uint
	ECigPerDayPast            *uint
	PipePerDayPast            *uint
	ChapoghPerDayPast         *uint
	SmokedOpiumPerDayPast     *uint
	ChewedOpiumPerDayPast     *uint
	HookahPerWeekPast         *uint
	SecondhandSmoke           bool
	SecondhandSmokeLocation   *string
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

	DrinksAlcohol             *bool
	CupsPerWeek               *string
	LastMonthSabzijatMeal     *string
	LastMonthSabzijatWeight   *string
	MediumActivityMonthInYear *uint
	MediumActivityHourInWeek  *string
	HardActivityMonthInYear   *uint
	HardActivityHourInWeek    *string
	SmokeAtLeast100           *bool
	SmokingAge                *uint
	SmokingNow                *bool
	LeaveSmokingAge           *uint
	CountSmokingDaily         *string
	CountGheliandaily         *string
	CountSmokingDailyPast     *string
	CountGheliandailyPast     *string
}
type UpdateMamographyRequest struct {
	UserID uint
	FormID uint

	GhaedeAge                    *uint
	HasChildren                  *bool
	NumberOfChildren             *uint
	AgeOfFirstBirth              *uint
	MenopausalStatus             *uint
	MenopauseAge                 *string
	HRT                          *bool
	HRTUseLength                 *uint
	LastFiveYearsHRTUse          *bool
	CurrentHRTUse                *bool
	IntendedHRTUse               *uint
	HRTType                      *string
	Oral                         *bool
	OralDuration                 *string
	OralTwoLastYears             *bool
	MamoGraphy                   *bool
	Falop                        *bool
	Andometrioz                  *bool
	LeavePestan                  *bool
	LeaveTokhmdan                *bool
	LaDeColon                    *bool
	LaDePol                      *bool
	AspLaMo                      *bool
	NsaiDLaMo                    *bool
	LastFiveYearBloodTestInStool *bool
	NumberOfBreastBiopsies       *uint
	HyperplasiaInBiopsy          *uint
}
type UpdateContactRequest struct {
	UserID uint
	FormID uint

	Name                 *string
	TestGen              *bool
	TestGenPicture       *multipart.FileHeader
	FmTestGen            *bool
	FatherTestGenPicture *multipart.FileHeader
	MotherTestGenPicture *multipart.FileHeader
	CallExpert           *bool
	BirthCountry         *string
	Province             *string
	City                 *string
	Country              *string
	Address              *string
	PostalCode           *string
}
type UpdateLungCancerRequest struct {
	UserID uint
	FormID uint

	InsuranceStatus           *string
	SupplementaryInsurances   *string
	Hypertension              *bool
	HypertensionTreatment     *bool
	HeartDisease              *bool
	HeartDiseaseTreatment     *bool
	Diabetes                  *bool
	DiabetesTreatment         *bool
	ChronicLungDisease        *bool
	ChronicLungDiseaseType    *string
	LungCancerHistory         *bool
	OtherCancerHistory        *bool
	OtherCancerType           *uint
	LungCancerFamily          *bool
	LungCancerFamilyRelation  *string
	OtherCancerFamily         *bool
	OtherCancerFamilyType     *uint
	OtherCancerFamilyRelation *string
	OccupationalExposure      *string
	CurrentSmoking            *bool
	SmokingStartAgeCurrent    *uint
	SmokingTypesCurrent       *string
	CigarettesPerDayCurrent   *uint
	CigarPerDayCurrent        *uint
	ECigPerDayCurrent         *uint
	PipePerDayCurrent         *uint
	ChapoghPerDayCurrent      *uint
	SmokedOpiumPerDayCurrent  *uint
	ChewedOpiumPerDayCurrent  *uint
	HookahPerWeekCurrent      *uint
	PastSmoking               *string
	SmokingStartAgePast       *uint
	SmokingTypesPast          *string
	CigarettesPerDayPast      *uint
	CigarPerDayPast           *uint
	ECigPerDayPast            *uint
	PipePerDayPast            *uint
	ChapoghPerDayPast         *uint
	SmokedOpiumPerDayPast     *uint
	ChewedOpiumPerDayPast     *uint
	HookahPerWeekPast         *uint
	SecondhandSmoke           *bool
	SecondhandSmokeLocation   *string
}
