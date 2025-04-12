package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func CreateTables(ctx *gin.Context) {
	uid, _ := strconv.ParseInt(ctx.GetHeader("UserId"), 10, 64) //第二个参数 10 表示字符串是十进制格式。第三个参数 64 表示转换结果的类型为 int64。
	var table1 = models.TableYanchendao1{ColumnBenjin: "5000", ColumnYongJin: "0.95", ColumnMean: "0.08", ColumnRestartIdx: "1", ColumnLiushuiIdx: "1", Uid: uid}
	var table2 models.TableYanchendao2
	//AutoMigrate自动迁移 没有这个的表的时候，用于自动创建数据库表或更新表的结构(不会插入数据)。
	err := global.Db.AutoMigrate(&table1)
	if err != nil {
		panic("failed to migrate database：" + err.Error())
	}
	var count int64 = 0
	global.Db.Model(&table1).Where("uid=?", uid).Count(&count)
	if count <= 0 {
		global.Db.Create(&table1)
	}
	global.Db.AutoMigrate(&table2)
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "删除本页数据成功", Data: table1})
}
func GetTable1(ctx *gin.Context) {
	var tableYanchendao1s []models.TableYanchendao1
	UserId := ctx.GetHeader("UserId")
	if err := global.Db.Where("uid=?", UserId).Find(&tableYanchendao1s).Error; err != nil {
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
	// 获取指定 Header 字段的值
	//userAgent := ctx.GetHeader("User-Agent")
	UserId := ctx.GetHeader("UserId")
	fmt.Println(UserId)
	var tableYanchendao2s []models.TableYanchendao2
	if err := global.Db.Where("user_id=?", UserId).Last(&tableYanchendao2s).Error; err != nil {
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
	//手机上是从零开始的，所以减掉1
	//for i, tableYanchendao2 := range tableYanchendao2s {
	//	tableYanchendao2s[i].ID = tableYanchendao2.ID - 1
	//}
	//time.Sleep(100 * time.Millisecond)
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

var mu sync.Mutex // 定义在函数外，作为全局锁
func InsertTable2(ctx *gin.Context) {
	var tableYanchendao2 models.TableYanchendao2
	mu.Lock()                                                     // 加锁
	defer mu.Unlock()                                             // 确保函数退出时解锁
	if err := ctx.ShouldBindJSON(&tableYanchendao2); err != nil { //移动端不传某个字段这里也不会报错，在结构体里需要加binding:"required"才会报错
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("测试%+v\n", tableYanchendao2)
	tableYanchendao2.ID = 0 //解决运行的过程中会自动给赋值
	//使用你提供的主键值，而不是数据库的自增值 Session(&gorm.Session{FullSaveAssociations: true})（gorm默认会忽略传的值），mysql数据库的特性也是下标从时1开始。 例如我删除一个，再插入一个值，这时候的主键自增的就会少一个值
	//现在继续使用自增的（从数据里可以看出来删除了哪个数据）
	if err := global.Db. /*.Session(&gorm.Session{FullSaveAssociations: true})*/ Create(&tableYanchendao2).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    err.Error(),
			Data:   gin.H{},
		})
		return
	}
	//手机上是从零开始计算的，所以减掉1
	//tableYanchendao2.ID = tableYanchendao2.ID - 1
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

// 重启系统（需要记录重启的位置（行））
func Restart(ctx *gin.Context) {
	uid := ctx.GetHeader("UserId")
	var tableYanchendao1 models.TableYanchendao1
	var tableYanchendao2 models.TableYanchendao2
	// 重启时，清除消数列数据（colmun_shuyingzhi_d=""）
	// 将所有记录的 colmun_shuyingzhi_d 列清空, 必须要加 Where("1 = 1")这个条件
	result := global.Db.Model(&tableYanchendao2).Where("1 = 1").Update("colmun_shuyingzhi_d", "")
	global.Db.Find(&tableYanchendao1)
	if result.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    result.Error.Error(),
			Data:   gin.H{},
		})
		return
	}
	//if err := global.Db.Last(&tableYanchendao1).Error; err != nil {
	//	Fail(ctx, ResponseJson{
	//		Status: http.StatusInternalServerError,
	//		Code:   1,
	//		Msg:    err.Error(),
	//		Data:   gin.H{},
	//	})
	//	return
	//}
	// 从上下文中绑定 JSON 数据
	//var value ValueX
	//if err := ctx.ShouldBindJSON(&value); err != nil {
	//	ctx.JSON(400, gin.H{"error": "Invalid JSON format"})
	//	return
	//}
	//tableYanchendao1.ColumnRestartIdx = value.Index //这个又一直传过来的是空值（这个还要看一下原因）
	//tableYanchendao1.ColumnRestartIdx = strconv.FormatInt(result.RowsAffected, 10) //假如本来就是空串，不会有影响行数

	// 查询表格总行数
	var count int64
	if err := global.Db.Model(&tableYanchendao2).Count(&count).Error; err != nil {
		fmt.Println("Failed to count rows:", err)
		return
	}
	//tableYanchendao1.ColumnRestartIdx = strconv.FormatInt(count, 10)
	//tableYanchendao1.ID = tableYanchendao1.ID + 1
	//if err := global.Db.Create(&tableYanchendao1).Error; err != nil {
	//	Fail(ctx, ResponseJson{
	//		Status: http.StatusInternalServerError,
	//		Code:   1,
	//		Msg:    err.Error(),
	//		Data:   gin.H{},
	//	})
	//	return
	//}
	//E := global.Db.Model(&tableYanchendao1).Where("uid = ?", uid).Update("column_restart_index", strconv.FormatInt(count, 10)) //这个后面后带一个ID的条件， UPDATE `table_yanchendao1` SET `column_restart_index`='739' WHERE uid = '1852251920824012800' AND `id` = 1
	//E := global.Db.Table("table_yanchendao1").Where("uid = ?", uid).Updates(map[string]interface{}{"column_restart_index": (tableYanchendao2)})

	global.Db.Last(&tableYanchendao2)
	E := global.Db.Table("table_yanchendao1").Where("uid = ?", uid).Updates(map[string]interface{}{"column_restart_index": tableYanchendao2.ID})
	if E.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    E.Error.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "重启成功", Data: tableYanchendao1})
}

