package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	. "exchangeapp/utils"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建/初始化用户数据表
func CreateTables(ctx *gin.Context) {
	// 获取用户ID（添加错误处理）
	uid, err := strconv.ParseInt(ctx.GetHeader("UserId"), 10, 64) //第二个参数 10 表示字符串是十进制格式。第三个参数 64 表示转换结果的类型为 int64。
	if err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "无效的用户ID: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	var table1 = models.TableYanchendao1{
		ColumnBenjin:       "5000",
		ColumnYongJin:      "0.95",
		ColumnMean:         "0.08",
		ColumnRestartIdx:   "1",
		ColumnLiushuiIdx:   "1",
		ColumnZhuangZhanBi: 50, // 默认庄占比50%
		Uid:                uid,
	}
	var table2 models.TableYanchendao2

	// AutoMigrate自动迁移：没有这个表的时候，用于自动创建数据库表或更新表的结构(不会插入数据)
	err = global.Db.AutoMigrate(&table1)
	if err != nil {
		panic("failed to migrate database：" + err.Error())
	}

	// 检查该用户是否已有数据（包括已删除的记录）
	// 使用 Unscoped() 查询包括已软删除的记录，避免重复创建
	var count int64 = 0
	global.Db.Unscoped().Model(&table1).Where("uid=?", uid).Count(&count)
	if count <= 0 {
		global.Db.Create(&table1) //把初始数据插入到数据库中
	}

	// 迁移 table2 表结构
	err = global.Db.AutoMigrate(&table2)
	if err != nil {
		panic("failed to migrate table2：" + err.Error())
	}

	Ok(ctx, ResponseJson{
		Code:   0,
		Status: http.StatusOK,
		Msg:    "重置数据成功",
		Data:   table1,
	})
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

