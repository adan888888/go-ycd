package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      获取买币记录列表
// @Tags         买币记录
// @Accept       json
// @Produce      json
// @Param        currency query string false "币种筛选"
// @Success      200  {object}  models.JSONResult{data=[]models.BuyRecordResponse}
// @Router       /api/buyRecords [get]
func GetBuyRecords(ctx *gin.Context) {
	// 获取筛选参数
	currency := ctx.Query("currency")

	// 构建查询
	query := global.Db

	if currency != "" {
		query = query.Where("currency = ?", currency)
	}

	// 查询所有数据
	var buyRecords []models.BuyRecord
	if err := query.Order("created_at ASC").Find(&buyRecords).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   0,
			Msg:    "查询买币记录失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 转换为响应格式
	var responses []models.BuyRecordResponse
	for _, record := range buyRecords {
		response := models.BuyRecordResponse{
			ID:        record.ID,
			Currency:  record.Currency,
			BuyPrice:  record.BuyPrice,
			BuyAmount: record.BuyAmount,
			BuyTime:   record.BuyTime,
			CreatedAt: record.CreatedAt,
		}
		responses = append(responses, response)
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   1,
		Msg:    "查询成功",
		Data:   responses,
	})
}

// @Summary      删除买币记录
// @Tags         买币记录
// @Accept       json
// @Produce      json
// @Param        id path int true "记录ID"
// @Success      200  {object}  models.JSONResult
// @Router       /api/buyRecords/{id} [delete]
func DeleteBuyRecord(ctx *gin.Context) {
	// 获取记录ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   0,
			Msg:    "无效的记录ID",
			Data:   gin.H{},
		})
		return
	}

	// 查询记录
	var buyRecord models.BuyRecord
	if err := global.Db.First(&buyRecord, id).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusNotFound,
			Code:   0,
			Msg:    "买币记录不存在",
			Data:   gin.H{},
		})
		return
	}

	// 软删除记录
	if err := global.Db.Delete(&buyRecord).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   0,
			Msg:    "删除买币记录失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   1,
		Msg:    "删除成功",
		Data:   gin.H{},
	})
}
