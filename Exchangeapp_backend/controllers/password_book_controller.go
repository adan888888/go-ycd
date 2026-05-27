package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// 全局数据库实例，从 global 包获取

// CreatePasswordItem 创建密码项
// @Summary 创建密码项
// @Description 创建新的密码项
// @Tags 密码本
// @Accept json
// @Produce json
// @Param request body models.PasswordItemRequest true "密码项信息"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book [post]
func CreatePasswordItem(c *gin.Context) {
	uid, ok := requireLoginUID(c)
	if !ok {
		return
	}

	var req models.PasswordItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "请求参数错误: " + err.Error(),
			Data:   nil,
		})
		return
	}

	// 创建密码项
	passwordItem := models.PasswordItem{
		Uid:      uid,
		Title:    req.Title,
		Username: req.Username,
		Password: req.Password,
		Website:  req.Website,
		Notes:    req.Notes,
	}

	if err := global.Db.Create(&passwordItem).Error; err != nil {
		ServerFail(c, ResponseJson{
			Status: 500,
			Code:   1,
			Msg:    "创建密码项失败: " + err.Error(),
			Data:   nil,
		})
		return
	}

	// 返回响应
	response := models.PasswordItemResponse(passwordItem)

	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "创建成功",
		Data:   response,
	})
}

// GetPasswordItems 获取密码列表
// @Summary 获取密码列表
// @Description 获取所有密码项，支持搜索
// @Tags 密码本
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book [get]
func GetPasswordItems(c *gin.Context) {
	// 获取搜索关键词
	keyword := c.Query("keyword")

	// 超管查全部；普通用户仅本人
	query := global.Db.Model(&models.PasswordItem{})
	var ok bool
	query, ok = applyOwnedDataScope(c, query, "uid")
	if !ok {
		return
	}

	// 添加搜索条件
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR username LIKE ? OR website LIKE ? OR notes LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// 获取数据
	var passwordItems []models.PasswordItem
	if err := query.Order("updated_at ASC").Find(&passwordItems).Error; err != nil {
		ServerFail(c, ResponseJson{
			Status: 500,
			Code:   1,
			Msg:    "查询失败: " + err.Error(),
			Data:   nil,
		})
		return
	}

	// 转换为响应格式
	var items []models.PasswordItemResponse
	for _, item := range passwordItems {
		items = append(items, models.PasswordItemResponse(item))
	}
	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "查询成功",
		Data:   items,
	})
}

// GetPasswordItem 获取单个密码项
// @Summary 获取单个密码项
// @Description 根据ID获取密码项详情
// @Tags 密码本
// @Produce json
// @Param id path int true "密码项ID"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book/{id} [get]
func GetPasswordItem(c *gin.Context) {
	// 获取密码项ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "无效的ID",
			Data:   nil,
		})
		return
	}

	// 查询密码项
	var passwordItem models.PasswordItem
	if err := global.Db.Where("id = ?", id).First(&passwordItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Fail(c, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "密码项不存在",
				Data:   nil,
			})
		} else {
			ServerFail(c, ResponseJson{
				Status: 500,
				Code:   1,
				Msg:    "查询失败: " + err.Error(),
				Data:   nil,
			})
		}
		return
	}

	// 返回响应
	response := models.PasswordItemResponse(passwordItem)

	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "查询成功",
		Data:   response,
	})
}

// UpdatePasswordItem 更新密码项
// @Summary 更新密码项
// @Description 更新密码项信息
// @Tags 密码本
// @Accept json
// @Produce json
// @Param id path int true "密码项ID"
// @Param request body models.PasswordItemRequest true "密码项信息"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book/{id} [put]
func UpdatePasswordItem(c *gin.Context) {
	// 获取密码项ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "无效的ID",
			Data:   nil,
		})
		return
	}

	// 绑定请求数据
	var req models.PasswordItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "请求参数错误: " + err.Error(),
			Data:   nil,
		})
		return
	}

	// 查询密码项是否存在
	var passwordItem models.PasswordItem
	if err := global.Db.Where("id = ?", id).First(&passwordItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Fail(c, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "密码项不存在",
				Data:   nil,
			})
		} else {
			ServerFail(c, ResponseJson{
				Status: 500,
				Code:   1,
				Msg:    "查询失败: " + err.Error(),
				Data:   nil,
			})
		}
		return
	}

	// 更新密码项
	passwordItem.Title = req.Title
	passwordItem.Username = req.Username
	passwordItem.Password = req.Password
	passwordItem.Website = req.Website
	passwordItem.Notes = req.Notes

	if err := global.Db.Save(&passwordItem).Error; err != nil {
		ServerFail(c, ResponseJson{
			Status: 500,
			Code:   1,
			Msg:    "更新失败: " + err.Error(),
			Data:   nil,
		})
		return
	}

	// 返回响应
	response := models.PasswordItemResponse(passwordItem)

	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "更新成功",
		Data:   response,
	})
}

// DeletePasswordItem 删除密码项
// @Summary 删除密码项
// @Description 根据ID删除密码项
// @Tags 密码本
// @Produce json
// @Param id path int true "密码项ID"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book/{id} [delete]
func DeletePasswordItem(c *gin.Context) {
	// 获取密码项ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "无效的ID",
			Data:   nil,
		})
		return
	}

	// 删除密码项
	result := global.Db.Where("id = ?", id).Delete(&models.PasswordItem{})
	if result.Error != nil {
		ServerFail(c, ResponseJson{
			Status: 500,
			Code:   1,
			Msg:    "删除失败: " + result.Error.Error(),
			Data:   nil,
		})
		return
	}

	if result.RowsAffected == 0 {
		Fail(c, ResponseJson{
			Status: http.StatusNotFound,
			Code:   1,
			Msg:    "密码项不存在",
			Data:   nil,
		})
		return
	}

	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "删除成功",
		Data:   nil,
	})
}

// BatchDeletePasswordItems 批量删除密码项
// @Summary 批量删除密码项
// @Description 根据ID列表批量删除密码项
// @Tags 密码本
// @Accept json
// @Produce json
// @Param request body []int true "密码项ID列表"
// @Success 200 {object} models.JSONResult
// @Router /api/password-book/batch-delete [post]
func BatchDeletePasswordItems(c *gin.Context) {
	// 绑定请求数据
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "请求参数错误: " + err.Error(),
			Data:   nil,
		})
		return
	}

	if len(ids) == 0 {
		Fail(c, ResponseJson{
			Status: 400,
			Code:   1,
			Msg:    "ID列表不能为空",
			Data:   nil,
		})
		return
	}

	// 批量删除
	result := global.Db.Where("id IN ?", ids).Delete(&models.PasswordItem{})
	if result.Error != nil {
		ServerFail(c, ResponseJson{
			Status: 500,
			Code:   1,
			Msg:    "批量删除失败: " + result.Error.Error(),
			Data:   nil,
		})
		return
	}

	Ok(c, ResponseJson{
		Status: 200,
		Code:   0,
		Msg:    "批量删除成功",
		Data: gin.H{
			"deleted_count": result.RowsAffected,
		},
	})
}
