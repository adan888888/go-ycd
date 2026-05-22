package subscription

import (
	"exchangeapp/models"
	"time"
)

const SuperAdminUsername = "Admin"
const YcdExpiredMsg = "请充值"

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

// FormatExpiresAt 返回展示用到期时间
func FormatExpiresAt(user models.User) (display string, expiresAt *time.Time, isPermanent bool) {
	if IsPermanentUser(user.Username) {
		return "永久", nil, true
	}
	if user.ExpiresAt == nil {
		return "未设置", nil, false
	}
	return user.ExpiresAt.Format("2006-01-02 15:04:05"), user.ExpiresAt, false
}
