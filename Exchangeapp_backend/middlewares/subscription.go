package middlewares

import (
	"exchangeapp/subscription"
	"net/http"

	"github.com/gin-gonic/gin"
)

// YcdSubscriptionMiddleware 拦截已到期普通用户的 ycd 接口
func YcdSubscriptionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if v, ok := ctx.Get("isSuperAdmin"); ok && v == true {
			ctx.Next()
			return
		}
		if v, ok := ctx.Get("ycdAllowed"); ok && v == true {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": subscription.CodeYcdExpired,
			"msg":  subscription.YcdExpiredMsg,
			"data": gin.H{},
		})
	}
}
