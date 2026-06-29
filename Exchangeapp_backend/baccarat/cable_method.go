package baccarat

const MaxCableLayer = 13

// 十三太保缆法下注表（层 1～13）
var (
	cableCol1 = []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377}
	cableCol2 = []int{0, 0, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
	cableCol3 = []int{0, 0, 4, 6, 10, 16, 26, 42, 68, 110, 178, 288, 466}
)

// CableState 缆法当前位置：层 + 列（1=第1列，2=必打，3=选打）
type CableState struct {
	Layer int `json:"layer"`
	Col   int `json:"col"`
}

func NewCableState() CableState {
	return CableState{Layer: 1, Col: 1}
}

// BetAmount 当前应下注单位数
func (s CableState) BetAmount() int {
	if s.Layer < 1 {
		s.Layer = 1
	}
	if s.Layer > MaxCableLayer {
		s.Layer = MaxCableLayer
	}
	idx := s.Layer - 1
	switch s.Col {
	case 2:
		if s.Layer <= 2 {
			return cableCol1[idx]
		}
		return cableCol2[idx]
	case 3:
		return cableCol3[idx]
	default:
		return cableCol1[idx]
	}
}

// cableOutcome 碰撞结果：win / loss / tie
func cableOutcome(pick, actual string) string {
	if actual == "和" {
		return "tie"
	}
	if pick == actual {
		return "win"
	}
	return "loss"
}

// cableProfit 缆法单局盈亏：庄赢扣 5% 手续费（净赢 0.95），闲赢 1:1，输全损，和局 0
func cableProfit(bet int, pick, outcome string) float64 {
	switch outcome {
	case "tie":
		return 0
	case "loss":
		return float64(-bet)
	default:
		if pick == "庄" {
			return float64(bet) * 0.95
		}
		return float64(bet)
	}
}

// advanceCable 根据中/不中推进缆法；和局不变层
func advanceCable(state CableState, outcome, pick string) (CableState, float64, bool) {
	if outcome == "tie" {
		return state, 0, false
	}

	bet := state.BetAmount()
	profit := cableProfit(bet, pick, outcome)
	win := outcome == "win"

	next := state
	bursted := false

	if state.Layer <= 2 {
		if win {
			next = CableState{Layer: 1, Col: 1}
		} else if state.Layer >= MaxCableLayer {
			bursted = true
			next = CableState{Layer: 1, Col: 1}
		} else {
			next = CableState{Layer: state.Layer + 1, Col: 1}
		}
		return next, profit, bursted
	}

	switch state.Col {
	case 1:
		if win {
			next = CableState{Layer: state.Layer, Col: 2}
		} else if state.Layer >= MaxCableLayer {
			bursted = true
			next = CableState{Layer: 1, Col: 1}
		} else {
			next = CableState{Layer: state.Layer + 1, Col: 1}
		}
	case 2:
		if win {
			next = CableState{Layer: 1, Col: 1}
		} else {
			next = CableState{Layer: state.Layer, Col: 3}
		}
	case 3:
		if win {
			next = CableState{Layer: 1, Col: 1}
		} else if state.Layer >= MaxCableLayer {
			bursted = true
			next = CableState{Layer: 1, Col: 1}
		} else {
			next = CableState{Layer: state.Layer + 1, Col: 1}
		}
	}

	return next, profit, bursted
}

// CableTableRow 缆法表一行（供前端展示）
type CableTableRow struct {
	Layer int `json:"layer"`
	Col1  int `json:"col1"`
	Col2  int `json:"col2"`
	Col3  int `json:"col3"`
}

// CableTable 返回完整缆法表
func CableTable() []CableTableRow {
	rows := make([]CableTableRow, MaxCableLayer)
	for i := 0; i < MaxCableLayer; i++ {
		rows[i] = CableTableRow{
			Layer: i + 1,
			Col1:  cableCol1[i],
			Col2:  cableCol2[i],
			Col3:  cableCol3[i],
		}
	}
	return rows
}
