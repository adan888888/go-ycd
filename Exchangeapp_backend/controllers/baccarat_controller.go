package controllers

import (
	"exchangeapp/apicode"
	"exchangeapp/baccarat"

	"github.com/gin-gonic/gin"
)

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
