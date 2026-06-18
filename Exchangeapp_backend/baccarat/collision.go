package baccarat

import "math/rand"

// CollisionStats 随机庄闲与实际开奖碰撞统计
type CollisionStats struct {
	UserID       int64   `json:"userId"`
	ZhuangZhanBi int     `json:"zhuangZhanBi"`
	WinCount     int     `json:"winCount"`
	LossCount    int     `json:"lossCount"`
	TieCount     int     `json:"tieCount"`
	SettledCount int     `json:"settledCount"`
	WinRate      float64 `json:"winRate"`
	NetWin                         int     `json:"netWin"`
	RandomBanker                   int     `json:"randomBanker"`
	RandomPlayer                   int     `json:"randomPlayer"`
	RandomBankerPlayerDiff         int     `json:"randomBankerPlayerDiff"`
	RandomBankerMinusPlayerPer10k  float64 `json:"randomBankerMinusPlayerPer10k"`
}

// PickRandomBankerPlayer 与 GetRandomBankerPlayer 相同规则：1~100 随机数 <= 庄占比则为庄
func PickRandomBankerPlayer(zhuangZhanBi int, rnd *rand.Rand) string {
	if zhuangZhanBi <= 0 {
		zhuangZhanBi = 50
	}
	if rnd.Intn(100)+1 <= zhuangZhanBi {
		return "庄"
	}
	return "闲"
}

func winnerToSide(winner string) string {
	switch winner {
	case WinnerBanker:
		return "庄"
	case WinnerPlayer:
		return "闲"
	default:
		return "和"
	}
}

func finalizeCollision(stats *CollisionStats) {
	stats.SettledCount = stats.WinCount + stats.LossCount
	stats.NetWin = stats.WinCount - stats.LossCount
	stats.RandomBankerPlayerDiff = stats.RandomBanker - stats.RandomPlayer
	totalRandom := stats.RandomBanker + stats.RandomPlayer
	if totalRandom > 0 {
		stats.RandomBankerMinusPlayerPer10k = float64(stats.RandomBankerPlayerDiff) * 10000 / float64(totalRandom)
	}
	if stats.SettledCount > 0 {
		stats.WinRate = float64(stats.WinCount) * 100 / float64(stats.SettledCount)
	}
}

func recordCollision(stats *CollisionStats, pick, actual string) {
	if pick == "庄" {
		stats.RandomBanker++
	} else {
		stats.RandomPlayer++
	}
	switch actual {
	case "和":
		stats.TieCount++
	case pick:
		stats.WinCount++
	default:
		stats.LossCount++
	}
}
