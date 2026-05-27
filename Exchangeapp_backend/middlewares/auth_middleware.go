package middlewares

import (
	"exchangeapp/apicode"
	"exchangeapp/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(utils.TOKEN_NAME)
		if token == "" || !strings.HasPrefix(token, utils.TOKEN_PREFIX) {
			apicode.Abort(ctx, apicode.CodeUnauthorized, "Missing Authorization Header")
			return
		}
		username, err := utils.ParseJWT(token)

		if err != nil {
			apicode.Abort(ctx, apicode.CodeUnauthorized, "token已过期或无效")
			return
		}

		ctx.Set("username", username)
		AttachLoginUser(ctx, username)
		ctx.Next()
	}
}
