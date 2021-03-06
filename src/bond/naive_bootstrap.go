package bond

import (
	m "../measures"
	"fmt"
	"gonum.org/v1/gonum/optimize"
	"log"
	"math"
	"sort"
)

func NaiveBootstrapFromZCYields(yields []m.Rate, ttms []m.Time) *FixedForwardRateCurve {
	if len(yields) != len(ttms) {
		log.Fatalf("yields and times must have same length! %d != %d\n", len(yields), len(ttms))
	}
	t0 := m.Time(0)

	bonds := make([]*ZeroCouponBond, len(ttms))
	for i := 0; i < len(bonds); i++ {
		bonds[i] = &ZeroCouponBond{Expirable{ttms[i]}}
	}

	return &FixedForwardRateCurve{
		Tenors: ttms,
		Rates:  fwdRatesFromZCBonds(yields, bonds, t0)}
}

func fwdRatesFromZCBonds(yields []m.Rate, bonds []*ZeroCouponBond, t0 m.Time) []m.Rate {
	fwd := make([]m.Rate, len(yields))
	for i := range yields {
		Z_i := bonds[i].Price(0, yields[i])
		if i == 0 {
			fwd[i] = bonds[i].YieldToMaturity(t0, Z_i)
		} else {
			Z_i_1 := float64(bonds[i-1].Price(0, yields[i-1]))
			fwd[i] = m.Rate(-math.Log(float64(Z_i)/Z_i_1) / float64(bonds[i].Maturity-bonds[i-1].Maturity))
		}
	}
	return fwd
}

type timeArray []m.Time

func (s timeArray) Len() int {
	return len(s)
}
func (s timeArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s timeArray) Less(i, j int) bool {
	return s[i] < s[j]
}

func NaiveBootstrapFromFixedCoupon(quotedYields []m.Rate, quotedBonds []*FixedCouponBond, t0 m.Time) *FixedForwardRateCurve {
	if len(quotedYields) != len(quotedBonds) {
		log.Fatalf("quotedYields and times must have same length! %d != %d\n", len(quotedYields), len(quotedBonds))
	}

	ttmMap := map[m.Time]int{}
	for i, bond := range quotedBonds {
		ttmMap[bond.Maturity] = i
	}

	ttms := make([]m.Time, 0)
	for k := range ttmMap {
		ttms = append(ttms, k)
	}

	sort.Sort(timeArray(ttms))

	bonds := make([]*FixedCouponBond, len(ttms))
	yields := make([]m.Rate, len(ttms))
	for i, ttm := range ttms {
		bonds[i] = quotedBonds[ttmMap[ttm]]
		yields[i] = quotedYields[ttmMap[ttm]]
	}

	result, err := fwdRatesFromFCBonds(yields, bonds, t0)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func fwdRatesFromFCBonds(yields []m.Rate, bonds []*FixedCouponBond, t0 m.Time) (*FixedForwardRateCurve, error) {
	if len(yields) != len(bonds) {
		return nil, fmt.Errorf("yields and bonds must be of same count: %d != %d", len(yields), len(bonds))
	}

	curve := &FixedForwardRateCurve{
		Tenors: []m.Time{bonds[0].Maturity - t0},
		Rates:  []m.Rate{yields[0]},
	}

	for i, quotedYield := range yields {
		if i == 0 {
			continue
		}

		nextMaturity := bonds[i].Maturity - t0
		quotedPrice := bonds[i].Price(t0, quotedYield)
		problem := optimize.Problem{
			Func: func(x []float64) float64 {
				augmented, _ := extendedWithRate(curve, nextMaturity, m.Rate(x[0]))
				return math.Abs(float64(bonds[i].PriceByDF(t0, DFByConstantRateInterpolation(augmented)) - quotedPrice))
			},
		}
		result, err := optimize.Minimize(problem, []float64{float64(yields[i])}, nil, &optimize.NelderMead{})
		if err != nil {
			return nil, err
		}

		curve, err = extendedWithRate(curve, nextMaturity, m.Rate(result.X[0]))
	}
	return curve, nil
}

func extendedWithRate(curve *FixedForwardRateCurve, nextT m.Time, fwdRate m.Rate) (*FixedForwardRateCurve, error) {
	nPts := len(curve.Tenors)
	lastMaturity := curve.Tenors[nPts-1]
	if lastMaturity >= nextT {
		return nil, fmt.Errorf("cannot augment curve with point that is after last maturity! %f >= %f", lastMaturity, nextT)
	}
	maturities := make([]m.Time, nPts+1)
	copy(maturities, curve.Tenors)
	maturities[nPts] = nextT

	rates := make([]m.Rate, nPts+1)
	copy(rates, curve.Rates)
	rates[nPts] = fwdRate

	return &FixedForwardRateCurve{Tenors: maturities, Rates: rates}, nil
}
