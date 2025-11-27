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
		PermissionViewPatients,
		PermissionCreatePatients,
		PermissionUpdatePatients,
		PermissionDeletePatients,
		PermissionViewForms,
		PermissionCreateForms,
		PermissionUpdateForms,
		PermissionDeleteForms,
		PermissionHandleForms,
		PermissionEnterData,
		PermissionUpdateOwnData,
		PermissionDeleteOwnData,
		PermissionSetPassword,
		PermissionHandleOperators,
		PermissionCreateFormForUser,
	},
	Operator: {
		PermissionViewPatients,
		PermissionCreatePatients,
		PermissionUpdatePatients,
		PermissionDeletePatients,
		PermissionViewForms,
		PermissionCreateForms,
		PermissionUpdateForms,
		PermissionDeleteForms,
		PermissionHandleForms,
		PermissionEnterData,
		PermissionUpdateOwnData,
		PermissionDeleteOwnData,
		PermissionSetPassword,
		PermissionCreateFormForUser,
	},
	Patient: {
		PermissionEnterData,
		PermissionUpdateOwnData,
		PermissionDeleteOwnData,
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
