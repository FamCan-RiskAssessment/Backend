package entity

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type Form struct {
	database.Model
	Status             enum.FormStatus
	UserID             uint  `gorm:"not null;index"`
	User               User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OperatorID         *uint `gorm:"type:int"`
	FilledByOperatorID *uint `gorm:"type:int;index"`
}

type BasicInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 1
	Gender               enum.Gender `gorm:"not null"`
	BirthDate            time.Time   `gorm:"not null;type:date"`
	IsAtba               bool        `gorm:"not null;default:false"`
	SocialSecurityNumber string      `gorm:"not null"`
	Height               float64     `gorm:"type:decimal(5,2);not null"`
	Weight               float64     `gorm:"type:decimal(5,2);not null"`
}

type GeneralHealthInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	//page 2
	DrinksAlcohol             *bool   `gorm:"type:boolean"`
	CupsPerWeek               *string `gorm:"type:varchar(50)"`
	LastMonthSabzijatMeal     string  `gorm:"type:varchar(50);not null"`
	LastMonthSabzijatWeight   string  `gorm:"type:varchar(50);not null"`
	MediumActivityMonthInYear uint    `gorm:"not null"`
	MediumActivityHourInWeek  string  `gorm:"type:varchar(50);not null"`
	HardActivityMonthInYear   uint    `gorm:"not null"`
	HardActivityHourInWeek    string  `gorm:"type:varchar(50);not null"`
	SmokeAtLeast100           *bool   `gorm:"type:boolean"`
	SmokingAge                *uint   `gorm:"type:int"`
	SmokingNow                bool    `gorm:"not null;default:false"`
	LeaveSmokingAge           *uint   `gorm:"type:int"`
	CountSmokingDaily         *string `gorm:"type:varchar(50)"`
	CountGheliandaily         *string `gorm:"type:varchar(50)"`
	CountSmokingDailyPast     *string `gorm:"type:varchar(50)"`
	CountGheliandailyPast     *string `gorm:"type:varchar(50)"`
}

type MamoGraphyInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 3
	GhaedeAge                    uint                            `gorm:"not null"`
	HasChildren                  bool                            `gorm:"not null;default:false"`
	NumberOfChildren             *uint                           `gorm:"type:int"`
	AgeOfFirstBirth              *uint                           `gorm:"type:int"`
	MenopausalStatus             enum.MenopausalStatus           `gorm:"not null"`
	MenopauseAge                 *string                         `gorm:"type:varchar(50)"`
	HRT                          *bool                           `gorm:"type:boolean"`
	HRTUseLength                 *uint                           `gorm:"type:int"`
	LastFiveYearsHRTUse          bool                            `gorm:"not null;default:false"`
	CurrentHRTUse                *bool                           `gorm:"type:boolean"`
	IntendedHRTUse               *uint                           `gorm:"type:int"`
	HRTType                      *string                         `gorm:"type:varchar(50)"`
	Oral                         *bool                           `gorm:"type:boolean"`
	OralDuration                 *string                         `gorm:"type:varchar(50)"`
	OralTwoLastYears             *bool                           `gorm:"type:boolean"`
	MamoGraphy                   *bool                           `gorm:"type:boolean"`
	MamoGraphyPicturePath        *string                         `gorm:"type:varchar(255)"`
	Falop                        *bool                           `gorm:"type:boolean"`
	Andometrioz                  *bool                           `gorm:"type:boolean"`
	LeavePestan                  bool                            `gorm:"not null;default:false"`
	LeaveTokhmdan                bool                            `gorm:"not null;default:false"`
	LaDeColon                    *bool                           `gorm:"type:boolean"`
	LaDePol                      *bool                           `gorm:"type:boolean"`
	AspLaMo                      *bool                           `gorm:"type:boolean"`
	NsaiDLaMo                    *bool                           `gorm:"type:boolean"`
	LastFiveYearBloodTestInStool *bool                           `gorm:"type:boolean"`
	NumberOfBreastBiopsies       *uint                           `gorm:"type:int"`
	HyperplasiaInBiopsy          *enum.HyperplasiaInBiopsyStatus `gorm:"type:int"`
}

type CancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`

	CancerAge   uint
	CancerType  enum.CancerType
	PicturePath *string `gorm:"type:varchar(255)"`
}

type FamilyCancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`

	Relative         enum.Relative
	RelativeRelation *string `gorm:"type:varchar(127)"`
	Name             *string `gorm:"type:varchar(127)"`
	LifeStatus       *enum.LifeStatus
	CancerAge        uint
	CancerType       enum.CancerType
	PicturePath      *string `gorm:"type:varchar(255)"`
}

type ContactInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 6
	Name    string `gorm:"type:varchar(255);not null"`
	TestGen *bool  `gorm:"type:boolean"`
	// Pics
	FmTestGen *bool `gorm:"type:boolean"`
	// Pics
	CallExpert   bool    `gorm:"not null;default:true"`
	BirthCountry *string `gorm:"type:varchar(50)"`
	Province     *string `gorm:"type:varchar(50)"`
	City         *string `gorm:"type:varchar(50)"`
	Country      *string `gorm:"type:varchar(50)"`
	Address      string  `gorm:"type:text;not null"`
	PostalCode   string  `gorm:"not null"`
}

type LungCancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 7
	InsuranceStatus           *string `gorm:"type:varchar(127)"`
	SupplementaryInsurances   *string `gorm:"type:varchar(50)"`
	Hypertension              bool    `gorm:"not null;default:false"`
	HypertensionTreatment     *bool   `gorm:"type:boolean"`
	HeartDisease              bool    `gorm:"not null;default:false"`
	HeartDiseaseTreatment     *bool   `gorm:"type:boolean"`
	Diabetes                  bool    `gorm:"not null;default:false"`
	DiabetesTreatment         *bool   `gorm:"type:boolean"`
	ChronicLungDisease        *bool   `gorm:"type:boolean"`
	ChronicLungDiseaseType    *string `gorm:"type:varchar(50)"`
	LungCancerHistory         bool    `gorm:"not null;default:false"`
	OtherCancerHistory        bool    `gorm:"not null;default:false"`
	OtherCancerType           *enum.CancerType
	LungCancerFamily          *bool   `gorm:"type:boolean"`
	LungCancerFamilyRelation  *string `gorm:"type:varchar(50)"`
	OtherCancerFamily         *bool   `gorm:"type:boolean"`
	OtherCancerFamilyType     *enum.CancerType
	OtherCancerFamilyRelation *string `gorm:"type:varchar(50)"`
	OccupationalExposure      *string `gorm:"type:varchar(255)"`
	CurrentSmoking            bool    `gorm:"not null;default:false"`
	SmokingStartAgeCurrent    *uint   `gorm:"type:int"`
	SmokingTypesCurrent       *string `gorm:"type:varchar(50)"`
	CigarettesPerDayCurrent   *uint   `gorm:"type:int"`
	CigarPerDayCurrent        *uint   `gorm:"type:int"`
	ECigPerDayCurrent         *uint   `gorm:"type:int"`
	PipePerDayCurrent         *uint   `gorm:"type:int"`
	ChapoghPerDayCurrent      *uint   `gorm:"type:int"`
	SmokedOpiumPerDayCurrent  *uint   `gorm:"type:int"`
	ChewedOpiumPerDayCurrent  *uint   `gorm:"type:int"`
	HookahPerWeekCurrent      *uint   `gorm:"type:int"`
	PastSmoking               *string `gorm:"type:varchar(127)"`
	SmokingStartAgePast       *uint   `gorm:"type:int"`
	SmokingTypesPast          *string `gorm:"type:varchar(50)"`
	CigarettesPerDayPast      *uint   `gorm:"type:int"`
	CigarPerDayPast           *uint   `gorm:"type:int"`
	ECigPerDayPast            *uint   `gorm:"type:int"`
	PipePerDayPast            *uint   `gorm:"type:int"`
	ChapoghPerDayPast         *uint   `gorm:"type:int"`
	SmokedOpiumPerDayPast     *uint   `gorm:"type:int"`
	ChewedOpiumPerDayPast     *uint   `gorm:"type:int"`
	HookahPerWeekPast         *uint   `gorm:"type:int"`
	SecondhandSmoke           bool    `gorm:"not null;default:false"`
	SecondhandSmokeLocation   *string `gorm:"type:varchar(50)"`
}

type Premm5Result struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:onDelete:CASCADE"`

	// Gene-specific mutation probabilities
	MLH1Probability float64 `gorm:"type:decimal(10,6);not null"`
	MSH2Probability float64 `gorm:"type:decimal(10,6);not null"`
	MSH6Probability float64 `gorm:"type:decimal(10,6);not null"`
	PMS2Probability float64 `gorm:"type:decimal(10,6);not null"`

	// Overall PREMM5 scores
	PAny  float64 `gorm:"type:decimal(10,6);not null"` // Overall PREMM5 score (probability of any mutation)
	PNone float64 `gorm:"type:decimal(10,6);not null"` // Probability of no mutation
}

type BCRAResult struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:onDelete:CASCADE"`

	AbsRisk    float64 `gorm:"type:decimal(10,6);not null"` // Absolute risk percentage
	AbsRiskAvg float64 `gorm:"type:decimal(10,6);not null"` // Absolute risk average
	RRStar1    float64 `gorm:"type:decimal(10,6);not null"` // Relative risk star 1
	RRStar2    float64 `gorm:"type:decimal(10,6);not null"` // Relative risk star 2
	ProjIntvl  float64 `gorm:"type:decimal(10,6);not null"` // Projection interval
}

type GailResult struct {
	database.Model
	FormID uint `gorm:"not null"`
	Form   Form `gorm:"foreignKey:FormID;constraint:onDelete:CASCADE"`

	AbsoluteRisk float64 `gorm:"not null"`
	RelativeRisk *float64
}

type PLCOResult struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:onDelete:CASCADE"`

	PLCOM20126YrRisk     float64  `gorm:"type:decimal(10,6);not null"`
	PLCOM20123YrRisk     float64  `gorm:"type:decimal(10,6);not null"`
	PLCOM2012RiskPercent float64  `gorm:"type:decimal(10,6);not null"`
	PLCO2012Results3Yr   *float64 `gorm:"type:decimal(10,6)"`
	PLCO2012Results6Yr   *float64 `gorm:"type:decimal(10,6)"`
}
