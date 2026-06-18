package baccarat

import "math/rand"

// BulkSimulateResult 批量靴模拟结果
type BulkSimulateResult struct {
	ShoeCount               int             `json:"shoeCount"`
	TotalHands              int             `json:"totalHands"`
	EffectiveHands          int             `json:"effectiveHands"`
	ExcludeTie              bool            `json:"excludeTie"`
	Stats                   WinStats        `json:"stats"`
	BankerMinusPlayerPer10k float64         `json:"bankerMinusPlayerPer10k"`
	NetWin                  int             `json:"netWin"`
	Collision               *CollisionStats `json:"collision,omitempty"`
}

func calcBankerMinusPlayerPer10k(banker, player, base int) float64 {
	if base <= 0 {
		return 0
	}
	return float64(banker-player) * 10000 / float64(base)
}

// SimulateShoes 连续模拟多靴，每靴洗牌、随机切牌后发至需换靴
// collisionZhuangZhanBi > 0 时同步做随机庄闲碰撞统计
func SimulateShoes(shoeCount int, excludeTie bool, collisionUserID int64, collisionZhuangZhanBi int) BulkSimulateResult {
	result := BulkSimulateResult{
		ShoeCount:  shoeCount,
		ExcludeTie: excludeTie,
	}
	rnd := rand.New(rand.NewSource(rand.Int63()))

	var collision *CollisionStats
	if collisionUserID > 0 && collisionZhuangZhanBi >= 0 {
		collision = &CollisionStats{
			UserID:       collisionUserID,
			ZhuangZhanBi: collisionZhuangZhanBi,
		}
		if collision.ZhuangZhanBi == 0 {
			collision.ZhuangZhanBi = 50
		}
	}

	for i := 0; i < shoeCount; i++ {
		shoe := NewShoe(rnd)
		shoe.Shuffle()
		shoe.RandomizeCutCard()

		for !shoe.NeedsReshuffle() {
			hand, err := PrepareHand(shoe)
			if err != nil {
				break
			}
			result.TotalHands++
			switch hand.Winner {
			case WinnerPlayer:
				result.Stats.Player++
			case WinnerBanker:
				result.Stats.Banker++
			case WinnerTie:
				result.Stats.Tie++
			}
			if collision != nil {
				pick := PickRandomBankerPlayer(collision.ZhuangZhanBi, rnd)
				recordCollision(collision, pick, winnerToSide(hand.Winner))
			}
		}
	}

	result.EffectiveHands = result.Stats.Player + result.Stats.Banker
	base := result.TotalHands
	if excludeTie {
		base = result.EffectiveHands
	}
	result.BankerMinusPlayerPer10k = calcBankerMinusPlayerPer10k(
		result.Stats.Banker,
		result.Stats.Player,
		base,
	)
	result.NetWin = result.Stats.Banker - result.Stats.Player
	if collision != nil {
		finalizeCollision(collision)
		result.Collision = collision
	}
	return result
}