// 对消数列进行排序
func SortXiaoShu(ctx *gin.Context) {
	var tableYanchendao2s []models.TableYanchendao2
	if err := global.Db.Find(&tableYanchendao2s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{Code: 1, Status: http.StatusNotFound, Msg: err.Error(), Data: gin.H{}})
			return
		}
	}
	// 提取 colmun_shuyingzhi_d 列的数据
	var floats []float64
	// 遍历字符串切片，将每个字符串转换为 float64
	for _, s := range tableYanchendao2s {
		num, err := strconv.ParseFloat(s.ColmunShuyingzhiD, 64) //把数字类型的添加到floats切片中
		if err != nil {
			fmt.Println("Error converting string to float:", err)
			continue
		}
		floats = append(floats, num)
	}
	if len(floats) > 0 {
		// 对浮点数切片进行排序
		sort.Float64s(floats)
		var slice1 []float64
		for i := 0; i < len(tableYanchendao2s)-len(floats); i++ {
			slice1 = append(slice1, 1234567.8)
		}
		// 在开头插入元素
		floats = append(slice1, floats...)
		//更新
		for i, _ := range tableYanchendao2s {
			if floats[i] == 1234567.8 {
				tableYanchendao2s[i].ColmunShuyingzhiD = "" //如果是“”空的字符串，db.update不会起效
			} else {
				tableYanchendao2s[i].ColmunShuyingzhiD = strconv.FormatFloat(floats[i], 'f', -1, 64) //转换为科学计数法字符串,-1 表示保留尽可能多的位数。'E': 科学计数法（大写 E）。
			}
		}
	}

	// 更新数据库中的数据
	for _, v := range tableYanchendao2s {
		global.Db.Model(&models.TableYanchendao2{}).Select("colmun_shuyingzhi_d").Where("id=?", v.ID).Updates(v) //要使用Select指定，空值才会更新
		//global.Db .Save(&v) //，这个方法不稳，感觉还是key造成的 或者数据太多操作的太快 底层判断不过来要加事务，Save 方法会更新结构体的所有字段 如果key相同就是update如果没有就是插入数据
	}
	Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "更新数据成功", Data: gin.H{}})
	//tableYanchendao2s[0].ColmunShuyingzhiD = "测试"
	//global.Db.Save(&tableYanchendao2s[0])// 总体测试下来，是需要自动生成的id才可以更新
}

