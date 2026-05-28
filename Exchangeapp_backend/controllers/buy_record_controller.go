package controllers

import (
	"exchangeapp/apicode"
	"exchangeapp/global"
	"exchangeapp/models"
	"strconv"
	"strings"
	"time"

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

	// 仅当前登录用户（JWT uid）可查本人持币记录
	query := global.Db.Model(&models.BuyRecord{})
	var ok bool
	query, ok = applyCurrentLoginDataScope(ctx, query, "uid")
	if !ok {
		return
	}

	if currency != "" {
		trimmed := strings.TrimSpace(currency)
		// 与库里写法（如 btc / BTC）无关，不区分大小写
		query = query.Where("LOWER(currency) = LOWER(?)", trimmed)
	}

	// 查询所有数据
	var buyRecords []models.BuyRecord
	if err := query.Order("created_at ASC").Find(&buyRecords).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, "查询买币记录失败: " + err.Error())
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

	Success(ctx, "查询成功", responses)
}

// CreateBuyRecord 录入一条买币/购买记录
func CreateBuyRecord(ctx *gin.Context) {
	uid, ok := requireLoginUID(ctx)
	if !ok {
		return
	}

	var body struct {
		Currency  string  `json:"currency"`
		BuyPrice  float64 `json:"buy_price"`
		BuyAmount float64 `json:"buy_amount"`
		BuyTime   string  `json:"buy_time"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, apicode.CodeParamInvalid, "参数无效: " + err.Error())
		return
	}
	body.Currency = strings.TrimSpace(body.Currency)
	if body.Currency == "" {
		Fail(ctx, apicode.CodeParamInvalid, "币种不能为空")
		return
	}
	if body.BuyTime == "" {
		Fail(ctx, apicode.CodeParamInvalid, "买入时间不能为空")
		return
	}
	var buyAt time.Time
	var err error
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		buyAt, err = time.ParseInLocation(layout, strings.TrimSpace(body.BuyTime), time.Local)
		if err == nil {
			break
		}
	}
	if err != nil {
		Fail(ctx, apicode.CodeParamInvalid, "买入时间格式无效，请使用日期时间或 ISO8601")
		return
	}

	record := models.BuyRecord{
		Uid:       uid,
		Currency:  body.Currency,
		BuyPrice:  body.BuyPrice,
		BuyAmount: body.BuyAmount,
		BuyTime:   buyAt,
	}
	if err := global.Db.Create(&record).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, "保存失败: " + err.Error())
		return
	}

	resp := models.BuyRecordResponse{
		ID:        record.ID,
		Currency:  record.Currency,
		BuyPrice:  record.BuyPrice,
		BuyAmount: record.BuyAmount,
		BuyTime:   record.BuyTime,
		CreatedAt: record.CreatedAt,
	}
	Success(ctx, "录入成功", resp)
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
		Fail(ctx, apicode.CodeParamInvalid, "无效的记录ID")
		return
	}

	// 仅当前登录用户可删本人记录
	var buyRecord models.BuyRecord
	if err := global.Db.First(&buyRecord, id).Error; err != nil {
		Fail(ctx, apicode.CodeNotFound, "买币记录不存在")
		return
	}
	loginUID, ok := requireLoginUID(ctx)
	if !ok {
		return
	}
	if buyRecord.Uid != loginUID {
		Fail(ctx, apicode.CodeForbidden, "无权限删除该记录")
		return
	}

	// 软删除记录
	if err := global.Db.Delete(&buyRecord).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, "删除买币记录失败: " + err.Error())
		return
	}

	Success(ctx, "删除成功", gin.H{})
}
