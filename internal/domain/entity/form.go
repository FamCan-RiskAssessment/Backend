package entity

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type Form struct {
	database.Model
	Status     enum.FormStatus
	UserID     uint  `gorm:"not null;index"`
	User       User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OperatorID *uint `gorm:"type:int"`
}

type BasicInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 1
	Gender               enum.Gender `gorm:"not null"`
	BirthYear            uint        `gorm:"not null"`
	BirthMonth           string      `gorm:"type:varchar(10);not null"`
	BirthDay             uint        `gorm:"not null"`
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
	// page 4
	Cancer     bool  `gorm:"not null;default:false"`
	CancerAge  *uint `gorm:"type:int"`
	CancerType *enum.CancerType
	// Pics
}

type FamilyCancerInfo struct {
	database.Model
	FormID uint `gorm:"not null;index"`
	Form   Form `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE"`
	// page 5
	ChildCancer     bool    `gorm:"not null;default:false"`
	ChildName       *string `gorm:"type:varchar(50)"`
	ChildCancerType *enum.CancerType
	ChildCancerAge  *uint   `gorm:"type:int"`
	ChildLifeStatus *string `gorm:"type:varchar(50)"`
	// Pics
	MotherCancer     bool    `gorm:"not null;default:false"`
	MotherName       *string `gorm:"type:varchar(50)"`
	MotherLifeStatus *string `gorm:"type:varchar(50)"`
	MotherCancerType *enum.CancerType
	MotherCancerAge  *uint `gorm:"type:int"`
	// Pics
	FatherCancer     bool    `gorm:"not null;default:false"`
	FatherName       *string `gorm:"type:varchar(50)"`
	FatherLifeStatus *string `gorm:"type:varchar(50)"`
	FatherCancerType *enum.CancerType
	FatherCancerAge  *uint `gorm:"type:int"`
	// Pics
	SiblingCancer     bool    `gorm:"not null;default:false"`
	SiblingName       *string `gorm:"type:varchar(50)"`
	SiblingLifeStatus *string `gorm:"type:varchar(50)"`
	SiblingCancerType *enum.CancerType
	SiblingCancerAge  *uint `gorm:"type:int"`
	// Pics
	AmeAmoCancer     bool    `gorm:"not null;default:false"`
	AmeAmoName       *string `gorm:"type:varchar(50)"`
	AmeAmoLifeStatus *string `gorm:"type:varchar(50)"`
	AmeAmoCancerType *enum.CancerType
	AmeAmoCancerAge  *uint `gorm:"type:int"`
	// Pics
	KhaleDaeiCancer     bool    `gorm:"not null;default:false"`
	KhaleDaeiName       *string `gorm:"type:varchar(50)"`
	KhaleDaeiLifeStatus *string `gorm:"type:varchar(50)"`
	KhaleDaeiCancerType *enum.CancerType
	KhaleDaeiCancerAge  *uint `gorm:"type:int"`
	// Pics
	OtherRelativeCancer     *bool   `gorm:"type:boolean"`
	OtherRelativeName       *string `gorm:"type:varchar(50)"`
	OtherRelativeRelation   *string `gorm:"type:varchar(50)"`
	OtherRelativeLifeStatus *string `gorm:"type:varchar(50)"`
	OtherRelativeCancerType *enum.CancerType
	OtherRelativeCancerAge  *uint `gorm:"type:int"`
	// Pics
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
