package enum

type ActionType uint

const (
	ActionTypeFormAssigned ActionType = iota + 1
	ActionTypeFormUnAssigned
	ActionTypeUserRoleUpdated
)

func (at ActionType) String() string {
	switch at {
	case ActionTypeFormAssigned:
		return "اساین شدن فرم"
	case ActionTypeUserRoleUpdated:
		return "بروزرسانی شدن رول های یوزر"
	case ActionTypeFormUnAssigned:
		return "معلق شدن فرم"
	}
	return "نامشخص"
}

func GetAllActionTypes() []ActionType {
	return []ActionType{
		ActionTypeFormAssigned,
		ActionTypeUserRoleUpdated,
		ActionTypeFormUnAssigned,
	}
}
