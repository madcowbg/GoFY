package option

import "math"

type EuropeanOption struct {
	Strike Money
	T      Time
}

type EuropeanCallOption struct {
	EuropeanOption
}

func (option *EuropeanCallOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(spot-option.Strike)))
}

func (option *EuropeanOption) Maturity() Time {
	return option.T
}
