package baccarat

import "math/rand"

// CableHandDetail 单局缆法明细
type CableHandDetail struct {
	HandIndex        int    `json:"handIndex"`
	ShoeIndex        int    `json:"shoeIndex"`
	Pick             string `json:"pick"`
	Actual           string `json:"actual"`
	Outcome          string `json:"outcome"`
	Bet              int    `json:"bet"`
	Profit           float64 `json:"profit"`
	Layer            int     `json:"layer"`
	Col              int     `json:"col"`
	NextLayer        int     `json:"nextLayer"`
	NextCol          int     `json:"nextCol"`
	CumulativeProfit float64 `json:"cumulativeProfit"`
	Bursted          bool   `json:"bursted"`
}

// CableShoeSummary 单靴缆法汇总
type CableShoeSummary struct {
	ShoeIndex   int `json:"shoeIndex"`
	Hands       int `json:"hands"`
	Profit      float64 `json:"profit"`
	MaxLayer    int `json:"maxLayer"`
	BurstCount  int `json:"burstCount"`
}

// CableSummary 缆法模拟汇总
type CableSummary struct {
	TotalProfit   float64 `json:"totalProfit"`
	MaxLayer       int `json:"maxLayer"`
	BurstCount     int `json:"burstCount"`
	SettledHands   int `json:"settledHands"`
	TieHands       int `json:"tieHands"`
	WinHands       int `json:"winHands"`
	LossHands      int `json:"lossHands"`
	TotalBetUnits  int `json:"totalBetUnits"`
}

// CableSimulateResult 缆法 + 八副牌 + 庄闲碰撞模拟结果
type CableSimulateResult struct {
	ShoeCount               int                `json:"shoeCount"`
	TotalHands              int                `json:"totalHands"`
	Stats                   WinStats           `json:"stats"`
	BankerMinusPlayerPer10k float64            `json:"bankerMinusPlayerPer10k"`
	Collision               *CollisionStats    `json:"collision"`
	Cable                   CableSummary       `json:"cable"`
	CableTable              []CableTableRow    `json:"cableTable"`
	ShoeSummaries           []CableShoeSummary `json:"shoeSummaries"`
	Details                 []CableHandDetail  `json:"details"`
	DetailTruncated         bool               `json:"detailTruncated"`
}

const defaultCableMaxDetails = 500

// SimulateCableMethod 八副牌发牌 + 庄闲碰撞 + 十三太保缆法下注
func SimulateCableMethod(shoeCount int, collisionUserID int64, collisionZhuangZhanBi int, maxDetails int) CableSimulateResult {
	if maxDetails <= 0 {
		maxDetails = defaultCableMaxDetails
	}
	if maxDetails > 5000 {
		maxDetails = 5000
	}

	result := CableSimulateResult{
		ShoeCount:  shoeCount,
		CableTable: CableTable(),
	}
	rnd := rand.New(rand.NewSource(rand.Int63()))

	collision := &CollisionStats{
		UserID:       collisionUserID,
		ZhuangZhanBi: collisionZhuangZhanBi,
	}
	if collision.ZhuangZhanBi <= 0 {
		collision.ZhuangZhanBi = 50
	}

	cableState := NewCableState()
	cumulative := 0.0
	handIndex := 0
	detailTruncated := false

	for shoeIdx := 0; shoeIdx < shoeCount; shoeIdx++ {
		shoeSummary := CableShoeSummary{ShoeIndex: shoeIdx + 1}
		shoeStartProfit := cumulative

		shoe := NewShoe(rnd)
		shoe.Shuffle()
		shoe.RandomizeCutCard()

		for !shoe.NeedsReshuffle() {
			hand, err := PrepareHand(shoe)
			if err != nil {
				break
			}
			result.TotalHands++
			handIndex++
			shoeSummary.Hands++

			switch hand.Winner {
			case WinnerPlayer:
				result.Stats.Player++
			case WinnerBanker:
				result.Stats.Banker++
			case WinnerTie:
				result.Stats.Tie++
			}

			pick := PickRandomBankerPlayer(collision.ZhuangZhanBi, rnd)
			actual := winnerToSide(hand.Winner)
			recordCollision(collision, pick, actual)

			outcome := cableOutcome(pick, actual)
			bet := cableState.BetAmount()
			layerBefore := cableState.Layer
			colBefore := cableState.Col

			if cableState.Layer > result.Cable.MaxLayer {
				result.Cable.MaxLayer = cableState.Layer
			}
			if cableState.Layer > shoeSummary.MaxLayer {
				shoeSummary.MaxLayer = cableState.Layer
			}

			nextState, profit, bursted := advanceCable(cableState, outcome, pick)
			cumulative += profit

			switch outcome {
			case "tie":
				result.Cable.TieHands++
			case "win":
				result.Cable.WinHands++
				result.Cable.SettledHands++
				result.Cable.TotalBetUnits += bet
			case "loss":
				result.Cable.LossHands++
				result.Cable.SettledHands++
				result.Cable.TotalBetUnits += bet
			}
			if bursted {
				result.Cable.BurstCount++
				shoeSummary.BurstCount++
			}

			if len(result.Details) < maxDetails {
				result.Details = append(result.Details, CableHandDetail{
					HandIndex:        handIndex,
					ShoeIndex:        shoeIdx + 1,
					Pick:             pick,
					Actual:           actual,
					Outcome:          outcome,
					Bet:              bet,
					Profit:           profit,
					Layer:            layerBefore,
					Col:              colBefore,
					NextLayer:        nextState.Layer,
					NextCol:          nextState.Col,
					CumulativeProfit: cumulative,
					Bursted:          bursted,
				})
			} else {
				detailTruncated = true
			}

			cableState = nextState
		}

		shoeSummary.Profit = cumulative - shoeStartProfit
		result.ShoeSummaries = append(result.ShoeSummaries, shoeSummary)
	}

	result.Cable.TotalProfit = cumulative
	result.BankerMinusPlayerPer10k = calcBankerMinusPlayerPer10k(
		result.Stats.Banker,
		result.Stats.Player,
		result.TotalHands,
	)
	finalizeCollision(collision)
	result.Collision = collision
	result.DetailTruncated = detailTruncated
	return result
}
