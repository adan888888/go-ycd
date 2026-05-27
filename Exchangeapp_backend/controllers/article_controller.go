package controllers

import (
	"exchangeapp/apicode"
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var article models.Article

	if err := ctx.ShouldBindJSON(&article); err != nil {
		Fail(ctx, apicode.CodeParamInvalid, err.Error())
		return
	}

	if err := global.Db.AutoMigrate(&article); err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	}

	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	}

	Success(ctx, "创建成功", article)
}

func GetArticles(ctx *gin.Context) {
	cachedData, err := global.RedisDB.Get(cacheKey).Result()

	if err == redis.Nil {
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				Fail(ctx, apicode.CodeNotFound, err.Error())
			} else {
				Fail(ctx, apicode.CodeServerError, err.Error())
			}
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			Fail(ctx, apicode.CodeServerError, err.Error())
			return
		}

		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			Fail(ctx, apicode.CodeServerError, err.Error())
			return
		}

		Success(ctx, "查询成功", articles)

	} else if err != nil {
		Fail(ctx, apicode.CodeServerError, err.Error())
		return
	} else {
		var articles []models.Article

		if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
			Fail(ctx, apicode.CodeServerError, err.Error())
			return
		}
		utils.Logger.Errorf("=%+v", utils.RemoveEscapeChars1(cachedData))
		Success(ctx, "查询成功", articles)
	}
}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article

	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, apicode.CodeNotFound, "文章不存在")
		} else {
			Fail(ctx, apicode.CodeServerError, err.Error())
		}
		return
	}

	Success(ctx, "查询成功", article)
}
