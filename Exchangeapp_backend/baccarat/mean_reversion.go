package baccarat

import "math/rand"

// MeanReversionZhuangZhanBi 根据净胜 S 返回庄占比（均值回归策略）
// S < -50 → 80%；-50~-20 → 65%；-20~20 → 55%；20~50 → 45%；≥50 → 35%
func MeanReversionZhuangZhanBi(netWin int) int {
	switch {
	case netWin < -50:
		return 80
	case netWin < -20:
		return 65
	case netWin < 20:
		return 55
	case netWin < 50:
		return 45
	default:
		return 35
	}
}

// SimulateShoesMeanReversion 均值回归策略碰撞模拟：按累计净胜动态调整庄占比
func SimulateShoesMeanReversion(shoeCount int, collisionUserID int64, initialNetWin int) BulkSimulateResult {
	result := BulkSimulateResult{
		ShoeCount:  shoeCount,
		ExcludeTie: false,
	}
	rnd := rand.New(rand.NewSource(rand.Int63()))

	collision := &CollisionStats{
		UserID:       collisionUserID,
		ZhuangZhanBi: MeanReversionZhuangZhanBi(initialNetWin),
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

			runningS := initialNetWin + collision.WinCount - collision.LossCount
			collision.ZhuangZhanBi = MeanReversionZhuangZhanBi(runningS)
			pick := PickRandomBankerPlayer(collision.ZhuangZhanBi, rnd)
			recordCollision(collision, pick, winnerToSide(hand.Winner))
		}
	}

	result.EffectiveHands = result.Stats.Player + result.Stats.Banker
	result.BankerMinusPlayerPer10k = calcBankerMinusPlayerPer10k(
		result.Stats.Banker,
		result.Stats.Player,
		result.TotalHands,
	)
	result.NetWin = result.Stats.Banker - result.Stats.Player
	finalizeCollision(collision)
	result.Collision = collision
	return result
}
