package bond

import (
	m "../measures"
	"fmt"
	"gonum.org/v1/gonum/optimize"
	"log"
	"math"
	"sort"
)

type DiscountFactor func(t m.Time) float64

type FixedForwardRateCurve struct {
	Maturities []m.Time
	Rates      []m.Rate
}

func Yield(discountFactor DiscountFactor, t m.Time) m.Rate {
	return (&ZeroCouponBond{Expirable{t}}).YieldToMaturity(0, m.Money(discountFactor(t)))
}

func (curve *FixedForwardRateCurve) DiscountFactor(t m.Time) float64 {
	if t < 0 {
		return math.NaN()
	}

	df := 1.0
	appliedTime := m.Time(0.0)
	for i := 0; i < len(curve.Maturities); i++ {
		if appliedTime >= t {
			break
		}

		timeToDiscount := m.Time(math.Min(float64(t), float64(curve.Maturities[i]))) - appliedTime
		df *= dfFor(curve.Rates[i], timeToDiscount)

		appliedTime += timeToDiscount
	}

	if appliedTime < t {
		df *= dfFor(curve.Rates[len(curve.Rates)-1], t-appliedTime)
	}
	return df
}

func dfFor(r m.Rate, t m.Time) float64 {
	return math.Exp(-float64(r) * float64(t))
}

func BootstrapForwardRates(yields []m.Rate, ttms []m.Time) *FixedForwardRateCurve {
	if len(yields) != len(ttms) {
		log.Fatalf("yields and times must have same length! %d != %d\n", len(yields), len(ttms))
	}
	t0 := m.Time(0)

	bonds := make([]*ZeroCouponBond, len(ttms))
	for i := 0; i < len(bonds); i++ {
		bonds[i] = &ZeroCouponBond{Expirable{ttms[i]}}
	}

	return &FixedForwardRateCurve{
		Maturities: ttms,
		Rates:      fwdRatesFromZCBonds(yields, bonds, t0)}
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

type TimeArray []m.Time

func (s TimeArray) Len() int {
	return len(s)
}
func (s TimeArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TimeArray) Less(i, j int) bool {
	return s[i] < s[j]
}

func BootstrapForwardRatesFromFixedCoupon(quotedYields []m.Rate, quotedBonds []*FixedCouponBond) *FixedForwardRateCurve {
	if len(quotedYields) != len(quotedBonds) {
		log.Fatalf("quotedYields and times must have same length! %d != %d\n", len(quotedYields), len(quotedBonds))
	}
	t0 := m.Time(0)

	ttmMap := map[m.Time]int{}
	for i, bond := range quotedBonds {
		ttmMap[bond.Maturity] = i
	}

	ttms := make([]m.Time, 0)
	for k := range ttmMap {
		ttms = append(ttms, k)
	}

	sort.Sort(TimeArray(ttms))

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
		Maturities: []m.Time{bonds[0].Maturity},
		Rates:      []m.Rate{yields[0]},
	}

	for i, quotedYield := range yields {
		if i == 0 {
			continue
		}

		nextMaturity := bonds[i].Maturity
		quotedPrice := bonds[i].Price(t0, quotedYield)
		problem := optimize.Problem{
			Func: func(x []float64) float64 {
				augmented, _ := AugmentedCurve(curve, nextMaturity, m.Rate(x[0]))
				return math.Abs(float64(bonds[i].PriceByDF(t0, augmented.DiscountFactor) - quotedPrice))
			},
		}
		result, err := optimize.Minimize(problem, []float64{float64(yields[i])}, nil, &optimize.NelderMead{})
		if err != nil {
			return nil, err
		}

		curve, err = AugmentedCurve(curve, nextMaturity, m.Rate(result.X[0]))
	}
	return curve, nil
}

func AugmentedCurve(curve *FixedForwardRateCurve, nextT m.Time, fwdRate m.Rate) (*FixedForwardRateCurve, error) {
	nPts := len(curve.Maturities)
	lastMaturity := curve.Maturities[nPts-1]
	if lastMaturity >= nextT {
		return nil, fmt.Errorf("cannot augment curve with point that is after last maturity! %f >= %f", lastMaturity, nextT)
	}
	maturities := make([]m.Time, nPts+1)
	copy(maturities, curve.Maturities)
	maturities[nPts] = nextT

	rates := make([]m.Rate, nPts+1)
	copy(rates, curve.Rates)
	rates[nPts] = fwdRate

	return &FixedForwardRateCurve{Maturities: maturities, Rates: rates}, nil
}
