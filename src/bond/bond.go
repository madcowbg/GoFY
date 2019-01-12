package bond

import (
	m "../measures"
	"math"
)

type Expirable struct {
	Maturity m.Time
}

func (expirable Expirable) TimeToExpiration(t m.Time) float64 {
	return math.Max(0, float64(expirable.Maturity-t))
}

type ZeroCouponBond struct {
	Expirable
}

func (bond ZeroCouponBond) CurrentYield(t m.Time, yield m.Rate) m.Rate {
	return 0.0
}

func (bond ZeroCouponBond) Price(t m.Time, yield m.Rate) m.Money {
	return m.Money(math.Exp(-bond.TimeToExpiration(t) * float64(yield)))
}

func (bond ZeroCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	return m.Rate(-math.Log(float64(price)) / bond.TimeToExpiration(t))
}
