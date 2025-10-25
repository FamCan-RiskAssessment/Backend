package enum

type MenopausalStatus uint

const (
	MenopausalStatusPreMenopausal MenopausalStatus = iota + 1
	MenopausalStatusPeriMenopausal
	MenopausalStatusPostMenopausal
	MenopausalStatusUnknown
)

func (ms MenopausalStatus) String() string {
	switch ms {
	case MenopausalStatusPreMenopausal:
		return "پیش از یائسگی"
	case MenopausalStatusPeriMenopausal:
		return "دوران یائسگی"
	case MenopausalStatusPostMenopausal:
		return "پس از یائسگی"
	case MenopausalStatusUnknown:
		return "نامشخص"
	}
	return "نامشخص"
}

func GetAllMenopausalStatuses() []MenopausalStatus {
	return []MenopausalStatus{
		MenopausalStatusPreMenopausal,
		MenopausalStatusPeriMenopausal,
		MenopausalStatusPostMenopausal,
		MenopausalStatusUnknown,
	}
}
