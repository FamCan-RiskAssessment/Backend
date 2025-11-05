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

type SendFormToGBRRequest struct {
	Age          int
	LaterAge     int
	HorizonYears int
	MenarcheAge  int
	NumBiopsies  int
	FLBAge       int
	NumRelatives int
	Race         int
	ShowRR       bool
}
