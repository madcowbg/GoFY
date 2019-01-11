package option

import (
	m "../measures"
	"github.com/phil-mansfield/gotetra/math/interpolate"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type Pricing func(option Option, spot m.Money, t m.Time) m.Money
type Greek func(option Option, spot m.Money, t m.Time) float64

type Decision interface {
	EarlyExcercise(spot m.Money, nonExcercisedValue m.Money) m.Money
}

type Option interface {
	Decision
	Expiration() m.Time
	Payoff(spot m.Money) m.Money

	Strike() m.Money // FIXME deprecated ... does not generalize!
}

type PricingParameters struct {
	Sigma m.Return
	R     m.Rate
}

func BinomialPricing(parameters PricingParameters) Pricing {
	return func(option Option, spot m.Money, t m.Time) m.Money {
		return BinomialModel(
			option,
			spot,
			t,
			parameters.Sigma,
			parameters.R)
	}
}

func GridPricing(parameters PricingParameters) Pricing {
	return func(option Option, spot m.Money, t m.Time) m.Money {
		SInf := math.Max(2.0*float64(option.Strike()), 1.1*float64(spot))
		NAS := 200

		S, V := FiniteDifferenceGrid(NAS, SInf)(
			option,
			t,
			parameters.Sigma,
			parameters.R)

		V0 := make([]float64, len(S))
		for i := 0; i < len(S); i++ {
			V0[i] = V[i][len(V[0])-1]
		}

		return m.Money(interpolate.NewLinear(S, V0).Eval(float64(spot)))
	}
}

func EuropeanMCPricing(parameters PricingParameters) Pricing {
	return func(option Option, spot m.Money, t m.Time) m.Money {
		return NoEarlyExcerciseMonteCarloModel(100000, 100, -1)(option, spot, t, parameters.Sigma, parameters.R)
	}
}

func diff(f func(x float64) float64, x, d float64) float64 {
	return (f(x+d) - f(x-d)) / (2 * d)
}

func diff2nd(f func(x float64) float64, x, d float64) float64 {
	return (f(x+d) - 2*f(x) + f(x-d)) / (d * d)
}

func Delta(pricing Pricing) Greek {
	return func(option Option, spot m.Money, t m.Time) float64 {
		return diff(
			func(x float64) float64 { return float64(pricing(option, m.Money(x), t)) },
			float64(spot),
			0.01*float64(spot))
	}
}

func Gamma(pricing Pricing) Greek {
	return func(option Option, spot m.Money, t m.Time) float64 {
		return diff2nd(
			func(x float64) float64 { return float64(pricing(option, m.Money(x), t)) },
			float64(spot),
			0.01*float64(spot))
	}
}

func Theta(pricing Pricing) Greek {
	return func(option Option, spot m.Money, t m.Time) float64 {
		return diff(
			func(x float64) float64 { return float64(pricing(option, spot, m.Time(x))) },
			float64(t),
			0.001)
	}
}

func Rho(pricingFromParameters func(parameters PricingParameters) Pricing, parameters PricingParameters) Greek {
	return func(option Option, spot m.Money, t m.Time) float64 {
		return diff(
			func(r float64) float64 {
				tweakedParameters := parameters
				tweakedParameters.R = m.Rate(r)
				return float64(pricingFromParameters(tweakedParameters)(option, spot, t))
			},
			float64(parameters.R),
			0.0001)
	}
}

func ImplyVol(pricingMethod func(PricingParameters) Pricing, R m.Rate) func(option Option, spot m.Money, t m.Time) func(price m.Money) (float64, error) {
	return func(option Option, spot m.Money, t m.Time) func(price m.Money) (float64, error) {
		return func(price m.Money) (float64, error) {
			problem := optimize.Problem{
				Func: func(x []float64) float64 {
					parameters := PricingParameters{Sigma: m.Return(x[0]), R: R}
					pricing := pricingMethod(parameters)
					return math.Abs(float64(pricing(option, spot, t) - price))
				},
			}
			result, err := optimize.Minimize(problem, []float64{0.2}, nil, &optimize.NelderMead{})
			if err != nil {
				return math.NaN(), err
			}
			return result.X[0], nil
		}
	}
}
