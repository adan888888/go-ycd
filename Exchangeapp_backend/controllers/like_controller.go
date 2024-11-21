package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}

// 多层嵌套
// @Summary      获取Banner图列表
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Success      200  {object} models.JSONResult{data=models.Banners{bannerx=[]models.Banner}} "成功响应"
// @Router       /api/banners [get]
func GetBanners(ctx *gin.Context) {
	var banners []models.Banner
	if err := global.Db.Debug().Find(&banners).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("没有查到数据")
			// 这里可以自定义错误信息
			err = errors.New("没有查到数据")
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var banners1 []string
	for _, url := range banners {
		banners1 = append(banners1, url.Url)
	}
	/*
		ctx.JSON(http.StatusOK,
		gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": gin.H{
				"banners": banners1,
			},
		})*/

	Ok(ctx, ResponseJson{
		Status: http.StatusOK, //不传话主不会返回给客户端 因为有  `json:"-"`  这个标签
		Code:   10000,
		Msg:    "查询成功",
		Data: map[string]any{
			"banners": banners1,
		},
	})
}

// @Summary      获取热门图列表
// @Tags         接口文档
// @Accept       json
// @Produce      json
// @Success      200  {object}  []string
// @Router       /api/hotgames [get]
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

	Ok(ctx, ResponseJson{
		Status: http.StatusOK, //不传话主不会返回给客户端 因为有  `json:"-"`  这个标签
		Code:   10000,
		Msg:    "查询成功",
		Data: map[string]any{
			"hotgames": hotgames,
		},
	})
}
