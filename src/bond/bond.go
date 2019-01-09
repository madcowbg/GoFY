package bond

import (
	m "../measures"
	"math"
)

type ZeroCouponBond struct {
	Maturity m.Time
}

func (bond ZeroCouponBond) Price(t m.Time, rate m.Rate) m.Money {
	timeToExpiration := math.Max(0, float64(bond.Maturity-t))
	return m.Money(math.Exp(-timeToExpiration * float64(rate)))
}
