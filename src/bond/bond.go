package bond

import (
	m "../measures"
	"gonum.org/v1/gonum/diff/fd"
	"math"
)

type Bond interface {
	TimeToExpiration(t m.Time) m.Time
	CurrentYield(t m.Time, yield m.Rate) m.Rate
	Price(t m.Time, rate m.Rate) m.Money
	PriceByDF(t m.Time, df m.DiscountFactor) m.Money
	YieldToMaturity(t m.Time, price m.Money) m.Rate
}

type Expirable struct {
	Maturity m.Time
}

func (expirable Expirable) TimeToExpiration(t m.Time) m.Time {
	return m.Time(math.Max(0, float64(expirable.Maturity-t)))
}

func Duration(bond Bond, t m.Time, rate m.Rate) float64 {
	return fd.Derivative(
		func(yield float64) float64 { return float64(bond.Price(t, m.Rate(yield))) },
		float64(rate),
		&fd.Settings{Formula: fd.Central})
}

func MacaulayDuration(bond Bond, t m.Time, rate m.Rate) float64 {
	return -Duration(bond, t, rate) / float64(bond.Price(t, rate))
}

func DollarConvexity(bond Bond, t m.Time, rate m.Rate) float64 {
	return fd.Derivative(
		func(yield float64) float64 { return float64(bond.Price(t, m.Rate(yield))) },
		float64(rate),
		&fd.Settings{Formula: fd.Central2nd})
}

func Convexity(bond Bond, t m.Time, rate m.Rate) float64 {
	return DollarConvexity(bond, t, rate) / float64(bond.Price(t, rate))
}
