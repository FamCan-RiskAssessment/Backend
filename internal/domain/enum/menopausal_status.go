package enum

type MenopausalStatus uint

const (
	MenopausalStatusPreMenopausal MenopausalStatus = iota + 1
	MenopausalStatusPostMenopausal
	MenopausalStatusMedicationInduced
	MenopausalStatusUnknown
)

func (ms MenopausalStatus) String() string {
	switch ms {
	case MenopausalStatusPreMenopausal:
		return "خیر (Premenopausal)"
	case MenopausalStatusPostMenopausal:
		return "بله (Postmenopausal)"
	case MenopausalStatusMedicationInduced:
		return "در اثر مصرف موقتا متوقف شده است"
	case MenopausalStatusUnknown:
		return "اطلاع ندارم"
	}
	return "نامشخص"
}

func GetAllMenopausalStatuses() []MenopausalStatus {
	return []MenopausalStatus{
		MenopausalStatusPreMenopausal,
		MenopausalStatusPostMenopausal,
		MenopausalStatusMedicationInduced,
		MenopausalStatusUnknown,
	}
}
