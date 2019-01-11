package bond

import (
	m "../measures"
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
)

func absCmp(tol float64) cmp.Option {
	return cmp.Comparer(func(x, y float64) bool {
		return math.Abs(x-y) < tol
	})
}

func TestZeroCouponBondPrice(t *testing.T) {
	bond := ZeroCouponBond{Maturity: 2}

	prices := []m.Money{bond.Price(1, 0.02), bond.Price(1, 0.03), bond.Price(0, 0.02), bond.Price(1.5, 0.02)}
	expectedPrices := []m.Money{0.9801986733067553, 0.9704455335485082, 0.9607894391523232, 0.9900498337491681}
	if !cmp.Equal(prices, expectedPrices, absCmp(1e-14)) {
		t.Errorf("wrong bond prices:\n got %v\n expected %v\n", prices, expectedPrices)
	}
}
