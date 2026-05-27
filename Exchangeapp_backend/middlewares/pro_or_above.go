package middlewares

import (
	"exchangeapp/apicode"
	"exchangeapp/models"

	"github.com/gin-gonic/gin"
)

// ProOrAboveMiddleware 仅允许专业版及以上访问（需先经过 AuthMiddleWare）
func ProOrAboveMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("loginRole")
		roleStr, _ := role.(string)
		if !models.IsProOrAboveRole(roleStr) {
			apicode.Abort(ctx, apicode.CodeForbidden, "无权限，该功能需专业版及以上")
			return
		}
		ctx.Next()
	}
}
