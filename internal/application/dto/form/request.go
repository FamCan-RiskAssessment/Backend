package formdto

import "time"

type CreateFormRequest struct {
	UserID               uint
	Name                 string    `json:"name" binding:"required"`
	DateOfBirth          time.Time `json:"date_of_birth" binding:"required"`
	Address              string    `json:"address" binding:"required"`
	PostalCode           string    `json:"postal_code" binding:"required"`
	SocialSecurityNumber string    `json:"social_security_number" binding:"required"`
	Gender               string    `json:"gender" binding:"required"`
	IsAtba               bool      `json:"is_atba"`
	Height               float64   `json:"height" binding:"required"`
	Weight               float64   `json:"weight" binding:"required"`

	DrinksAlcohol             *bool   `json:"drinks_alcohol,omitempty"`
	CupsPerWeek               *string `json:"cups_per_week,omitempty"`
	LastMonthSabzijatMeal     string  `json:"last_month_sabzijat_meal" binding:"required"`
	LastMonthSabzijatWeight   string  `json:"last_month_sabzijat_weight" binding:"required"`
	MediumActivityMonthInYear uint    `json:"medium_activity_month_in_year" binding:"required"`
	MediumActivityHourInWeek  string  `json:"medium_activity_hour_in_week" binding:"required"`
	HardActivityMonthInYear   uint    `json:"hard_activity_month_in_year" binding:"required"`
	HardActivityHourInWeek    string  `json:"hard_activity_hour_in_week" binding:"required"`
	SmokeAtLeast100           *bool   `json:"smoke_at_least_100,omitempty"`
	SmokingAge                uint    `json:"smoking_age" binding:"required"`
	SmokingNow                bool    `json:"smoking_now"`
	LeaveSmokingAge           *uint   `json:"leave_smoking_age,omitempty"`
	CountSmokingDaily         *uint   `json:"count_smoking_daily,omitempty"`
	CountGheliandaily         *uint   `json:"count_gheliandaily,omitempty"`
	CountSmokingDailyPast     *uint   `json:"count_smoking_daily_past,omitempty"`
	CountGheliandailyPast     *uint   `json:"count_gheliandaily_past,omitempty"`

	GhaedeAge                    uint    `json:"ghaede_age" binding:"required"`
	HasChildren                  bool    `json:"has_children"`
	NumberOfChildren             *uint   `json:"number_of_children,omitempty"`
	AgeOfFirstBirth              *uint   `json:"age_of_first_birth,omitempty"`
	MenopausalStatus             string  `json:"menopausal_status" binding:"required"`
	MenopauseAge                 string  `json:"menopause_age" binding:"required"`
	HRT                          *bool   `json:"hrt,omitempty"`
	HRTUseLength                 *uint   `json:"hrt_use_length,omitempty"`
	LastFiveYearsHRTUse          bool    `json:"last_five_years_hrt_use"`
	CurrentHRTUse                *bool   `json:"current_hrt_use,omitempty"`
	IntendedHRTUse               *uint   `json:"intended_hrt_use,omitempty"`
	HRTType                      *string `json:"hrt_type,omitempty"`
	Oral                         *bool   `json:"oral,omitempty"`
	OralDuration                 *string `json:"oral_duration,omitempty"`
	OralTwoLastYears             *bool   `json:"oral_two_last_years,omitempty"`
	MamoGraphy                   *bool   `json:"mamo_graphy,omitempty"`
	Falop                        *bool   `json:"falop,omitempty"`
	Andometrioz                  *bool   `json:"andometrioz,omitempty"`
	LeavePestan                  bool    `json:"leave_pestan"`
	LeaveTokhmdan                bool    `json:"leave_tokhmdan"`
	LaDeColon                    *bool   `json:"la_de_colon,omitempty"`
	LaDePol                      *bool   `json:"la_de_pol,omitempty"`
	AspLaMo                      *bool   `json:"asp_la_mo,omitempty"`
	NsaiDLaMo                    *bool   `json:"nsai_d_la_mo,omitempty"`
	LastFiveYearBloodTestInStool *bool   `json:"last_five_year_blood_test_in_stool,omitempty"`

	Cancer     bool    `json:"cancer"`
	CancerType *string `json:"cancer_type,omitempty"`
	CancerAge  *uint   `json:"cancer_age,omitempty"`

	ChildCancer     bool    `json:"child_cancer"`
	ChildName       *string `json:"child_name,omitempty"`
	ChildCancerType *string `json:"child_cancer_type,omitempty"`
	ChildCancerAge  *string `json:"child_cancer_age,omitempty"`
	ChildLifeStatus *string `json:"child_life_status,omitempty"`

	MotherCancer     bool    `json:"mother_cancer"`
	MotherName       *string `json:"mother_name,omitempty"`
	MotherLifeStatus *string `json:"mother_life_status,omitempty"`
	MotherCancerType *string `json:"mother_cancer_type,omitempty"`
	MotherCancerAge  *string `json:"mother_cancer_age,omitempty"`

	FatherCancer     bool    `json:"father_cancer"`
	FatherName       *string `json:"father_name,omitempty"`
	FatherLifeStatus *string `json:"father_life_status,omitempty"`
	FatherCancerType *string `json:"father_cancer_type,omitempty"`
	FatherCancerAge  *string `json:"father_cancer_age,omitempty"`

	SiblingCancer     bool    `json:"sibling_cancer"`
	SiblingName       *string `json:"sibling_name,omitempty"`
	SiblingLifeStatus *string `json:"sibling_life_status,omitempty"`
	SiblingCancerType *string `json:"sibling_cancer_type,omitempty"`
	SiblingCancerAge  *string `json:"sibling_cancer_age,omitempty"`

	AmeAmoCancer     bool    `json:"ame_amo_cancer"`
	AmeAmoName       *string `json:"ame_amo_name,omitempty"`
	AmeAmoLifeStatus *string `json:"ame_amo_life_status,omitempty"`
	AmeAmoCancerType *string `json:"ame_amo_cancer_type,omitempty"`
	AmeAmoCancerAge  *string `json:"ame_amo_cancer_age,omitempty"`

	KhaleDaeiCancer     bool    `json:"khale_daei_cancer"`
	KhaleDaeiName       *string `json:"khale_daei_name,omitempty"`
	KhaleDaeiLifeStatus *string `json:"khale_daei_life_status,omitempty"`
	KhaleDaeiCancerType *string `json:"khale_daei_cancer_type,omitempty"`
	KhaleDaeiCancerAge  *string `json:"khale_daei_cancer_age,omitempty"`

	OtherRelativeCancer     *bool   `json:"other_relative_cancer,omitempty"`
	OtherRelativeName       *string `json:"other_relative_name,omitempty"`
	OtherRelativeRelation   *string `json:"other_relative_relation,omitempty"`
	OtherRelativeLifeStatus *string `json:"other_relative_life_status,omitempty"`
	OtherRelativeCancerType *string `json:"other_relative_cancer_type,omitempty"`
	OtherRelativeCancerAge  *string `json:"other_relative_cancer_age,omitempty"`

	TestGen   *bool `json:"test_gen,omitempty"`
	FmTestGen *bool `json:"fm_test_gen,omitempty"`

	CallExpert   bool    `json:"call_expert"`
	BirthCountry *string `json:"birth_country,omitempty"`
	Province     *string `json:"province,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`

	InsuranceStatus           string  `json:"insurance_status" binding:"required"`
	SupplementaryInsurances   *string `json:"supplementary_insurances,omitempty"`
	Hypertension              bool    `json:"hypertension"`
	HypertensionTreatment     *bool   `json:"hypertension_treatment,omitempty"`
	HeartDisease              bool    `json:"heart_disease"`
	HeartDiseaseTreatment     *bool   `json:"heart_disease_treatment,omitempty"`
	Diabetes                  bool    `json:"diabetes"`
	DiabetesTreatment         *bool   `json:"diabetes_treatment,omitempty"`
	ChronicLungDisease        *bool   `json:"chronic_lung_disease,omitempty"`
	ChronicLungDiseaseType    *string `json:"chronic_lung_disease_type,omitempty"`
	LungCancerHistory         bool    `json:"lung_cancer_history"`
	OtherCancerHistory        bool    `json:"other_cancer_history"`
	OtherCancerType           *string `json:"other_cancer_type,omitempty"`
	LungCancerFamily          *bool   `json:"lung_cancer_family,omitempty"`
	LungCancerFamilyRelation  *string `json:"lung_cancer_family_relation,omitempty"`
	OtherCancerFamily         *bool   `json:"other_cancer_family,omitempty"`
	OtherCancerFamilyType     *string `json:"other_cancer_family_type,omitempty"`
	OtherCancerFamilyRelation *string `json:"other_cancer_family_relation,omitempty"`
	OccupationalExposure      *string `json:"occupational_exposure,omitempty"`
	CurrentSmoking            bool    `json:"current_smoking"`
	SmokingStartAgeCurrent    *uint   `json:"smoking_start_age_current,omitempty"`
	SmokingTypesCurrent       *string `json:"smoking_types_current,omitempty"`
	CigarettesPerDayCurrent   *uint   `json:"cigarettes_per_day_current,omitempty"`
	CigarPerDayCurrent        *uint   `json:"cigar_per_day_current,omitempty"`
	ECigPerDayCurrent         *uint   `json:"e_cig_per_day_current,omitempty"`
	PipePerDayCurrent         *uint   `json:"pipe_per_day_current,omitempty"`
	ChapoghPerDayCurrent      *uint   `json:"chapogh_per_day_current,omitempty"`
	SmokedOpiumPerDayCurrent  *uint   `json:"smoked_opium_per_day_current,omitempty"`
	ChewedOpiumPerDayCurrent  *uint   `json:"chewed_opium_per_day_current,omitempty"`
	HookahPerWeekCurrent      *uint   `json:"hookah_per_week_current,omitempty"`
	PastSmoking               *string `json:"past_smoking,omitempty"`
	SmokingStartAgePast       *uint   `json:"smoking_start_age_past,omitempty"`
	SmokingTypesPast          *string `json:"smoking_types_past,omitempty"`
	CigarettesPerDayPast      *uint   `json:"cigarettes_per_day_past,omitempty"`
	CigarPerDayPast           *uint   `json:"cigar_per_day_past,omitempty"`
	ECigPerDayPast            *uint   `json:"e_cig_per_day_past,omitempty"`
	PipePerDayPast            *uint   `json:"pipe_per_day_past,omitempty"`
	ChapoghPerDayPast         *uint   `json:"chapogh_per_day_past,omitempty"`
	SmokedOpiumPerDayPast     *uint   `json:"smoked_opium_per_day_past,omitempty"`
	ChewedOpiumPerDayPast     *uint   `json:"chewed_opium_per_day_past,omitempty"`
	HookahPerWeekPast         *uint   `json:"hookah_per_week_past,omitempty"`
	SecondhandSmoke           bool    `json:"secondhand_smoke"`
	SecondhandSmokeLocation   *string `json:"secondhand_smoke_location,omitempty"`
}