// InsertTable2 插入Table2数据
// 优化：移除了全局锁，依赖MySQL数据库的ACID特性和并发控制机制
// MySQL的InnoDB引擎提供了行级锁和事务隔离，可以安全地处理并发插入
func InsertTable2(ctx *gin.Context) {
	var tableYanchendao2 models.TableYanchendao2

	// JSON解析（不需要锁保护）
	if err := ctx.ShouldBindJSON(&tableYanchendao2); err != nil { //移动端不传某个字段这里也不会报错，在结构体里需要加binding:"required"才会报错
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("测试%+v\n", tableYanchendao2)
	tableYanchendao2.ID = 0 //解决运行的过程中会自动给赋值
	//使用你提供的主键值，而不是数据库的自增值 Session(&gorm.Session{FullSaveAssociations: true})（gorm默认会忽略传的值），mysql数据库的特性也是下标从时1开始。 例如我删除一个，再插入一个值，这时候的主键自增的就会少一个值
	//现在继续使用自增的（从数据里可以看出来删除了哪个数据）

	// 直接执行数据库插入，依赖数据库的并发控制
	// MySQL会自动处理并发插入，使用行级锁保证数据一致性
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

// 重启系统（需要记录重启的位置（行））
func Restart(ctx *gin.Context) {
	uid, _ := strconv.ParseInt(ctx.GetHeader("UserId"), 10, 64)
	var tableYanchendao1 models.TableYanchendao1
	var tableYanchendao2 models.TableYanchendao2
	// 重启时，清除消数列数据（colmun_shuyingzhi_d=""）
	// 将当前用户所有未删除记录的 colmun_shuyingzhi_d 列清空为空字符串
	// GORM 会自动添加 deleted_at IS NULL 条件，无需手动添加
	result := global.Db.Model(&tableYanchendao2).Where("user_id = ?", uid).Update("colmun_shuyingzhi_d", "")
	if result.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    result.Error.Error(),
			Data:   gin.H{},
		})
		return
	}
	global.Db.Last(&tableYanchendao2)
	//E := global.Db.Table("table_yanchendao1").Where("uid = ?", uid).Updates(map[string]interface{}{"column_restart_index": tableYanchendao2.ID})
	//改变需求，要存起来，不是修改，将来要看的见重启的历史
	// 从现有记录复制所有字段值，确保新字段也有正确的值
	if err := global.Db.Where("uid=?", uid).Last(&tableYanchendao1).Error; err != nil {
		// 如果查询失败，使用默认值
		tableYanchendao1.ColumnZhuangZhanBi = 50 // 默认庄占比50%
	}
	tableYanchendao1.Uid = uid
	tableYanchendao1.TempIndex = "-1"
	tableYanchendao1.ColumnRestartIdx = strconv.Itoa(tableYanchendao2.ID)
	// 如果新字段是0（旧数据没有这个字段），使用默认值
	if tableYanchendao1.ColumnZhuangZhanBi == 0 {
		tableYanchendao1.ColumnZhuangZhanBi = 50
	}
	E := global.Db.Table("table_yanchendao1").Omit("id").Create(tableYanchendao1) //.Omit忽略id插入数据
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
	// 按创建时间正序查询，确保最新的记录在数组最后
	if err := global.Db.Where("user_id=?", ctx.GetHeader("UserId")).Order("created_at ASC").Find(&tableYanchendao2s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{Code: 1, Status: http.StatusNotFound, Msg: err.Error(), Data: gin.H{}})
			return
		}
	}

	// 提取 colmun_shuyingzhi_d 列的有效数据（非空且能转换为浮点数）
	var floats []float64
	for _, s := range tableYanchendao2s {
		if s.ColmunShuyingzhiD == "" {
			continue // 跳过空值
		}
		num, err := strconv.ParseFloat(s.ColmunShuyingzhiD, 64)
		if err != nil {
			fmt.Println("Error converting string to float:", err)
			continue
		}
		floats = append(floats, num)
	}

	// 如果有有效值，进行排序
	if len(floats) > 0 {
		// 对浮点数切片进行排序（从小到大）
		sort.Float64s(floats)

		// 先清空所有记录的 colmun_shuyingzhi_d
		for i := range tableYanchendao2s {
			tableYanchendao2s[i].ColmunShuyingzhiD = ""
		}

		// 从最新的记录（数组最后）开始，倒序写入排序后的值
		// 例如：如果有100条记录，其中50条有值，排序后的值应该写入到索引50-99（最新的50条记录）
		floatsIndex := len(floats) - 1 // 从排序后的最后一个值开始
		for i := len(tableYanchendao2s) - 1; i >= 0 && floatsIndex >= 0; i-- {
			tableYanchendao2s[i].ColmunShuyingzhiD = strconv.FormatFloat(floats[floatsIndex], 'f', -1, 64)
			floatsIndex--
		}
	}

	// 使用事务批量更新数据库
	tx := global.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		Fail(ctx, ResponseJson{
			Code:   1,
			Status: http.StatusInternalServerError,
			Msg:    "开启事务失败: " + tx.Error.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 使用 CASE WHEN 语句批量更新，避免 N 次数据库查询
	if len(tableYanchendao2s) > 0 {
		// 构建 CASE WHEN SQL 语句
		var caseWhenSQL strings.Builder
		var ids []interface{}

		caseWhenSQL.WriteString("CASE id ")
		for _, v := range tableYanchendao2s {
			caseWhenSQL.WriteString("WHEN ? THEN ? ")
			ids = append(ids, v.ID, v.ColmunShuyingzhiD)
		}
		caseWhenSQL.WriteString("END")

		// 构建 WHERE 条件
		var whereIDs []interface{}
		var placeholders []string
		for _, v := range tableYanchendao2s {
			placeholders = append(placeholders, "?")
			whereIDs = append(whereIDs, v.ID)
		}

		// 执行批量更新
		userID := ctx.GetHeader("UserId")
		sql := fmt.Sprintf(
			"UPDATE table_yanchendao2 SET colmun_shuyingzhi_d = %s WHERE user_id = ? AND id IN (%s) AND deleted_at IS NULL",
			caseWhenSQL.String(),
			strings.Join(placeholders, ","),
		)

		// 合并所有参数：CASE WHEN 的参数 + user_id + WHERE IN 的参数
		args := append(ids, userID)
		args = append(args, whereIDs...)

		if err := tx.Exec(sql, args...).Error; err != nil {
			tx.Rollback()
			Fail(ctx, ResponseJson{
				Code:   1,
				Status: http.StatusInternalServerError,
				Msg:    "批量更新失败: " + err.Error(),
				Data:   gin.H{},
			})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		Fail(ctx, ResponseJson{
			Code:   1,
			Status: http.StatusInternalServerError,
			Msg:    "提交事务失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 构建排序后的序列（字符串数组，用于返回）
	var sortedSequence []string
	if len(floats) > 0 {
		sortedSequence = make([]string, len(floats))
		for i, v := range floats {
			sortedSequence[i] = strconv.FormatFloat(v, 'f', -1, 64)
		}
	}

	Ok(ctx, ResponseJson{
		Code:   0,
		Status: http.StatusOK,
		Msg:    "排序成功",
		Data: gin.H{
			"sorted_sequence": sortedSequence,      // 排序后的序列（字符串数组）
			"count":           len(sortedSequence), // 排序后的数量
		},
	})
	//tableYanchendao2s[0].ColmunShuyingzhiD = "测试"
	//global.Db.Save(&tableYanchendao2s[0])// 总体测试下来，是需要自动生成的id才可以更新
}

// 消数（单个）
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
		// 只更新colmun_shuyingzhi_d这一列，传入的是空字符串""也会起效
		global.Db.Model(&tableYanchendao2).Select("colmun_shuyingzhi_d").Where("id=?", tableYanchendao2.ID).Updates(tableYanchendao2)
		Ok(ctx, ResponseJson{Code: 0, Status: http.StatusOK, Msg: "消数据成功", Data: gin.H{}})
	}

}

// 删除本页（清空用户所有数据并重新初始化）
func DeleteAll(ctx *gin.Context) {
	UserId := ctx.GetHeader("UserId")

	// 物理删除：真正清空用户的所有数据（包括已软删除的记录）
	// 使用 Unscoped() 跳过软删除机制，执行真正的 DELETE 操作
	result := global.Db.Unscoped().Where("uid=?", UserId).Delete(&models.TableYanchendao1{})
	if result.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "删除数据失败: " + result.Error.Error(),
			Data:   gin.H{},
		})
		return
	}

	result1 := global.Db.Unscoped().Where("user_id=?", UserId).Delete(&models.TableYanchendao2{})
	if result1.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "删除数据失败: " + result1.Error.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 输出受影响的行数
	println("Deleted rows:", result.RowsAffected, result1.RowsAffected)

	// 重新初始化数据
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
	// 如果新字段是0（旧数据没有这个字段），使用默认值，避免覆盖
	if tableYanchendao1.ColumnZhuangZhanBi == 0 {
		tableYanchendao1.ColumnZhuangZhanBi = 50
	}
	tx := global.Db.Save(&tableYanchendao1) //Save 会保存所有的字段，即使字段是零值.必须要保证有主键id，否则是新增数据
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
func UpdateQiWangValue(ctx *gin.Context) {
	type TempValuse struct {
		//* 表示该字段是指针类型；不加 * 则表示该字段是值类型
		Mean *string `json:"mean"` //ResetIndex一定要大写要不然赋不了值
	}
	var temp TempValuse
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		return
	}
	if temp.Mean != nil {
		fmt.Printf("前端传递的 Mean 值为: %s\n", *temp.Mean)
	} else {
		fmt.Println("Mean 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("user_id=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnMean = *temp.Mean
	// 如果新字段是0（旧数据没有这个字段），使用默认值，避免覆盖
	if tableYanchendao1.ColumnZhuangZhanBi == 0 {
		tableYanchendao1.ColumnZhuangZhanBi = 50
	}
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
		fmt.Printf("前端传递的 Benjin 值为: %s\n", *temp.Odds)
	} else {
		fmt.Println("Benjin 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("user_id=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	tableYanchendao1.ColumnYongJin = *temp.Odds
	// 如果新字段是0（旧数据没有这个字段），使用默认值，避免覆盖
	if tableYanchendao1.ColumnZhuangZhanBi == 0 {
		tableYanchendao1.ColumnZhuangZhanBi = 50
	}
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
		fmt.Printf("前端传递的 Benjin 值为: %s\n", *temp.Benjin)
	} else {
		fmt.Println("Benjin 是结构体默认值")
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户数据失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}
	tableYanchendao1.ColumnBenjin = *temp.Benjin
	// 如果新字段是0（旧数据没有这个字段），使用默认值，避免覆盖
	if tableYanchendao1.ColumnZhuangZhanBi == 0 {
		tableYanchendao1.ColumnZhuangZhanBi = 50
	}
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

// 修改庄占比
func UpdateZhuangZhanBi(ctx *gin.Context) {
	type TempValues struct {
		ZhuangZhanBi *int `json:"zhuangZhanBi"` // 庄占比，范围0-100
	}
	var temp TempValues
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "参数错误: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}
	if temp.ZhuangZhanBi == nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "庄占比不能为空",
			Data:   gin.H{},
		})
		return
	}
	// 验证范围 0-100
	if *temp.ZhuangZhanBi < 0 || *temp.ZhuangZhanBi > 100 {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "庄占比必须在0-100之间",
			Data:   gin.H{},
		})
		return
	}
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", ctx.GetHeader("UserId")).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "未找到用户数据",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户数据失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}
	tableYanchendao1.ColumnZhuangZhanBi = *temp.ZhuangZhanBi
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "更新失败: " + tx.Error.Error(),
			Data:   gin.H{},
		})
		return
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改庄占比成功",
		Data:   gin.H{"zhuangZhanBi": tableYanchendao1.ColumnZhuangZhanBi},
	})
}

