package controllers

import (
	"errors"
	"exchangeapp/apicode"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/subscription"
	"exchangeapp/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		FailMsg(ctx, err.Error())
		return
	}
	var existing models.User
	if err := global.Db.Unscoped().Where("username = ?", user.Username).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			if err := global.Db.Unscoped().Delete(&existing).Error; err != nil {
				ServerFailMsg(ctx, "清理已删除用户失败: "+err.Error())
				return
			}
		} else {
			utils.Logger.Errorln("Username is already taken:", existing.Username)
			FailMsg(ctx, "该用户已注册过")
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ServerFailMsg(ctx, err.Error())
		return
	}
	hashedPwd, err := utils.HashPassword(user.Password)

	if err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	user.Password = hashedPwd
	user.Uid = utils.GetUid()
	user.Role = models.RoleUser
	exp := time.Now().AddDate(0, 0, 30)
	user.ExpiresAt = &exp
	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	if err := global.Db.AutoMigrate(&user); err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	OkMsg(ctx, "注册成功", gin.H{
		"token":          token,
		"userId":         strconv.FormatInt(user.Uid, 10),
		"nickname":       user.Username,
		"role":           user.Role,
		"is_super_admin": false,
	})
}

/*
// Tags         json  //放在哪个类里面
// Accept       json  //接收
// Produce      json //返回
*/
// @Summary      登录
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Param        data body models.UserBody true "json"
// @Success      200  {object} models.JSONResult{data=models.User} "成功响应"
// @Router       /api/auth/login [post]
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		FailMsg(ctx, err.Error())
		return
	}

	var user models.User

	if err := global.Db.Select("uid", "username", "password", "role", "expires_at").
		Where(usernameCaseSensitiveSQL, input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidDB) {
			Fail(ctx, ResponseJson{
				Code: apicode.CodeFail,
				Msg:  fmt.Errorf("数据库连接失败: %w", err).Error(),
				Data: gin.H{},
			})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			FailMsg(ctx, "用户名或者密码错误！")
			return
		}

		ServerFailMsg(ctx, err.Error())
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		FailMsg(ctx, "密码错误！")
		return
	}

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	data := gin.H{
		"token":          token,
		"userId":         strconv.FormatInt(user.Uid, 10),
		"nickname":       user.Username,
		"role":           models.NormalizeUserRole(user.Role),
		"is_super_admin": models.IsSuperAdminRole(user.Role),
	}
	for k, v := range subscriptionPayload(user) {
		data[k] = v
	}
	OkMsg(ctx, "登录成功", data)
	ctx.SetCookie(
		"token", user.Username,
		3600,
		"/api/auth/", "", true, false)
}

func subscriptionPayload(user models.User) gin.H {
	display, expiresAt, isPermanent := subscription.FormatExpiresAt(user)
	payload := gin.H{
		"is_permanent": isPermanent,
		"ycd_allowed":  subscription.IsYcdAllowed(user),
	}
	if isPermanent {
		payload["expires_at"] = "永久"
	} else if expiresAt != nil {
		payload["expires_at"] = display
	} else {
		payload["expires_at"] = ""
	}
	return payload
}
