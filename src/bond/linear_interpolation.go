package bond

import (
	m "../measures"
	"math"
	. "sort"
)

func DFByConstantRateInterpolation(curve *FixedForwardRateCurve) m.DiscountFactor {
	return func(t m.Time) m.Money {
		if t < 0 {
			return m.Money(math.NaN())
		}

		df := m.Money(1.0)
		appliedTime := m.Time(0.0)
		for i := 0; i < len(curve.Tenors); i++ {
			if appliedTime >= t {
				break
			}

			timeToDiscount := m.Time(math.Min(float64(t), float64(curve.Tenors[i]))) - appliedTime
			df *= asDiscountFactor(curve.Rates[i], timeToDiscount)

			appliedTime += timeToDiscount
		}

		if appliedTime < t {
			df *= asDiscountFactor(curve.Rates[len(curve.Rates)-1], t-appliedTime)
		}
		return df
	}
}

func SpotCurveByConstantRateInterpolation(curve *FixedForwardRateCurve) *FixedSpotCurve {
	rate := AsSpotRate(DFByConstantRateInterpolation(curve))

	rates := make([]m.Rate, len(curve.Tenors))
	for i, ttm := range curve.Tenors {
		rates[i] = rate(ttm)
	}

	return &FixedSpotCurve{
		Tenors: curve.Tenors,
		Rates:  rates,
	}
}

func ConstantSpotRateInterpolation(curve *FixedSpotCurve, t m.Time) m.Rate {
	i := Search(len(curve.Rates), func(i int) bool { return curve.Tenors[i] >= t })
	if i >= len(curve.Rates) {
		return curve.Rates[i-1]
	} else {
		return curve.Rates[i]
	}
}
