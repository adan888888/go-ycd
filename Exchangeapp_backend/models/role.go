package models

const (
	RoleSuperAdmin = "super_admin"
	RolePro        = "pro"
	RoleUser       = "user"
)

// LegacySuperAdminUsername 历史超管默认用户名，仅用于数据迁移
const LegacySuperAdminUsername = "Admin"

func IsSuperAdminRole(role string) bool {
	return role == RoleSuperAdmin
}

func IsProRole(role string) bool {
	return NormalizeUserRole(role) == RolePro
}

// IsProOrAboveRole 专业版及以上（含超级管理员）
func IsProOrAboveRole(role string) bool {
	r := NormalizeUserRole(role)
	return r == RoleSuperAdmin || r == RolePro
}

func NormalizeUserRole(role string) string {
	switch role {
	case RoleSuperAdmin:
		return RoleSuperAdmin
	case RolePro:
		return RolePro
	default:
		return RoleUser
	}
}
