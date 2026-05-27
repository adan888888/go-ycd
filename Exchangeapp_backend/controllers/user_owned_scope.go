package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func requireLoginUID(ctx *gin.Context) (int64, bool) {
	uid := loginUID(ctx)
	if uid == 0 {
		forbidScope(ctx, http.StatusUnauthorized, "无法识别当前登录用户")
		return 0, false
	}
	return uid, true
}

// applyOwnedDataScope 超管查全部；普通用户仅查本人（买币/密码本接口当前仅超管可访问）
func applyOwnedDataScope(ctx *gin.Context, query *gorm.DB, column string) (*gorm.DB, bool) {
	if isSuperAdmin(ctx) {
		return query, true
	}
	uid, ok := requireLoginUID(ctx)
	if !ok {
		return nil, false
	}
	return query.Where(column+" = ?", uid), true
}
