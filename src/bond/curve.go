package bond

import (
	m "../measures"
	"math"
)

func AsSpotRate(discountFactor m.DiscountFactor) m.SpotRate {
	return func(t m.Time) m.Rate {
		return asRate(m.Money(discountFactor(t)), t)
	}
}

func AsDiscountFactor(rate m.SpotRate) m.DiscountFactor {
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
	Tenors []m.Time
	Rates  []m.Rate
}

type FixedSpotCurve struct {
	Tenors []m.Time
	Rates  []m.Rate
}

func InterpolateOnArray(interpolator func(m.Time) m.Rate, tenors []m.Time) []m.Rate {
	result := make([]m.Rate, len(tenors))
	for i, t := range tenors {
		result[i] = interpolator(t)
	}
	return result
}
