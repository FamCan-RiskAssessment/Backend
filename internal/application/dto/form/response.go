package formdto

import "time"

type FormResponse struct {
	ID                   uint    `json:"id"`
	UserID               uint    `json:"user_id"`
	Name                 string  `json:"name"`
	BirthDay             string  `json:"birth_day"`
	BirthMonth           string  `json:"birth_month"`
	BirthYear            string  `json:"birth_year"`
	Address              string  `json:"address"`
	PostalCode           string  `json:"postal_code"`
	SocialSecurityNumber string  `json:"social_security_number"`
	Gender               string  `json:"gender"`
	IsAtba               bool    `json:"is_atba"`
	Height               float64 `json:"height"`
	Weight               float64 `json:"weight"`

	DrinksAlcohol             *bool   `json:"drinks_alcohol"`
	CupsPerWeek               *string `json:"cups_per_week"`
	LastMonthSabzijatMeal     string  `json:"last_month_sabzijat_meal"`
	LastMonthSabzijatWeight   string  `json:"last_month_sabzijat_weight"`
	MediumActivityMonthInYear uint    `json:"medium_activity_month_in_year"`
	MediumActivityHourInWeek  string  `json:"medium_activity_hour_in_week"`
	HardActivityMonthInYear   uint    `json:"hard_activity_month_in_year"`
	HardActivityHourInWeek    string  `json:"hard_activity_hour_in_week"`
	SmokeAtLeast100           *bool   `json:"smoke_at_least_100"`
	SmokingAge                uint    `json:"smoking_age"`
	SmokingNow                bool    `json:"smoking_now"`
	LeaveSmokingAge           *uint   `json:"leave_smoking_age"`
	CountSmokingDaily         *uint   `json:"count_smoking_daily"`
	CountGheliandaily         *uint   `json:"count_gheliandaily"`
	CountSmokingDailyPast     *uint   `json:"count_smoking_daily_past"`
	CountGheliandailyPast     *uint   `json:"count_gheliandaily_past"`

	GhaedeAge                    uint    `json:"ghaede_age"`
	HasChildren                  bool    `json:"has_children"`
	NumberOfChildren             *uint   `json:"number_of_children"`
	AgeOfFirstBirth              *uint   `json:"age_of_first_birth"`
	MenopausalStatus             string  `json:"menopausal_status"`
	MenopauseAge                 string  `json:"menopause_age"`
	HRT                          *bool   `json:"hrt"`
	HRTUseLength                 *uint   `json:"hrt_use_length"`
	LastFiveYearsHRTUse          bool    `json:"last_five_years_hrt_use"`
	CurrentHRTUse                *bool   `json:"current_hrt_use"`
	IntendedHRTUse               *uint   `json:"intended_hrt_use"`
	HRTType                      *string `json:"hrt_type"`
	Oral                         *bool   `json:"oral"`
	OralDuration                 *string `json:"oral_duration"`
	OralTwoLastYears             *bool   `json:"oral_two_last_years"`
	MamoGraphy                   *bool   `json:"mamo_graphy"`
	Falop                        *bool   `json:"falop"`
	Andometrioz                  *bool   `json:"andometrioz"`
	LeavePestan                  bool    `json:"leave_pestan"`
	LeaveTokhmdan                bool    `json:"leave_tokhmdan"`
	LaDeColon                    *bool   `json:"la_de_colon"`
	LaDePol                      *bool   `json:"la_de_pol"`
	AspLaMo                      *bool   `json:"asp_la_mo"`
	NsaiDLaMo                    *bool   `json:"nsai_d_la_mo"`
	LastFiveYearBloodTestInStool *bool   `json:"last_five_year_blood_test_in_stool"`

	Cancer     bool    `json:"cancer"`
	CancerType *string `json:"cancer_type"`
	CancerAge  *uint   `json:"cancer_age"`

	ChildCancer     bool    `json:"child_cancer"`
	ChildName       *string `json:"child_name"`
	ChildCancerType *string `json:"child_cancer_type"`
	ChildCancerAge  *string `json:"child_cancer_age"`
	ChildLifeStatus *string `json:"child_life_status"`

	MotherCancer     bool    `json:"mother_cancer"`
	MotherName       *string `json:"mother_name"`
	MotherLifeStatus *string `json:"mother_life_status"`
	MotherCancerType *string `json:"mother_cancer_type"`
	MotherCancerAge  *string `json:"mother_cancer_age"`

	FatherCancer     bool    `json:"father_cancer"`
	FatherName       *string `json:"father_name"`
	FatherLifeStatus *string `json:"father_life_status"`
	FatherCancerType *string `json:"father_cancer_type"`
	FatherCancerAge  *string `json:"father_cancer_age"`

	SiblingCancer     bool    `json:"sibling_cancer"`
	SiblingName       *string `json:"sibling_name"`
	SiblingLifeStatus *string `json:"sibling_life_status"`
	SiblingCancerType *string `json:"sibling_cancer_type"`
	SiblingCancerAge  *string `json:"sibling_cancer_age"`

	AmeAmoCancer     bool    `json:"ame_amo_cancer"`
	AmeAmoName       *string `json:"ame_amo_name"`
	AmeAmoLifeStatus *string `json:"ame_amo_life_status"`
	AmeAmoCancerType *string `json:"ame_amo_cancer_type"`
	AmeAmoCancerAge  *string `json:"ame_amo_cancer_age"`

	KhaleDaeiCancer     bool    `json:"khale_daei_cancer"`
	KhaleDaeiName       *string `json:"khale_daei_name"`
	KhaleDaeiLifeStatus *string `json:"khale_daei_life_status"`
	KhaleDaeiCancerType *string `json:"khale_daei_cancer_type"`
	KhaleDaeiCancerAge  *string `json:"khale_daei_cancer_age"`

	OtherRelativeCancer     *bool   `json:"other_relative_cancer"`
	OtherRelativeName       *string `json:"other_relative_name"`
	OtherRelativeRelation   *string `json:"other_relative_relation"`
	OtherRelativeLifeStatus *string `json:"other_relative_life_status"`
	OtherRelativeCancerType *string `json:"other_relative_cancer_type"`
	OtherRelativeCancerAge  *string `json:"other_relative_cancer_age"`

	TestGen   *bool `json:"test_gen"`
	FmTestGen *bool `json:"fm_test_gen"`

	CallExpert   bool    `json:"call_expert"`
	BirthCountry *string `json:"birth_country"`
	Province     *string `json:"province"`
	City         *string `json:"city"`
	Country      *string `json:"country"`

	InsuranceStatus           string  `json:"insurance_status"`
	SupplementaryInsurances   *string `json:"supplementary_insurances"`
	Hypertension              bool    `json:"hypertension"`
	HypertensionTreatment     *bool   `json:"hypertension_treatment"`
	HeartDisease              bool    `json:"heart_disease"`
	HeartDiseaseTreatment     *bool   `json:"heart_disease_treatment"`
	Diabetes                  bool    `json:"diabetes"`
	DiabetesTreatment         *bool   `json:"diabetes_treatment"`
	ChronicLungDisease        *bool   `json:"chronic_lung_disease"`
	ChronicLungDiseaseType    *string `json:"chronic_lung_disease_type"`
	LungCancerHistory         bool    `json:"lung_cancer_history"`
	OtherCancerHistory        bool    `json:"other_cancer_history"`
	OtherCancerType           *string `json:"other_cancer_type"`
	LungCancerFamily          *bool   `json:"lung_cancer_family"`
	LungCancerFamilyRelation  *string `json:"lung_cancer_family_relation"`
	OtherCancerFamily         *bool   `json:"other_cancer_family"`
	OtherCancerFamilyType     *string `json:"other_cancer_family_type"`
	OtherCancerFamilyRelation *string `json:"other_cancer_family_relation"`
	OccupationalExposure      *string `json:"occupational_exposure"`
	CurrentSmoking            bool    `json:"current_smoking"`
	SmokingStartAgeCurrent    *uint   `json:"smoking_start_age_current"`
	SmokingTypesCurrent       *string `json:"smoking_types_current"`
	CigarettesPerDayCurrent   *uint   `json:"cigarettes_per_day_current"`
	CigarPerDayCurrent        *uint   `json:"cigar_per_day_current"`
	ECigPerDayCurrent         *uint   `json:"e_cig_per_day_current"`
	PipePerDayCurrent         *uint   `json:"pipe_per_day_current"`
	ChapoghPerDayCurrent      *uint   `json:"chapogh_per_day_current"`
	SmokedOpiumPerDayCurrent  *uint   `json:"smoked_opium_per_day_current"`
	ChewedOpiumPerDayCurrent  *uint   `json:"chewed_opium_per_day_current"`
	HookahPerWeekCurrent      *uint   `json:"hookah_per_week_current"`
	PastSmoking               *string `json:"past_smoking"`
	SmokingStartAgePast       *uint   `json:"smoking_start_age_past"`
	SmokingTypesPast          *string `json:"smoking_types_past"`
	CigarettesPerDayPast      *uint   `json:"cigarettes_per_day_past"`
	CigarPerDayPast           *uint   `json:"cigar_per_day_past"`
	ECigPerDayPast            *uint   `json:"e_cig_per_day_past"`
	PipePerDayPast            *uint   `json:"pipe_per_day_past"`
	ChapoghPerDayPast         *uint   `json:"chapogh_per_day_past"`
	SmokedOpiumPerDayPast     *uint   `json:"smoked_opium_per_day_past"`
	ChewedOpiumPerDayPast     *uint   `json:"chewed_opium_per_day_past"`
	HookahPerWeekPast         *uint   `json:"hookah_per_week_past"`
	SecondhandSmoke           bool    `json:"secondhand_smoke"`
	SecondhandSmokeLocation   *string `json:"secondhand_smoke_location"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
