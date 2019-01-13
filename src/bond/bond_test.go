package bond

import (
	m "../measures"
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
)

func absCmp(tol float64) cmp.Option {
	return cmp.Comparer(func(x, y float64) bool {
		return (math.IsNaN(x) && math.IsNaN(y)) || math.Abs(x-y) < tol
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

func compareSingleCashflow(tol float64) cmp.Option {
	return cmp.Comparer(func(x, y Cashflow) bool {
		return math.Abs(float64(x.Amount-y.Amount)) < tol && math.Abs(float64(x.Time-y.Time)) < tol
	})
}

func TestCouponsAndMaturity(t *testing.T) {
	coupon := FixedCouponTerm{Frequency: 2, PerAnnum: 0.05}

	nextCoupons := []m.Time{
		coupon.NextCoupon(0, 0),
		coupon.NextCoupon(0, -0.0001),
		coupon.NextCoupon(0, 0.0001),
	}
	expectedNextCoupons := []m.Time{0, 0, 0.5}
	if !cmp.Equal(nextCoupons, expectedNextCoupons, absCmp(1e-14)) {
		t.Errorf("wrong next coupons:\n got %v\n expected %v\n", nextCoupons, expectedNextCoupons)
	}

	bond := FixedCouponBond{
		Expirable: Expirable{Maturity: 3},
		IssueTime: 0,
		Coupon:    coupon}

	nextCashflows := [][]Cashflow{
		bond.RemainingCashflows(-1),
		bond.RemainingCashflows(0.001),
		bond.RemainingCashflows(0.499999),
		bond.RemainingCashflows(0.5),
		bond.RemainingCashflows(0.500001),
		bond.RemainingCashflows(2.999999),
		bond.RemainingCashflows(3),
		bond.RemainingCashflows(3.000001),
	}

	expectedNextCashflows := [][]Cashflow{
		{
			{Time: 0, Amount: 0.025},
			{Time: 0.5, Amount: 0.025},
			{Time: 1, Amount: 0.025},
			{Time: 1.5, Amount: 0.025},
			{Time: 2, Amount: 0.025},
			{Time: 2.5, Amount: 0.025},
			{Time: 3, Amount: 1},
		}, {
			{Time: 0.5, Amount: 0.025},
			{Time: 1, Amount: 0.025},
			{Time: 1.5, Amount: 0.025},
			{Time: 2, Amount: 0.025},
			{Time: 2.5, Amount: 0.025},
			{Time: 3, Amount: 1},
		}, {
			{Time: 0.5, Amount: 0.025},
			{Time: 1, Amount: 0.025},
			{Time: 1.5, Amount: 0.025},
			{Time: 2, Amount: 0.025},
			{Time: 2.5, Amount: 0.025},
			{Time: 3, Amount: 1},
		}, {
			{Time: 0.5, Amount: 0.025},
			{Time: 1, Amount: 0.025},
			{Time: 1.5, Amount: 0.025},
			{Time: 2, Amount: 0.025},
			{Time: 2.5, Amount: 0.025},
			{Time: 3, Amount: 1},
		}, {
			{Time: 1, Amount: 0.025},
			{Time: 1.5, Amount: 0.025},
			{Time: 2, Amount: 0.025},
			{Time: 2.5, Amount: 0.025},
			{Time: 3, Amount: 1}},
		{
			{Time: 3, Amount: 1},
		}, {
			{Time: 3, Amount: 1},
		}, {},
	}
	if !cmp.Equal(nextCashflows, expectedNextCashflows, compareSingleCashflow(1e-12)) {
		t.Errorf(
			"wrong next cashflows:\n got %v\n expected %v\n%s",
			nextCashflows,
			expectedNextCashflows,
			cmp.Diff(nextCashflows, expectedNextCashflows, compareSingleCashflow(1e-12)))
	}
}

func TestFixedCouponBondPriceAndYield(t *testing.T) {
	var fixedBond Bond = &FixedCouponBond{
		Expirable: Expirable{Maturity: 3},
		IssueTime: 0,
		Coupon:    FixedCouponTerm{Frequency: 2, PerAnnum: 0.05}}

	bondPrices := []m.Money{
		fixedBond.Price(-1, 0.02),
		fixedBond.Price(0, 0.02),
		fixedBond.Price(1, 0.05),
		fixedBond.Price(3, 0.07),
		fixedBond.Price(3.1, 0.02),
	}
	expectedFixedBondPrices := []m.Money{1.0665368819865981, 1.0880823561906854, 1.0011944886073996, 1, 0}

	if !cmp.Equal(bondPrices, expectedFixedBondPrices, absCmp(1e-14)) {
		t.Errorf("wrong fixed coupon bond prices:\n got %v\n expected %v\n", bondPrices, expectedFixedBondPrices)
	}

	bondYields := []float64{
		float64(fixedBond.YieldToMaturity(0, 0.8)),
		float64(fixedBond.YieldToMaturity(0, 1.0)),
		float64(fixedBond.YieldToMaturity(0, 1.4)),
		float64(fixedBond.YieldToMaturity(-5, 1.4)),
		float64(fixedBond.YieldToMaturity(5, 1.4)),
	}

	expectedBondYields := []float64{0.1324619, 0.050635, -0.070590, -0.025293, math.NaN()}
	if !cmp.Equal(bondYields, expectedBondYields, absCmp(1e-5)) {
		t.Errorf(
			"wrong fixed coupon bond yields:\n got %v\n expected %v\n%s\n",
			bondYields, expectedBondYields, cmp.Diff(bondYields, expectedBondYields, absCmp(1e-5)))
	}
}
