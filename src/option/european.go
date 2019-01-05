package option

import "math"

type EuropeanOption struct {
	VanillaOption
}

func (option *EuropeanOption) EarlyExcercise(excercisePayoff Money, nonExcercisedValue Money) Money {
	return nonExcercisedValue
}

type EuropeanCallOption struct {
	EuropeanOption
}

func (option *EuropeanCallOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(spot-option.Strike)))
}

type EuropeanPutOption struct {
	EuropeanOption
}

func (option *EuropeanPutOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(option.Strike-spot)))
}
