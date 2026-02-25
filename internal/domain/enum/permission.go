package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1

	PermissionSetPassword

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

	PermissionViewForms
	PermissionCreateForms
	PermissionUpdateForms
	PermissionDeleteForms
	PermissionHandleForms
	PermissionCreateFormForUser

	PermissionHandleOperators

	PermissionViewLogs

	PermissionFillProfile
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryDataEntry
	CategoryPatientManagement
	CategoryUserManagement
	CategoryAccessManagement
	CategoryFormManagement
	CategoryOperatorManagement
	CategoryLogManagement
	CategoryProfileManagement
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

	PermissionViewRoles:   CategoryAccessManagement,
	PermissionCreateRoles: CategoryAccessManagement,
	PermissionUpdateRoles: CategoryAccessManagement,
	PermissionDeleteRoles: CategoryAccessManagement,

	PermissionViewForms:   CategoryFormManagement,
	PermissionCreateForms: CategoryFormManagement,
	PermissionUpdateForms: CategoryFormManagement,
	PermissionDeleteForms: CategoryFormManagement,
	PermissionHandleForms: CategoryFormManagement,

	PermissionHandleOperators: CategoryOperatorManagement,
	PermissionViewLogs:        CategoryLogManagement,
	PermissionFillProfile:     CategoryProfileManagement,
}

var permissionNames = map[PermissionType]string{
	PermissionAll: "دسترسی کامل",

	PermissionEnterData:     "درج داده",
	PermissionUpdateOwnData: "ویرایش داده خود",
	PermissionDeleteOwnData: "حذف داده خود",

	PermissionViewPatients:   "مشاهده مراجعه کنندگان",
	PermissionCreatePatients: "ایجاد مراجعه کننده",
	PermissionUpdatePatients: "ویرایش مراجعه کننده",
	PermissionDeletePatients: "حذف مراجعه کننده",

	PermissionViewUsers:   "مشاهده کاربران",
	PermissionCreateUsers: "ایجاد کاربر",
	PermissionUpdateUsers: "ویرایش کاربر",
	PermissionDeleteUsers: "حذف کاربر",

	PermissionViewRoles:   "مشاهده نقش ها",
	PermissionCreateRoles: "ایجاد نقش",
	PermissionUpdateRoles: "ویرایش نقش",
	PermissionDeleteRoles: "حذف نقش",

	PermissionSetPassword: "تنظیم پسورد",

	PermissionViewForms:         "مشاهده فرم ها",
	PermissionCreateForms:       "ایجاد فرم",
	PermissionUpdateForms:       "ویرایش فرم",
	PermissionDeleteForms:       "حذف فرم",
	PermissionHandleForms:       "مدیریت فرم ها",
	PermissionCreateFormForUser: "ایجاد فرم برای کاربر",

	PermissionHandleOperators: "مدیریت اپراتور ها",

	PermissionViewLogs: "مشاهده لاگ ها",

	PermissionFillProfile: "تکمیل پروفایل",
}

var CategoryNames = map[PermissionCategory]string{
	CategoryGeneral:            "عمومی",
	CategoryDataEntry:          "درج داده",
	CategoryUserManagement:     "مدیریت کاربران",
	CategoryPatientManagement:  "مدیریت مراجعه کنندگان",
	CategoryAccessManagement:   "مدیریت دسترسی",
	CategoryFormManagement:     "مدیریت فرم ها",
	CategoryOperatorManagement: "مدیریت اپراتور ها",
	CategoryLogManagement:      "مدیریت لاگ ها",
	CategoryProfileManagement:  "مدیریت پروفایل",
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

		PermissionViewForms,
		PermissionCreateForms,
		PermissionUpdateForms,
		PermissionDeleteForms,
		PermissionHandleForms,
		PermissionCreateFormForUser,

		PermissionSetPassword,
		PermissionHandleOperators,
		PermissionViewLogs,
		PermissionFillProfile,
	}
}

func GetAllPermissionCategories() []PermissionCategory {
	return []PermissionCategory{
		CategoryGeneral,
		CategoryDataEntry,
		CategoryPatientManagement,
		CategoryUserManagement,
		CategoryAccessManagement,
		CategoryFormManagement,
		CategoryLogManagement,
		CategoryProfileManagement,
	}
}