// 消数
func Xiaoshu(ctx *gin.Context) {
	//if ctx.Request.ContentLength == 0 { //ShouldBindJSON如果不传这里也不会报错，1.所以要加这个判断， 2.另外加binding:"required"
	//	Fail(ctx, ResponseJson{
	//		Code:   1,
	//		Status: http.StatusBadRequest,
	//		Msg:    "请求体不能为空",
	//		Data:   gin.H{},
	//	})
	//	return
	//}
	var tableYanchendao2 models.TableYanchendao2
	if err := ctx.ShouldBindJSON(&tableYanchendao2); err != nil { // ShouldBindJSON如果不传这里也不会报错
		Fail(ctx, ResponseJson{
			Code:   1,
			Status: http.StatusInternalServerError,
			Msg:    "输入数据错误",
			Data:   gin.H{},
		})
		return
	}
	if tableYanchendao2.ColmunShuyingzhiD == "" && tableYanchendao2.ColumnXiazhujine != "" {
		// 更新数据库中的数据
		global.Db.Model(&tableYanchendao2).Select("colmun_shuyingzhi_d").Where("id=?", tableYanchendao2.ID).Updates(tableYanchendao2)
		Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "更新数据成功", Data: gin.H{}})
	}

}

// 删除本页
func DeleteAll(ctx *gin.Context) {
	UserId := ctx.GetHeader("UserId")
	// 调用 Delete 方法并传入一个空的 TableYanchendao1 结构体指针，这会删除 user 表中的所有记录。
	//1. 删除 TableYanchendao1 表中的所有数据
	result := global.Db.Where("uid=?", UserId).Delete(&models.TableYanchendao1{})
	if result.Error != nil {
		panic(result.Error)
	}

	result1 := global.Db.Where("user_id=?", UserId).Delete(&models.TableYanchendao2{})
	if result1.Error != nil {
		panic(result1.Error)
	}
	// 输出受影响的行数
	println("Deleted rows:", result.RowsAffected, result1.RowsAffected)

	//2.使用sql语句
	//result1 := global.Db.Exec("DELETE FROM table_yanchendao1")
	//result2 := global.Db.Exec("DELETE FROM table_yanchendao2")
	// 输出受影响的行数
	//println("Deleted rows:", result1.RowsAffected)
	//println("Deleted rows:", result2.RowsAffected)

	//直接删除表
	/*if err := global.Db.Migrator().DropTable(&models.TableYanchendao1{}); err != nil {
		panic(err)
	}
	if err := global.Db.Migrator().DropTable(&models.TableYanchendao2{}); err != nil {
		panic(err)
	}*/

	//time.Sleep(2 * time.Second)
	CreateTables(ctx)
}

// 重置流水
func ResetLiushui(ctx *gin.Context) {
	type TempValuse struct {
		//* 表示该字段是指针类型；不加 * 则表示该字段是值类型
		ResetIndex *int `json:"resetIndex"` //ResetIndex一定要大写要不然赋不了值
	}
	var temp TempValuse
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		return
	}
	if temp.ResetIndex != nil {
		fmt.Printf("前端传递的 age 值为: %d\n", *temp.ResetIndex)
	} else {
		fmt.Println("Mean 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnLiushuiIdx = strconv.Itoa(*temp.ResetIndex)
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		panic(tx.Error)
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "重置流水",
		Data:   nil,
	})
}

