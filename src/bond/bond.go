package bond

import (
	m "../measures"
	"math"
)

type Bond interface {
	TimeToExpiration(t m.Time) m.Time
	CurrentYield(t m.Time, yield m.Rate) m.Rate
	Price(t m.Time, rate m.Rate) m.Money
	YieldToMaturity(t m.Time, price m.Money) m.Rate
}

type Expirable struct {
	Maturity m.Time
}

func (expirable Expirable) TimeToExpiration(t m.Time) m.Time {
	return m.Time(math.Max(0, float64(expirable.Maturity-t)))
}
