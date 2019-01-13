package bond

import (
	m "../measures"
	"fmt"
	"log"
	"math"
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
		Z_i := float64(bonds[i].Price(0, yields[i]))
		if i == 0 {
			fwd[i] = m.Rate(-math.Log(Z_i) / float64(bonds[i].Maturity-t0))
		} else {
			Z_i_1 := float64(bonds[i-1].Price(0, yields[i-1]))
			fwd[i] = m.Rate(-math.Log(Z_i/Z_i_1) / float64(bonds[i].Maturity-bonds[i-1].Maturity))
		}
	}
	return fwd
}
