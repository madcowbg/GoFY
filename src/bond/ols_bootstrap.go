package bond

import (
	m "../measures"
	"gonum.org/v1/gonum/optimize"
	"math"
)

func OLSBootstrapFromFixedCoupon(
	interpolator func(terms []m.Time, rates []m.Rate) m.SpotRate,
	quotedYields []m.Rate,
	quotedBonds []*FixedCouponBond,
	t0 m.Time,
	terms []m.Time) *FixedSpotCurve {

	rates0 := naivelyFitted(quotedYields, quotedBonds, t0, terms)

	dirtyPrices := make([]float64, len(quotedBonds))
	for i, bond := range quotedBonds {
		dirtyPrices[i] = float64(bond.Price(t0, quotedYields[i]))
	}

	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			rates := make([]m.Rate, len(terms))
			for i := range rates {
				rates[i] = m.Rate(x[i] / 100)
			}

			spotRate := interpolator(terms, rates)
			ssq := 0.0
			ss := make([]float64, len(quotedBonds))
			for i, bond := range quotedBonds {

				ss[i] = math.Abs(dirtyPrices[i] - float64(bond.PriceByDF(t0, AsDiscountFactor(spotRate))))
				ssq += ss[i]
			}
			return ssq
		},
	}

	result, err := optimize.Minimize(problem, rates0, nil, &optimize.NelderMead{})
	if err != nil {
		panic(err)
	}

	rates := make([]m.Rate, len(terms))
	for i := range rates {
		rates[i] = m.Rate(result.X[i] / 100)
	}

	return &FixedSpotCurve{Tenors: terms, Rates: rates}
}

func naivelyFitted(quotedYields []m.Rate, quotedBonds []*FixedCouponBond, t0 m.Time, terms []m.Time) []float64 {
	curve0 := NaiveBootstrapFromFixedCoupon(quotedYields, quotedBonds, t0)
	spotCurve0 := SpotCurveByConstantRateInterpolation(curve0)

	rates0 := make([]float64, len(terms))

	for i := range rates0 {
		rates0[i] = 100 * float64(ConstantSpotRateInterpolation(spotCurve0, terms[i]))
	}
	return rates0
}
