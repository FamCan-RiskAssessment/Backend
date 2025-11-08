package enum

type FormStatus uint

const (
	FormStatusPending FormStatus = iota + 1
	FormStatusAssigned
	FormStatusApproved
	FormStatusRejected
	FormStatusReady
	FormStatusSentToCalc
)

func (fs FormStatus) String() string {
	switch fs {
	case FormStatusPending:
		return "در حال برسی"
	case FormStatusAssigned:
		return "تخصیص شده"
	case FormStatusApproved:
		return "قبول شده"
	case FormStatusRejected:
		return "رد شده"
	case FormStatusReady:
		return "ارسال شده"
	case FormStatusSentToCalc:
		return "ارسال شده به مدل"
	}
	return "unknown"
}

func GetAllFormStatuses() []FormStatus {
	return []FormStatus{
		FormStatusPending,
		FormStatusAssigned,
		FormStatusApproved,
		FormStatusRejected,
		FormStatusReady,
		FormStatusSentToCalc,
	}
}
