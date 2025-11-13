package calcdto

type SendFormToCalcRequest struct {
	UserID uint
	FormID uint
	CalcID uint
}

type SendFormToPremm5Request struct {
	Sex        int `json:"sex"`
	CurrentAge int `json:"current_age"`

	PersonalCrcCount       int `json:"personal_crc_count"` // NEW: 0, 1, or 2+
	PersonalCrcYoungestAge int `json:"personal_crc_youngest_age"`
	PersonalEc             int `json:"personal_ec"`
	PersonalEcAge          int `json:"personal_ec_age"`
	PersonalOtherLs        int `json:"personal_other_ls"`

	// NEW FIELDS - FDR (First Degree Relatives)
	NumFdrCrc         int `json:"num_fdr_crc"`
	YoungestFdrCrcAge int `json:"youngest_fdr_crc_age"`
	NumFdrEc          int `json:"num_fdr_ec"`
	YoungestFdrEcAge  int `json:"youngest_fdr_ec_age"`

	// NEW FIELDS - SDR (Second Degree Relatives)
	NumSdrCrc         int `json:"num_sdr_crc"`
	YoungestSdrCrcAge int `json:"youngest_sdr_crc_age"`
	NumSdrEc          int `json:"num_sdr_ec"`
	YoungestSdrEcAge  int `json:"youngest_sdr_ec_age"`

	HasFdrOtherLs int `json:"has_fdr_other_ls"`
	HasSdrOtherLs int `json:"has_sdr_other_ls"`
}

type SendFormToBCRARequest struct {
	T1      float64 `json:"T1"`      // Current age
	T2      float64 `json:"T2"`      // Projection age (hardcoded)
	N_Biop  int     `json:"N_Biop"`  // Number of breast biopsies
	HypPlas int     `json:"HypPlas"` // Hyperplasia in biopsy (0=no, 1=yes, 99=unknown)
	AgeMen  int     `json:"AgeMen"`  // Age at menarche/first period
	Age1st  int     `json:"Age1st"`  // Age at first live birth (98=nulliparous, 99=unknown)
	N_Rels  int     `json:"N_Rels"`  // Number of first-degree relatives with breast cancer
	Race    int     `json:"Race"`    // Race code (hardcoded)
}

type SendFormToGailRequest struct {
	Age int `json:"age"`
	// LaterAge     *int   `json:"later_age"`
	HorizonYears int    `json:"horizon_years"`
	MenarcheAge  int    `json:"menarche_age"`
	NumBiopsies  int    `json:"num_biopsies"`
	FLBAge       int    `json:"flb_age"`
	NumRelatives int    `json:"num_relatives"`
	Race         string `json:"race"` // Changed from int to string
	ShowRR       bool   `json:"show_rr"`
}

type SendFormToPLCORequest struct {
	Age                   int     `json:"age"`
	Education             int     `json:"education"`
	BMI                   float64 `json:"bmi"`
	COPD                  int     `json:"copd"`
	PersonalCancerHistory int     `json:"personal_cancer_history"`
	FamilyLungCancer      int     `json:"family_lung_cancer"`
	RaceWhite             int     `json:"race_white"`
	RaceBlack             int     `json:"race_black"`
	RaceHispanic          int     `json:"race_hispanic"`
	RaceAsian             int     `json:"race_asian"`
	RaceNHPI              int     `json:"race_nhpi"`
	RaceAmericanIndian    int     `json:"race_american_indian"`
	SmokingStatus         int     `json:"smoking_status"`
	CigarettesPerDay      float64 `json:"cigarettes_per_day"`
	SmokingDuration       int     `json:"smoking_duration"`
	YearsQuit             int     `json:"years_quit"`
	ScreeningResult       *string `json:"screening_result,omitempty"`
}
