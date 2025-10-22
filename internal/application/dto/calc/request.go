package calcdto

type SendFormToCalcRequest struct {
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
}
