package controllers

import (
	"exchangeapp/apicode"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		Fail(ctx, apicode.CodeParamInvalid, err.Error())
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	}

	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	}

	Success(ctx, "创建成功", exchangeRate)
}

func GetExchangeRates(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate

	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, apicode.CodeNotFound, err.Error())
		} else {
			Fail(ctx, apicode.CodeServerError, err.Error())
		}
		return
	}
	Success(ctx, "查询成功", exchangeRates)
}
