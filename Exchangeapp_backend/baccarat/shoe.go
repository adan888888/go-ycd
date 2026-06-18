package baccarat

import (
	"fmt"
	"math/rand"
)

const (
	DeckCount            = 8
	CardsPerDeck         = 52
	TotalCards           = DeckCount * CardsPerDeck
	CutCardMinRemaining  = 12
	CutCardMaxRemaining  = 20
	MinCardsPerHand      = 6
	BigRoadRows          = 6
	BigRoadCols          = 120
)

var (
	suits = []string{"♠", "♥", "♦", "♣"}
	ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
)

// Card 单张扑克牌
type Card struct {
	Suit    string `json:"suit"`
	Rank    string `json:"rank"`
	Value   int    `json:"value"`
	Display string `json:"display"`
}

// Shoe 8 副牌百家乐牌靴（416 张，无放回发牌，切牌后换靴）
type Shoe struct {
	rand              *rand.Rand
	cards             []Card
	drawIndex         int
	cutCardRemaining  int
}

// NewShoe 创建牌靴
func NewShoe(r *rand.Rand) *Shoe {
	if r == nil {
		r = rand.New(rand.NewSource(rand.Int63()))
	}
	return &Shoe{
		rand:             r,
		cutCardRemaining: CutCardMaxRemaining,
	}
}

func (s *Shoe) Remaining() int {
	return len(s.cards) - s.drawIndex
}

func (s *Shoe) CutCardRemaining() int {
	return s.cutCardRemaining
}

func (s *Shoe) HasActiveShoe() bool {
	return len(s.cards) > 0 && s.Remaining() > 0
}

func (s *Shoe) NeedsReshuffle() bool {
	return !s.HasActiveShoe() || s.Remaining() < MinCardsPerHand || s.Remaining() <= s.cutCardRemaining
}

func (s *Shoe) RandomizeCutCard() {
	s.cutCardRemaining = CutCardMinRemaining + s.rand.Intn(CutCardMaxRemaining-CutCardMinRemaining+1)
}

func (s *Shoe) Shuffle() {
	s.cards = s.cards[:0]
	s.drawIndex = 0
	for deck := 0; deck < DeckCount; deck++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				s.cards = append(s.cards, makeCard(suit, rank))
			}
		}
	}
	s.rand.Shuffle(len(s.cards), func(i, j int) {
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	})
}

func (s *Shoe) Draw() (Card, error) {
	if s.Remaining() <= 0 {
		return Card{}, fmt.Errorf("牌靴已空")
	}
	card := s.cards[s.drawIndex]
	s.drawIndex++
	return card, nil
}

func makeCard(suit, rank string) Card {
	var value int
	switch rank {
	case "A":
		value = 1
	case "10", "J", "Q", "K":
		value = 0
	default:
		fmt.Sscanf(rank, "%d", &value)
	}
	return Card{
		Suit:    suit,
		Rank:    rank,
		Value:   value,
		Display: rank + suit,
	}
}
