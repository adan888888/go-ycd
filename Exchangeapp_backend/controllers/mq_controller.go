package controllers

import (
	"exchangeapp/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendRabbitMsg(ctx *gin.Context) {
	msg := ctx.Param("msg")
	config.InitRabbitMQ(msg)
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   1,
		Msg:    "发送成功",
		Data: map[string]interface{}{
			"data": msg,
		}})
}
