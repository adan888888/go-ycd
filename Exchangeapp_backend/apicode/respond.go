package apicode

import (
	"github.com/gin-gonic/gin"
)

// Abort 中间件/拦截器统一写响应：HTTP 状态由 code 推导，body 固定结构。
func Abort(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(HTTPStatusForCode(code), gin.H{
		"code": code,
		"msg":  msg,
		"data": gin.H{},
	})
}
