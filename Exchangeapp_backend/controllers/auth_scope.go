package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const superAdminUsername = "Admin"

// usernameCaseSensitiveSQL MySQL 下按字节区分大小写匹配用户名
const usernameCaseSensitiveSQL = "BINARY username = ?"

func isSuperAdmin(ctx *gin.Context) bool {
	v, ok := ctx.Get("isSuperAdmin")
	return ok && v == true
}

func loginUID(ctx *gin.Context) int64 {
	v, ok := ctx.Get("loginUid")
	if !ok {
		return 0
	}
	switch uid := v.(type) {
	case int64:
		return uid
	case int:
		return int64(uid)
	case float64:
		return int64(uid)
	default:
		return 0
	}
}

// resolveListUserScope 解析列表/统计类接口的 user_id 范围。
// 超级管理员 Admin：user_id 为空查全部，有值查指定用户。
// 普通用户：只能查自己的 uid，传入其他 user_id 返回 403。
func resolveListUserScope(ctx *gin.Context, requestedUserIDStr string) (scopedUID int64, queryAll bool, httpStatus int, errMsg string) {
	if isSuperAdmin(ctx) {
		if requestedUserIDStr == "" {
			return 0, true, 0, ""
		}
		uid, err := strconv.ParseInt(requestedUserIDStr, 10, 64)
		if err != nil {
			return 0, false, http.StatusBadRequest, "用户ID格式错误: " + err.Error()
		}
		return uid, false, 0, ""
	}

	uid := loginUID(ctx)
	if uid == 0 {
		return 0, false, http.StatusUnauthorized, "无法识别当前登录用户"
	}
	if requestedUserIDStr != "" {
		reqUID, err := strconv.ParseInt(requestedUserIDStr, 10, 64)
		if err != nil {
			return 0, false, http.StatusBadRequest, "用户ID格式错误: " + err.Error()
		}
		if reqUID != uid {
			return 0, false, http.StatusForbidden, "无权查看其他用户数据"
		}
	}
	return uid, false, 0, ""
}

func requireSuperAdmin(ctx *gin.Context) bool {
	if isSuperAdmin(ctx) {
		return true
	}
	forbidScope(ctx, http.StatusForbidden, "仅超级管理员 Admin 可执行此操作")
	return false
}

func forbidScope(ctx *gin.Context, status int, msg string) {
	code := 1
	if status == http.StatusForbidden {
		code = 403
	}
	Fail(ctx, ResponseJson{
		Status: status,
		Code:   code,
		Msg:    msg,
		Data:   gin.H{},
	})
}
