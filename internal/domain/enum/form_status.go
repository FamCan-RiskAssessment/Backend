package enum

type FormStatus uint

const (
	FormStatusDraft FormStatus = iota + 1
	FormStatusSubmitted
	FormStatusWaitingForPatientResponse
	FormStatusWaitingForDocuments
	FormStatusRejected
	FormStatusReadyForCalculation
	FormStatusCalculated
)

func (fs FormStatus) String() string {
	switch fs {
	case FormStatusDraft:
		return "پیش‌نویس"
	case FormStatusSubmitted:
		return "ارسال شده"
	case FormStatusWaitingForPatientResponse:
		return "در انتظار پاسخ بیمار"
	case FormStatusWaitingForDocuments:
		return "در انتظار مدارک"
	case FormStatusRejected:
		return "رد شده"
	case FormStatusReadyForCalculation:
		return "آماده محاسبه"
	case FormStatusCalculated:
		return "محاسبه شده"
	}
	return "unknown"
}

func GetAllFormStatuses() []FormStatus {
	return []FormStatus{
		FormStatusDraft,
		FormStatusSubmitted,
		FormStatusWaitingForPatientResponse,
		FormStatusWaitingForDocuments,
		FormStatusRejected,
		FormStatusReadyForCalculation,
		FormStatusCalculated,
	}
}
