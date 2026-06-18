package controllers

import (
	"errors"
	"exchangeapp/apicode"
	"exchangeapp/baccarat"
	"exchangeapp/global"
	"exchangeapp/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultCollisionUserID int64 = 1907650735441448960

func baccaratSession(ctx *gin.Context) (*baccarat.Session, bool) {
	uid := loginUID(ctx)
	if uid == 0 {
		Fail(ctx, apicode.CodeUnauthorized, "无法识别当前登录用户")
		return nil, false
	}
	return baccarat.DefaultSessions.Get(uid), true
}

// GetBaccaratState 获取百家乐模拟状态
// @Summary 获取百家乐模拟状态
// @Tags baccarat
// @Router /api/baccarat/state [get]
func GetBaccaratState(ctx *gin.Context) {
	session, ok := baccaratSession(ctx)
	if !ok {
		return
	}
	if !session.Shoe.HasActiveShoe() && !session.AwaitingCutCard {
		session.Shuffle()
	}
	Success(ctx, "ok", session.ToDTO())
}

// ShuffleBaccaratShoe 洗牌并开始等待切牌
// @Summary 百家乐洗牌
// @Tags baccarat
// @Router /api/baccarat/shuffle [post]
func ShuffleBaccaratShoe(ctx *gin.Context) {
	session, ok := baccaratSession(ctx)
	if !ok {
		return
	}
	session.Shuffle()
	session.ClearHistory()
	Success(ctx, "洗牌完成，请随机切牌位", session.ToDTO())
}

// CutBaccaratCard 随机切牌位
// @Summary 百家乐随机切牌
// @Tags baccarat
// @Router /api/baccarat/cut-card [post]
func CutBaccaratCard(ctx *gin.Context) {
	session, ok := baccaratSession(ctx)
	if !ok {
		return
	}
	if err := session.CutCard(); err != nil {
		Fail(ctx, apicode.CodeParamInvalid, err.Error())
		return
	}
	Success(ctx, "切牌完成，可以发牌", session.ToDTO())
}

// DealBaccaratHand 发一局牌
// @Summary 百家乐发一局
// @Tags baccarat
// @Router /api/baccarat/deal [post]
func DealBaccaratHand(ctx *gin.Context) {
	session, ok := baccaratSession(ctx)
	if !ok {
		return
	}
	if !session.Shoe.HasActiveShoe() && !session.AwaitingCutCard {
		session.Shuffle()
		Fail(ctx, apicode.CodeParamInvalid, "请先随机切牌位", session.ToDTO())
		return
	}
	if session.AwaitingCutCard {
		Fail(ctx, apicode.CodeParamInvalid, "请先随机切牌位", session.ToDTO())
		return
	}
	if session.Shoe.NeedsReshuffle() {
		session.Shuffle()
		Fail(ctx, apicode.CodeParamInvalid, "牌靴需换靴，已自动洗牌，请重新切牌", session.ToDTO())
		return
	}

	hand, err := session.Deal()
	if err != nil {
		Fail(ctx, apicode.CodeParamInvalid, err.Error())
		return
	}

	data := session.ToDTO()
	data.LastSteps = hand.Steps
	Success(ctx, hand.ResultText, data)
}

// ResetBaccaratSession 重置会话（清空历史并洗牌）
// @Summary 重置百家乐模拟
// @Tags baccarat
// @Router /api/baccarat/reset [post]
func ResetBaccaratSession(ctx *gin.Context) {
	uid := loginUID(ctx)
	if uid == 0 {
		Fail(ctx, apicode.CodeUnauthorized, "无法识别当前登录用户")
		return
	}
	session := baccarat.DefaultSessions.Reset(uid)
	session.Shuffle()
	Success(ctx, "已重置", session.ToDTO())
}

// BulkSimulateBaccarat 批量模拟多靴并统计庄闲和
// @Summary 百家乐批量靴模拟
// @Tags baccarat
// @Router /api/baccarat/bulk-simulate [post]
func BulkSimulateBaccarat(ctx *gin.Context) {
	if loginUID(ctx) == 0 {
		Fail(ctx, apicode.CodeUnauthorized, "无法识别当前登录用户")
		return
	}

	var input struct {
		ShoeCount  int  `json:"shoeCount"`
		ExcludeTie bool `json:"excludeTie"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil || input.ShoeCount <= 0 {
		input.ShoeCount = 1000
	}
	if input.ShoeCount > 10000 {
		Fail(ctx, apicode.CodeParamInvalid, "单次最多模拟 10000 靴")
		return
	}

	result := baccarat.SimulateShoes(input.ShoeCount, input.ExcludeTie, 0, -1)
	msg := "模拟完成"
	if input.ExcludeTie {
		msg = "模拟完成（统计已去掉和局）"
	}
	Success(ctx, msg, result)
}

func loadUserZhuangZhanBi(uid int64) (int, error) {
	var table models.TableYanchendao1
	err := global.Db.Where("uid = ?", uid).Last(&table).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 50, nil
	}
	if err != nil {
		return 0, err
	}
	if table.ZhuangZhanBi == 0 {
		return 50, nil
	}
	return table.ZhuangZhanBi, nil
}

// BulkCollisionSimulate 千靴开奖与用户随机庄闲碰撞测胜率
// @Summary 百家乐随机庄闲碰撞
// @Tags baccarat
// @Router /api/baccarat/bulk-collision [post]
func BulkCollisionSimulate(ctx *gin.Context) {
	if loginUID(ctx) == 0 {
		Fail(ctx, apicode.CodeUnauthorized, "无法识别当前登录用户")
		return
	}

	var input struct {
		ShoeCount int    `json:"shoeCount"`
		UserID    string `json:"userId"`
	}
	_ = ctx.ShouldBindJSON(&input)
	if input.ShoeCount <= 0 {
		input.ShoeCount = 1000
	}
	if input.ShoeCount > 10000 {
		Fail(ctx, apicode.CodeParamInvalid, "单次最多模拟 10000 靴")
		return
	}

	uid := defaultCollisionUserID
	if input.UserID != "" {
		parsed, err := strconv.ParseInt(input.UserID, 10, 64)
		if err != nil {
			Fail(ctx, apicode.CodeParamInvalid, "用户ID格式错误")
			return
		}
		uid = parsed
	}

	zhuangZhanBi, err := loadUserZhuangZhanBi(uid)
	if err != nil {
		Fail(ctx, apicode.CodeServerError, "查询用户庄占比失败: "+err.Error())
		return
	}

	result := baccarat.SimulateShoes(input.ShoeCount, false, uid, zhuangZhanBi)
	Success(ctx, "碰撞统计完成", result)
}
