package bond

import (
	m "../measures"
	"math"
)

type ZeroCouponBond struct {
	Expirable
}

func (bond *ZeroCouponBond) CurrentYield(t m.Time, rate m.Rate) m.Rate {
	return 0.0
}

func (bond *ZeroCouponBond) Price(t m.Time, rate m.Rate) m.Money {
	return m.Money(math.Exp(-float64(bond.TimeToExpiration(t)) * float64(rate)))
}

func (bond *ZeroCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	return m.Rate(-math.Log(float64(price)) / float64(bond.TimeToExpiration(t)))
}
