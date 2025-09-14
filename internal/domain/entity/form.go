package entity

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type Form struct {
	database.Model
	UserID               uint      `gorm:"not null;index"`
	User                 User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name                 string    `gorm:"type:varchar(255);not null"`
	DateOfBirth          time.Time `gorm:"type:date;not null"`
	Address              string    `gorm:"type:text;not null"`
	PostalCode           string    `gorm:"type:varchar(20);not null"`
	SocialSecurityNumber string    `gorm:"type:varchar(20);not null"`
	Gender               string    `gorm:"type:varchar(10);not null"`
	IsAtba               bool      `gorm:"not null;default:false"`
	Height               float64   `gorm:"type:decimal(5,2);not null"`
	Weight               float64   `gorm:"type:decimal(5,2);not null"`

	DrinksAlcohol             *bool   `gorm:"type:boolean"`
	CupsPerWeek               *string `gorm:"type:varchar(50)"`
	LastMonthSabzijatMeal     string  `gorm:"type:varchar(50);not null"`
	LastMonthSabzijatWeight   string  `gorm:"type:varchar(50);not null"`
	MediumActivityMonthInYear uint    `gorm:"not null"`
	MediumActivityHourInWeek  string  `gorm:"type:varchar(50);not null"`
	HardActivityMonthInYear   uint    `gorm:"not null"`
	HardActivityHourInWeek    string  `gorm:"type:varchar(50);not null"`
	SmokeAtLeast100           *bool   `gorm:"type:boolean"`
	SmokingAge                uint    `gorm:"not null"`
	SmokingNow                bool    `gorm:"not null;default:false"`
	LeaveSmokingAge           *uint   `gorm:"type:int"`
	CountSmokingDaily         *uint   `gorm:"type:int"`
	CountGheliandaily         *uint   `gorm:"type:int"`
	CountSmokingDailyPast     *uint   `gorm:"type:int"`
	CountGheliandailyPast     *uint   `gorm:"type:int"`

	GhaedeAge                    uint    `gorm:"not null"`
	HasChildren                  bool    `gorm:"not null;default:false"`
	NumberOfChildren             *uint   `gorm:"type:int"`
	AgeOfFirstBirth              *uint   `gorm:"type:int"`
	MenopausalStatus             string  `gorm:"type:varchar(50);not null"`
	MenopauseAge                 string  `gorm:"type:varchar(50);not null"`
	HRT                          *bool   `gorm:"type:boolean"`
	HRTUseLength                 *uint   `gorm:"type:int"`
	LastFiveYearsHRTUse          bool    `gorm:"not null;default:false"`
	CurrentHRTUse                *bool   `gorm:"type:boolean"`
	IntendedHRTUse               *uint   `gorm:"type:int"`
	HRTType                      *string `gorm:"type:varchar(50)"`
	Oral                         *bool   `gorm:"type:boolean"`
	OralDuration                 *string `gorm:"type:varchar(50)"`
	OralTwoLastYears             *bool   `gorm:"type:boolean"`
	MamoGraphy                   *bool   `gorm:"type:boolean"`
	Falop                        *bool   `gorm:"type:boolean"`
	Andometrioz                  *bool   `gorm:"type:boolean"`
	LeavePestan                  bool    `gorm:"not null;default:false"`
	LeaveTokhmdan                bool    `gorm:"not null;default:false"`
	LaDeColon                    *bool   `gorm:"type:boolean"`
	LaDePol                      *bool   `gorm:"type:boolean"`
	AspLaMo                      *bool   `gorm:"type:boolean"`
	NsaiDLaMo                    *bool   `gorm:"type:boolean"`
	LastFiveYearBloodTestInStool *bool   `gorm:"type:boolean"`

	Cancer     bool    `gorm:"not null;default:false"`
	CancerType *string `gorm:"type:varchar(50)"`
	CancerAge  *uint   `gorm:"type:int"`
	// Pics

	ChildCancer     bool    `gorm:"not null;default:false"`
	ChildName       *string `gorm:"type:varchar(50)"`
	ChildCancerType *string `gorm:"type:varchar(50)"`
	ChildCancerAge  *string `gorm:"type:varchar(50)"`
	ChildLifeStatus *string `gorm:"type:varchar(50)"`
	// Pics
	MotherCancer     bool    `gorm:"not null;default:false"`
	MotherName       *string `gorm:"type:varchar(50)"`
	MotherLifeStatus *string `gorm:"type:varchar(50)"`
	MotherCancerType *string `gorm:"type:varchar(50)"`
	MotherCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics
	FatherCancer     bool    `gorm:"not null;default:false"`
	FatherName       *string `gorm:"type:varchar(50)"`
	FatherLifeStatus *string `gorm:"type:varchar(50)"`
	FatherCancerType *string `gorm:"type:varchar(50)"`
	FatherCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics
	SiblingCancer     bool    `gorm:"not null;default:false"`
	SiblingName       *string `gorm:"type:varchar(50)"`
	SiblingLifeStatus *string `gorm:"type:varchar(50)"`
	SiblingCancerType *string `gorm:"type:varchar(50)"`
	SiblingCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics
	AmeAmoCancer     bool    `gorm:"not null;default:false"`
	AmeAmoName       *string `gorm:"type:varchar(50)"`
	AmeAmoLifeStatus *string `gorm:"type:varchar(50)"`
	AmeAmoCancerType *string `gorm:"type:varchar(50)"`
	AmeAmoCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics
	KhaleDaeiCancer     bool    `gorm:"not null;default:false"`
	KhaleDaeiName       *string `gorm:"type:varchar(50)"`
	KhaleDaeiLifeStatus *string `gorm:"type:varchar(50)"`
	KhaleDaeiCancerType *string `gorm:"type:varchar(50)"`
	KhaleDaeiCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics
	OtherRelativeCancer     *bool   `gorm:"type:boolean"`
	OtherRelativeName       *string `gorm:"type:varchar(50)"`
	OtherRelativeRelation   *string `gorm:"type:varchar(50)"`
	OtherRelativeLifeStatus *string `gorm:"type:varchar(50)"`
	OtherRelativeCancerType *string `gorm:"type:varchar(50)"`
	OtherRelativeCancerAge  *string `gorm:"type:varchar(50)"`
	// Pics

	TestGen *bool `gorm:"type:boolean"`
	// Pics
	FmTestGen *bool `gorm:"type:boolean"`
	// Pics
	CallExpert   bool    `gorm:"not null;default:true"`
	BirthCountry *string `gorm:"type:varchar(50)"`
	Province     *string `gorm:"type:varchar(50)"`
	City         *string `gorm:"type:varchar(50)"`
	Country      *string `gorm:"type:varchar(50)"`

	InsuranceStatus           string  `gorm:"type:varchar(50);not null"`
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
	OtherCancerType           *string `gorm:"type:varchar(50)"`
	LungCancerFamily          *bool   `gorm:"type:boolean"`
	LungCancerFamilyRelation  *string `gorm:"type:varchar(50)"`
	OtherCancerFamily         *bool   `gorm:"type:boolean"`
	OtherCancerFamilyType     *string `gorm:"type:varchar(50)"`
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
	PastSmoking               *string `gorm:"type:varchar(50)"`
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