// 获取用户庄占比（无需认证，通过user_id参数）
// @Summary      获取用户庄占比
// @Tags         ycd投注记录
// @Accept       json
// @Produce      json
// @Param        user_id query string true "用户ID"
// @Success      200  {object}  ResponseJson{data=object}
// @Router       /api/ycd/zhuangzhanbi [get]
func GetZhuangZhanBi(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "用户ID不能为空",
			Data:   gin.H{},
		})
		return
	}

	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", userIDStr).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "未找到用户数据",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户数据失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	zhuangZhanBi := tableYanchendao1.ColumnZhuangZhanBi
	if zhuangZhanBi == 0 {
		zhuangZhanBi = 50 // 默认值
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data: gin.H{
			"zhuangZhanBi": zhuangZhanBi,
			"user_id":      userIDStr,
		},
	})
}

// 更新用户庄占比（无需认证，通过user_id参数）
// @Summary      更新用户庄占比
// @Tags         ycd投注记录
// @Accept       json
// @Produce      json
// @Param        user_id query string true "用户ID"
// @Param        zhuangZhanBi body int true "庄占比(0-100)"
// @Success      200  {object}  ResponseJson{data=object}
// @Router       /api/ycd/zhuangzhanbi [post]
func UpdateZhuangZhanBiPublic(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "用户ID不能为空",
			Data:   gin.H{},
		})
		return
	}

	type TempValues struct {
		ZhuangZhanBi *int `json:"zhuangZhanBi"` // 庄占比，范围0-100
	}
	var temp TempValues
	if err := ctx.ShouldBindJSON(&temp); err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "参数错误: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}
	if temp.ZhuangZhanBi == nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "庄占比不能为空",
			Data:   gin.H{},
		})
		return
	}
	// 验证范围 0-100
	if *temp.ZhuangZhanBi < 0 || *temp.ZhuangZhanBi > 100 {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "庄占比必须在0-100之间",
			Data:   gin.H{},
		})
		return
	}

	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", userIDStr).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "未找到用户数据",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户数据失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	tableYanchendao1.ColumnZhuangZhanBi = *temp.ZhuangZhanBi
	tx := global.Db.Save(&tableYanchendao1)
	if tx.Error != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "更新失败: " + tx.Error.Error(),
			Data:   gin.H{},
		})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "修改庄占比成功",
		Data: gin.H{
			"zhuangZhanBi": tableYanchendao1.ColumnZhuangZhanBi,
			"user_id":      userIDStr,
		},
	})
}

