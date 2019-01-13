package bond

import (
	m "../measures"
	"log"
	"math"
)

func BootstrapForwardRates(yields []m.Rate, ttms []m.Time) []m.Rate {
	if len(yields) != len(ttms) {
		log.Fatalf("yields and times must have same length! %d != %d\n", len(yields), len(ttms))
	}
	t0 := m.Time(0)

	bonds := make([]*ZeroCouponBond, len(ttms))
	for i := 0; i < len(bonds); i++ {
		bonds[i] = &ZeroCouponBond{Expirable{ttms[i]}}
	}
	return fwdRatesFromZCBonds(yields, bonds, t0)
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
