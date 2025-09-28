package enum

type FormStatus uint

const (
	FormStatusPending FormStatus = iota + 1
	FormStatusApproved
	FormStatusRejected
)

func (fs FormStatus) String() string {
	switch fs {
	case FormStatusPending:
		return "درحال بررسی"
	case FormStatusApproved:
		return "قبول شده"
	case FormStatusRejected:
		return "رد شده"
	}
	return "unknown"
}

func GetAllFormStatuses() []FormStatus {
	return []FormStatus{
		FormStatusPending,
		FormStatusApproved,
		FormStatusRejected,
	}
}
