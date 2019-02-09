package bond

import (
	m "../measures"
	"math"
)

type DiscountFactor func(t m.Time) float64
type SpotRate func(t m.Time) m.Rate

func AsRate(discountFactor DiscountFactor) SpotRate {
	return func(t m.Time) m.Rate {
		return (&ZeroCouponBond{Expirable{t}}).YieldToMaturity(0, m.Money(discountFactor(t)))
	}
}

func AsDiscountFactor(rate SpotRate) DiscountFactor {
	return func(t m.Time) float64 {
		return discountFactor(rate(t), t)
	}
}

func discountFactor(rate m.Rate, time m.Time) float64 {
	return math.Exp(-float64(rate) * float64(time))
}

type FixedForwardRateCurve struct {
	Maturities []m.Time
	Rates      []m.Rate
}

type FixedSpotCurve struct {
	Maturities []m.Time
	Rates      []m.Rate
}
