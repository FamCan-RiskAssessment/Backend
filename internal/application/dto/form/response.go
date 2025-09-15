package formdto

import "time"

type FormResponse struct {
	ID                   uint    `json:"id"`
	UserID               uint    `json:"user_id"`
	Name                 string  `json:"name"`
	BirthDay             uint    `json:"birthDay"`
	BirthMonth           string  `json:"birthMonth"`
	BirthYear            uint    `json:"birthYear"`
	Address              string  `json:"address"`
	PostalCode           uint    `json:"postalCode"`
	SocialSecurityNumber uint    `gorm:"not null"`
	Gender               string  `json:"gender"`
	IsAtba               bool    `json:"isAtba"`
	Height               float64 `json:"height"`
	Weight               float64 `json:"weight"`

	DrinksAlcohol             *bool   `json:"drinksAlcohol"`
	CupsPerWeek               *string `json:"cupsPerWeek"`
	LastMonthSabzijatMeal     string  `json:"lastMonthSabzijatMeal"`
	LastMonthSabzijatWeight   string  `json:"lastMonthSabzijatWeight"`
	MediumActivityMonthInYear uint    `json:"mediumActivityMonthInYear"`
	MediumActivityHourInWeek  string  `json:"mediumActivityHourInWeek"`
	HardActivityMonthInYear   uint    `json:"hardActivityMonthInYear"`
	HardActivityHourInWeek    string  `json:"hardActivityHourInWeek"`
	SmokeAtLeast100           *bool   `json:"smokeAtLeast100"`
	SmokingAge                *uint   `json:"smokingAge"`
	SmokingNow                bool    `json:"smokingNow"`
	LeaveSmokingAge           *uint   `json:"leaveSmokingAge"`
	CountSmokingDaily         *string `json:"countSmokingDaily"`
	CountGheliandaily         *string `json:"countGhelianDaily"`
	CountSmokingDailyPast     *string `json:"countSmokingDailyPast"`
	CountGheliandailyPast     *string `json:"countGhelianDailyPast"`

	GhaedeAge                    uint    `json:"ghaedeAge"`
	HasChildren                  bool    `json:"hasChildren"`
	NumberOfChildren             *uint   `json:"numberOfChildren"`
	AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth"`
	MenopausalStatus             string  `json:"menopausalStatus"`
	MenopauseAge                 *string `json:"menopauseAge"`
	HRT                          *bool   `json:"hrt"`
	HRTUseLength                 *uint   `json:"hrtUseLength"`
	LastFiveYearsHRTUse          bool    `json:"lastFiveYearsHrtUse"`
	CurrentHRTUse                *bool   `json:"currentHrtUse"`
	IntendedHRTUse               *uint   `json:"intendedHrtUse"`
	HRTType                      *string `json:"hrtType"`
	Oral                         *bool   `json:"oral"`
	OralDuration                 *string `json:"oralDuration"`
	OralTwoLastYears             *bool   `json:"oralTwoLastYears"`
	MamoGraphy                   *bool   `json:"mamoGraphy"`
	Falop                        *bool   `json:"falop"`
	Andometrioz                  *bool   `json:"andometrioz"`
	LeavePestan                  bool    `json:"leavePestan"`
	LeaveTokhmdan                bool    `json:"leaveTokhmdan"`
	LaDeColon                    *bool   `json:"laDeColon"`
	LaDePol                      *bool   `json:"laDePol"`
	AspLaMo                      *bool   `json:"aspLaMo"`
	NsaiDLaMo                    *bool   `json:"nsaiDLaMo"`
	LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool"`

	Cancer     bool    `json:"cancer"`
	CancerType *string `json:"cancerType"`
	CancerAge  *uint   `json:"cancerAge"`

	ChildCancer     bool    `json:"childCancer"`
	ChildName       *string `json:"childName"`
	ChildCancerType *string `json:"childCancerType"`
	ChildCancerAge  *uint   `json:"childCancerAge"`
	ChildLifeStatus *string `json:"childLifeStatus"`

	MotherCancer     bool    `json:"motherCancer"`
	MotherName       *string `json:"motherName"`
	MotherLifeStatus *string `json:"motherLifeStatus"`
	MotherCancerType *string `json:"motherCancerType"`
	MotherCancerAge  *uint   `json:"motherCancerAge"`

	FatherCancer     bool    `json:"fatherCancer"`
	FatherName       *string `json:"fatherName"`
	FatherLifeStatus *string `json:"fatherLifeStatus"`
	FatherCancerType *string `json:"fatherCancerType"`
	FatherCancerAge  *uint   `json:"fatherCancerAge"`

	SiblingCancer     bool    `json:"siblingCancer"`
	SiblingName       *string `json:"siblingName"`
	SiblingLifeStatus *string `json:"siblingLifeStatus"`
	SiblingCancerType *string `json:"siblingCancerType"`
	SiblingCancerAge  *uint   `json:"siblingCancerAge"`

	AmeAmoCancer     bool    `json:"ameAmoCancer"`
	AmeAmoName       *string `json:"ameAmoName"`
	AmeAmoLifeStatus *string `json:"ameAmoLifeStatus"`
	AmeAmoCancerType *string `json:"ameAmoCancerType"`
	AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge"`

	KhaleDaeiCancer     bool    `json:"khaleDaeiCancer"`
	KhaleDaeiName       *string `json:"khaleDaeiName"`
	KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus"`
	KhaleDaeiCancerType *string `json:"khaleDaeiCancerType"`
	KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge"`

	OtherRelativeCancer     *bool   `json:"otherRelativeCancer"`
	OtherRelativeName       *string `json:"otherRelativeName"`
	OtherRelativeRelation   *string `json:"otherRelativeRelation"`
	OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus"`
	OtherRelativeCancerType *string `json:"otherRelativeCancerType"`
	OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge"`

	TestGen   *bool `json:"testGen"`
	FmTestGen *bool `json:"fmTestGen"`

	CallExpert   bool    `json:"callExpert"`
	BirthCountry *string `json:"birthCountry"`
	Province     *string `json:"province"`
	City         *string `json:"city"`
	Country      *string `json:"country"`

	InsuranceStatus           *string `json:"insuranceStatus"`
	SupplementaryInsurances   *string `json:"supplementaryInsurances"`
	Hypertension              bool    `json:"hypertension"`
	HypertensionTreatment     *bool   `json:"hypertensionTreatment"`
	HeartDisease              bool    `json:"heartDisease"`
	HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment"`
	Diabetes                  bool    `json:"diabetes"`
	DiabetesTreatment         *bool   `json:"diabetesTreatment"`
	ChronicLungDisease        *bool   `json:"chronicLungDisease"`
	ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType"`
	LungCancerHistory         bool    `json:"lungCancerHistory"`
	OtherCancerHistory        bool    `json:"otherCancerHistory"`
	OtherCancerType           *string `json:"otherCancerType"`
	LungCancerFamily          *bool   `json:"lungCancerFamily"`
	LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation"`
	OtherCancerFamily         *bool   `json:"otherCancerFamily"`
	OtherCancerFamilyType     *string `json:"otherCancerFamilyType"`
	OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation"`
	OccupationalExposure      *string `json:"occupationalExposure"`
	CurrentSmoking            bool    `json:"currentSmoking"`
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
	SecondhandSmoke           bool    `json:"secondhandSmoke"`
	SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateFormResponse struct {
	Form    FormResponse `json:"form"`
	Message string       `json:"message"`
}

type UpdateFormResponse struct {
	Form    FormResponse `json:"form"`
	Message string       `json:"message"`
}

type GetUserFormsResponse struct {
	Forms []FormResponse `json:"forms"`
	Total int            `json:"total"`
}

type DeleteFormResponse struct {
	Message string `json:"message"`
}
