package main

import (
	"fmt"
	"math"
)

type Time float64
type Price float64
type Return float64
type Rate float64

type Option interface {
	payoff(spot Price) Price
	price(spot Price, t Time) Price
}

type OptionParameters struct {
	sigma Return
	r     Rate
}

type CallOption struct {
	parameters OptionParameters
	strike     Price
	T          Time
}

func (option *CallOption) payoff(spot Price) Price {
	return Price(math.Max(0.0, float64(spot-option.strike)))
}

func (option *CallOption) price(spot Price, t Time) Price {
	if t >= option.T {
		return option.payoff(spot)
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

	S := make([]Price, nsteps+1)
	S[0] = spot

	for n := 1; n <= nsteps; n++ {
		for j := n; j >= 1; j-- {
			S[j] = Price(u * float64(S[j-1]))
		}
		S[0] = Price(d * float64(S[0]))
	}

	V := make([]Price, nsteps+1)
	for j := 0; j <= nsteps; j++ {
		V[j] = option.payoff(S[j])
	}

	for n := nsteps; n >= 1; n-- {
		for j := 0; j < n; j++ {
			V[j] = Price((p*float64(V[j+1]) + (1-p)*float64(V[j])) * discountFactor)
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

func delta(option Option, spot Price, t Time) float64 {
	return diff(
		func(x float64) float64 { return float64(option.price(Price(x), t)) },
		float64(spot),
		0.01*float64(spot))
}

func gamma(option Option, spot Price, t Time) float64 {
	return diff2nd(
		func(x float64) float64 { return float64(option.price(Price(x), t)) },
		float64(spot),
		0.01*float64(spot))
}

func (option *CallOption) rho(spot Price, t Time) float64 {
	return diff(
		func(r float64) float64 {
			tweaked := *option
			tweaked.parameters.r = Rate(r)
			return float64(tweaked.price(spot, t))
		},
		float64(option.parameters.r),
		0.0001)
}

func theta(option Option, spot Price, t Time) float64 {
	return diff(
		func(x float64) float64 { return float64(option.price(spot, Time(x))) },
		float64(t),
		0.001)
}

func main() {
	// dec := json.NewDecoder(os.Stdin)
	// enc := json.NewEncoder(os.Stdout)
	opt := &CallOption{OptionParameters{0.2, 0.02}, 100, 1}
	for i := 1; i < 20; i++ {
		spot := Price(i * 10.0)
		fmt.Printf("S=%f V(T)=%f\n", float64(spot), opt.payoff(spot))
		fmt.Printf("V(0)=%f\n", opt.price(spot, 0))
		fmt.Printf("Delta(0)=%f\n", delta(opt, spot, 0))
		fmt.Printf("Gamma(0)=%f\n", gamma(opt, spot, 0))
		fmt.Printf("Rho(0)=%f\n", opt.rho(spot, 0))
		fmt.Printf("Theta(0)=%f\n", theta(opt, spot, 0.0))
	}
}
