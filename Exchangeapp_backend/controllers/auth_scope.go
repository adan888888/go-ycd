package controllers

import (
	"exchangeapp/apicode"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
// 超级管理员：user_id 为空查全部，有值查指定用户。
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

// listUserScope 列表/统计接口的用户范围（支持多选 user_ids）
type listUserScope struct {
	QueryAll bool
	UserIDs  []int64
}

func (s listUserScope) userIDLabel() string {
	if s.QueryAll || len(s.UserIDs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(s.UserIDs))
	for _, id := range s.UserIDs {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}

// resolveListUserIDsScope 解析 user_id / user_ids（逗号分隔）。
// 超管：都不传则查全部；传 user_ids 则查多个用户综合数据。
// 普通用户：只能查本人。
func resolveListUserIDsScope(ctx *gin.Context) (listUserScope, int, string) {
	userIDsParam := strings.TrimSpace(ctx.Query("user_ids"))
	userIDParam := strings.TrimSpace(ctx.Query("user_id"))

	parseIDs := func(raw string) ([]int64, int, string) {
		if raw == "" {
			return nil, 0, ""
		}
		parts := strings.Split(raw, ",")
		ids := make([]int64, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			uid, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, http.StatusBadRequest, "用户ID格式错误: " + err.Error()
			}
			ids = append(ids, uid)
		}
		return ids, 0, ""
	}

	if !isSuperAdmin(ctx) {
		uid := loginUID(ctx)
		if uid == 0 {
			return listUserScope{}, http.StatusUnauthorized, "无法识别当前登录用户"
		}
		if userIDsParam != "" {
			ids, status, msg := parseIDs(userIDsParam)
			if status != 0 {
				return listUserScope{}, status, msg
			}
			for _, id := range ids {
				if id != uid {
					return listUserScope{}, http.StatusForbidden, "无权查看其他用户数据"
				}
			}
		} else if userIDParam != "" {
			reqUID, err := strconv.ParseInt(userIDParam, 10, 64)
			if err != nil {
				return listUserScope{}, http.StatusBadRequest, "用户ID格式错误: " + err.Error()
			}
			if reqUID != uid {
				return listUserScope{}, http.StatusForbidden, "无权查看其他用户数据"
			}
		}
		return listUserScope{QueryAll: false, UserIDs: []int64{uid}}, 0, ""
	}

	if userIDsParam != "" {
		ids, status, msg := parseIDs(userIDsParam)
		if status != 0 {
			return listUserScope{}, status, msg
		}
		if len(ids) == 0 {
			return listUserScope{QueryAll: true}, 0, ""
		}
		return listUserScope{QueryAll: false, UserIDs: ids}, 0, ""
	}
	if userIDParam != "" {
		uid, err := strconv.ParseInt(userIDParam, 10, 64)
		if err != nil {
			return listUserScope{}, http.StatusBadRequest, "用户ID格式错误: " + err.Error()
		}
		return listUserScope{QueryAll: false, UserIDs: []int64{uid}}, 0, ""
	}
	return listUserScope{QueryAll: true}, 0, ""
}

func applyUserScope(query *gorm.DB, scope listUserScope, column string) *gorm.DB {
	if scope.QueryAll || len(scope.UserIDs) == 0 {
		return query
	}
	if len(scope.UserIDs) == 1 {
		return query.Where(column+" = ?", scope.UserIDs[0])
	}
	return query.Where(column+" IN ?", scope.UserIDs)
}

func buildUserIDInClause(scope listUserScope) (string, []interface{}) {
	if scope.QueryAll || len(scope.UserIDs) == 0 {
		return "", nil
	}
	if len(scope.UserIDs) == 1 {
		return columnEqPlaceholder("user_id"), []interface{}{scope.UserIDs[0]}
	}
	placeholders := make([]string, len(scope.UserIDs))
	args := make([]interface{}, len(scope.UserIDs))
	for i, id := range scope.UserIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	return "user_id IN (" + strings.Join(placeholders, ",") + ")", args
}

func columnEqPlaceholder(column string) string {
	return column + " = ?"
}

func requireSuperAdmin(ctx *gin.Context) bool {
	if isSuperAdmin(ctx) {
		return true
	}
	forbidScope(ctx, http.StatusForbidden, "仅超级管理员可执行此操作")
	return false
}

func forbidScope(ctx *gin.Context, httpStatus int, msg string) {
	Fail(ctx, apicode.CodeFromHTTPStatus(httpStatus), msg)
}
