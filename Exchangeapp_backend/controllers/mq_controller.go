package controllers

import (
	"exchangeapp/config"

	"github.com/gin-gonic/gin"
)

func SendRabbitMsg(ctx *gin.Context) {
	msg := ctx.Param("msg")
	config.InitRabbitMQ(msg)
	Success(ctx, "发送成功", gin.H{"data": msg})
}
