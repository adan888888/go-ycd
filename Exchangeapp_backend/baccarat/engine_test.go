package baccarat

import (
	"math/rand"
	"testing"
)

func TestCalculateTotal(t *testing.T) {
	cards := []Card{
		{Value: 9},
		{Value: 8},
	}
	if got := CalculateTotal(cards); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestPrepareHandNatural(t *testing.T) {
	shoe := NewShoe(rand.New(rand.NewSource(1)))
	shoe.Shuffle()
	hand, err := PrepareHand(shoe)
	if err != nil {
		t.Fatal(err)
	}
	if len(hand.PlayerCards) < 2 || len(hand.BankerCards) < 2 {
		t.Fatalf("invalid hand: %+v", hand)
	}
	if hand.Winner == "" {
		t.Fatal("winner should not be empty")
	}
}

func TestUpdateBigRoad(t *testing.T) {
	state := NewBigRoadState()
	UpdateBigRoad(state, WinnerPlayer)
	UpdateBigRoad(state, WinnerPlayer)
	UpdateBigRoad(state, WinnerBanker)

	if state.BigRoad[0][0] != WinnerPlayer {
		t.Fatalf("expected player at [0][0], got %s", state.BigRoad[0][0])
	}
	if state.BigRoad[1][0] != WinnerPlayer {
		t.Fatalf("expected player at [1][0], got %s", state.BigRoad[1][0])
	}
	if state.BigRoad[0][1] != WinnerBanker {
		t.Fatalf("expected banker at [0][1], got %s", state.BigRoad[0][1])
	}
}

func TestShoeDraw(t *testing.T) {
	shoe := NewShoe(rand.New(rand.NewSource(2)))
	shoe.Shuffle()
	if shoe.Remaining() != TotalCards {
		t.Fatalf("expected %d cards, got %d", TotalCards, shoe.Remaining())
	}
	_, err := shoe.Draw()
	if err != nil {
		t.Fatal(err)
	}
	if shoe.Remaining() != TotalCards-1 {
		t.Fatalf("expected %d remaining", TotalCards-1)
	}
}
