package middlewares

import (
	"exchangeapp/apicode"

	"github.com/gin-gonic/gin"
)

// SuperAdminMiddleware 仅允许超级管理员访问（需先经过 AuthMiddleWare）
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get("isSuperAdmin")
		if !ok || v != true {
			apicode.Abort(ctx, apicode.CodeForbidden, "无权限，请联系管理员")
			return
		}
		ctx.Next()
	}
}
