package middlewares

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/subscription"

	"github.com/gin-gonic/gin"
)

// AttachLoginUser 根据 JWT 用户名加载 uid、role，并标记是否超级管理员
func AttachLoginUser(ctx *gin.Context, username string) {
	ctx.Set("loginUsername", username)
	ctx.Set("isSuperAdmin", false)
	ctx.Set("loginRole", models.RoleUser)

	var user models.User
	if err := global.Db.Select("uid", "username", "expires_at", "role").
		Where("BINARY username = ?", username).First(&user).Error; err == nil && user.Uid != 0 {
		role := models.NormalizeUserRole(user.Role)
		ctx.Set("loginRole", role)
		ctx.Set("isSuperAdmin", models.IsSuperAdminRole(role))
		ctx.Set("loginUid", user.Uid)
		ctx.Set("jsqAllowed", subscription.IsJsqAllowed(user))
		if user.ExpiresAt != nil {
			ctx.Set("expiresAt", *user.ExpiresAt)
		}
		return
	}
	ctx.Set("loginUid", int64(0))
	ctx.Set("jsqAllowed", false)
}
