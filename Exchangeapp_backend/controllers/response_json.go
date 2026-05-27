package controllers

import (
	"exchangeapp/apicode"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// code 0 成功，1 普通失败；401/403/2202 等为全局业务码，见 apicode 包
type ResponseJson struct {
	Status int    `json:"-"`             //是忽略的 (系统的状态码,传的话，就是自己写的，不传系统会默认传)
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

// 若 resp.Status 为 0：Ok 固定 200；Fail/ServerFail 按 apicode.HTTPStatusForCode(resp.Code) 推导
func buildStatus(resp ResponseJson, nDefaultStatus int) int {
	if resp.Status != 0 {
		return resp.Status
	}
	if nDefaultStatus == http.StatusOK {
		return http.StatusOK
	}
	if resp.Code == 0 {
		return nDefaultStatus
	}
	return apicode.HTTPStatusForCode(resp.Code)
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

func FailMsg(ctx *gin.Context, msg string) {
	Fail(ctx, ResponseJson{Code: apicode.CodeFail, Msg: msg, Data: gin.H{}})
}

func ServerFailMsg(ctx *gin.Context, msg string) {
	ServerFail(ctx, ResponseJson{Code: apicode.CodeServerError, Msg: msg, Data: gin.H{}})
}

func NotFoundMsg(ctx *gin.Context, msg string) {
	Fail(ctx, ResponseJson{Code: apicode.CodeNotFound, Msg: msg, Data: gin.H{}})
}

func OkMsg(ctx *gin.Context, msg string, data any) {
	if data == nil {
		data = gin.H{}
	}
	Ok(ctx, ResponseJson{Code: apicode.CodeOK, Msg: msg, Data: data})
}
