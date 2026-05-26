package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/subscription"
	"exchangeapp/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminListUsers 超级管理员获取全部用户列表
func AdminListUsers(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	var users []models.User
	if err := global.Db.Unscoped().Model(&models.User{}).
		Where("uid IS NOT NULL").
		Order("deleted_at IS NULL DESC, username ASC").
		Find(&users).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户列表失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	items := make([]adminUserItem, 0, len(users))
	for _, u := range users {
		items = append(items, buildAdminUserItem(u))
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

	user, ok := findAdminUserByUID(uid)
	if !ok {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户不存在", Data: gin.H{}})
		return
	}
	if user.DeletedAt.Valid {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户已删除", Data: gin.H{}})
		return
	}

	if models.IsSuperAdminRole(user.Role) {
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
		// 用户表软删除（设置 deleted_at，不物理删行）
		return tx.Delete(&user).Error
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

// AdminRestoreUser 超级管理员恢复已软删除的用户
func AdminRestoreUser(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	uid, err := parsePathUID(ctx.Param("uid"))
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: err.Error(), Data: gin.H{}})
		return
	}

	user, ok := findAdminUserByUID(uid)
	if !ok {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户不存在", Data: gin.H{}})
		return
	}
	if !user.DeletedAt.Valid {
		Fail(ctx, ResponseJson{Code: 1, Msg: "用户未删除，无需恢复", Data: gin.H{}})
		return
	}
	if models.IsSuperAdminRole(user.Role) {
		Fail(ctx, ResponseJson{Code: 1, Msg: "不能操作超级管理员账号", Data: gin.H{}})
		return
	}

	if err := global.Db.Unscoped().Model(&user).Update("deleted_at", nil).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "恢复用户失败: " + err.Error(), Data: gin.H{}})
		return
	}
	if err := global.Db.Unscoped().Select("uid", "username", "expires_at", "deleted_at").
		Where("uid = ?", uid).First(&user).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "读取恢复后用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "恢复成功",
		Data:   buildAdminUserItem(user),
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

// AdminUpdateExpiresAt 超级管理员修改用户到期时间（超管账号永久，不可改）
func AdminUpdateExpiresAt(ctx *gin.Context) {
	if !requireSuperAdmin(ctx) {
		return
	}

	uid, err := parsePathUID(ctx.Param("uid"))
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: err.Error(), Data: gin.H{}})
		return
	}

	var body struct {
		ExpiresAt string `json:"expires_at" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: "参数错误: " + err.Error(), Data: gin.H{}})
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

	if models.IsSuperAdminRole(user.Role) {
		Fail(ctx, ResponseJson{Code: 1, Msg: "超级管理员账号为永久有效，无需设置到期时间", Data: gin.H{}})
		return
	}

	exp, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(body.ExpiresAt), time.Local)
	if err != nil {
		Fail(ctx, ResponseJson{Code: 1, Msg: "到期时间格式应为 YYYY-MM-DD HH:mm:ss", Data: gin.H{}})
		return
	}

	if err := global.Db.Model(&user).Select("expires_at").Updates(map[string]interface{}{
		"expires_at": exp,
	}).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "修改到期时间失败: " + err.Error(), Data: gin.H{}})
		return
	}
	// 写库后重新读取，确保返回与 ycd 校验一致
	if err := global.Db.Select("uid", "username", "expires_at").Where("uid = ?", uid).First(&user).Error; err != nil {
		ServerFail(ctx, ResponseJson{Code: 1, Msg: "读取更新后用户失败: " + err.Error(), Data: gin.H{}})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "到期时间修改成功",
		Data: gin.H{
			"user_id":     strconv.FormatInt(uid, 10),
			"expires_at":  subscription.FormatExpiresAtForAdmin(user),
			"ycd_allowed": subscription.IsYcdAllowed(user),
		},
	})
}

type adminUserItem struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
	IsPermanent bool   `json:"is_permanent"`
	YcdAllowed  bool   `json:"ycd_allowed"`
	IsDeleted   bool   `json:"is_deleted"`
	Status      string `json:"status"`
}

func buildAdminUserItem(u models.User) adminUserItem {
	_, _, isPermanent := subscription.FormatExpiresAt(u)
	deleted := u.DeletedAt.Valid
	status := "正常"
	ycdAllowed := subscription.IsYcdAllowed(u)
	if deleted {
		status = "已删除"
		ycdAllowed = false
	}
	return adminUserItem{
		UserID:      strconv.FormatInt(u.Uid, 10),
		Username:    u.Username,
		Role:        models.NormalizeUserRole(u.Role),
		CreatedAt:   u.CreatedAt.Format("2006-01-02 15:04:05"),
		ExpiresAt:   subscription.FormatExpiresAtForAdmin(u),
		IsPermanent: isPermanent,
		YcdAllowed:  ycdAllowed,
		IsDeleted:   deleted,
		Status:      status,
	}
}

func findAdminUserByUID(uid int64) (models.User, bool) {
	var user models.User
	if err := global.Db.Unscoped().Where("uid = ?", uid).First(&user).Error; err != nil {
		return user, false
	}
	return user, true
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
