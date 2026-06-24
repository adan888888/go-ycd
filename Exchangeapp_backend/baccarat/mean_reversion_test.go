package baccarat

import "testing"

func TestMeanReversionZhuangZhanBi(t *testing.T) {
	cases := []struct {
		netWin int
		want   int
	}{
		{-51, 80},
		{-50, 65},
		{-21, 65},
		{-20, 55},
		{0, 55},
		{19, 55},
		{20, 45},
		{49, 45},
		{50, 35},
		{100, 35},
	}
	for _, c := range cases {
		if got := MeanReversionZhuangZhanBi(c.netWin); got != c.want {
			t.Fatalf("S=%d: expected %d, got %d", c.netWin, c.want, got)
		}
	}
}

func TestSimulateShoesMeanReversion(t *testing.T) {
	result := SimulateShoesMeanReversion(5, 1907650735441448960, 0)
	if result.ShoeCount != 5 {
		t.Fatalf("expected 5 shoes, got %d", result.ShoeCount)
	}
	if result.Collision == nil {
		t.Fatal("expected collision stats")
	}
	if result.Collision.ZhuangZhanBi < 35 || result.Collision.ZhuangZhanBi > 80 {
		t.Fatalf("unexpected final zhuang zhan bi: %d", result.Collision.ZhuangZhanBi)
	}
}
