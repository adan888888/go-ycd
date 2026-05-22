package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// @Summary      注册
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Param        data body models.UserBody true "传json数据"
// @Success      200  {object}  models.User
// @Router       /api/exchangeRates/articles [post]
func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询是否已注册（含软删除记录，便于释放用户名后重新注册）
	var existing models.User
	if err := global.Db.Unscoped().Where("username = ?", user.Username).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			if err := global.Db.Unscoped().Delete(&existing).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "清理已删除用户失败: " + err.Error()})
				return
			}
		} else {
			utils.Logger.Errorln("Username is already taken:", existing.Username)
			ctx.JSON(http.StatusNotImplemented, gin.H{"error": "该用户已注册过"})
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	//验证“用户名”
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidDB) {
			Fail(ctx, ResponseJson{
				Status: http.StatusUnauthorized,
				Code:   0,
				Msg:    fmt.Errorf("数据库连接失败: %w", err).Error(),
				Data:   gin.H{},
			})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) /*|| user == (models.User{})*/ {
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

	//验证“密码” --用户传过来的密码和数据库查询的该用户的密码的加密值对比
	if !utils.CheckPassword(input.Password, user.Password) {
		Fail(ctx, ResponseJson{
			Status: http.StatusUnauthorized,
			Code:   0,
			Msg:    "密码错误！",
			Data:   gin.H{},
		})
		return
	}

	//生成token
	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "登录成功",
		Data:   gin.H{"token": token, "userId": strconv.FormatInt(user.Uid, 10), "nickname": user.Username},
	})
	ctx.SetCookie(
		"token", user.Username,
		3600, //3600秒=1小时
		"/api/auth/", "", true, false)
}
