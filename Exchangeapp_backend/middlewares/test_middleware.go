package middlewares

import (
	"exchangeapp/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckUser(context *gin.Context) {
	name, err := context.Cookie("name")
	if err != nil || name == "" {
		controllers.Ok(context, controllers.ResponseJson{
			Status: http.StatusOK,
			Code:   0,
			Msg:    "您尚未登录",
		})
		//在网页上，可以重定向到登录页面
		context.Redirect(http.StatusMovedPermanently, "/login")
		context.Abort()
	}
	context.Next()
}
