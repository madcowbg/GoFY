package option

import (
	m "../measures"
	"math"
)

type AmericanOption struct {
	VanillaOption
}

func (option *AmericanOption) EarlyExcercise(excercisePayoff m.Money, nonExcercisedValue m.Money) m.Money {
	return m.Money(math.Max(float64(excercisePayoff), float64(nonExcercisedValue)))
}

type AmericanCallOption struct {
	AmericanOption
}

func (option *AmericanCallOption) Payoff(spot m.Money) m.Money {
	return m.Money(math.Max(0.0, float64(spot-option.S)))
}

type AmericanPutOption struct {
	AmericanOption
}

func (option *AmericanPutOption) Payoff(spot m.Money) m.Money {
	return m.Money(math.Max(0.0, float64(option.S-spot)))
}