// GetTable1List 获取table_yanchendao1数据列表（无需认证，支持分页）
func GetTable1List(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// 限制每页最大数量
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	var tableYanchendao1s []models.TableYanchendao1
	var total int64
	var query *gorm.DB

	// 如果 user_id 为空，查询所有用户的记录；否则查询指定用户的记录
	if userIDStr == "" {
		// 查询所有用户的记录，按创建时间倒序排列
		query = global.Db.Model(&models.TableYanchendao1{}).Where("deleted_at IS NULL")
	} else {
		// 查询该用户的所有记录，按创建时间倒序排列
		// 将字符串转换为 int64，避免大整数查询问题
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			Fail(ctx, ResponseJson{
				Status: http.StatusBadRequest,
				Code:   1,
				Msg:    "用户ID格式错误: " + err.Error(),
				Data:   gin.H{},
			})
			return
		}
		query = global.Db.Model(&models.TableYanchendao1{}).Where("uid = ? AND deleted_at IS NULL", userID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询总数失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 分页查询（按ID正序排列，保持数据库中的顺序）
	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&tableYanchendao1s).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 获取所有唯一的用户ID
	uidMap := make(map[int64]string) // uid -> username
	for _, item := range tableYanchendao1s {
		if _, exists := uidMap[item.Uid]; !exists {
			uidMap[item.Uid] = ""
		}
	}

	// 批量查询用户名
	if len(uidMap) > 0 {
		var uids []int64
		for uid := range uidMap {
			uids = append(uids, uid)
		}
		var users []models.User
		if err := global.Db.Select("uid, username").Where("uid IN ?", uids).Find(&users).Error; err == nil {
			for _, user := range users {
				uidMap[user.Uid] = user.Username
				// 调试日志：检查 uid 和 username 的对应关系
				fmt.Printf("GetTable1List 调试: uid=%d, username=%s\n", user.Uid, user.Username)
			}
		} else {
			fmt.Printf("GetTable1List 错误: 查询用户名失败: %v\n", err)
		}
	}

	// 构建返回数据，添加用户名，并确保 uid 是字符串
	var resultList []gin.H
	for _, item := range tableYanchendao1s {
		username := uidMap[item.Uid]
		// 调试日志：检查返回数据中的 uid 和 username 对应关系
		if userIDStr != "" {
			fmt.Printf("GetTable1List 返回数据: uid=%d, username=%s, 查询的user_id=%s\n", item.Uid, username, userIDStr)
		}
		resultList = append(resultList, gin.H{
			"id":                    item.ID,
			"uid":                   strconv.FormatInt(item.Uid, 10), // 确保 uid 是字符串
			"username":              username,
			"column_benjin":         item.ColumnBenjin,
			"column_yongJin":        item.ColumnYongJin,
			"column_mean":           item.ColumnMean,
			"column_restart_index":  item.ColumnRestartIdx,
			"column_liushui_index":  item.ColumnLiushuiIdx,
			"column_zhuang_zhan_bi": item.ColumnZhuangZhanBi,
			"temp_index":            item.TempIndex,
			"created_at":            item.CreatedAt,
			"deleted_at":            item.DeletedAt,
		})
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data: gin.H{
			"list":      resultList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTable2List 获取table_yanchendao2数据列表（无需认证，支持分页）
func GetTable2List(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// 限制每页最大数量
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	var tableYanchendao2s []models.TableYanchendao2
	var total int64
	var query *gorm.DB

	// 如果 user_id 为空，查询所有用户的记录；否则查询指定用户的记录
	if userIDStr == "" {
		// 查询所有用户的记录，按创建时间倒序排列
		query = global.Db.Model(&models.TableYanchendao2{}).Where("deleted_at IS NULL")
	} else {
		// 查询该用户的所有记录，按创建时间倒序排列
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			Fail(ctx, ResponseJson{
				Status: http.StatusBadRequest,
				Code:   1,
				Msg:    "用户ID格式错误: " + err.Error(),
				Data:   gin.H{},
			})
			return
		}
		query = global.Db.Model(&models.TableYanchendao2{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询总数失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 分页查询（按创建时间正序排列，最早的数据在前）
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&tableYanchendao2s).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 获取所有唯一的用户ID
	uidMap := make(map[int64]string) // user_id -> username
	for _, item := range tableYanchendao2s {
		if _, exists := uidMap[item.UserID]; !exists {
			uidMap[item.UserID] = ""
		}
	}

	// 批量查询用户名
	if len(uidMap) > 0 {
		var uids []int64
		for uid := range uidMap {
			uids = append(uids, uid)
		}
		var users []models.User
		if err := global.Db.Select("uid, username").Where("uid IN ?", uids).Find(&users).Error; err == nil {
			for _, user := range users {
				uidMap[user.Uid] = user.Username
			}
		}
	}

	// 构建返回数据，添加用户名，并确保 user_id 是字符串
	var resultList []gin.H
	for _, item := range tableYanchendao2s {
		resultList = append(resultList, gin.H{
			"id":                  item.ID,
			"user_id":             strconv.FormatInt(item.UserID, 10), // 确保 user_id 是字符串
			"username":            uidMap[item.UserID],
			"column_xiazhujine":   item.ColumnXiazhujine,
			"colmun_shuyingzhi":   item.ColmunShuyingzhi,
			"colmun_shuyingzhi_d": item.ColmunShuyingzhiD,
			"colmun_shengfulu":    item.ColmunShengfulu,
			"colmun_zx":           item.ColmunZX,
			"colmun_remark":       item.ColmunRemark,
			"column_current_jin":  item.ColumnCurrentJin,
			"created_at":          item.CreatedAt,
			"deleted_at":          item.DeletedAt,
		})
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data: gin.H{
			"list":      resultList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateTable2Config 更新table_yanchendao2的输赢值和消数后输赢值
func UpdateTable2Config(ctx *gin.Context) {
	type UpdateTable2ConfigRequest struct {
		ID                int    `json:"id" binding:"required"` // 记录ID
		ColmunShuyingzhi  string `json:"colmun_shuyingzhi"`     // 输赢值
		ColmunShuyingzhiD string `json:"colmun_shuyingzhi_d"`   // 消数后输赢值
	}

	var req UpdateTable2ConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "参数错误: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	var tableYanchendao2 models.TableYanchendao2
	if err := global.Db.Where("id = ? AND deleted_at IS NULL", req.ID).First(&tableYanchendao2).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "记录不存在",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	updates := make(map[string]interface{})
	if req.ColmunShuyingzhi != "" {
		updates["colmun_shuyingzhi"] = req.ColmunShuyingzhi
	}
	if req.ColmunShuyingzhiD != "" {
		updates["colmun_shuyingzhi_d"] = req.ColmunShuyingzhiD
	}

	if len(updates) > 0 {
		if err := global.Db.Model(&tableYanchendao2).Updates(updates).Error; err != nil {
			Fail(ctx, ResponseJson{
				Status: http.StatusInternalServerError,
				Code:   1,
				Msg:    "更新失败: " + err.Error(),
				Data:   gin.H{},
			})
			return
		}
	}

	// 重新查询以返回最新数据
	global.Db.Where("id = ?", req.ID).First(&tableYanchendao2)

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "更新成功",
		Data: gin.H{
			"id":                  tableYanchendao2.ID,
			"colmun_shuyingzhi":   tableYanchendao2.ColmunShuyingzhi,
			"colmun_shuyingzhi_d": tableYanchendao2.ColmunShuyingzhiD,
		},
	})
}

// 更新 table_yanchendao1 的临时索引和重启位置
func UpdateTable1Config(ctx *gin.Context) {
	type UpdateTable1ConfigRequest struct {
		ID           int    `json:"id" binding:"required"` // 记录ID
		TempIndex    string `json:"temp_index"`            // 临时索引
		RestartIndex string `json:"restart_index"`         // 重启位置
	}

	var req UpdateTable1ConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "参数错误: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 查询记录是否存在
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("id = ? AND deleted_at IS NULL", req.ID).First(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "记录不存在",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 更新字段（只更新传入的字段）
	updates := make(map[string]interface{})
	if req.TempIndex != "" {
		updates["temp_index"] = req.TempIndex
	}
	if req.RestartIndex != "" {
		updates["column_restart_index"] = req.RestartIndex
	}

	// 如果有要更新的字段，执行更新
	if len(updates) > 0 {
		if err := global.Db.Model(&tableYanchendao1).Updates(updates).Error; err != nil {
			Fail(ctx, ResponseJson{
				Status: http.StatusInternalServerError,
				Code:   1,
				Msg:    "更新失败: " + err.Error(),
				Data:   gin.H{},
			})
			return
		}
	}

	// 重新查询更新后的数据
	global.Db.Where("id = ?", req.ID).First(&tableYanchendao1)

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "更新成功",
		Data: gin.H{
			"id":                   tableYanchendao1.ID,
			"temp_index":           tableYanchendao1.TempIndex,
			"column_restart_index": tableYanchendao1.ColumnRestartIdx,
		},
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
            WHERE user_id = ? AND deleted_at IS NULL
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
            WHERE id < ? AND user_id = ? AND deleted_at IS NULL
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

// tempIndex -10000 app第一次进来
// tempIndex -1 取消局部平衡/重启...
// tempIndex -2 确保不会破坏局部平衡-->每次下注记录输赢/改变数据库里面的值等 ...
// tempIndex >2 点击进行局部平衡
func GetStatisticalAreasData(ctx *gin.Context) {
	var CurrentTempIndex int64             //  这个全局变量主要是用来区分有没有点局部平衡
	statisticalAreas := make([]string, 32) // 定义一个空的字符串切片，类似于 Dart 中的空字符串列表
	var restartIndex int64
	var tableYanchendao1 = models.TableYanchendao1{}
	var tableYanchendao2s []models.TableYanchendao2
	UserId := ctx.GetHeader("UserId")

	//从app传过来的 tempIndex
	tempIndex, err := strconv.ParseInt(ctx.Query("tempIndex"), 10, 64)
	if err != nil {
		println("解释错误", err.Error())
	}

	//查询
	if tx := global.Db.Where("uid=?", UserId).Last(&tableYanchendao1); tx.Error != nil {
		println(tx.Error)
		return
	}

	newRecord := tableYanchendao1
	newRecord.ID = 0 // 重置 ID，让数据库自动生成新 ID
	newRecord.TempIndex = strconv.FormatInt(tempIndex, 10)
	//防止取消局部平衡的时候再次存一次
	if (tempIndex == -1 && newRecord.TempIndex != tableYanchendao1.TempIndex) || (tempIndex > 2 && newRecord.TempIndex != tableYanchendao1.TempIndex) {
		if err := global.Db.Save(&newRecord).Error; err != nil {
			println("创建新记录失败:", err.Error())
			return
		}
	}

	// -10000(app退出应用再进来)
	// -2（正常打）
	if (tempIndex == -10000 && tableYanchendao1.TempIndex != "") || (tempIndex == -2 && tableYanchendao1.TempIndex != "") {
		parseInt, _ := strconv.ParseInt(tableYanchendao1.TempIndex, 10, 64)
		tempIndex = parseInt
	}
	CurrentTempIndex = tempIndex

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
	var zhuangCount = 0
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
			zhuangCount++
		}
	}
	statisticalAreas[5] = strconv.Itoa(zt_y)
	//胜
	if len(tableYanchendao2s) == 0 {
		statisticalAreas[9] = ""
	} else {
		statisticalAreas[9] = fmt.Sprintf("%.2f%%", float64(zt_y)/float64(len(tableYanchendao2s))*100) //胜率 ,保留两位小数点. %%两个表示一个
	}
	//winRate := float64(jb_y) / float64(jb_count) * 100

	statisticalAreas[13] = fmt.Sprintf("%d", IntAbs(zt_y)-IntAbs(zt_s)) //净胜~须多少手回到50%
	statisticalAreas[17] = fmt.Sprintf("%.2f", zt_syz)                  //一共输赢多少钱

	//计算平均赢
	if statisticalAreas[13] == "0" {
		statisticalAreas[21] = "-"
	} else {
		// 移除中文字符
		numStr := RemoveChineseCharacters(statisticalAreas[13])
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
	p := IntAbs(IntAbs(zt_y) - IntAbs(zt_s))
	var result string
	if statisticalAreas[13] == "0" {
		result = "-"
	} else {
		if zt_syz < 0 {
			value := (math.Abs(zt_syz) + d) / float64(p)
			formattedValue := strconv.FormatFloat(value, 'f', 0, 64)
			result = fmt.Sprintf(" 须%sx%d ", formattedValue, p)
		} else {
			value := (math.Abs(zt_syz) - d) / float64(p)
			formattedValue := strconv.FormatFloat(value, 'f', 0, 64)
			result = fmt.Sprintf(" 可负%sx%d ", formattedValue, p)
		}
	}
	statisticalAreas[25] = result //还需要多少 加到50%的时候
	// 处理重启位置
	if len(tableYanchendao2s) > 0 {
		//statisticalAreas[29] = tableYanchendao1.ColumnRestartIdx
		statisticalAreas[29] = "" //不要这个值了吧
	}
	// 处理本金使用
	statisticalAreas[8] = strconv.Itoa(IntAbs(benUse1))
	// 清空索引 16 的值
	statisticalAreas[16] = ""
	// 处理当前金额
	if num, err := strconv.ParseFloat(statisticalAreas[0], 64); err == nil {
		result := num + zt_syz
		statisticalAreas[4] = strconv.FormatFloat(result, 'f', 2, 64)
	}

	// 局部
	// 计算重启位置

	//局部平衡是必须是>2的数，如果是重启传过来是-1，每次下注计算是-2（防止打一手就把局部平衡破坏了）
	if CurrentTempIndex > 2 {
		restartIndex = CurrentTempIndex
	} else {
		if restartIdx, err := strconv.ParseInt(tableYanchendao1.ColumnRestartIdx, 10, 64); err == nil {
			restartIndex = restartIdx
		}
	}
	fmt.Println("===========CurrentTempIndex", CurrentTempIndex)

	jb_y := 0
	jb_s := 0
	jb_syz := 0.0
	jb_count := 0
	// 遍历 table2List 计算局部数据
	for i := 0; i < len(tableYanchendao2s); i++ {
		if CurrentTempIndex == -2 /*正常打*/ || CurrentTempIndex == -1 /*重启*/ {
			///不是局部平衡的时候，要从重启的位置的下一行计算
			if tableYanchendao2s[i].ID > int(restartIndex) {
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
			///有局部平衡标志的时候（点某一行的时候）要从当前行计算
			if tableYanchendao2s[i].ID >= int(restartIndex) {
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
	statisticalAreas[14] = strconv.Itoa(jb_y - IntAbs(jb_s))
	statisticalAreas[18] = fmt.Sprintf("%.3f", jb_syz)
	if statisticalAreas[14] == "0" {
		statisticalAreas[22] = "-"
	} else {
		parseStr := RemoveChineseCharacters(statisticalAreas[14])
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
			statisticalAreas[26] = fmt.Sprintf(" 须%.1fx%d ", (jb_syz*-1+dJ)/float64(parse), parse)
		}
	} else {
		if parse == 0 {
			statisticalAreas[26] = ""
		} else {
			statisticalAreas[26] = fmt.Sprintf(" 可负%.1fx%d ", (jb_syz-dJ)/float64(parse), parse)
		}
	}
	// 填充第四列数据
	statisticalAreas[3] = fmt.Sprintf("流水%.0f", runningWater)
	if len(tableYanchendao2s) > 0 {
		statisticalAreas[7] = fmt.Sprintf("均利%.2f", zt_syz/float64(len(tableYanchendao2s)))
	}
	//连胜负
	statisticalAreas[11] = fmt.Sprintf("%d", countLianShengFu)
	xianCount := len(tableYanchendao2s) - zhuangCount
	statisticalAreas[15] = fmt.Sprintf("%d/%d/%d", zhuangCount, xianCount, zhuangCount-xianCount) //统计庄闲差
	// 保存佣金值，用于后续计算
	var yongJinValue string
	if len(tableYanchendao2s) > 0 {
		yongJinValue = tableYanchendao1.ColumnYongJin //扣水（庄扣5%）
	}
	// 计算期望值（原27的值），直接放到23位置
	if statisticalAreas[14] == "0" {
		statisticalAreas[23] = "" //期望值
	} else if statisticalAreas[21] == "-" {
		statisticalAreas[23] = "" //期望值
	} else {
		parts := strings.Split(RemoveChineseCharacters(statisticalAreas[25]), "x")
		if len(parts) > 0 {
			num25, _ := strconv.ParseFloat(parts[0], 64)
			num23, _ := strconv.ParseFloat(yongJinValue, 64)
			statisticalAreas[23] = fmt.Sprintf("%.2f", num25/num23) //期望值
		}
	}
	// 计算原31的值，直接放到27位置
	if statisticalAreas[14] == "0" {
		statisticalAreas[27] = ""
	} else if statisticalAreas[22] == "-" {
		statisticalAreas[27] = ""
	} else {
		parts := strings.Split(RemoveChineseCharacters(statisticalAreas[26]), "x")
		if len(parts) > 0 {
			num26, _ := strconv.ParseFloat(parts[0], 64)
			num23, _ := strconv.ParseFloat(yongJinValue, 64)
			statisticalAreas[27] = fmt.Sprintf("%.2f", num26/num23)
		}
	}
	// 佣金值直接放到31位置
	statisticalAreas[31] = yongJinValue //佣金

	// 预测平均值. 手机上做
	//textEditingControllerText := ""
	//if textEditingControllerText != "" {
	//	statisticalAreas[20] = pVal1(tableYanchendao1)
	//	statisticalAreas[24] = pVal2()
	//}

	if CurrentTempIndex > 2 {
		statisticalAreas[30] = fmt.Sprintf("%d", CurrentTempIndex)
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "统计数据",
		Data:   statisticalAreas,
	})
}

// 折线图数据
func LinechartData(ctx *gin.Context) {
	var arr [75]string
	uid, _ := strconv.ParseInt(ctx.GetHeader("UserId"), 10, 64) //第二个参数 10 表示字符串是十进制格式。第三个参数 64 表示转换结果的类型为 int64。
	global.Db.Model(&models.TableYanchendao2{}).Where("user_id=?", uid).Order("id DESC").Limit(75).Pluck("column_current_jin", &arr)

	//// 反转切片
	//reversedArr := make([]string, len(arr))
	//for i, j := 0, len(arr)-1; i < len(arr); i, j = i+1, j-1 {
	//	reversedArr[i] = arr[j]
	//}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "折线图数据",
		Data:   arr,
	})
}

// 清除数据（消数列数据全部清除）
func CleanDataD(ctx *gin.Context) {
	UserId := ctx.GetHeader("UserId")
	result := global.Db.Model(&models.TableYanchendao2{}).Where("user_id = ?", UserId).Update("colmun_shuyingzhi_d", "")
	if result.Error != nil {
		return
	}
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "清除数据成功",
		//Data:   [...]int64{result.RowsAffected},
		Data: [1]int64{result.RowsAffected},
		//Data: result.RowsAffected,
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

// 获取用户列表（所有用户，用于选择查询）
// @Summary      获取用户列表
// @Tags         ycd投注记录
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseJson{data=[]object}
// @Router       /api/ycd/today/users [get]
func GetTodayBettingUsers(ctx *gin.Context) {
	// 查询所有用户列表
	var users []models.User

	// 使用User模型查询所有用户（包括软删除的记录）
	if err := global.Db.Unscoped().Model(&models.User{}).
		Where("uid IS NOT NULL").
		Order("username ASC").
		Find(&users).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户列表失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 转换为响应格式，将user_id转换为字符串
	// 手动构建完整的JSON响应，确保user_id是字符串类型，避免JavaScript大整数精度丢失
	var jsonBuilder strings.Builder
	jsonBuilder.WriteString(`{"code":0,"msg":"查询成功","data":[`)
	for i, user := range users {
		if i > 0 {
			jsonBuilder.WriteString(",")
		}
		// 确保user_id是字符串（用引号包裹）
		jsonBuilder.WriteString(fmt.Sprintf(`{"user_id":"%s","username":"%s"}`,
			strconv.FormatInt(user.Uid, 10),
			strings.ReplaceAll(user.Username, `"`, `\"`)))
	}
	jsonBuilder.WriteString("]}")

	// 调试日志
	fmt.Printf("查询到 %d 个用户，JSON: %s\n", len(users), jsonBuilder.String())

	// 直接返回JSON字符串，避免gin的自动序列化
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(jsonBuilder.String()))
}

// 查询投注流水（总金额）- 支持日期范围查询
// @Summary      查询投注流水
// @Tags         ycd投注记录
// @Accept       json
// @Produce      json
// @Param        user_id query string false "用户ID，不传则查询所有用户"
// @Param        start_date query string false "开始日期，格式：YYYY-MM-DD，不传则默认为今天"
// @Param        end_date query string false "结束日期，格式：YYYY-MM-DD，不传则默认为开始日期"
// @Success      200  {object}  ResponseJson{data=object}
// @Router       /api/ycd/today/amount [get]
func GetTodayBettingAmount(ctx *gin.Context) {
	// 获取日期参数
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	// 如果没有提供日期，默认使用今天
	if startDateStr == "" {
		startDateStr = time.Now().Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = startDateStr
	}

	// 获取可选的用户ID参数
	userIDStr := ctx.Query("user_id")

	// 构建查询条件：日期范围
	query := global.Db.Model(&models.TableYanchendao2{}).
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDateStr, endDateStr)

	// 如果提供了用户ID，则按用户筛选
	if userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	var records []models.TableYanchendao2
	if err := query.Find(&records).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   0,
			Msg:    "查询流水失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 计算总金额（column_xiazhujine 是字符串，需要转换）
	var totalAmount float64
	for _, record := range records {
		if amount, err := strconv.ParseFloat(record.ColumnXiazhujine, 64); err == nil {
			totalAmount += amount
		}
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data: gin.H{
			"total_amount": totalAmount,
			"start_date":   startDateStr,
			"end_date":     endDateStr,
			"count":        len(records),
			"user_id":      userIDStr,
		},
	})
}

// 查询下注次数 - 支持日期范围查询
// @Summary      查询下注次数
// @Tags         ycd投注记录
// @Accept       json
// @Produce      json
// @Param        user_id query string false "用户ID，不传则查询所有用户"
// @Param        start_date query string false "开始日期，格式：YYYY-MM-DD，不传则默认为今天"
// @Param        end_date query string false "结束日期，格式：YYYY-MM-DD，不传则默认为开始日期"
// @Success      200  {object}  ResponseJson{data=object}
// @Router       /api/ycd/today/count [get]
func GetTodayBettingCount(ctx *gin.Context) {
	// 获取日期参数
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	// 如果没有提供日期，默认使用今天
	if startDateStr == "" {
		startDateStr = time.Now().Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = startDateStr
	}

	// 获取可选的用户ID参数
	userIDStr := ctx.Query("user_id")

	// 构建查询条件：日期范围
	query := global.Db.Model(&models.TableYanchendao2{}).
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDateStr, endDateStr)

	// 如果提供了用户ID，则按用户筛选
	if userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   0,
			Msg:    "查询下注次数失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "查询成功",
		Data: gin.H{
			"count":      count,
			"start_date": startDateStr,
			"end_date":   endDateStr,
			"user_id":    userIDStr,
		},
	})
}

// 随机庄闲接口
func GetRandomBankerPlayer(ctx *gin.Context) {
	// 获取用户ID
	UserId := ctx.GetHeader("UserId")
	if UserId == "" {
		Fail(ctx, ResponseJson{
			Status: http.StatusBadRequest,
			Code:   1,
			Msg:    "用户ID不能为空",
			Data:   gin.H{},
		})
		return
	}

	// 查询用户的 tableYanchendao1 记录，获取庄占比
	var tableYanchendao1 models.TableYanchendao1
	if err := global.Db.Where("uid=?", UserId).Last(&tableYanchendao1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(ctx, ResponseJson{
				Status: http.StatusNotFound,
				Code:   1,
				Msg:    "未找到用户数据，请先初始化",
				Data:   gin.H{},
			})
			return
		}
		Fail(ctx, ResponseJson{
			Status: http.StatusInternalServerError,
			Code:   1,
			Msg:    "查询用户数据失败: " + err.Error(),
			Data:   gin.H{},
		})
		return
	}

	// 获取庄占比，如果为0则使用默认值50
	zhuangZhanBi := tableYanchendao1.ColumnZhuangZhanBi
	if zhuangZhanBi == 0 {
		zhuangZhanBi = 50
	}

	// 使用工具类生成1-100之间的随机数 >=1. <=100
	randomValue := GenerateRandomValue(1, 100)
	// 根据庄占比判断结果
	// 如果随机数 <= 庄占比，则为庄；否则为闲
	var result string
	var resultValue int
	if randomValue <= zhuangZhanBi {
		result = "庄"
		resultValue = 0
	} else {
		result = "闲"
		resultValue = 1
	}
	fmt.Println(result, zhuangZhanBi)
	// 返回结果
	Ok(ctx, ResponseJson{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "获取随机庄闲成功",
		Data: gin.H{
			"result":      result,
			"value":       resultValue,
			"randomValue": randomValue,  // 实际随机数
			"biasValue":   zhuangZhanBi, // 庄占比（0-100）
			"timestamp":   time.Now().Unix(),
		},
	})
}
