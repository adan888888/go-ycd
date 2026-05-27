package apicode

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Abort 中间件/拦截器统一写响应：HTTP 固定 200，body 为 { code, msg, data }。
func Abort(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": gin.H{},
	})
}
