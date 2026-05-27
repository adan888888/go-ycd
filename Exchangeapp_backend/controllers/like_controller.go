package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	OkMsg(ctx, "点赞成功", gin.H{})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ServerFailMsg(ctx, err.Error())
		return
	}

	OkMsg(ctx, "查询成功", gin.H{"likes": likes})
}

func GetBanners(ctx *gin.Context) {
	var banners []models.Banner
	if err := global.Db.Debug().Find(&banners).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("没有查到数据")
			err = errors.New("没有查到数据")
			return
		}
		ServerFailMsg(ctx, err.Error())
		return
	}
	var banners1 []string
	for _, url := range banners {
		banners1 = append(banners1, url.Url)
	}

	OkMsg(ctx, "查询成功", gin.H{"banners": banners1})
}

func GetHotgames(ctx *gin.Context) {
	var hotgames = []string{
		"https://9f.com/images/game/551201.jpg",
		"https://9f.com/images/game/551205.jpg",
		"https://9f.com/images/game/551206.jpg",
		"https://9f.com/images/game/551208.jpg",
		"https://9f.com/images/game/551209.jpg",
		"https://9f.com/images/game/551210.jpg",
		"https://9f.com/images/game/551212.jpg",
		"https://9f.com/images/game/551216.jpg",
		"https://9f.com/images/game/551301.jpg",
		"https://9f.com/images/game/551338.jpg",
		"https://9f.com/images/game/551339.jpg",
	}

	OkMsg(ctx, "查询成功", gin.H{"hotgames": hotgames})
}
