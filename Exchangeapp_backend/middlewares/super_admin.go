package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuperAdminMiddleware 仅允许超级管理员访问（需先经过 AuthMiddleWare）
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get("isSuperAdmin")
		if !ok || v != true {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"code":   403,
				"msg":    "仅超级管理员可访问",
				"data":   gin.H{},
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
