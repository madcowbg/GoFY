package bond

import (
	m "../measures"
	"math"
)

type DiscountFactor func(t m.Time) m.Money
type SpotRate func(t m.Time) m.Rate

func AsRate(discountFactor DiscountFactor) SpotRate {
	return func(t m.Time) m.Rate {
		return asRate(m.Money(discountFactor(t)), t)
	}
}

func AsDiscountFactor(rate SpotRate) DiscountFactor {
	return func(t m.Time) m.Money {
		return asDiscountFactor(rate(t), t)
	}
}

func asDiscountFactor(rate m.Rate, time m.Time) m.Money {
	return m.Money(math.Exp(-float64(rate) * float64(time)))
}

func asRate(price m.Money, ttm m.Time) m.Rate {
	return m.Rate(-math.Log(float64(price)) / float64(ttm))
}

type FixedForwardRateCurve struct {
	Maturities []m.Time
	Rates      []m.Rate
}

type FixedSpotCurve struct {
	Maturities []m.Time
	Rates      []m.Rate
}