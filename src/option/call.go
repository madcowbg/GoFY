package option

import "math"

type CallOption struct {
	parameters OptionParameters
	strike     Money
	T          Time
}

func Call(sigma Return, r Rate, strike Money, T Time) *CallOption {
	return &CallOption{OptionParameters{sigma, r}, strike, T}
}

func (option *CallOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(spot-option.strike)))
}

func (option *CallOption) Maturity() Time {
	return option.T
}

func (option *CallOption) Parameters() OptionParameters {
	return option.parameters
}

func (option *CallOption) Rho(spot Money, t Time) float64 {
	return diff(
		func(r float64) float64 {
			tweaked := *option
			tweaked.parameters.r = Rate(r)
			return float64(Price(&tweaked, spot, t))
		},
		float64(option.parameters.r),
		0.0001)
}
