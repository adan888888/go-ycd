package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ResponseJson struct {
	Status int    `json:"-"`             //是忽略的 (系统的状态码)
	Code   int    `json:"code"`          //`json:"code,omitempty"` omitempty 如果不存在就不返回给前端  假如设置为o，这个字段不会返回给前端
	Msg    string `json:"msg,omitempty"` //描述
	Data   any    `json:"data,omitempty"`
}

func (m ResponseJson) IsEmpty() bool {
	return reflect.DeepEqual(m, ResponseJson{})
}
func HttpResponse(ctx *gin.Context, status int, resp ResponseJson) {
	if resp.IsEmpty() {
		ctx.AbortWithStatus(status)
		return
	}
	//AbortWithStatusJSON这个接口请求完了，停止后续动作（也就是不会返回 两个ctx json去前端）
	ctx.AbortWithStatusJSON(status, resp)
}
func buildStatus(resp ResponseJson, nDefaultStatus int) int {
	if 0 == resp.Status {
		return nDefaultStatus
	}
	return resp.Status
}
func Ok(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}
func Fail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusBadRequest), resp)
}
func ServerFail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusInternalServerError), resp) //500服务器错误
}
