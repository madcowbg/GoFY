package bond

import (
	m "../measures"
	"math"
)

func DFByConstantRateInterpolation(curve *FixedForwardRateCurve) DiscountFactor {
	return func(t m.Time) m.Money {
		if t < 0 {
			return m.Money(math.NaN())
		}

		df := m.Money(1.0)
		appliedTime := m.Time(0.0)
		for i := 0; i < len(curve.Maturities); i++ {
			if appliedTime >= t {
				break
			}

			timeToDiscount := m.Time(math.Min(float64(t), float64(curve.Maturities[i]))) - appliedTime
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
	rates := make([]m.Rate, len(curve.Maturities))
	for i, ttm := range curve.Maturities {
		bond := &ZeroCouponBond{Expirable{ttm}}
		rates[i] = bond.YieldToMaturity(0, bond.PriceByDF(0, DFByConstantRateInterpolation(curve)))
	}

	return &FixedSpotCurve{
		Maturities: curve.Maturities,
		Rates:      rates,
	}
}
