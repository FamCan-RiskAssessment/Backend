package formdto

type CreateBasicFormRequest struct {
	UserID uint

	// page 1
	Name                 string  `json:"name" binding:"required"`
	BirthDay             uint    `json:"birthDay" binding:"required"`
	BirthMonth           string  `json:"birthMonth" binding:"required"`
	BirthYear            uint    `json:"birthYear" binding:"required"`
	SocialSecurityNumber string  `json:"socialSecurityNumber" binding:"required"`
	Gender               string  `json:"gender" binding:"required"`
	IsAtba               bool    `json:"isAtba"`
	Height               float64 `json:"height" binding:"required"`
	Weight               float64 `json:"weight" binding:"required"`
}

type UpdateBasicFormRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	// page 1
	Name                 *string  `json:"name,omitempty"`
	BirthDay             *uint    `json:"birthDay,omitempty"`
	BirthMonth           *string  `json:"birthMonth,omitempty"`
	BirthYear            *uint    `json:"birthYear,omitempty"`
	SocialSecurityNumber *string  `json:"socialSecurityNumber,omitempty"`
	Gender               *string  `json:"gender,omitempty"`
	IsAtba               *bool    `json:"isAtba,omitempty"`
	Height               *float64 `json:"height,omitempty"`
	Weight               *float64 `json:"weight,omitempty"`
}

type GetUserFormsRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Offset int  `json:"offset"`
	Limit  int  `json:"limit"`
}

type GetFormRequest struct {
	FormID uint `json:"form_id" binding:"required"`
}

type DeleteFormRequest struct {
	FormID uint `json:"form_id" binding:"required"`
}

type UpsertGeneralHealthRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	DrinksAlcohol             *bool   `json:"drinksAlcohol,omitempty"`
	CupsPerWeek               *string `json:"cupsPerWeek,omitempty"`
	LastMonthSabzijatMeal     string  `json:"lastMonthSabzijatMeal" binding:"required"`
	LastMonthSabzijatWeight   string  `json:"lastMonthSabzijatWeight" binding:"required"`
	MediumActivityMonthInYear uint    `json:"mediumActivityMonthInYear" binding:"required"`
	MediumActivityHourInWeek  string  `json:"mediumActivityHourInWeek" binding:"required"`
	HardActivityMonthInYear   uint    `json:"hardActivityMonthInYear" binding:"required"`
	HardActivityHourInWeek    string  `json:"hardActivityHourInWeek" binding:"required"`
	SmokeAtLeast100           *bool   `json:"smokeAtLeast100,omitempty"`
	SmokingAge                *uint   `json:"smokingAge,omitempty"`
	SmokingNow                bool    `json:"smokingNow"`
	LeaveSmokingAge           *uint   `json:"leaveSmokingAge,omitempty"`
	CountSmokingDaily         *string `json:"countSmokingDaily,omitempty"`
	CountGheliandaily         *string `json:"countGheliandaily,omitempty"`
	CountSmokingDailyPast     *string `json:"countSmokingDailyPast,omitempty"`
	CountGheliandailyPast     *string `json:"countGheliandailyPast,omitempty"`
}

type UpsertMamographyRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	GhaedeAge                    uint    `json:"ghaedeAge" binding:"required"`
	HasChildren                  bool    `json:"hasChildren"`
	NumberOfChildren             *uint   `json:"numberOfChildren,omitempty"`
	AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth,omitempty"`
	MenopausalStatus             string  `json:"menopausalStatus" binding:"required"`
	MenopauseAge                 *string `json:"menopauseAge,omitempty"`
	HRT                          *bool   `json:"hrt,omitempty"`
	HRTUseLength                 *uint   `json:"hrtUseLength,omitempty"`
	LastFiveYearsHRTUse          bool    `json:"lastFiveYearsHrtUse"`
	CurrentHRTUse                *bool   `json:"currentHrtUse,omitempty"`
	IntendedHRTUse               *uint   `json:"intendedHrtUse,omitempty"`
	HRTType                      *string `json:"hrtType,omitempty"`
	Oral                         *bool   `json:"oral,omitempty"`
	OralDuration                 *string `json:"oralDuration,omitempty"`
	OralTwoLastYears             *bool   `json:"oralTwoLastYears,omitempty"`
	MamoGraphy                   *bool   `json:"mamoGraphy,omitempty"`
	Falop                        *bool   `json:"falop,omitempty"`
	Andometrioz                  *bool   `json:"andometrioz,omitempty"`
	LeavePestan                  bool    `json:"leavePestan"`
	LeaveTokhmdan                bool    `json:"leaveTokhmdan"`
	LaDeColon                    *bool   `json:"laDeColon,omitempty"`
	LaDePol                      *bool   `json:"laDePol,omitempty"`
	AspLaMo                      *bool   `json:"aspLaMo,omitempty"`
	NsaiDLaMo                    *bool   `json:"nsaiDLaMo,omitempty"`
	LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool,omitempty"`
}

type UpsertCancerRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	Cancer     bool    `json:"cancer"`
	CancerType *string `json:"cancerType,omitempty"`
	CancerAge  *uint   `json:"cancerAge,omitempty"`
}

type UpsertFamilyCancerRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	ChildCancer     bool    `json:"childCancer"`
	ChildName       *string `json:"childName,omitempty"`
	ChildCancerType *string `json:"childCancerType,omitempty"`
	ChildCancerAge  *uint   `json:"childCancerAge,omitempty"`
	ChildLifeStatus *string `json:"childLifeStatus,omitempty"`

	MotherCancer     bool    `json:"motherCancer"`
	MotherName       *string `json:"motherName,omitempty"`
	MotherLifeStatus *string `json:"motherLifeStatus,omitempty"`
	MotherCancerType *string `json:"motherCancerType,omitempty"`
	MotherCancerAge  *uint   `json:"motherCancerAge,omitempty"`

	FatherCancer     bool    `json:"fatherCancer"`
	FatherName       *string `json:"fatherName,omitempty"`
	FatherLifeStatus *string `json:"fatherLifeStatus,omitempty"`
	FatherCancerType *string `json:"fatherCancerType,omitempty"`
	FatherCancerAge  *uint   `json:"fatherCancerAge,omitempty"`

	SiblingCancer     bool    `json:"siblingCancer"`
	SiblingName       *string `json:"siblingName,omitempty"`
	SiblingLifeStatus *string `json:"siblingLifeStatus,omitempty"`
	SiblingCancerType *string `json:"siblingCancerType,omitempty"`
	SiblingCancerAge  *uint   `json:"siblingCancerAge,omitempty"`

	AmeAmoCancer     bool    `json:"ameAmoCancer"`
	AmeAmoName       *string `json:"ameAmoName,omitempty"`
	AmeAmoLifeStatus *string `json:"ameAmoLifeStatus,omitempty"`
	AmeAmoCancerType *string `json:"ameAmoCancerType,omitempty"`
	AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge,omitempty"`

	KhaleDaeiCancer     bool    `json:"khaleDaeiCancer"`
	KhaleDaeiName       *string `json:"khaleDaeiName,omitempty"`
	KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus,omitempty"`
	KhaleDaeiCancerType *string `json:"khaleDaeiCancerType,omitempty"`
	KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge,omitempty"`

	OtherRelativeCancer     *bool   `json:"otherRelativeCancer,omitempty"`
	OtherRelativeName       *string `json:"otherRelativeName,omitempty"`
	OtherRelativeRelation   *string `json:"otherRelativeRelation,omitempty"`
	OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus,omitempty"`
	OtherRelativeCancerType *string `json:"otherRelativeCancerType,omitempty"`
	OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge,omitempty"`
}

type UpsertContactRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	TestGen      *bool   `json:"testGen,omitempty"`
	FmTestGen    *bool   `json:"fmTestGen,omitempty"`
	CallExpert   bool    `json:"callExpert"`
	BirthCountry *string `json:"birthCountry,omitempty"`
	Province     *string `json:"province,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`
	Address      string  `json:"address" binding:"required"`
	PostalCode   string  `json:"postalCode" binding:"required"`
}