// 修改期望值
func Updateqiwangvalue(ctx *gin.Context) {
	type TempValuse struct {
		//* 表示该字段是指针类型；不加 * 则表示该字段是值类型
		Mean *string `json:"mean"` //ResetIndex一定要大写要不然赋不了值
	}
	var temp TempValuse
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		return
	}
	if temp.Mean != nil {
		fmt.Printf("前端传递的 age 值为: %d\n", *temp.Mean)
	} else {
		fmt.Println("Mean 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("user_id=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnMean = *temp.Mean
	/*UPDATE `table_yanchendao1` SET `column_benjin`='5000',`column_yongjin`='0.95',`column_mean`='1',`column_restart_index`='0',`column_liushui_index`='26',`created_at`='2025-03-26 15:11:32' WHERE `id` = 1*/
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		panic(tx.Error)
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改期望值成功",
		Data:   gin.H{"mean": *temp.Mean},
	})
}

// 修改赔率
func UpdateOdds(ctx *gin.Context) {
	type TempValuse struct {
		//* 表示该字段是指针类型；不加 * 则表示该字段是值类型
		Odds *string `json:"odds"` //ResetIndex一定要大写要不然赋不了值
	}
	var temp TempValuse
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		return
	}
	if temp.Odds != nil {
		fmt.Printf("前端传递的 Benjin 值为: %d\n", *temp.Odds)
	} else {
		fmt.Println("Benjin 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("user_id=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnYongJin = *temp.Odds
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		panic(tx.Error)
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改赔率成功",
		Data:   gin.H{"odds": *temp.Odds},
	})
}

// 修改本金
func UpdateBenjin(ctx *gin.Context) {
	type TempValuse struct {
		//* 表示该字段是指针类型；不加 * 则表示该字段是值类型
		Benjin *string `json:"benjin"` //ResetIndex一定要大写要不然赋不了值
	}
	var temp TempValuse
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		return
	}
	if temp.Benjin != nil {
		fmt.Printf("前端传递的 Benjin 值为: %d\n", *temp.Benjin)
	} else {
		fmt.Println("Benjin 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnBenjin = *temp.Benjin
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		panic(tx.Error)
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改本金成功",
		Data:   tableYanchendao1,
	})
}

