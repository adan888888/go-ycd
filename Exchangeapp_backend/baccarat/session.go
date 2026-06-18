package baccarat

import (
	"fmt"
	"sync"
)

// Session 单用户百家乐模拟会话
type Session struct {
	Shoe              *Shoe
	BigRoad           *BigRoadState
	AwaitingCutCard   bool
	ShoeCutCardChosen bool
	GameHistory       []GameRecord
	LastHand          *HandResult
}

// SessionManager 按用户隔离的会话管理
type SessionManager struct {
	mu       sync.Mutex
	sessions map[int64]*Session
}

var DefaultSessions = &SessionManager{
	sessions: make(map[int64]*Session),
}

func (m *SessionManager) Get(uid int64) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[uid]
	if !ok {
		s = newSession()
		m.sessions[uid] = s
	}
	return s
}

func (m *SessionManager) Reset(uid int64) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := newSession()
	m.sessions[uid] = s
	return s
}

func newSession() *Session {
	return &Session{
		Shoe:    NewShoe(nil),
		BigRoad: NewBigRoadState(),
	}
}

// StateDTO 对外返回的状态
type StateDTO struct {
	ShoeRemaining      int          `json:"shoeRemaining"`
	ShoeTotalCards     int          `json:"shoeTotalCards"`
	ShoeCutCardRemaining int        `json:"shoeCutCardRemaining"`
	AwaitingCutCard    bool         `json:"awaitingCutCard"`
	ShoeCutCardChosen  bool         `json:"shoeCutCardChosen"`
	NeedsReshuffle     bool         `json:"needsReshuffle"`
	HasActiveShoe      bool         `json:"hasActiveShoe"`
	CurrentResult      string       `json:"currentResult"`
	PlayerCards        []Card       `json:"playerCards"`
	BankerCards        []Card       `json:"bankerCards"`
	PlayerTotal        int          `json:"playerTotal"`
	BankerTotal        int          `json:"bankerTotal"`
	Winner             string       `json:"winner"`
	BigRoad            [][]string   `json:"bigRoad"`
	GameHistory        []GameRecord `json:"gameHistory"`
	Stats              WinStats     `json:"stats"`
	LastSteps          []DealStep   `json:"lastSteps,omitempty"`
}

func (s *Session) ToDTO() StateDTO {
	dto := StateDTO{
		ShoeRemaining:        s.Shoe.Remaining(),
		ShoeTotalCards:       TotalCards,
		ShoeCutCardRemaining: s.Shoe.CutCardRemaining(),
		AwaitingCutCard:      s.AwaitingCutCard,
		ShoeCutCardChosen:    s.ShoeCutCardChosen,
		NeedsReshuffle:       s.Shoe.NeedsReshuffle(),
		HasActiveShoe:        s.Shoe.HasActiveShoe(),
		BigRoad:              s.BigRoad.BigRoad,
		GameHistory:          s.GameHistory,
		Stats:                CalcWinStats(s.GameHistory),
	}
	if s.LastHand != nil {
		dto.CurrentResult = s.LastHand.ResultText
		dto.PlayerCards = s.LastHand.PlayerCards
		dto.BankerCards = s.LastHand.BankerCards
		dto.PlayerTotal = s.LastHand.PlayerTotal
		dto.BankerTotal = s.LastHand.BankerTotal
		dto.Winner = s.LastHand.Winner
		dto.LastSteps = s.LastHand.Steps
	}
	if dto.GameHistory == nil {
		dto.GameHistory = []GameRecord{}
	}
	return dto
}

func (s *Session) Shuffle() {
	s.Shoe.Shuffle()
	s.AwaitingCutCard = true
	s.ShoeCutCardChosen = false
	s.LastHand = nil
}

func (s *Session) CutCard() error {
	if !s.AwaitingCutCard {
		return fmt.Errorf("当前不需要切牌")
	}
	s.Shoe.RandomizeCutCard()
	s.AwaitingCutCard = false
	s.ShoeCutCardChosen = true
	return nil
}

func (s *Session) Deal() (*HandResult, error) {
	if s.AwaitingCutCard {
		return nil, fmt.Errorf("请先随机切牌位")
	}
	if s.Shoe.NeedsReshuffle() {
		return nil, fmt.Errorf("牌靴需换靴，请先洗牌")
	}

	hand, err := PrepareHand(s.Shoe)
	if err != nil {
		return nil, err
	}

	s.LastHand = hand
	UpdateBigRoad(s.BigRoad, hand.Winner)

	record := GameRecord{
		PlayerCards: cardsToDisplay(hand.PlayerCards),
		BankerCards: cardsToDisplay(hand.BankerCards),
		PlayerTotal: hand.PlayerTotal,
		BankerTotal: hand.BankerTotal,
		Winner:      hand.Winner,
	}
	s.GameHistory = append([]GameRecord{record}, s.GameHistory...)
	if len(s.GameHistory) > 20 {
		s.GameHistory = s.GameHistory[:20]
	}

	return hand, nil
}

func (s *Session) ClearHistory() {
	s.BigRoad = NewBigRoadState()
	s.GameHistory = nil
	s.LastHand = nil
}
