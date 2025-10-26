package formdto

type CreateBasicFormRequest struct {
	UserID uint

	// page 1
	BirthDay             uint    `json:"birthDay" binding:"required"`
	BirthMonth           string  `json:"birthMonth" binding:"required"`
	BirthYear            uint    `json:"birthYear" binding:"required"`
	SocialSecurityNumber string  `json:"socialSecurityNumber" binding:"required"`
	Gender               uint    `json:"gender" binding:"required"`
	IsAtba               bool    `json:"isAtba"`
	Height               float64 `json:"height" binding:"required"`
	Weight               float64 `json:"weight" binding:"required"`
}

type UpdateBasicFormRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	// page 1
	BirthDay             *uint    `json:"birthDay"`
	BirthMonth           *string  `json:"birthMonth"`
	BirthYear            *uint    `json:"birthYear"`
	SocialSecurityNumber *string  `json:"socialSecurityNumber"`
	Gender               *uint    `json:"gender"`
	IsAtba               *bool    `json:"isAtba"`
	Height               *float64 `json:"height"`
	Weight               *float64 `json:"weight"`
}

type GetUserFormsRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Offset int  `json:"offset"`
	Limit  int  `json:"limit"`
}

type GetFormRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`
}

type GetPartialFormRequest struct {
	UserID uint
	FormID uint
}

type DeleteFormRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`
}

type UpsertGeneralHealthRequest struct {
	UserID uint
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
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	GhaedeAge                    uint    `json:"ghaedeAge" binding:"required"`
	HasChildren                  bool    `json:"hasChildren"`
	NumberOfChildren             *uint   `json:"numberOfChildren,omitempty"`
	AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth,omitempty"`
	MenopausalStatus             uint    `json:"menopausalStatus" binding:"required"`
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
	NumberOfBreastBiopsies       *uint   `json:"numberOfBreastBiopsies,omitempty"`
	HyperplasiaInBiopsy          *uint   `json:"hyperplasiaInBiopsy,omitempty"`
}

type UpsertCancerRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	Cancer     bool  `json:"cancer"`
	CancerType *uint `json:"cancerType,omitempty"`
	CancerAge  *uint `json:"cancerAge,omitempty"`
}

type UpsertFamilyCancerRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	ChildCancer     bool    `json:"childCancer"`
	ChildName       *string `json:"childName,omitempty"`
	ChildCancerType *uint   `json:"childCancerType,omitempty"`
	ChildCancerAge  *uint   `json:"childCancerAge,omitempty"`
	ChildLifeStatus *string `json:"childLifeStatus,omitempty"`

	MotherCancer     bool    `json:"motherCancer"`
	MotherName       *string `json:"motherName,omitempty"`
	MotherLifeStatus *string `json:"motherLifeStatus,omitempty"`
	MotherCancerType *uint   `json:"motherCancerType,omitempty"`
	MotherCancerAge  *uint   `json:"motherCancerAge,omitempty"`

	FatherCancer     bool    `json:"fatherCancer"`
	FatherName       *string `json:"fatherName,omitempty"`
	FatherLifeStatus *string `json:"fatherLifeStatus,omitempty"`
	FatherCancerType *uint   `json:"fatherCancerType,omitempty"`
	FatherCancerAge  *uint   `json:"fatherCancerAge,omitempty"`

	SiblingCancer     bool    `json:"siblingCancer"`
	SiblingName       *string `json:"siblingName,omitempty"`
	SiblingLifeStatus *string `json:"siblingLifeStatus,omitempty"`
	SiblingCancerType *uint   `json:"siblingCancerType,omitempty"`
	SiblingCancerAge  *uint   `json:"siblingCancerAge,omitempty"`

	AmeAmoCancer     bool    `json:"ameAmoCancer"`
	AmeAmoName       *string `json:"ameAmoName,omitempty"`
	AmeAmoLifeStatus *string `json:"ameAmoLifeStatus,omitempty"`
	AmeAmoCancerType *uint   `json:"ameAmoCancerType,omitempty"`
	AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge,omitempty"`

	KhaleDaeiCancer     bool    `json:"khaleDaeiCancer"`
	KhaleDaeiName       *string `json:"khaleDaeiName,omitempty"`
	KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus,omitempty"`
	KhaleDaeiCancerType *uint   `json:"khaleDaeiCancerType,omitempty"`
	KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge,omitempty"`

	OtherRelativeCancer     *bool   `json:"otherRelativeCancer,omitempty"`
	OtherRelativeName       *string `json:"otherRelativeName,omitempty"`
	OtherRelativeRelation   *string `json:"otherRelativeRelation,omitempty"`
	OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus,omitempty"`
	OtherRelativeCancerType *uint   `json:"otherRelativeCancerType,omitempty"`
	OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge,omitempty"`
}

type UpsertContactRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	Name         string  `json:"name" binding:"required"`
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
	UserID uint
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
	OtherCancerType           *uint   `json:"otherCancerType,omitempty"`
	LungCancerFamily          *bool   `json:"lungCancerFamily,omitempty"`
	LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation,omitempty"`
	OtherCancerFamily         *bool   `json:"otherCancerFamily,omitempty"`
	OtherCancerFamilyType     *uint   `json:"otherCancerFamilyType,omitempty"`
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
	UserID uint
	FormID uint `json:"form_id" binding:"required"`
}

type AssignOperatorRequest struct {
	UserID     uint
	FormID     uint `json:"form_id" binding:"required"`
	OperatorID uint `json:"operator_id" binding:"required"`
}

type UnassignOperatorRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`
}

type UpdateGeneralHealthRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	DrinksAlcohol             *bool   `json:"drinksAlcohol"`
	CupsPerWeek               *string `json:"cupsPerWeek"`
	LastMonthSabzijatMeal     *string `json:"lastMonthSabzijatMeal"`
	LastMonthSabzijatWeight   *string `json:"lastMonthSabzijatWeight"`
	MediumActivityMonthInYear *uint   `json:"mediumActivityMonthInYear"`
	MediumActivityHourInWeek  *string `json:"mediumActivityHourInWeek"`
	HardActivityMonthInYear   *uint   `json:"hardActivityMonthInYear"`
	HardActivityHourInWeek    *string `json:"hardActivityHourInWeek"`
	SmokeAtLeast100           *bool   `json:"smokeAtLeast100"`
	SmokingAge                *uint   `json:"smokingAge"`
	SmokingNow                *bool   `json:"smokingNow"`
	LeaveSmokingAge           *uint   `json:"leaveSmokingAge"`
	CountSmokingDaily         *string `json:"countSmokingDaily"`
	CountGheliandaily         *string `json:"countGheliandaily"`
	CountSmokingDailyPast     *string `json:"countSmokingDailyPast"`
	CountGheliandailyPast     *string `json:"countGheliandailyPast"`
}
type UpdateMamographyRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	GhaedeAge                    *uint   `json:"ghaedeAge"`
	HasChildren                  *bool   `json:"hasChildren"`
	NumberOfChildren             *uint   `json:"numberOfChildren"`
	AgeOfFirstBirth              *uint   `json:"ageOfFirstBirth"`
	MenopausalStatus             *uint   `json:"menopausalStatus"`
	MenopauseAge                 *string `json:"menopauseAge"`
	HRT                          *bool   `json:"hrt"`
	HRTUseLength                 *uint   `json:"hrtUseLength"`
	LastFiveYearsHRTUse          *bool   `json:"lastFiveYearsHrtUse"`
	CurrentHRTUse                *bool   `json:"currentHrtUse"`
	IntendedHRTUse               *uint   `json:"intendedHrtUse"`
	HRTType                      *string `json:"hrtType"`
	Oral                         *bool   `json:"oral"`
	OralDuration                 *string `json:"oralDuration"`
	OralTwoLastYears             *bool   `json:"oralTwoLastYears"`
	MamoGraphy                   *bool   `json:"mamoGraphy"`
	Falop                        *bool   `json:"falop"`
	Andometrioz                  *bool   `json:"andometrioz"`
	LeavePestan                  *bool   `json:"leavePestan"`
	LeaveTokhmdan                *bool   `json:"leaveTokhmdan"`
	LaDeColon                    *bool   `json:"laDeColon"`
	LaDePol                      *bool   `json:"laDePol"`
	AspLaMo                      *bool   `json:"aspLaMo"`
	NsaiDLaMo                    *bool   `json:"nsaiDLaMo"`
	LastFiveYearBloodTestInStool *bool   `json:"lastFiveYearBloodTestInStool"`
	NumberOfBreastBiopsies       *uint   `json:"numberOfBreastBiopsies,omitempty"`
	HyperplasiaInBiopsy          *uint   `json:"hyperplasiaInBiopsy,omitempty"`
}
type UpdateCancerRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	Cancer     *bool `json:"cancer"`
	CancerType *uint `json:"cancerType"`
	CancerAge  *uint `json:"cancerAge"`
}
type UpdateFamilyCancerRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	ChildCancer     *bool   `json:"childCancer"`
	ChildName       *string `json:"childName"`
	ChildCancerType *uint   `json:"childCancerType"`
	ChildCancerAge  *uint   `json:"childCancerAge"`
	ChildLifeStatus *string `json:"childLifeStatus"`

	MotherCancer     *bool   `json:"motherCancer"`
	MotherName       *string `json:"motherName"`
	MotherLifeStatus *string `json:"motherLifeStatus"`
	MotherCancerType *uint   `json:"motherCancerType"`
	MotherCancerAge  *uint   `json:"motherCancerAge"`

	FatherCancer     *bool   `json:"fatherCancer"`
	FatherName       *string `json:"fatherName"`
	FatherLifeStatus *string `json:"fatherLifeStatus"`
	FatherCancerType *uint   `json:"fatherCancerType"`
	FatherCancerAge  *uint   `json:"fatherCancerAge"`

	SiblingCancer     *bool   `json:"siblingCancer"`
	SiblingName       *string `json:"siblingName"`
	SiblingLifeStatus *string `json:"siblingLifeStatus"`
	SiblingCancerType *uint   `json:"siblingCancerType"`
	SiblingCancerAge  *uint   `json:"siblingCancerAge"`

	AmeAmoCancer     *bool   `json:"ameAmoCancer"`
	AmeAmoName       *string `json:"ameAmoName"`
	AmeAmoLifeStatus *string `json:"ameAmoLifeStatus"`
	AmeAmoCancerType *uint   `json:"ameAmoCancerType"`
	AmeAmoCancerAge  *uint   `json:"ameAmoCancerAge"`

	KhaleDaeiCancer     *bool   `json:"khaleDaeiCancer"`
	KhaleDaeiName       *string `json:"khaleDaeiName"`
	KhaleDaeiLifeStatus *string `json:"khaleDaeiLifeStatus"`
	KhaleDaeiCancerType *uint   `json:"khaleDaeiCancerType"`
	KhaleDaeiCancerAge  *uint   `json:"khaleDaeiCancerAge"`

	OtherRelativeCancer     *bool   `json:"otherRelativeCancer"`
	OtherRelativeName       *string `json:"otherRelativeName"`
	OtherRelativeRelation   *string `json:"otherRelativeRelation"`
	OtherRelativeLifeStatus *string `json:"otherRelativeLifeStatus"`
	OtherRelativeCancerType *uint   `json:"otherRelativeCancerType"`
	OtherRelativeCancerAge  *uint   `json:"otherRelativeCancerAge"`
}
type UpdateContactRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	Name         *string `json:"name"`
	TestGen      *bool   `json:"testGen"`
	FmTestGen    *bool   `json:"fmTestGen"`
	CallExpert   *bool   `json:"callExpert"`
	BirthCountry *string `json:"birthCountry"`
	Province     *string `json:"province"`
	City         *string `json:"city"`
	Country      *string `json:"country"`
	Address      *string `json:"address"`
	PostalCode   *string `json:"postalCode"`
}
type UpdateLungCancerRequest struct {
	UserID uint
	FormID uint `json:"form_id" binding:"required"`

	InsuranceStatus           *string `json:"insuranceStatus"`
	SupplementaryInsurances   *string `json:"supplementaryInsurances"`
	Hypertension              *bool   `json:"hypertension"`
	HypertensionTreatment     *bool   `json:"hypertensionTreatment"`
	HeartDisease              *bool   `json:"heartDisease"`
	HeartDiseaseTreatment     *bool   `json:"heartDiseaseTreatment"`
	Diabetes                  *bool   `json:"diabetes"`
	DiabetesTreatment         *bool   `json:"diabetesTreatment"`
	ChronicLungDisease        *bool   `json:"chronicLungDisease"`
	ChronicLungDiseaseType    *string `json:"chronicLungDiseaseType"`
	LungCancerHistory         *bool   `json:"lungCancerHistory"`
	OtherCancerHistory        *bool   `json:"otherCancerHistory"`
	OtherCancerType           *uint   `json:"otherCancerType"`
	LungCancerFamily          *bool   `json:"lungCancerFamily"`
	LungCancerFamilyRelation  *string `json:"lungCancerFamilyRelation"`
	OtherCancerFamily         *bool   `json:"otherCancerFamily"`
	OtherCancerFamilyType     *uint   `json:"otherCancerFamilyType"`
	OtherCancerFamilyRelation *string `json:"otherCancerFamilyRelation"`
	OccupationalExposure      *string `json:"occupationalExposure"`
	CurrentSmoking            *bool   `json:"currentSmoking"`
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
	SecondhandSmoke           *bool   `json:"secondhandSmoke"`
	SecondhandSmokeLocation   *string `json:"secondhandSmokeLocation"`
}
