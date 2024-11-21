package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// @Summary      注册
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Param        data body models.UserBody true "传json数据"
// @Success      200  {object}  models.User
// @Router       /api/exchangeRates/articles [post]
func Register(ctx *gin.Context) {
	var user1 models.UserBody
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//先查询是否注册过
	global.Db.Debug().Where("username = ?", user.Username).First(&user1)
	if len(user1.Username) != 0 || user1.Username != "" {
		utils.Logger.Errorln("Username is already taken:", user1)
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "该用户已注册过"})
		return
	}
	hashedPwd, err := utils.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashedPwd
	user.Uid = utils.GetUid() //雪花算法
	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

/*
// Tags         json  //放在哪个类里面
// Accept       json  //接收
// Produce      json //返回
*/
// @Summary      登录
// 不要描述 // @Description  描述
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Param        data body models.UserBody true "json"  #models.UserBody里面的字段一定要大写 要不然生成不了
// @Success      200  {object} models.JSONResult{data=models.User} "成功响应"
// @Router       /api/auth/login [post]
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == (models.User{}) {
			Fail(ctx, ResponseJson{
				Status: http.StatusUnauthorized,
				Code:   0,
				Msg:    "用户名或者密码错误！",
				Data:   gin.H{},
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}
	//生成token
	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//ctx.JSON(http.StatusOK, gin.H{"token": token})
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   1,
		Msg:    "登录成功",
		Data:   gin.H{"token": token, "Uid": user.Uid, "Username": user.Username},
	})
	ctx.SetCookie(
		"token", user.Username,
		3600, //3600秒=1小时
		"/api/auth/", "", true, false)
}
