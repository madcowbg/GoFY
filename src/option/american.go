package option

import (
	"math"
)

type AmericanOption struct {
	VanillaOption
}

func (option *AmericanOption) EarlyExcercise(excercisePayoff Money, nonExcercisedValue Money) Money {
	return Money(math.Max(float64(excercisePayoff), float64(nonExcercisedValue)))
}

type AmericanCallOption struct {
	AmericanOption
}

func (option *AmericanCallOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(spot-option.Strike)))
}

type AmericanPutOption struct {
	AmericanOption
}

func (option *AmericanPutOption) Payoff(spot Money) Money {
	return Money(math.Max(0.0, float64(option.Strike-spot)))
}
