package option

import "math"
import m "../measures"

type EuropeanOption struct {
	VanillaOption
}

func (option *EuropeanOption) EarlyExcercise(excercisePayoff m.Money, nonExcercisedValue m.Money) m.Money {
	return nonExcercisedValue
}

type EuropeanCallOption struct {
	EuropeanOption
}

func (option *EuropeanCallOption) Payoff(spot m.Money) m.Money {
	return m.Money(math.Max(0.0, float64(spot-option.S)))
}

type EuropeanPutOption struct {
	EuropeanOption
}

func (option *EuropeanPutOption) Payoff(spot m.Money) m.Money {
	return m.Money(math.Max(0.0, float64(option.S-spot)))
}
