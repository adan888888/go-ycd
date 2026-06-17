package controllers

import (
	"exchangeapp/apicode"
	"exchangeapp/global"
	"exchangeapp/models"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	dataStatsMaxRecords  = 50000
	dataStatsTrendSample = 360
)

type dataStatsRecord struct {
	UserID           int64     `gorm:"column:user_id"`
	ZX         string    `gorm:"column:zx"`
	Shengfulu  string    `gorm:"column:shengfulu"`
	Remark     string    `gorm:"column:remark"`
	Shuyingzhi float64   `gorm:"column:shuyingzhi"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

type dataStatsZx struct {
	ZhuangCount int     `json:"zhuang_count"`
	XianCount   int     `json:"xian_count"`
	Diff        int     `json:"diff"`
	ZhuangRate  float64 `json:"zhuang_rate"`
	XianRate    float64 `json:"xian_rate"`
}

type dataStatsUser struct {
	UserID       string      `json:"user_id"`
	Username     string      `json:"username"`
	Total        int         `json:"total"`
	Random       dataStatsZx `json:"random"`
	Draw         dataStatsZx `json:"draw"`
	WinCount     int         `json:"win_count"`
	LossCount    int         `json:"loss_count"`
	SettledCount int         `json:"settled_count"`
	WinRate      float64     `json:"win_rate"`
	LatestAt     string      `json:"latest_at"`
}

type dataStatsTrendPoint struct {
	Index       int     `json:"index"`
	WinRate     float64 `json:"win_rate"`
	WinLossDiff int     `json:"win_loss_diff"`
}

type dataStatsUserAcc struct {
	userID       int64
	username     string
	total        int
	randomZhuang int
	randomXian   int
	drawZhuang   int
	drawXian     int
	winCount     int
	lossCount    int
	latestAt     time.Time
	hasLatest    bool
}

// GetDataStats 数据统计页聚合接口：后端计算汇总、分用户统计与趋势点，避免前端分页拉全量明细。
func GetDataStats(ctx *gin.Context) {
	userScope, status, msg := resolveListUserIDsScope(ctx)
	if status != 0 {
		forbidScope(ctx, status, msg)
		return
	}

	countQuery := global.Db.Model(&models.TableYanchendao2{}).Where("deleted_at IS NULL")
	countQuery = applyUserScope(countQuery, userScope, "user_id")
	var totalRecords int64
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, "查询总数失败: "+err.Error())
		return
	}

	query := global.Db.Model(&models.TableYanchendao2{}).
		Select("user_id, zx, shengfulu, remark, shuyingzhi, created_at").
		Where("deleted_at IS NULL").
		Order("created_at ASC")
	query = applyUserScope(query, userScope, "user_id")
	query = query.Limit(dataStatsMaxRecords)

	var records []dataStatsRecord
	if err := query.Find(&records).Error; err != nil {
		Fail(ctx, apicode.CodeServerError, "查询失败: "+err.Error())
		return
	}

	userMap := make(map[int64]*dataStatsUserAcc)
	ensureUser := func(userID int64) *dataStatsUserAcc {
		if acc, ok := userMap[userID]; ok {
			return acc
		}
		acc := &dataStatsUserAcc{userID: userID}
		userMap[userID] = acc
		return acc
	}

	win := 0
	loss := 0
	trendAll := make([]dataStatsTrendPoint, 0, len(records))
	for i, record := range records {
		acc := ensureUser(record.UserID)
		acc.total++

		zx := strings.TrimSpace(record.ZX)
		switch zx {
		case "庄":
			acc.drawZhuang++
		case "闲":
			acc.drawXian++
		}

		switch deriveRandomZx(zx, record.Shengfulu, record.Remark) {
		case "庄":
			acc.randomZhuang++
		case "闲":
			acc.randomXian++
		}

		switch parseShuyingzhiOutcome(record.Shuyingzhi) {
		case 1:
			acc.winCount++
			win++
		case -1:
			acc.lossCount++
			loss++
		}

		if !acc.hasLatest || record.CreatedAt.After(acc.latestAt) {
			acc.latestAt = record.CreatedAt
			acc.hasLatest = true
		}

		settled := win + loss
		winRate := 0.0
		if settled > 0 {
			winRate = float64(win) / float64(settled) * 100
		}
		trendAll = append(trendAll, dataStatsTrendPoint{
			Index:       i + 1,
			WinRate:     winRate,
			WinLossDiff: win - loss,
		})
	}

	userIDs := make([]int64, 0, len(userMap))
	for uid := range userMap {
		userIDs = append(userIDs, uid)
	}
	uidName := loadUsernames(userIDs)

	users := make([]dataStatsUser, 0, len(userMap))
	for _, acc := range userMap {
		username := uidName[acc.userID]
		users = append(users, buildDataStatsUser(acc, username))
	}
	sort.Slice(users, func(i, j int) bool {
		if users[i].Total == users[j].Total {
			return users[i].UserID < users[j].UserID
		}
		return users[i].Total > users[j].Total
	})

	aggregateAcc := &dataStatsUserAcc{}
	for _, acc := range userMap {
		mergeUserAcc(aggregateAcc, acc)
	}
	aggregate := buildDataStatsUser(aggregateAcc, "全部用户")
	aggregate.UserID = "all"

	Success(ctx, "查询成功", gin.H{
		"total_records": totalRecords,
		"used_records":  len(records),
		"limit_reached": totalRecords > int64(dataStatsMaxRecords),
		"max_records":   dataStatsMaxRecords,
		"aggregate":     aggregate,
		"users":         users,
		"trend":         sampleDataStatsTrend(trendAll, dataStatsTrendSample),
	})
}

func deriveRandomZx(zx, sfl, remark string) string {
	zx = strings.TrimSpace(zx)
	sfl = strings.TrimSpace(sfl)
	remark = strings.TrimSpace(remark)
	if zx != "庄" && zx != "闲" {
		return ""
	}
	if sfl != "正打" && sfl != "反打" {
		return ""
	}
	if remark == "1" || remark == "-1" {
		sameAsDraw := (sfl == "正打" && remark == "1") || (sfl == "反打" && remark == "-1")
		if sameAsDraw {
			return zx
		}
		if zx == "庄" {
			return "闲"
		}
		return "庄"
	}
	if sfl == "正打" {
		return zx
	}
	if zx == "庄" {
		return "闲"
	}
	return "庄"
}

func parseShuyingzhiOutcome(value float64) int {
	if value > 0 {
		return 1
	}
	if value < 0 {
		return -1
	}
	return 0
}

func calcDataStatsZx(zhuangCount, xianCount, total int) dataStatsZx {
	diff := zhuangCount - xianCount
	zhuangRate := 0.0
	xianRate := 0.0
	if total > 0 {
		zhuangRate = float64(zhuangCount) / float64(total) * 100
		xianRate = float64(xianCount) / float64(total) * 100
	}
	return dataStatsZx{
		ZhuangCount: zhuangCount,
		XianCount:   xianCount,
		Diff:        diff,
		ZhuangRate:  zhuangRate,
		XianRate:    xianRate,
	}
}

func buildDataStatsUser(acc *dataStatsUserAcc, username string) dataStatsUser {
	settled := acc.winCount + acc.lossCount
	winRate := 0.0
	if settled > 0 {
		winRate = float64(acc.winCount) / float64(settled) * 100
	}
	latestAt := ""
	if acc.hasLatest {
		latestAt = acc.latestAt.Format("2006-01-02 15:04:05")
	}
	return dataStatsUser{
		UserID:       strconv.FormatInt(acc.userID, 10),
		Username:     username,
		Total:        acc.total,
		Random:       calcDataStatsZx(acc.randomZhuang, acc.randomXian, acc.total),
		Draw:         calcDataStatsZx(acc.drawZhuang, acc.drawXian, acc.total),
		WinCount:     acc.winCount,
		LossCount:    acc.lossCount,
		SettledCount: settled,
		WinRate:      winRate,
		LatestAt:     latestAt,
	}
}

func mergeUserAcc(dst, src *dataStatsUserAcc) {
	dst.total += src.total
	dst.randomZhuang += src.randomZhuang
	dst.randomXian += src.randomXian
	dst.drawZhuang += src.drawZhuang
	dst.drawXian += src.drawXian
	dst.winCount += src.winCount
	dst.lossCount += src.lossCount
	if src.hasLatest && (!dst.hasLatest || src.latestAt.After(dst.latestAt)) {
		dst.latestAt = src.latestAt
		dst.hasLatest = true
	}
}

func loadUsernames(userIDs []int64) map[int64]string {
	result := make(map[int64]string, len(userIDs))
	if len(userIDs) == 0 {
		return result
	}
	var users []models.User
	if err := global.Db.Select("uid, username").Where("uid IN ?", userIDs).Find(&users).Error; err != nil {
		return result
	}
	for _, user := range users {
		result[user.Uid] = user.Username
	}
	return result
}

func sampleDataStatsTrend(points []dataStatsTrendPoint, maxCount int) []dataStatsTrendPoint {
	if len(points) <= maxCount {
		return points
	}
	if maxCount <= 1 {
		if len(points) == 0 {
			return points
		}
		return []dataStatsTrendPoint{points[len(points)-1]}
	}
	step := float64(len(points)-1) / float64(maxCount-1)
	sampled := make([]dataStatsTrendPoint, maxCount)
	for i := 0; i < maxCount; i++ {
		idx := int(math.Round(float64(i) * step))
		if idx >= len(points) {
			idx = len(points) - 1
		}
		sampled[i] = points[idx]
	}
	return sampled
}
