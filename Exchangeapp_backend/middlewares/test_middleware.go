package middlewares

import (
	"exchangeapp/controllers"
	"github.com/gin-gonic/gin"
	
	"net/http"
)

func CheckUser(context *gin.Context) {
	name, err := context.Cookie("name")
	if err != nil || name == "" {
		controllers.Success(context, "您尚未登录", gin.H{})
		//在网页上，可以重定向到登录页面
		context.Redirect(http.StatusMovedPermanently, "/login")
		context.Abort()
	}
	context.Next()
}