// 加载更多历史数据
func LoadMore(ctx *gin.Context) {
	//http://localhost:3000/api/LoadMore?last_id=836&c=10&uid=1852251920824012800
	//http://localhost:8080/user?name=张三&age=100&addr=广东  //这种传值用ctx.Query
	//http://localhost:3000/api/testmq/你好  //这种用Param id := ctx.Param("msg")

	lv := ctx.Query("last_id")
	c := ctx.Query("c")
	uid := ctx.Query("uid")
	var tableYanchendao2s []models.TableYanchendao2
	if lv == "-1" {
		global.Db.Raw(`
        SELECT *
        FROM (
            SELECT *
            FROM table_yanchendao2
            WHERE  user_id = ? 
            ORDER BY created_at DESC
            LIMIT ?
        ) AS subquery
        ORDER BY created_at ASC;`, uid, c).Scan(&tableYanchendao2s)
	} else {
		result := global.Db.Raw(`
        SELECT *
        FROM (
            SELECT *
            FROM table_yanchendao2
            WHERE id < ? AND user_id = ? 
            ORDER BY created_at DESC
            LIMIT ?
        ) AS subquery
        ORDER BY created_at ASC;`, lv, uid, c).Scan(&tableYanchendao2s)
		if result.Error != nil {
			fmt.Println("查询出错:", result.Error)
			return
		}
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "加载更多成功",
		Data:   tableYanchendao2s,
	})
}
func GetStatisticalAreasData(ctx *gin.Context) {
	var CurrentTempIndex int64
	var restartIndex int64
	a, err := strconv.ParseInt(ctx.Query("tempIndex"), 10, 64)
	if err != nil {
		println("解释错误", err.Error())
	}
	CurrentTempIndex = a
	var tableYanchendao1 models.TableYanchendao1
	var tableYanchendao2s []models.TableYanchendao2
	statisticalAreas := make([]string, 32) // 定义一个空的字符串切片，类似于 Dart 中的空字符串列表
	UserId := ctx.GetHeader("UserId")
	if tx := global.Db.Where("uid=?", UserId).First(&tableYanchendao1); tx.Error != nil {
		println(tx.Error)
		return
	}
	if err := global.Db.Where("user_id=?", UserId).Find(&tableYanchendao2s).Error; err != nil {
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
	statisticalAreas[0] = tableYanchendao1.ColumnBenjin
	statisticalAreas[1] = strconv.Itoa(len(tableYanchendao2s)) //一共打多少手
	statisticalAreas[19] = tableYanchendao1.ColumnMean         //期望值

	//总体
	var zt_y = 0 //总体赢的次数
	var zt_s = 0 //总体输的次数
	var zt_syz = 0.0
	var runningWater = 0.0
	var countLianShengFu = 1
	var zCount = 0
	var benUse1 = 0
	for index, element := range tableYanchendao2s {
		// 累加输赢值
		shuyingzhiStr := fmt.Sprintf("%v", element.ColmunShuyingzhi)
		shuyingzhi, _ := strconv.ParseFloat(shuyingzhiStr, 64)
		zt_syz += shuyingzhi
		if zt_syz < 0 && zt_syz < float64(benUse1) {
			benUse1 = int(zt_syz)
		}

		// 累加下注金额
		xiazhujineStr := fmt.Sprintf("%v", element.ColumnXiazhujine)
		xiazhujine, _ := strconv.ParseFloat(xiazhujineStr, 64)
		runningWater += xiazhujine

		// 根据备注判断 zt_s 和 zt_y
		if element.ColmunRemark != "" && element.ColmunRemark == "-1" {
			zt_s--
		} else {
			zt_y++
		}

		// 连胜负计算
		if len(tableYanchendao2s) > 1 && index-1 >= 0 {
			prevShuyingzhiStr := fmt.Sprintf("%v", tableYanchendao2s[index-1].ColmunShuyingzhi)
			prevShuyingzhi, _ := strconv.ParseFloat(prevShuyingzhiStr, 64)
			if (shuyingzhi > 0 && prevShuyingzhi > 0) || (shuyingzhi < 0 && prevShuyingzhi < 0) {
				countLianShengFu++
			} else {
				countLianShengFu = 1
			}
		}

		// 庄个数统计
		if element.ColmunZX == "庄" {
			zCount++
		}
	}
	statisticalAreas[5] = strconv.Itoa(zt_y)
	//胜
	statisticalAreas[9] = fmt.Sprintf("%.2f%%", float64(zt_y)/float64(len(tableYanchendao2s))*100) //胜率 ,保留两位小数点. %%两个表示一个
	//winRate := float64(jb_y) / float64(jb_count) * 100

	statisticalAreas[13] = fmt.Sprintf("%d", intAbs(zt_y)-intAbs(zt_s)) //净胜~须多少手回到50%
	statisticalAreas[17] = fmt.Sprintf("%.2f", zt_syz)                  //一共输赢多少钱

	//计算平均赢
	if statisticalAreas[13] == "0" {
		statisticalAreas[21] = "-"
	} else {
		// 移除中文字符
		numStr := removeChineseCharacters(statisticalAreas[13])
		// 将字符串转换为浮点数
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			fmt.Printf("字符串转换为浮点数出错: %v\n", err)
			return
		}
		// 计算平均值并保留两位小数
		averageWin := zt_syz / num
		statisticalAreas[21] = strconv.FormatFloat(averageWin, 'f', 2, 64)
	}
	f, err := strconv.ParseFloat(statisticalAreas[19], 64)
	if err != nil {
		return
	}
	d := float64(len(tableYanchendao2s)+1) * f //期望一共的值
	p := intAbs(intAbs(zt_y) - intAbs(zt_s))
	var result string
	if statisticalAreas[13] == "0" {
		result = "-"
	} else {
		if zt_syz < 0 {
			value := (math.Abs(zt_syz) + d) / float64(p)
			formattedValue := strconv.FormatFloat(value, 'f', 1, 64)
			result = fmt.Sprintf("须%sx%d", formattedValue, p)
		} else {
			value := (math.Abs(zt_syz) - d) / float64(p)
			formattedValue := strconv.FormatFloat(value, 'f', 1, 64)
			result = fmt.Sprintf("可负%sx%d", formattedValue, p)
		}
	}
	statisticalAreas[25] = result //还需要多少 加到50%的时候
	// 处理重启位置
	if len(tableYanchendao2s) > 0 {
		statisticalAreas[29] = tableYanchendao1.ColumnRestartIdx
	}
	// 处理流水索引
	statisticalAreas[8] = tableYanchendao1.ColumnLiushuiIdx
	// 处理本金使用
	statisticalAreas[12] = strconv.Itoa(intAbs(benUse1))
	// 清空索引 16 的值
	statisticalAreas[16] = ""
	// 处理当前金额
	if num, err := strconv.ParseFloat(statisticalAreas[0], 64); err == nil {
		result := num + zt_syz
		statisticalAreas[4] = strconv.FormatFloat(result, 'f', 2, 64)
	}

	//局部
	// 计算重启位置
	if CurrentTempIndex == -1 { //-1时说明没有传tempIndex
		if a, _ := strconv.ParseInt(tableYanchendao1.ColumnRestartIdx, 10, 64); err == nil {
			restartIndex = a
		}
	} else {
		restartIndex = CurrentTempIndex
	}

	jb_y := 0
	jb_s := 0
	jb_syz := 0.0
	jb_count := 0
	// 遍历 table2List 计算局部数据
	for i := 0; i < len(tableYanchendao2s); i++ {
		if CurrentTempIndex == -1 {
			if tableYanchendao2s[i].ID > int(restartIndex) { //重启的时候，要从下一行计算
				jb_count++
				shuyingzhiStr := fmt.Sprintf("%v", tableYanchendao2s[i].ColmunShuyingzhi)
				shuyingzhi, _ := strconv.ParseFloat(shuyingzhiStr, 64)
				jb_syz += shuyingzhi
				if tableYanchendao2s[i].ColmunRemark != "" && strings.HasPrefix(tableYanchendao2s[i].ColmunRemark, "-1") {
					jb_s--
				} else {
					jb_y++
				}
			}
		} else {
			if tableYanchendao2s[i].ID >= int(restartIndex) { //点某一行的时候，要从当前行计算
				jb_count++
				shuyingzhiStr := fmt.Sprintf("%v", tableYanchendao2s[i].ColmunShuyingzhi)
				shuyingzhi, _ := strconv.ParseFloat(shuyingzhiStr, 64)
				jb_syz += shuyingzhi
				if tableYanchendao2s[i].ColmunRemark != "" && strings.HasPrefix(tableYanchendao2s[i].ColmunRemark, "-1") {
					jb_s--
				} else {
					jb_y++
				}
			}
		}

	}
	// 计算一共打多少手
	statisticalAreas[2] = strconv.Itoa(jb_count)
	// 填充局部统计数据到 totalValue
	statisticalAreas[6] = strconv.Itoa(jb_y)
	if jb_count == 0 {
		statisticalAreas[10] = ""
	} else {
		winRate := float64(jb_y) / float64(jb_count) * 100
		statisticalAreas[10] = fmt.Sprintf("%.2f%%", winRate)
	}
	statisticalAreas[14] = strconv.Itoa(jb_y - jb_s)
	statisticalAreas[18] = fmt.Sprintf("%.3f", jb_syz)
	if statisticalAreas[14] == "0" {
		statisticalAreas[22] = "-"
	} else {
		parseStr := removeChineseCharacters(statisticalAreas[14])
		parse, _ := strconv.ParseFloat(parseStr, 64)
		statisticalAreas[22] = fmt.Sprintf("%.3f", jb_syz/parse)
	}
	// 计算期望一共的值
	num19, _ := strconv.ParseFloat(statisticalAreas[19], 64)
	dJ := float64(jb_count+1) * num19                                     //期望一共多少
	parse, _ := strconv.Atoi(strings.TrimLeft(statisticalAreas[14], "-")) //净胜

	if statisticalAreas[14] == "0" {
		statisticalAreas[26] = "-"
	} else if jb_syz < 0 {
		if parse == 0 {
			statisticalAreas[26] = ""
		} else {
			statisticalAreas[26] = fmt.Sprintf("须%.1fx%d", (jb_syz*-1+dJ)/float64(parse), parse)
		}
	} else {
		if parse == 0 {
			statisticalAreas[26] = ""
		} else {
			statisticalAreas[26] = fmt.Sprintf("可负%.1fx%d", (jb_syz-dJ)/float64(parse), parse)
		}
	}
	// 填充第四列数据
	statisticalAreas[3] = fmt.Sprintf("流水%.0f", runningWater)
	if len(tableYanchendao2s) > 0 {
		statisticalAreas[7] = fmt.Sprintf("均利%.2f", zt_syz/float64(len(tableYanchendao2s)))
	}
	statisticalAreas[11] = fmt.Sprintf("连胜负%d", countLianShengFu)
	num1, _ := strconv.Atoi(statisticalAreas[1])
	statisticalAreas[15] = fmt.Sprintf("%d/%d", zCount, num1)
	if len(tableYanchendao2s) > 0 {
		statisticalAreas[23] = tableYanchendao1.ColumnYongJin
	}
	if statisticalAreas[14] == "0" {
		statisticalAreas[27] = ""
	} else if statisticalAreas[21] == "-" {
		statisticalAreas[27] = ""
	} else {
		parts := strings.Split(removeChineseCharacters(statisticalAreas[25]), "x")
		if len(parts) > 0 {
			num25, _ := strconv.ParseFloat(parts[0], 64)
			num23, _ := strconv.ParseFloat(statisticalAreas[23], 64)
			statisticalAreas[27] = fmt.Sprintf("%.2f", num25/num23)
		}
	}
	if statisticalAreas[14] == "0" {
		statisticalAreas[31] = ""
	} else if statisticalAreas[22] == "-" {
		statisticalAreas[31] = ""
	} else {
		parts := strings.Split(removeChineseCharacters(statisticalAreas[26]), "x")
		if len(parts) > 0 {
			num26, _ := strconv.ParseFloat(parts[0], 64)
			num23, _ := strconv.ParseFloat(statisticalAreas[23], 64)
			statisticalAreas[31] = fmt.Sprintf("%.2f", num26/num23)
		}
	}

	// 预测平均值. 手机上做
	//textEditingControllerText := ""
	//if textEditingControllerText != "" {
	//	statisticalAreas[20] = pVal1(tableYanchendao1)
	//	statisticalAreas[24] = pVal2()
	//}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "统计数据",
		Data:   statisticalAreas,
	})
}

