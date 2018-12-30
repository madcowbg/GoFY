package option

import "math"

type Time float64
type Money float64
type Return float64
type Rate float64

type Option interface {
	Payoff(spot Money) Money
	Price(spot Money, t Time) Money
}

type OptionParameters struct {
	sigma Return
	r     Rate
}

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

func (option *CallOption) Price(spot Money, t Time) Money {
	if t >= option.T {
		return option.Payoff(spot)
	}

	// binomial model
	sigma := float64(option.parameters.sigma)
	r := float64(option.parameters.r)

	nsteps := 1000
	step := (float64(option.T) - float64(t)) / float64(nsteps)

	discountFactor := math.Exp(-r * step)

	temp2 := 0.5 * (discountFactor + math.Exp((r+sigma*sigma)*step))

	u := temp2 + math.Sqrt(math.Pow(temp2, 2)-1)
	d := 1 / u
	p := (math.Exp(r*step) - d) / (u - d)

	S := make([]Money, nsteps+1)
	S[0] = spot

	for n := 1; n <= nsteps; n++ {
		for j := n; j >= 1; j-- {
			S[j] = Money(u * float64(S[j-1]))
		}
		S[0] = Money(d * float64(S[0]))
	}

	V := make([]Money, nsteps+1)
	for j := 0; j <= nsteps; j++ {
		V[j] = option.Payoff(S[j])
	}

	for n := nsteps; n >= 1; n-- {
		for j := 0; j < n; j++ {
			V[j] = Money((p*float64(V[j+1]) + (1-p)*float64(V[j])) * discountFactor)
		}
	}
	return V[0]
}

func diff(f func(x float64) float64, x, d float64) float64 {
	return (f(x+d) - f(x-d)) / (2 * d)
}

func diff2nd(f func(x float64) float64, x, d float64) float64 {
	return (f(x+d) - 2*f(x) + f(x-d)) / (2 * d)
}

func Delta(option Option, spot Money, t Time) float64 {
	return diff(
		func(x float64) float64 { return float64(option.Price(Money(x), t)) },
		float64(spot),
		0.01*float64(spot))
}

func Gamma(option Option, spot Money, t Time) float64 {
	return diff2nd(
		func(x float64) float64 { return float64(option.Price(Money(x), t)) },
		float64(spot),
		0.01*float64(spot))
}

func (option *CallOption) Rho(spot Money, t Time) float64 {
	return diff(
		func(r float64) float64 {
			tweaked := *option
			tweaked.parameters.r = Rate(r)
			return float64(tweaked.Price(spot, t))
		},
		float64(option.parameters.r),
		0.0001)
}

func Theta(option Option, spot Money, t Time) float64 {
	return diff(
		func(x float64) float64 { return float64(option.Price(spot, Time(x))) },
		float64(t),
		0.001)
}
