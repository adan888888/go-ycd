package controllers

import (
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
		FailMsg(ctx, err.Error())
		return
	}

	if err := global.Db.AutoMigrate(&article); err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	OkMsg(ctx, "创建成功", article)
}

func GetArticles(ctx *gin.Context) {
	cachedData, err := global.RedisDB.Get(cacheKey).Result()

	if err == redis.Nil {
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFoundMsg(ctx, err.Error())
			} else {
				ServerFailMsg(ctx, err.Error())
			}
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			ServerFailMsg(ctx, err.Error())
			return
		}

		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			ServerFailMsg(ctx, err.Error())
			return
		}

		OkMsg(ctx, "查询成功", articles)

	} else if err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	} else {
		var articles []models.Article

		if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
			ServerFailMsg(ctx, err.Error())
			return
		}
		utils.Logger.Errorf("=%+v", utils.RemoveEscapeChars1(cachedData))
		OkMsg(ctx, "查询成功", articles)
	}
}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article

	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFoundMsg(ctx, "文章不存在")
		} else {
			ServerFailMsg(ctx, err.Error())
		}
		return
	}

	OkMsg(ctx, "查询成功", article)
}
