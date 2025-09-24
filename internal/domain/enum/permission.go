package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1

	PermissionEnterData
	PermissionUpdateOwnData
	PermissionDeleteOwnData

	PermissionViewPatients
	PermissionCreatePatients
	PermissionUpdatePatients
	PermissionDeletePatients

	PermissionViewUsers
	PermissionCreateUsers
	PermissionUpdateUsers
	PermissionDeleteUsers

	PermissionViewRoles
	PermissionCreateRoles
	PermissionUpdateRoles
	PermissionDeleteRoles
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryDataEntry
	CategoryPatientManagement
	CategoryUserManagement
)

var permissionCategories = map[PermissionType]PermissionCategory{
	PermissionAll: CategoryGeneral,

	PermissionEnterData:     CategoryDataEntry,
	PermissionUpdateOwnData: CategoryDataEntry,
	PermissionDeleteOwnData: CategoryDataEntry,

	PermissionViewPatients:   CategoryPatientManagement,
	PermissionCreatePatients: CategoryPatientManagement,
	PermissionUpdatePatients: CategoryPatientManagement,
	PermissionDeletePatients: CategoryPatientManagement,

	PermissionViewUsers:   CategoryUserManagement,
	PermissionCreateUsers: CategoryUserManagement,
	PermissionUpdateUsers: CategoryUserManagement,
	PermissionDeleteUsers: CategoryUserManagement,
}

var permissionNames = map[PermissionType]string{
	PermissionAll: "دسترسی کامل",

	PermissionEnterData:     "درج داده",
	PermissionUpdateOwnData: "ویرایش داده خود",
	PermissionDeleteOwnData: "حذف داده خود",

	PermissionViewPatients:   "مشاهده بیماران",
	PermissionCreatePatients: "ایجاد بیمار",
	PermissionUpdatePatients: "ویرایش بیمار",
	PermissionDeletePatients: "حذف بیمار",

	PermissionViewUsers:   "مشاهده کاربران",
	PermissionCreateUsers: "ایجاد کاربر",
	PermissionUpdateUsers: "ویرایش کاربر",
	PermissionDeleteUsers: "حذف کاربر",

	PermissionViewRoles:   "مشاهده نقش ها",
	PermissionCreateRoles: "ایجاد نقش",
	PermissionUpdateRoles: "ویرایش نقش",
	PermissionDeleteRoles: "حذف نقش",
}

var CategoryNames = map[PermissionCategory]string{
	CategoryGeneral:           "عمومی",
	CategoryDataEntry:         "درج داده",
	CategoryUserManagement:    "مدیریت کاربران",
	CategoryPatientManagement: "مدیریت بیماران",
}

func (p PermissionType) String() string {
	if name, ok := permissionNames[p]; ok {
		return name
	}
	return "نامشخص"
}

func (c PermissionCategory) String() string {
	if name, ok := CategoryNames[c]; ok {
		return name
	}
	return "نامشخص"
}

func (p PermissionType) Category() PermissionCategory {
	if category, ok := permissionCategories[p]; ok {
		return category
	}
	return CategoryGeneral
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll,

		PermissionEnterData,
		PermissionUpdateOwnData,
		PermissionDeleteOwnData,

		PermissionViewPatients,
		PermissionCreatePatients,
		PermissionUpdatePatients,
		PermissionDeletePatients,

		PermissionViewUsers,
		PermissionCreateUsers,
		PermissionUpdateUsers,
		PermissionDeleteUsers,

		PermissionViewRoles,
		PermissionCreateRoles,
		PermissionUpdateRoles,
		PermissionDeleteRoles,
	}
}

func GetAllPermissionCategories() []PermissionCategory {
	return []PermissionCategory{
		CategoryGeneral,
		CategoryDataEntry,
		CategoryPatientManagement,
		CategoryUserManagement,
	}
}
