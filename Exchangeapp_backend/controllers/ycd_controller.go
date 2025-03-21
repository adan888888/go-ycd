package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetTable1(ctx *gin.Context) {
	var tableYanchendao1s []models.TableYanchendao1
	if err := global.Db.Find(&tableYanchendao1s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   0,
				Msg:    err.Error(),
				Data:   gin.H{},
			})
			return
		} else {
			Fail(ctx, ResponseJson{
				Status: http.StatusInternalServerError,
				Code:   0,
				Msg:    err.Error(),
				Data:   gin.H{},
			})
		}
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "查询成功", Data: tableYanchendao1s})
}
func GetTable2(ctx *gin.Context) {
	var tableYanchendao2s []models.TableYanchendao2
	if err := global.Db.Find(&tableYanchendao2s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    err.Error(),
				Data:   gin.H{},
			})
			return
		} else {
			Fail(ctx, ResponseJson{
				Status: http.StatusInternalServerError,
				Code:   1,
				Msg:    err.Error(),
				Data:   gin.H{},
			})
		}
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "查询成功", Data: tableYanchendao2s})
}
func InsertTable1(ctx *gin.Context) {
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Create(&tableYanchendao1).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "插入数据成功", Data: tableYanchendao1})
}
func InsertTable2(ctx *gin.Context) {
	var tableYanchendao2 models.TableYanchendao2

	if err := ctx.ShouldBindJSON(&tableYanchendao2); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.Db.Create(&tableYanchendao2).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "插入数据成功", Data: tableYanchendao2})
}

// 删除最后一行
func DeleteLast(ctx *gin.Context) {
	var tableYanchendao2 models.TableYanchendao2

	if err := global.Db.Last(&tableYanchendao2).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	if err := global.Db.Delete(&tableYanchendao2).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "删除数据成功", Data: tableYanchendao2})
}

// 删除最后一行
func Restart(ctx *gin.Context) {
	var tableYanchendao1 models.TableYanchendao1
	var tableYanchendao2 models.TableYanchendao2
	// 重启时，清除消数列数据（colmun_shuyingzhi_d=""）
	// 将所有记录的 colmun_shuyingzhi_d 列清空, 必须要加 Where("1 = 1")这个条件
	result := global.Db.Debug().Model(&tableYanchendao2).Where("1 = 1").Update("colmun_shuyingzhi_d", "")
	if result.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    result.Error.Error(),
			Data:   gin.H{},
		})
		return
	}
	if err := global.Db.Last(&tableYanchendao1).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	// 从上下文中绑定 JSON 数据
	var value ValueX
	if err := ctx.ShouldBindJSON(&value); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}
	//tableYanchendao1.ColumnRestartIdx = strconv.FormatInt(result.RowsAffected, 10) //假如本来就是空串，不会有影响行数
	tableYanchendao1.ColumnRestartIdx = value.Index //这个又一直传过来的是空值
	// 查询表格总行数
	var count int64
	if err := global.Db.Model(&tableYanchendao2).Count(&count).Error; err != nil {
		fmt.Println("Failed to count rows:", err)
		return
	}
	tableYanchendao1.ColumnRestartIdx = strconv.FormatInt(count, 10)
	tableYanchendao1.ID = tableYanchendao1.ID + 1
	if err := global.Db.Create(&tableYanchendao1).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "删除数据成功", Data: tableYanchendao1})
}

type ValueX struct {
	Index string `json:"index"`
}
