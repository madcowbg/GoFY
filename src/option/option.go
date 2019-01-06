package option

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
	Expiration() Time
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
