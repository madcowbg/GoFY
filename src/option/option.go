package option

import (
	"math"
)

type Time float64
type Money float64
type Return float64
type Rate float64

type Pricing func(option Option, spot Money, t Time) Money
type Greek func(option Option, spot Money, t Time) float64

type Decision interface {
	EarlyExcercise(spot Money, nonExcercisedValue Money) Money
}

type Option interface {
	Decision
	Maturity() Time
	Payoff(spot Money) Money
}

type PricingParameters struct {
	Sigma Return
	R     Rate
}

func Price(parameters PricingParameters) Pricing {
	return func(option Option, spot Money, t Time) Money {
		return binomialModel(
			option,
			spot,
			t,
			parameters.Sigma,
			parameters.R)
	}
}

func binomialModel(option Option, spot Money, t Time, sigmaX Return, rf Rate) Money {
	if t >= option.Maturity() {
		return option.Payoff(spot)
	}

	sigma := float64(sigmaX)
	r := float64(rf)

	nsteps := 1000
	step := (float64(option.Maturity()) - float64(t)) / float64(nsteps)

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
			S[j] = Money(float64(S[j]) / d)
		}
		S[n] = Money(math.NaN())

		for j := 0; j < n; j++ {
			V[j] = option.EarlyExcercise(
				option.Payoff(S[j]),
				Money((p*float64(V[j+1])+(1-p)*float64(V[j]))*discountFactor))
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

func Delta(pricing Pricing) Greek {
	return func(option Option, spot Money, t Time) float64 {
		return diff(
			func(x float64) float64 { return float64(pricing(option, Money(x), t)) },
			float64(spot),
			0.01*float64(spot))
	}
}

func Gamma(pricing Pricing) Greek {
	return func(option Option, spot Money, t Time) float64 {
		return diff2nd(
			func(x float64) float64 { return float64(pricing(option, Money(x), t)) },
			float64(spot),
			0.01*float64(spot))
	}
}

func Theta(pricing Pricing) Greek {
	return func(option Option, spot Money, t Time) float64 {
		return diff(
			func(x float64) float64 { return float64(pricing(option, spot, Time(x))) },
			float64(t),
			0.001)
	}
}

func Rho(parameters PricingParameters) Greek {
	return func(option Option, spot Money, t Time) float64 {
		return diff(
			func(r float64) float64 {
				tweakedParameters := parameters
				tweakedParameters.R = Rate(r)
				return float64(Price(tweakedParameters)(option, spot, t))
			},
			float64(parameters.R),
			0.0001)
	}
}
