package enum

type ActionType uint

const (
	ActionTypeFormAssigned ActionType = iota + 1
	ActionTypeFormUnAssigned
	ActionTypeUserRoleUpdated
	ActionTypeFormSentToPremm5
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
		return "فرم به مدل پرم۵ فرستاده شد"
	}
	return "نامشخص"
}

func GetAllActionTypes() []ActionType {
	return []ActionType{
		ActionTypeFormAssigned,
		ActionTypeUserRoleUpdated,
		ActionTypeFormUnAssigned,
		ActionTypeFormSentToPremm5,
	}
}
