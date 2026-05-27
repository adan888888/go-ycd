package controllers

import (
	"exchangeapp/apicode"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseJson 统一响应体 { code, msg, data }；HTTP 状态码固定 200，只看 code。
type ResponseJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

// Success 成功响应，code 固定为 0。
func Success(ctx *gin.Context, msg string, data any) {
	write(ctx, apicode.CodeOK, msg, data)
}

// Fail 失败响应，code 为具体业务码，msg 由后端填写。
func Fail(ctx *gin.Context, code int, msg string, data ...any) {
	d := any(gin.H{})
	if len(data) > 0 && data[0] != nil {
		d = data[0]
	}
	write(ctx, code, msg, d)
}

func write(ctx *gin.Context, code int, msg string, data any) {
	if data == nil {
		data = gin.H{}
	}
	ctx.AbortWithStatusJSON(http.StatusOK, ResponseJson{Code: code, Msg: msg, Data: data})
}
