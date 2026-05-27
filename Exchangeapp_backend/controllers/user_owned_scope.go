package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func requireLoginUID(ctx *gin.Context) (int64, bool) {
	uid := loginUID(ctx)
	if uid == 0 {
		forbidScope(ctx, http.StatusUnauthorized, "无法识别当前登录用户")
		return 0, false
	}
	return uid, true
}
