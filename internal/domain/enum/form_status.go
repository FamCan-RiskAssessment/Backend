package enum

type FormStatus uint

const (
	FormStatusPending FormStatus = iota + 1
	FormStatusApproved
	FormStatusRejected
	FormStatusInComplete
	FormStatusReady
)

func (fs FormStatus) String() string {
	switch fs {
	case FormStatusPending:
		return "درحال بررسی"
	case FormStatusApproved:
		return "قبول شده"
	case FormStatusRejected:
		return "رد شده"
	case FormStatusInComplete:
		return "تکمیل نشده"
	case FormStatusReady:
		return "ارسال شده"
	}
	return "unknown"
}

func GetAllFormStatuses() []FormStatus {
	return []FormStatus{
		FormStatusPending,
		FormStatusApproved,
		FormStatusRejected,
		FormStatusInComplete,
		FormStatusReady,
	}
}
