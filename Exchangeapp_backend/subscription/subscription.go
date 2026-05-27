package subscription

import (
	"exchangeapp/models"
	"fmt"
	"time"
)

const JsqExpiredMsg = "服务已到期，请联系管理员"

// ExpiresAtNotSet 仅管理端展示，勿下发给 Flutter 等客户端
const ExpiresAtNotSet = "未设置"

// IsPermanentUser 超级管理员永久有效（按 role 判断）
func IsPermanentUser(user models.User) bool {
	return models.IsSuperAdminRole(user.Role)
}

// IsJsqAllowed 是否可使用 jsq（计数器）功能
func IsJsqAllowed(user models.User) bool {
	if IsPermanentUser(user) {
		return true
	}
	if user.ExpiresAt == nil {
		return false
	}
	return !time.Now().After(*user.ExpiresAt)
}

// JsqExpiredMessage 按用户到期信息生成提示，供 API 响应 msg 使用
func JsqExpiredMessage(user models.User) string {
	if IsPermanentUser(user) {
		return JsqExpiredMsg
	}
	display, _, _ := FormatExpiresAt(user)
	if display == "" {
		return "服务未开通，请联系管理员"
	}
	date := display
	if len(display) >= 10 {
		date = display[:10]
	}
	return fmt.Sprintf("服务已到期（%s），请联系管理员", date)
}

// FormatExpiresAt 客户端（登录等）展示：无到期时间为空字符串，不含「未设置」
func FormatExpiresAt(user models.User) (display string, expiresAt *time.Time, isPermanent bool) {
	if IsPermanentUser(user) {
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