// 一对多， 多表关联查询
func Getusers(ctx *gin.Context) {
	targetUid := int64(1852251920824012800)
	var user models.User
	// 预加载并指定查询条件
	result := global.Db.Preload("TableYanchendao1s"). /*, func(db *gorm.DB) *gorm.DB {
			return db.Where("uid = ?", targetUid)
		}*/Preload("TableYanchendao2s"). /*, func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ?", targetUid)
		}*/Where("uid = ?", targetUid).First(&user)

	if result.Error != nil {
		fmt.Printf("failed to query user: %v\n", result.Error)
		return
	}
	//还有一种自己拼接
	/*// 预加载关联数据进行查询，使用 Uid 进行关联
	  targetUid := int64(1852251920824012800)
	  var user models.User
	  // 手动构建查询逻辑
	  result := global.Db.Where("uid = ?", targetUid).First(&user)
	  if result.Error != nil {
	  	fmt.Printf("failed to query user: %v\n", result.Error)
	  	return
	  }
	  // 手动查询关联的 TableYanchendao1 数据
	  var tableYanchendao1s []models.TableYanchendao1
	  result = global.Db.Where("user_id = ?", targetUid).Find(&tableYanchendao1s)
	  if result.Error != nil {
	  	fmt.Printf("failed to query TableYanchendao1: %v\n", result.Error)
	  	return
	  }
	  user.TableYanchendao1s = tableYanchendao1s
	*/
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data:   user,
	})
}

// intAbs 自定义函数，用于求整数的绝对值
func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// removeChineseCharacters 移除字符串中的中文字符
func removeChineseCharacters(s string) string {
	re := regexp.MustCompile(`[\p{Han}]+`)
	return re.ReplaceAllString(s, "")
}
