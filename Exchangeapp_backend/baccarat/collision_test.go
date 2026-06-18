package baccarat

import (
	"math/rand"
	"testing"
)

func TestPickRandomBankerPlayer(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	banker := 0
	for i := 0; i < 1000; i++ {
		if PickRandomBankerPlayer(100, rnd) == "庄" {
			banker++
		}
	}
	if banker != 1000 {
		t.Fatalf("expected all banker, got %d", banker)
	}
}

func TestCollisionWinLoss(t *testing.T) {
	stats := &CollisionStats{}
	recordCollision(stats, "庄", "庄")
	recordCollision(stats, "庄", "闲")
	recordCollision(stats, "闲", "和")
	finalizeCollision(stats)
	if stats.WinCount != 1 || stats.LossCount != 1 || stats.TieCount != 1 || stats.NetWin != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	if stats.RandomBankerPlayerDiff != 1 {
		t.Fatalf("expected random diff 1, got %d", stats.RandomBankerPlayerDiff)
	}
	if stats.WinRate != 50 {
		t.Fatalf("expected 50%% win rate, got %v", stats.WinRate)
	}
}
