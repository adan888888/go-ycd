package baccarat

import "fmt"

const (
	WinnerPlayer = "闲家"
	WinnerBanker = "庄家"
	WinnerTie    = "和局"
)

// DealStep 发牌步骤
type DealStep struct {
	Side         string `json:"side"`
	Card         Card   `json:"card"`
	IsThirdCard  bool   `json:"isThirdCard"`
}

// HandResult 一局完整结果
type HandResult struct {
	PlayerCards []Card     `json:"playerCards"`
	BankerCards []Card     `json:"bankerCards"`
	PlayerTotal int        `json:"playerTotal"`
	BankerTotal int        `json:"bankerTotal"`
	Winner      string     `json:"winner"`
	ResultText  string     `json:"resultText"`
	Steps       []DealStep `json:"steps"`
}

// GameRecord 历史记录
type GameRecord struct {
	PlayerCards string `json:"playerCards"`
	BankerCards string `json:"bankerCards"`
	PlayerTotal int    `json:"playerTotal"`
	BankerTotal int    `json:"bankerTotal"`
	Winner      string `json:"winner"`
}

// WinStats 胜负统计
type WinStats struct {
	Player int `json:"player"`
	Banker int `json:"banker"`
	Tie    int `json:"tie"`
}

// BigRoadState 大路图状态
type BigRoadState struct {
	BigRoad           [][]string `json:"bigRoad"`
	CurrentRow        int        `json:"currentRow"`
	CurrentCol        int        `json:"currentCol"`
	LastWinner        string     `json:"lastWinner"`
	DragonStartCol    int        `json:"dragonStartCol"`
	DragonParallelRow int        `json:"dragonParallelRow"`
}

// NewBigRoadState 初始化大路图
func NewBigRoadState() *BigRoadState {
	bigRoad := make([][]string, BigRoadRows)
	for i := range bigRoad {
		bigRoad[i] = make([]string, BigRoadCols)
	}
	return &BigRoadState{
		BigRoad:           bigRoad,
		DragonStartCol:    -1,
		DragonParallelRow: -1,
	}
}

func CalculateTotal(cards []Card) int {
	total := 0
	for _, c := range cards {
		total += c.Value
	}
	return total % 10
}

func cardsToDisplay(cards []Card) string {
	if len(cards) == 0 {
		return ""
	}
	result := cards[0].Display
	for i := 1; i < len(cards); i++ {
		result += " " + cards[i].Display
	}
	return result
}

func containsInt(values []int, target int) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

// PrepareHand 按百家乐规则发一局牌
func PrepareHand(shoe *Shoe) (*HandResult, error) {
	playerCards := make([]Card, 0, 3)
	bankerCards := make([]Card, 0, 3)
	steps := make([]DealStep, 0, 6)

	dealToPlayer := func(isThird bool) error {
		card, err := shoe.Draw()
		if err != nil {
			return err
		}
		playerCards = append(playerCards, card)
		steps = append(steps, DealStep{Side: "player", Card: card, IsThirdCard: isThird})
		return nil
	}
	dealToBanker := func(isThird bool) error {
		card, err := shoe.Draw()
		if err != nil {
			return err
		}
		bankerCards = append(bankerCards, card)
		steps = append(steps, DealStep{Side: "banker", Card: card, IsThirdCard: isThird})
		return nil
	}

	if err := dealToPlayer(false); err != nil {
		return nil, err
	}
	if err := dealToBanker(false); err != nil {
		return nil, err
	}
	if err := dealToPlayer(false); err != nil {
		return nil, err
	}
	if err := dealToBanker(false); err != nil {
		return nil, err
	}

	playerTotal := CalculateTotal(playerCards)
	bankerTotal := CalculateTotal(bankerCards)

	if playerTotal >= 8 || bankerTotal >= 8 {
		return finalizeHand(playerCards, bankerCards, playerTotal, bankerTotal, steps), nil
	}

	playerGetsThird := playerTotal <= 5
	if playerGetsThird {
		if err := dealToPlayer(true); err != nil {
			return nil, err
		}
		playerTotal = CalculateTotal(playerCards)

		bankerGetsThird := false
		switch bankerTotal {
		case 0, 1, 2:
			bankerGetsThird = true
		case 3:
			bankerGetsThird = playerCards[2].Value != 8
		case 4:
			bankerGetsThird = containsInt([]int{2, 3, 4, 5, 6, 7}, playerCards[2].Value)
		case 5:
			bankerGetsThird = containsInt([]int{4, 5, 6, 7}, playerCards[2].Value)
		case 6:
			bankerGetsThird = containsInt([]int{6, 7}, playerCards[2].Value)
		}

		if bankerGetsThird {
			if err := dealToBanker(true); err != nil {
				return nil, err
			}
			bankerTotal = CalculateTotal(bankerCards)
		}
	} else if bankerTotal <= 5 {
		if err := dealToBanker(true); err != nil {
			return nil, err
		}
		bankerTotal = CalculateTotal(bankerCards)
	}

	return finalizeHand(playerCards, bankerCards, playerTotal, bankerTotal, steps), nil
}

func finalizeHand(playerCards, bankerCards []Card, playerTotal, bankerTotal int, steps []DealStep) *HandResult {
	var winner, resultText string
	if playerTotal > bankerTotal {
		winner = WinnerPlayer
		resultText = fmt.Sprintf("闲家胜 (%d vs %d)", playerTotal, bankerTotal)
	} else if bankerTotal > playerTotal {
		winner = WinnerBanker
		resultText = fmt.Sprintf("庄家胜 (%d vs %d)", bankerTotal, playerTotal)
	} else {
		winner = WinnerTie
		resultText = fmt.Sprintf("和局 (%d vs %d)", playerTotal, bankerTotal)
	}
	return &HandResult{
		PlayerCards: playerCards,
		BankerCards: bankerCards,
		PlayerTotal: playerTotal,
		BankerTotal: bankerTotal,
		Winner:      winner,
		ResultText:  resultText,
		Steps:       steps,
	}
}

// UpdateBigRoad 根据百家乐大路规则更新大路图
func UpdateBigRoad(state *BigRoadState, winner string) {
	if winner == WinnerTie {
		return
	}

	if state.LastWinner == "" {
		state.BigRoad[state.CurrentRow][state.CurrentCol] = winner
		state.CurrentCol++
	} else if state.LastWinner != winner {
		state.DragonStartCol = -1
		state.DragonParallelRow = -1
		state.CurrentRow = 0
		state.BigRoad[state.CurrentRow][state.CurrentCol] = winner
		state.CurrentCol++
	} else {
		state.CurrentRow++
		col := state.CurrentCol - 1

		occupiedBelow := state.CurrentRow < BigRoadRows && state.BigRoad[state.CurrentRow][col] != ""
		exceedRows := state.CurrentRow > BigRoadRows-1
		if occupiedBelow || exceedRows {
			state.DragonStartCol++
			state.BigRoad[state.DragonParallelRow][state.DragonStartCol] = winner
		} else {
			state.BigRoad[state.CurrentRow][state.CurrentCol-1] = winner
			state.DragonParallelRow = state.CurrentRow
			state.DragonStartCol = state.CurrentCol - 1
		}
	}

	state.LastWinner = winner
}

func CalcWinStats(history []GameRecord) WinStats {
	var stats WinStats
	for _, r := range history {
		switch r.Winner {
		case WinnerPlayer:
			stats.Player++
		case WinnerBanker:
			stats.Banker++
		case WinnerTie:
			stats.Tie++
		}
	}
	return stats
}