type UpdateFormRequest struct {
	FormID               uint       `json:"form_id" binding:"required"`
	Name                 *string    `json:"name,omitempty"`
	DateOfBirth          *time.Time `json:"date_of_birth,omitempty"`
	Address              *string    `json:"address,omitempty"`
	PostalCode           *string    `json:"postal_code,omitempty"`
	SocialSecurityNumber *string    `json:"social_security_number,omitempty"`
	Gender               *string    `json:"gender,omitempty"`
	IsAtba               *bool      `json:"is_atba,omitempty"`
	Height               *float64   `json:"height,omitempty"`
	Weight               *float64   `json:"weight,omitempty"`

	DrinksAlcohol             *bool   `json:"drinks_alcohol,omitempty"`
	CupsPerWeek               *string `json:"cups_per_week,omitempty"`
	LastMonthSabzijatMeal     *string `json:"last_month_sabzijat_meal,omitempty"`
	LastMonthSabzijatWeight   *string `json:"last_month_sabzijat_weight,omitempty"`
	MediumActivityMonthInYear *uint   `json:"medium_activity_month_in_year,omitempty"`
	MediumActivityHourInWeek  *string `json:"medium_activity_hour_in_week,omitempty"`
	HardActivityMonthInYear   *uint   `json:"hard_activity_month_in_year,omitempty"`
	HardActivityHourInWeek    *string `json:"hard_activity_hour_in_week,omitempty"`
	SmokeAtLeast100           *bool   `json:"smoke_at_least_100,omitempty"`
	SmokingAge                *uint   `json:"smoking_age,omitempty"`
	SmokingNow                *bool   `json:"smoking_now,omitempty"`
	LeaveSmokingAge           *uint   `json:"leave_smoking_age,omitempty"`
	CountSmokingDaily         *uint   `json:"count_smoking_daily,omitempty"`
	CountGheliandaily         *uint   `json:"count_gheliandaily,omitempty"`
	CountSmokingDailyPast     *uint   `json:"count_smoking_daily_past,omitempty"`
	CountGheliandailyPast     *uint   `json:"count_gheliandaily_past,omitempty"`

	GhaedeAge                    *uint   `json:"ghaede_age,omitempty"`
	HasChildren                  *bool   `json:"has_children,omitempty"`
	NumberOfChildren             *uint   `json:"number_of_children,omitempty"`
	AgeOfFirstBirth              *uint   `json:"age_of_first_birth,omitempty"`
	MenopausalStatus             *string `json:"menopausal_status,omitempty"`
	MenopauseAge                 *string `json:"menopause_age,omitempty"`
	HRT                          *bool   `json:"hrt,omitempty"`
	HRTUseLength                 *uint   `json:"hrt_use_length,omitempty"`
	LastFiveYearsHRTUse          *bool   `json:"last_five_years_hrt_use,omitempty"`
	CurrentHRTUse                *bool   `json:"current_hrt_use,omitempty"`
	IntendedHRTUse               *uint   `json:"intended_hrt_use,omitempty"`
	HRTType                      *string `json:"hrt_type,omitempty"`
	Oral                         *bool   `json:"oral,omitempty"`
	OralDuration                 *string `json:"oral_duration,omitempty"`
	OralTwoLastYears             *bool   `json:"oral_two_last_years,omitempty"`
	MamoGraphy                   *bool   `json:"mamo_graphy,omitempty"`
	Falop                        *bool   `json:"falop,omitempty"`
	Andometrioz                  *bool   `json:"andometrioz,omitempty"`
	LeavePestan                  *bool   `json:"leave_pestan,omitempty"`
	LeaveTokhmdan                *bool   `json:"leave_tokhmdan,omitempty"`
	LaDeColon                    *bool   `json:"la_de_colon,omitempty"`
	LaDePol                      *bool   `json:"la_de_pol,omitempty"`
	AspLaMo                      *bool   `json:"asp_la_mo,omitempty"`
	NsaiDLaMo                    *bool   `json:"nsai_d_la_mo,omitempty"`
	LastFiveYearBloodTestInStool *bool   `json:"last_five_year_blood_test_in_stool,omitempty"`

	Cancer     *bool   `json:"cancer,omitempty"`
	CancerType *string `json:"cancer_type,omitempty"`
	CancerAge  *uint   `json:"cancer_age,omitempty"`

	ChildCancer     *bool   `json:"child_cancer,omitempty"`
	ChildName       *string `json:"child_name,omitempty"`
	ChildCancerType *string `json:"child_cancer_type,omitempty"`
	ChildCancerAge  *string `json:"child_cancer_age,omitempty"`
	ChildLifeStatus *string `json:"child_life_status,omitempty"`

	MotherCancer     *bool   `json:"mother_cancer,omitempty"`
	MotherName       *string `json:"mother_name,omitempty"`
	MotherLifeStatus *string `json:"mother_life_status,omitempty"`
	MotherCancerType *string `json:"mother_cancer_type,omitempty"`
	MotherCancerAge  *string `json:"mother_cancer_age,omitempty"`

	FatherCancer     *bool   `json:"father_cancer,omitempty"`
	FatherName       *string `json:"father_name,omitempty"`
	FatherLifeStatus *string `json:"father_life_status,omitempty"`
	FatherCancerType *string `json:"father_cancer_type,omitempty"`
	FatherCancerAge  *string `json:"father_cancer_age,omitempty"`

	SiblingCancer     *bool   `json:"sibling_cancer,omitempty"`
	SiblingName       *string `json:"sibling_name,omitempty"`
	SiblingLifeStatus *string `json:"sibling_life_status,omitempty"`
	SiblingCancerType *string `json:"sibling_cancer_type,omitempty"`
	SiblingCancerAge  *string `json:"sibling_cancer_age,omitempty"`

	AmeAmoCancer     *bool   `json:"ame_amo_cancer,omitempty"`
	AmeAmoName       *string `json:"ame_amo_name,omitempty"`
	AmeAmoLifeStatus *string `json:"ame_amo_life_status,omitempty"`
	AmeAmoCancerType *string `json:"ame_amo_cancer_type,omitempty"`
	AmeAmoCancerAge  *string `json:"ame_amo_cancer_age,omitempty"`

	KhaleDaeiCancer     *bool   `json:"khale_daei_cancer,omitempty"`
	KhaleDaeiName       *string `json:"khale_daei_name,omitempty"`
	KhaleDaeiLifeStatus *string `json:"khale_daei_life_status,omitempty"`
	KhaleDaeiCancerType *string `json:"khale_daei_cancer_type,omitempty"`
	KhaleDaeiCancerAge  *string `json:"khale_daei_cancer_age,omitempty"`

	OtherRelativeCancer     *bool   `json:"other_relative_cancer,omitempty"`
	OtherRelativeName       *string `json:"other_relative_name,omitempty"`
	OtherRelativeRelation   *string `json:"other_relative_relation,omitempty"`
	OtherRelativeLifeStatus *string `json:"other_relative_life_status,omitempty"`
	OtherRelativeCancerType *string `json:"other_relative_cancer_type,omitempty"`
	OtherRelativeCancerAge  *string `json:"other_relative_cancer_age,omitempty"`

	TestGen   *bool `json:"test_gen,omitempty"`
	FmTestGen *bool `json:"fm_test_gen,omitempty"`

	CallExpert   *bool   `json:"call_expert,omitempty"`
	BirthCountry *string `json:"birth_country,omitempty"`
	Province     *string `json:"province,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`

	InsuranceStatus           *string `json:"insurance_status,omitempty"`
	SupplementaryInsurances   *string `json:"supplementary_insurances,omitempty"`
	Hypertension              *bool   `json:"hypertension,omitempty"`
	HypertensionTreatment     *bool   `json:"hypertension_treatment,omitempty"`
	HeartDisease              *bool   `json:"heart_disease,omitempty"`
	HeartDiseaseTreatment     *bool   `json:"heart_disease_treatment,omitempty"`
	Diabetes                  *bool   `json:"diabetes,omitempty"`
	DiabetesTreatment         *bool   `json:"diabetes_treatment,omitempty"`
	ChronicLungDisease        *bool   `json:"chronic_lung_disease,omitempty"`
	ChronicLungDiseaseType    *string `json:"chronic_lung_disease_type,omitempty"`
	LungCancerHistory         *bool   `json:"lung_cancer_history,omitempty"`
	OtherCancerHistory        *bool   `json:"other_cancer_history,omitempty"`
	OtherCancerType           *string `json:"other_cancer_type,omitempty"`
	LungCancerFamily          *bool   `json:"lung_cancer_family,omitempty"`
	LungCancerFamilyRelation  *string `json:"lung_cancer_family_relation,omitempty"`
	OtherCancerFamily         *bool   `json:"other_cancer_family,omitempty"`
	OtherCancerFamilyType     *string `json:"other_cancer_family_type,omitempty"`
	OtherCancerFamilyRelation *string `json:"other_cancer_family_relation,omitempty"`
	OccupationalExposure      *string `json:"occupational_exposure,omitempty"`
	CurrentSmoking            *bool   `json:"current_smoking,omitempty"`
	SmokingStartAgeCurrent    *uint   `json:"smoking_start_age_current,omitempty"`
	SmokingTypesCurrent       *string `json:"smoking_types_current,omitempty"`
	CigarettesPerDayCurrent   *uint   `json:"cigarettes_per_day_current,omitempty"`
	CigarPerDayCurrent        *uint   `json:"cigar_per_day_current,omitempty"`
	ECigPerDayCurrent         *uint   `json:"e_cig_per_day_current,omitempty"`
	PipePerDayCurrent         *uint   `json:"pipe_per_day_current,omitempty"`
	ChapoghPerDayCurrent      *uint   `json:"chapogh_per_day_current,omitempty"`
	SmokedOpiumPerDayCurrent  *uint   `json:"smoked_opium_per_day_current,omitempty"`
	ChewedOpiumPerDayCurrent  *uint   `json:"chewed_opium_per_day_current,omitempty"`
	HookahPerWeekCurrent      *uint   `json:"hookah_per_week_current,omitempty"`
	PastSmoking               *string `json:"past_smoking,omitempty"`
	SmokingStartAgePast       *uint   `json:"smoking_start_age_past,omitempty"`
	SmokingTypesPast          *string `json:"smoking_types_past,omitempty"`
	CigarettesPerDayPast      *uint   `json:"cigarettes_per_day_past,omitempty"`
	CigarPerDayPast           *uint   `json:"cigar_per_day_past,omitempty"`
	ECigPerDayPast            *uint   `json:"e_cig_per_day_past,omitempty"`
	PipePerDayPast            *uint   `json:"pipe_per_day_past,omitempty"`
	ChapoghPerDayPast         *uint   `json:"chapogh_per_day_past,omitempty"`
	SmokedOpiumPerDayPast     *uint   `json:"smoked_opium_per_day_past,omitempty"`
	ChewedOpiumPerDayPast     *uint   `json:"chewed_opium_per_day_past,omitempty"`
	HookahPerWeekPast         *uint   `json:"hookah_per_week_past,omitempty"`
	SecondhandSmoke           *bool   `json:"secondhand_smoke,omitempty"`
	SecondhandSmokeLocation   *string `json:"secondhand_smoke_location,omitempty"`
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
