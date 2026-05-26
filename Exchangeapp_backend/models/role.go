package models

const (
	RoleSuperAdmin = "super_admin"
	RoleUser       = "user"
)

// LegacySuperAdminUsername 历史超管默认用户名，仅用于数据迁移
const LegacySuperAdminUsername = "Admin"

func IsSuperAdminRole(role string) bool {
	return role == RoleSuperAdmin
}

func NormalizeUserRole(role string) string {
	if IsSuperAdminRole(role) {
		return RoleSuperAdmin
	}
	return RoleUser
}
