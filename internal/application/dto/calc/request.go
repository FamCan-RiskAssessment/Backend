package calcdto

type SendFormToCalcRequest struct {
	UserID uint
	FormID uint
	CalcID uint
}

type SendFormToPremm5Request struct {
	Sex                    uint `json:"sex"`
	PersonalCrcOne         uint `json:"personal_crc_one"`
	PersonalCrcMultiple    uint `json:"personal_crc_multiple"`
	AgeCrcDx               uint `json:"age_crc_dx"`
	PersonalEndometrial    uint `json:"personal_endometrial"`
	AgeEcDx                uint `json:"age_ec_dx"`
	PersonalLsOther        uint `json:"personal_ls_other"`
	FirstDegreeCrcOne      uint `json:"first_degree_crc_one"`
	AgeYoungestRelativeCrc uint `json:"age_youngest_relative_crc"`
	CurrentAge             uint `json:"current_age"`
	FamilyLsOther          uint `json:"family_ls_other"`
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
