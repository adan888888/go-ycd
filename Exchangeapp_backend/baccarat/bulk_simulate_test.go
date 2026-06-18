package baccarat

import "testing"

func TestSimulateShoes(t *testing.T) {
	result := SimulateShoes(10, false, 0, -1)
	if result.ShoeCount != 10 {
		t.Fatalf("expected 10 shoes, got %d", result.ShoeCount)
	}
	if result.TotalHands <= 0 {
		t.Fatal("expected positive total hands")
	}
	total := result.Stats.Player + result.Stats.Banker + result.Stats.Tie
	if total != result.TotalHands {
		t.Fatalf("stats sum %d != total hands %d", total, result.TotalHands)
	}
	if result.EffectiveHands != result.Stats.Player+result.Stats.Banker {
		t.Fatalf("effective hands mismatch")
	}
}

func TestCalcBankerMinusPlayerPer10k(t *testing.T) {
	got := calcBankerMinusPlayerPer10k(4600, 4400, 10000)
	if got != 200 {
		t.Fatalf("expected 200, got %v", got)
	}
}

func TestSimulateShoesExcludeTieFlag(t *testing.T) {
	result := SimulateShoes(5, true, 0, -1)
	if !result.ExcludeTie {
		t.Fatal("expected excludeTie true")
	}
}
