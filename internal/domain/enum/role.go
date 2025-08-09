package enum

type RoleName uint

const (
	SuperAdmin RoleName = iota + 1
)

var rolePermissions = map[RoleName][]PermissionType{
	SuperAdmin: {
		PermissionAll,
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
	}
	return "unknown"
}

func GetAllRoleNames() []RoleName {
	return []RoleName{
		SuperAdmin,
	}
}