type UpsertLungCancerRequest struct {
	FormID uint `json:"form_id" binding:"required"`

	InsuranceStatus           *string `json:"insuranceStatus,omitempty"`
	SupplementaryInsurances   *string `json:"supplementaryInsurances,omitempty"`
	Hypertension              bool    `json:"hypertension"`
	HypertensionTreatment     *bool   `json:"hypertensionTreatment,omitempty"`
	HeartDisease              bool    `json:"heartDisease"`
	HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment,omitempty"`
	Diabetes                  bool    `json:"diabetes"`
	DiabetesTreatment         *bool   `json:"diabetesTreatment,omitempty"`
	ChronicLungDisease        *bool   `json:"chronicLungDisease,omitempty"`
	ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType,omitempty"`
	LungCancerHistory         bool    `json:"lungCancerHistory"`
	OtherCancerHistory        bool    `json:"otherCancerHistory"`
	OtherCancerType           *string `json:"otherCancerType,omitempty"`
	LungCancerFamily          *bool   `json:"lungCancerFamily,omitempty"`
	LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation,omitempty"`
	OtherCancerFamily         *bool   `json:"otherCancerFamily,omitempty"`
	OtherCancerFamilyType     *string `json:"otherCancerFamilyType,omitempty"`
	OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation,omitempty"`
	OccupationalExposure      *string `json:"occupationalExposure,omitempty"`
	CurrentSmoking            bool    `json:"currentSmoking"`
	SmokingStartAgeCurrent    *uint   `json:"smokingStartAgeCurrent,omitempty"`
	SmokingTypesCurrent       *string `json:"smokingTypesCurrent,omitempty"`
	CigarettesPerDayCurrent   *uint   `json:"cigarettesPerDayCurrent,omitempty"`
	CigarPerDayCurrent        *uint   `json:"cigarPerDayCurrent,omitempty"`
	ECigPerDayCurrent         *uint   `json:"eCigPerDayCurrent,omitempty"`
	PipePerDayCurrent         *uint   `json:"pipePerDayCurrent,omitempty"`
	ChapoghPerDayCurrent      *uint   `json:"chapoghPerDayCurrent,omitempty"`
	SmokedOpiumPerDayCurrent  *uint   `json:"smokedOpiumPerDayCurrent,omitempty"`
	ChewedOpiumPerDayCurrent  *uint   `json:"chewedOpiumPerDayCurrent,omitempty"`
	HookahPerWeekCurrent      *uint   `json:"hookahPerWeekCurrent,omitempty"`
	PastSmoking               *string `json:"pastSmoking,omitempty"`
	SmokingStartAgePast       *uint   `json:"smokingStartAgePast,omitempty"`
	SmokingTypesPast          *string `json:"smokingTypesPast,omitempty"`
	CigarettesPerDayPast      *uint   `json:"cigarettesPerDayPast,omitempty"`
	CigarPerDayPast           *uint   `json:"cigarPerDayPast,omitempty"`
	ECigPerDayPast            *uint   `json:"eCigPerDayPast,omitempty"`
	PipePerDayPast            *uint   `json:"pipePerDayPast,omitempty"`
	ChapoghPerDayPast         *uint   `json:"chapoghPerDayPast,omitempty"`
	SmokedOpiumPerDayPast     *uint   `json:"smokedOpiumPerDayPast,omitempty"`
	ChewedOpiumPerDayPast     *uint   `json:"chewedOpiumPerDayPast,omitempty"`
	HookahPerWeekPast         *uint   `json:"hookahPerWeekPast,omitempty"`
	SecondhandSmoke           bool    `json:"secondhandSmoke"`
	SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation,omitempty"`
}

type ChangeFormStatusRequest struct {
	FormID uint `json:"form_id" binding:"required"`
}
