package middlewares

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/subscription"

	"github.com/gin-gonic/gin"
)

const SuperAdminUsername = subscription.SuperAdminUsername

// AttachLoginUser 根据 JWT 用户名加载 uid，并标记是否超级管理员
func AttachLoginUser(ctx *gin.Context, username string) {
	ctx.Set("loginUsername", username)
	ctx.Set("isSuperAdmin", username == SuperAdminUsername)

	var user models.User
	if err := global.Db.Select("uid", "username", "expires_at").Where("BINARY username = ?", username).First(&user).Error; err == nil && user.Uid != 0 {
		ctx.Set("loginUid", user.Uid)
		ctx.Set("ycdAllowed", subscription.IsYcdAllowed(user))
		if user.ExpiresAt != nil {
			ctx.Set("expiresAt", *user.ExpiresAt)
		}
		return
	}
	ctx.Set("loginUid", int64(0))
	ctx.Set("ycdAllowed", false)
}
