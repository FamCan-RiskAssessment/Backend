package enum

type ActionType uint

const (
	ActionTypeFormAssigned ActionType = iota + 1
	ActionTypeFormUnAssigned
	ActionTypeUserRoleUpdated
	ActionTypeFormSentToPremm5
	ActionTypeFormSentToBCRA
	ActionTypeFormSentToGail
	ActionTypeOperatorCreatedForm
	ActionTypeOperatorUpdatedForm
)

func (at ActionType) String() string {
	switch at {
	case ActionTypeFormAssigned:
		return "اساین شدن فرم"
	case ActionTypeUserRoleUpdated:
		return "بروزرسانی شدن رول های یوزر"
	case ActionTypeFormUnAssigned:
		return "معلق شدن فرم"
	case ActionTypeFormSentToPremm5:
		return "فرم به مدل Premm5 فرستاده شد"
	case ActionTypeFormSentToBCRA:
		return "فرم به مدل BCRA فرستاده شد"
	case ActionTypeFormSentToGail:
		return "فرم به مدل Gail فرستاده شد"
	case ActionTypeOperatorCreatedForm:
		return "اپراتور فرم ایجاد کرد"
	case ActionTypeOperatorUpdatedForm:
		return "اپراتور فرم را بروزرسانی کرد"
	}
	return "نامشخص"
}

func GetAllActionTypes() []ActionType {
	return []ActionType{
		ActionTypeFormAssigned,
		ActionTypeUserRoleUpdated,
		ActionTypeFormUnAssigned,
		ActionTypeFormSentToPremm5,
		ActionTypeFormSentToBCRA,
		ActionTypeFormSentToGail,
		ActionTypeOperatorCreatedForm,
		ActionTypeOperatorUpdatedForm,
	}
}
