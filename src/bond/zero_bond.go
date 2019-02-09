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

func (bond *ZeroCouponBond) PriceByDF(t m.Time, df DiscountFactor) m.Money {
	return m.Money(df(bond.TimeToExpiration(t)))
}

func (bond *ZeroCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	return asRate(price, bond.TimeToExpiration(t))
}
