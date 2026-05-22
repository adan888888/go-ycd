package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminListUsers 超级管理员获取全部用户列表
func AdminListUsers(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	var users []models.User
	if err := global.Db.Unscoped().
		Model(&models.User{}).
		Where("uid IS NOT NULL").
		Order("username ASC").
		Find(&users).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户列表失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	type userItem struct {
		UserID    string `json:"user_id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	}
	items := make([]userItem, 0, len(users))
	for _, u := range users {
		items = append(items, userItem{
			UserID:    strconv.FormatInt(u.Uid, 10),
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data:   items,
	})
}

// AdminDeleteUser 超级管理员删除用户（不可删除 Admin）
func AdminDeleteUser(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	uid, err := parsePathUID(ctx.Param("uid"))
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: err.Error(), Data: gin.H{}})
		return
	}

	var user models.User
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{Code: 1, Msg: "用户不存在", Data: gin.H{}})
			return
		}
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "查询用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	if user.Username == superAdminUsername {
		Fail(ctx, ResponseJson{Code: 1, Msg: "不能删除超级管理员账号", Data: gin.H{}})
		return
	}

	err = global.Db.Transaction(func(tx *gorm.DB) error {
		//Unscoped忽略软删除 会物理删除
		if err := tx.Unscoped().Where("uid = ?", uid).Delete(&models.TableYanchendao1{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("user_id = ?", uid).Delete(&models.TableYanchendao2{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Where("uid = ?", uid).Delete(&models.User{}).Error
	})
	if err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "删除用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "删除成功",
		Data:   gin.H{"user_id": strconv.FormatInt(uid, 10)},
	})
}

// AdminUpdateUsername 超级管理员修改用户名
func AdminUpdateUsername(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	uid, err := parsePathUID(ctx.Param("uid"))
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: err.Error(), Data: gin.H{}})
		return
	}

	var body struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: "参数错误: " + err.Error(), Data: gin.H{}})
		return
	}
	newName := strings.TrimSpace(body.Username)
	if newName == "" {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户名不能为空", Data: gin.H{}})
		return
	}

	var user models.User
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{Code: 1, Msg: "用户不存在", Data: gin.H{}})
			return
		}
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "查询用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	if user.Username == superAdminUsername && newName != superAdminUsername {
		Fail(ctx, ResponseJson{Code: 1, Msg: "不能修改超级管理员账号的用户名", Data: gin.H{}})
		return
	}

	var exists models.User
	if err := global.Db.Where("username = ? AND uid <> ?", newName, uid).First(&exists).Error; err == nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户名已被占用", Data: gin.H{}})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "校验用户名失败: " + err.Error(), Data: gin.H{}})
		return
	}

	if err := global.Db.Model(&user).Update("username", newName).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "修改用户名失败: " + err.Error(), Data: gin.H{}})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改成功",
		Data: gin.H{
			"user_id":  strconv.FormatInt(uid, 10),
			"username": newName,
		},
	})
}

// AdminUpdatePassword 超级管理员修改用户密码
func AdminUpdatePassword(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	uid, err := parsePathUID(ctx.Param("uid"))
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: err.Error(), Data: gin.H{}})
		return
	}

	var body struct {
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: "参数错误: " + err.Error(), Data: gin.H{}})
		return
	}
	if len(strings.TrimSpace(body.Password)) < 4 {
		Fail(ctx, ResponseJson{Code: 1, Msg: "密码长度至少 4 位", Data: gin.H{}})
		return
	}

	var user models.User
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{Code: 1, Msg: "用户不存在", Data: gin.H{}})
			return
		}
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "查询用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	hashedPwd, err := utils.HashPassword(body.Password)
	if err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "加密密码失败: " + err.Error(), Data: gin.H{}})
		return
	}

	if err := global.Db.Model(&user).Update("password", hashedPwd).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "修改密码失败: " + err.Error(), Data: gin.H{}})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "密码修改成功",
		Data:   gin.H{"user_id": strconv.FormatInt(uid, 10)},
	})
}

func parsePathUID(uidStr string) (int64, error) {
	uidStr = strings.TrimSpace(uidStr)
	if uidStr == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("用户ID格式错误: %w", err)
	}
	return uid, nil
}
