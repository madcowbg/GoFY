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

func TestZeroCouponBondPriceAndYield(t *testing.T) {
	bond := ZeroCouponBond{Expirable{Maturity: 2}}

	prices := []m.Money{
		bond.Price(1, 0.02),
		bond.Price(1, 0.03),
		bond.Price(0, 0.02),
		bond.Price(1.5, 0.02)}
	expectedPrices := []m.Money{0.9801986733067553, 0.9704455335485082, 0.9607894391523232, 0.9900498337491681}
	if !cmp.Equal(prices, expectedPrices, absCmp(1e-14)) {
		t.Errorf("wrong bond prices:\n got %v\n expected %v\n", prices, expectedPrices)
	}

	yields := []m.Rate{
		bond.YieldToMaturity(1, expectedPrices[0]),
		bond.YieldToMaturity(1, expectedPrices[1]),
		bond.YieldToMaturity(0, expectedPrices[2]),
		bond.YieldToMaturity(1.5, expectedPrices[3])}
	expectedYields := []m.Rate{0.02, 0.03, 0.02, 0.02}
	if !cmp.Equal(prices, expectedPrices, absCmp(1e-14)) {
		t.Errorf("wrong bond yields:\n got %v\n expected %v\n", yields, expectedYields)
	}

}
