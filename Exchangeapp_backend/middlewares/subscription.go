package middlewares

import (
	"exchangeapp/apicode"
	"exchangeapp/models"
	"exchangeapp/subscription"
	"time"

	"github.com/gin-gonic/gin"
)

// JsqSubscriptionMiddleware 拦截已到期普通用户的 jsq 接口
func JsqSubscriptionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if v, ok := ctx.Get("isSuperAdmin"); ok && v == true {
			ctx.Next()
			return
		}
		if v, ok := ctx.Get("jsqAllowed"); ok && v == true {
			ctx.Next()
			return
		}
		msg := subscription.JsqExpiredMsg
		if exp, ok := ctx.Get("expiresAt"); ok {
			if t, ok := exp.(time.Time); ok {
				msg = subscription.JsqExpiredMessage(models.User{ExpiresAt: &t})
			}
		}
		apicode.Abort(ctx, apicode.CodeJsqExpired, msg)
	}
}
