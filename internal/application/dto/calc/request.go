package calc

type SendFormToCalcRequest struct {
	FormID uint
	CalcID uint
}

type SendFormToPremm5Request struct {
	Sex                    uint
	PersonalCrcOne         bool
	PersonalCrcMultiple    bool
	AgeCrcDx               uint
	PersonalEndometrial    bool
	AgeEcDx                uint
	PersonalLsOther        bool
	FirstDegreeCrcOne      bool
	AgeYoungestRelativeCrc uint
	CurrentAge             uint
	FamilyLsOther          bool
}
