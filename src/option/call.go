package option

import "math"

type EuropeanOption struct {
	parameters OptionParameters
	strike     Money
	T          Time
}

type EuropeanCallOption struct {
	EuropeanOption
}

func Call(sigma Return, r Rate, strike Money, T Time) Option {
	return &EuropeanCallOption{EuropeanOption{OptionParameters{sigma, r}, strike, T}}
}

func (option *EuropeanCallOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(spot-option.strike)))
}

func (option *EuropeanOption) Maturity() Time {
	return option.T
}

func (option *EuropeanOption) Parameters() OptionParameters {
	return option.parameters
}

func (option *EuropeanCallOption) Rho(spot Money, t Time) float64 {
	return diff(
		func(r float64) float64 {
			tweaked := *option
			tweaked.parameters.r = Rate(r)
			return float64(Price(&tweaked, spot, t))
		},
		float64(option.parameters.r),
		0.0001)
}
