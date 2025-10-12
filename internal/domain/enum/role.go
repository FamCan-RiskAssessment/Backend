package enum

type RoleName uint

const (
	SuperAdmin RoleName = iota + 1
	Supervisor
	Operator
	Patient
)

var rolePermissions = map[RoleName][]PermissionType{
	SuperAdmin: {
		PermissionAll,
		PermissionSetPassword,
	},
	Supervisor: {
		PermissionType(CategoryPatientManagement),
		PermissionType(CategoryFormManagement),
		PermissionType(CategoryDataEntry),
		PermissionSetPassword,
	},
	Operator: {
		PermissionType(CategoryPatientManagement),
		PermissionType(CategoryFormManagement),
		PermissionType(CategoryDataEntry),
		PermissionSetPassword,
	},
	Patient: {
		PermissionType(CategoryDataEntry),
	},
}

func (role RoleName) Permissions() []PermissionType {
	if permissions, ok := rolePermissions[role]; ok {
		return permissions
	}
	return nil
}
func (role RoleName) String() string {
	switch role {
	case SuperAdmin:
		return "سوپر ادمین"
	case Supervisor:
		return "مدیر"
	case Operator:
		return "اپراتور"
	case Patient:
		return "بیمار"
	}
	return "unknown"
}

func GetAllRoleNames() []RoleName {
	return []RoleName{
		SuperAdmin,
		Supervisor,
		Operator,
		Patient,
	}
}
