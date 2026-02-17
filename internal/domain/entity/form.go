package entity

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type Form struct {
	database.Model
	Status             enum.FormStatus
	UserID             uint          `gorm:"not null;index"`
	User               User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FormType           enum.FormType `gorm:"not null;default:1"`
	OperatorID         *uint         `gorm:"type:int"`
	FilledByOperatorID *uint         `gorm:"type:int;index"`
}

type BasicInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 1
	Gender               enum.Gender `gorm:"not null"`
	BirthDate            time.Time   `gorm:"not null;type:date"`
	IsAtba               bool        `gorm:"not null;default:false"`
	SocialSecurityNumber string      `gorm:"type:varchar(500);not null"` // Encrypted, larger size needed
	Height               float64     `gorm:"type:decimal(5,2);not null"`
	Weight               float64     `gorm:"type:decimal(5,2);not null"`
}

type GeneralHealthInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	//page 2
	DrinksAlcohol             *enum.Answer
	CupsPerWeek               *string      `gorm:"type:varchar(50)"`
	LastMonthSabzijatMeal     string       `gorm:"type:varchar(50);not null"`
	LastMonthSabzijatWeight   string       `gorm:"type:varchar(50);not null"`
	MediumActivityMonthInYear uint         `gorm:"not null"`
	MediumActivityHourInWeek  string       `gorm:"type:varchar(50);not null"`
	HardActivityMonthInYear   uint         `gorm:"not null"`
	HardActivityHourInWeek    string       `gorm:"type:varchar(50);not null"`
	SmokeAtLeast100           *enum.Answer `gorm:"type:int"`
	SmokingAge                *uint        `gorm:"type:int"`
	SmokingNow                *enum.Answer `gorm:"type:int"`
	YearSmoke                 *uint        `gorm:"type:int"`
	LeaveSmokingAge           *uint        `gorm:"type:int"`
	CountSmokingDaily         *string      `gorm:"type:varchar(50)"`
	CountGheliandaily         *string      `gorm:"type:varchar(50)"`
	CountSmokingDailyPast     *string      `gorm:"type:varchar(50)"`
	CountGheliandailyPast     *string      `gorm:"type:varchar(50)"`
}

type MamoGraphyInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 3
	GhaedeAge                    uint                            `gorm:"not null"`
	HasChildren                  *enum.Answer                    `gorm:"type:int"`
	NumberOfChildren             *uint                           `gorm:"type:int"`
	SonCount                     *uint                           `gorm:"type:int"`
	DaughterCount                *uint                           `gorm:"type:int"`
	AgeOfFirstBirth              *uint                           `gorm:"type:int"`
	MenopausalStatus             enum.MenopausalStatus           `gorm:"not null"`
	MenopauseAge                 *string                         `gorm:"type:varchar(50)"`
	HRT                          *enum.Answer                    `gorm:"type:int"`
	HRTUseLength                 *uint                           `gorm:"type:int"`
	LastFiveYearsHRTUse          *enum.Answer                    `gorm:"type:int"`
	CurrentHRTUse                *enum.Answer                    `gorm:"type:int"`
	IntendedHRTUse               *uint                           `gorm:"type:int"`
	HRTType                      *string                         `gorm:"type:varchar(50)"`
	Oral                         *enum.Answer                    `gorm:"type:int"`
	OralDuration                 *string                         `gorm:"type:varchar(50)"`
	OralTwoLastYears             *enum.Answer                    `gorm:"type:int"`
	MamoGraphy                   *enum.Answer                    `gorm:"type:int"`
	MamoGraphyPicturePaths       []string                        `gorm:"type:jsonb;serializer:json"`
	BreastDensity                *string                         `gorm:"type:varchar(50)"`
	Falop                        *enum.Answer                    `gorm:"type:int"`
	Andometrioz                  *enum.Answer                    `gorm:"type:int"`
	LeavePestan                  bool                            `gorm:"not null;default:false"`
	LeaveTokhmdan                bool                            `gorm:"not null;default:false"`
	LaDeColon                    *enum.Answer                    `gorm:"type:int"`
	LaDePol                      *enum.Answer                    `gorm:"type:int"`
	AspLaMo                      *enum.Answer                    `gorm:"type:int"`
	NsaiDLaMo                    *enum.Answer                    `gorm:"type:int"`
	LastFiveYearBloodTestInStool *enum.Answer                    `gorm:"type:int"`
	NumberOfBreastBiopsies       *uint                           `gorm:"type:int"`
	HyperplasiaInBiopsy          *enum.HyperplasiaInBiopsyStatus `gorm:"type:int"`
}

type CancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`

	CancerAge    uint
	CancerType   enum.CancerType
	PicturePaths []string `gorm:"type:jsonb;serializer:json"`
}

type FamilyCancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`

	Relative         enum.Relative
	RelativeRelation *string `gorm:"type:varchar(127)"`
	Name             *string `gorm:"type:varchar(127)"`
	NumberRelative   *uint
	LifeStatus       *enum.LifeStatus
	CancerAge        uint
	CancerType       enum.CancerType
	PicturePaths     []string `gorm:"type:jsonb;serializer:json"`
}

type ContactInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 6
	Name                      string       `gorm:"type:varchar(255);not null"`
	TestGen                   *enum.Answer `gorm:"type:int"`
	TestGenPicturePaths       []string     `gorm:"type:jsonb;serializer:json"`
	FmTestGen                 *enum.Answer `gorm:"type:int"`
	FatherTestGenPicturePaths []string     `gorm:"type:jsonb;serializer:json"`
	MotherTestGenPicturePaths []string     `gorm:"type:jsonb;serializer:json"`
	CallExpert                bool         `gorm:"not null;default:true"`
	BirthCountry              *string      `gorm:"type:varchar(50)"`
	Province                  *string      `gorm:"type:varchar(50)"`
	City                      *string      `gorm:"type:varchar(50)"`
	Country                   *string      `gorm:"type:varchar(50)"`
	Address                   string       `gorm:"type:text;not null"` // Encrypted, text type for variable length
	PostalCode                string       `gorm:"not null"`
	Education                 string       `gorm:"type:varchar(127)"`
	Phone2                    *string      `gorm:"type:varchar(11)"`
	Phone3                    *string      `gorm:"type:varchar(11)"`
}

type LungCancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 7
	InsuranceStatus              *string      `gorm:"type:varchar(127)"`
	SupplementaryInsuranceStatus *enum.Answer `gorm:"type:int"`
	SupplementaryInsurances      *string      `gorm:"type:varchar(50)"`
	ChronicLungDisease           *enum.Answer `gorm:"type:int"`
	ChronicLungDiseaseType       *string      `gorm:"type:varchar(50)"`
	LungCancerHistory            *enum.Answer `gorm:"type:int"`
	OtherCancerHistory           *enum.Answer `gorm:"type:int"`
	OtherCancerType              *enum.CancerType
	LungCancerFamily             *enum.Answer `gorm:"type:int"`
	LungCancerFamilyRelation     *string      `gorm:"type:varchar(50)"`
	OtherCancerFamily            *enum.Answer `gorm:"type:int"`
	OtherCancerFamilyType        *enum.CancerType
	OtherCancerFamilyRelation    *string      `gorm:"type:varchar(50)"`
	OccupationalExposure         *string      `gorm:"type:varchar(255)"`
	CurrentSmoking               *enum.Answer `gorm:"type:int"`
	SmokingStartAgeCurrent       *uint        `gorm:"type:int"`
	SmokingTypesCurrent          *string      `gorm:"type:varchar(50)"`
	CigarettesPerDayCurrent      *uint        `gorm:"type:int"`
	CigarPerDayCurrent           *uint        `gorm:"type:int"`
	ECigPerDayCurrent            *uint        `gorm:"type:int"`
	PipePerDayCurrent            *uint        `gorm:"type:int"`
	ChapoghPerDayCurrent         *uint        `gorm:"type:int"`
	SmokedOpiumPerDayCurrent     *uint        `gorm:"type:int"`
	ChewedOpiumPerDayCurrent     *uint        `gorm:"type:int"`
	HookahPerWeekCurrent         *uint        `gorm:"type:int"`
	PastSmoking                  *string      `gorm:"type:varchar(127)"`
	SmokePastAvg                 *uint        `gorm:"type:int"`
	SmokeCurrentAvg              *uint        `gorm:"type:int"`
	Bronchitis                   *enum.Answer `gorm:"type:int"`
	LungIll                      *enum.Answer `gorm:"type:int"`
	Fibrosis                     *enum.Answer `gorm:"type:int"`
	LeaveSmoke                   *uint        `gorm:"type:int"`
	SmokingStartAgePast          *uint        `gorm:"type:int"`
	SmokingTypesPast             *string      `gorm:"type:varchar(50)"`
	CigarettesPerDayPast         *uint        `gorm:"type:int"`
	CigarPerDayPast              *uint        `gorm:"type:int"`
	ECigPerDayPast               *uint        `gorm:"type:int"`
	PipePerDayPast               *uint        `gorm:"type:int"`
	ChapoghPerDayPast            *uint        `gorm:"type:int"`
	SmokedOpiumPerDayPast        *uint        `gorm:"type:int"`
	ChewedOpiumPerDayPast        *uint        `gorm:"type:int"`
	HookahPerWeekPast            *uint        `gorm:"type:int"`
	SecondhandSmoke              *enum.Answer `gorm:"type:int"`
	SecondhandSmokeLocation      *string      `gorm:"type:varchar(50)"`
	LungDiseaseHistory           *string      `gorm:"type:varchar(127)"`
}

type NavidInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 7
	InsuranceStatus              *string      `gorm:"type:varchar(127)"`
	SupplementaryInsuranceStatus *enum.Answer `gorm:"type:int"`
	SupplementaryInsurances      *string      `gorm:"type:varchar(50)"`
	Hypertension                 string       `gorm:"type:varchar(127)"`
	HypertensionTreatment        *enum.Answer `gorm:"type:int"`
	HeartDisease                 string       `gorm:"type:varchar(127)"`
	HeartDiseaseTreatment        *enum.Answer `gorm:"type:int"`
	Diabetes                     string       `gorm:"type:varchar(127)"`
	DiabetesTreatment            *enum.Answer `gorm:"type:int"`
	ChronicLungDisease           *enum.Answer `gorm:"type:int"`
	ChronicLungDiseaseType       *string      `gorm:"type:varchar(50)"`
	LungCancerHistory            *enum.Answer `gorm:"type:int"`
	OtherCancerHistory           *enum.Answer `gorm:"type:int"`
	OtherCancerType              *enum.CancerType
	LungCancerFamily             *enum.Answer `gorm:"type:int"`
	LungCancerFamilyRelation     *string      `gorm:"type:varchar(50)"`
	OtherCancerFamily            *enum.Answer `gorm:"type:int"`
	OtherCancerFamilyType        *enum.CancerType
	OtherCancerFamilyRelation    *string      `gorm:"type:varchar(50)"`
	OccupationalExposure         *string      `gorm:"type:varchar(255)"`
	CurrentSmoking               *enum.Answer `gorm:"type:int"`
	SmokingStartAgeCurrent       *uint        `gorm:"type:int"`
	SmokingTypesCurrent          *string      `gorm:"type:varchar(50)"`
	CigarettesPerDayCurrent      *uint        `gorm:"type:int"`
	CigarPerDayCurrent           *uint        `gorm:"type:int"`
	ECigPerDayCurrent            *uint        `gorm:"type:int"`
	PipePerDayCurrent            *uint        `gorm:"type:int"`
	ChapoghPerDayCurrent         *uint        `gorm:"type:int"`
	SmokedOpiumPerDayCurrent     *uint        `gorm:"type:int"`
	ChewedOpiumPerDayCurrent     *uint        `gorm:"type:int"`
	HookahPerWeekCurrent         *uint        `gorm:"type:int"`
	PastSmoking                  *string      `gorm:"type:varchar(127)"`
	LeaveSmoke                   *uint        `gorm:"type:int"`
	SmokingStartAgePast          *uint        `gorm:"type:int"`
	SmokingTypesPast             *string      `gorm:"type:varchar(50)"`
	CigarettesPerDayPast         *uint        `gorm:"type:int"`
	CigarPerDayPast              *uint        `gorm:"type:int"`
	ECigPerDayPast               *uint        `gorm:"type:int"`
	PipePerDayPast               *uint        `gorm:"type:int"`
	ChapoghPerDayPast            *uint        `gorm:"type:int"`
	SmokedOpiumPerDayPast        *uint        `gorm:"type:int"`
	ChewedOpiumPerDayPast        *uint        `gorm:"type:int"`
	HookahPerWeekPast            *uint        `gorm:"type:int"`
	SecondhandSmoke              *enum.Answer `gorm:"type:int"`
	SecondhandSmokeLocation      *string      `gorm:"type:varchar(50)"`
	LungDiseaseHistory           *string      `gorm:"type:varchar(127)"`
	CurrentCigaretteSmoking      *enum.Answer `gorm:"type:int"`
	CurrentRolledTobacco         *enum.Answer `gorm:"type:int"`
	CurrentPipeSmoking           *enum.Answer `gorm:"type:int"`
	CurrentHookahUse             *enum.Answer `gorm:"type:int"`
	CurrentChiboukSmoking        *enum.Answer `gorm:"type:int"`
	CurrentOpiumUse              *enum.Answer `gorm:"type:int"`
	FormerCigaretteSmoking       *enum.Answer `gorm:"type:int"`
	FormerRolledTobacco          *enum.Answer `gorm:"type:int"`
	FormerPipeSmoking            *enum.Answer `gorm:"type:int"`
	FormerHookahUse              *enum.Answer `gorm:"type:int"`
	FormerChiboukSmoking         *enum.Answer `gorm:"type:int"`
	FormerOpiumUse               *enum.Answer `gorm:"type:int"`
	PelecSig                     *enum.Answer `gorm:"type:int"`
	CelecSig                     *enum.Answer `gorm:"type:int"`
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

type AttentionQuestions struct {
	database.Model
	FormID               uint `gorm:"not null;uniqueIndex"`
	Form                 Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	GeneralHealthCorrect *bool
	MamographyCorrect    *bool
	LungCancerCorrect    *bool
}
