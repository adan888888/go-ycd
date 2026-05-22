package subscription

import (
	"exchangeapp/models"
	"time"
)

const SuperAdminUsername = "Admin"
const YcdExpiredMsg = "请充值"
// ExpiresAtNotSet 仅管理端展示，勿下发给 Flutter 等客户端
const ExpiresAtNotSet = "未设置"

// CodeYcdExpired 业务码：ycd 服务已到期
const CodeYcdExpired = 2202

// IsPermanentUser 超级管理员永久有效
func IsPermanentUser(username string) bool {
	return username == SuperAdminUsername
}

// IsYcdAllowed 是否可使用 ycd 功能
func IsYcdAllowed(user models.User) bool {
	if IsPermanentUser(user.Username) {
		return true
	}
	if user.ExpiresAt == nil {
		return false
	}
	return !time.Now().After(*user.ExpiresAt)
}

// FormatExpiresAt 客户端（登录等）展示：无到期时间为空字符串，不含「未设置」
func FormatExpiresAt(user models.User) (display string, expiresAt *time.Time, isPermanent bool) {
	if IsPermanentUser(user.Username) {
		return "永久", nil, true
	}
	if user.ExpiresAt == nil {
		return "", nil, false
	}
	return user.ExpiresAt.Format("2006-01-02 15:04:05"), user.ExpiresAt, false
}

// FormatExpiresAtForAdmin 管理端列表/编辑：无到期时间返回「未设置」
func FormatExpiresAtForAdmin(user models.User) string {
	display, _, _ := FormatExpiresAt(user)
	if display == "" {
		return ExpiresAtNotSet
	}
	return display
}
