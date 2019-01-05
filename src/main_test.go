package main

import (
	"./option"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"math"
	"strings"
	"testing"
)

func limitedStringRepresentation(vs []float64) string {
	reprs := make([]string, len(vs))
	for i := 0; i < len(vs); i++ {
		reprs[i] = fmt.Sprintf("%.5f", vs[i])
	}
	return "[]float64{" + strings.Join(reprs, ", ") + "}"
}

func absCmp(tol float64) cmp.Option {
	return cmp.Comparer(func(x, y float64) bool {
		return math.Abs(x-y) < tol
	})
}

func EqualTol(tol float64) func(t *testing.T, got, want []float64, msg string) {
	return func(t *testing.T, got, want []float64, msg string) {
		if !cmp.Equal(got, want, absCmp(tol)) {
			t.Fatalf(
				"%s:\n got %s\n wanted %s\n diff: %v\n",
				msg,
				limitedStringRepresentation(got),
				limitedStringRepresentation(want),
				cmp.Diff(got, want, absCmp(tol)))
		}
	}
}

func EqualCallTot(msg string, t *testing.T, xs []option.Money, f func(price option.Money) float64, want []float64, tol float64) {
	vals := make([]float64, len(xs))
	for i := 0; i < len(xs); i++ {
		vals[i] = f(xs[i])
	}
	EqualTol(tol)(t, vals, want, msg)
}

func TestEuropeanCallATMStatsAreCorrect(t *testing.T) {
	pricingParameters := option.PricingParameters{0.2, 0.02}
	pricing := option.Price(pricingParameters)

	var callOpt option.Option = &option.EuropeanCallOption{option.EuropeanOption{100, 1}}

	spots := make([]option.Money, 10)
	for i := 0; i < 10; i++ {
		spots[i] = option.Money(100.0 + float64(i-5)*15.0)
	}

	EqualCallTot(
		"payoffs",
		t, spots,
		func(price option.Money) float64 { return float64(callOpt.Payoff(price)) },
		[]float64{0, 0, 0, 0, 0, 0, 15, 30, 45, 60},
		1e-5)

	EqualCallTot(
		"prices",
		t, spots,
		func(price option.Money) float64 { return float64(pricing(callOpt, price, 0)) },
		[]float64{0.00000, 0.00001, 0.00816, 0.31394, 2.54571, 8.91424, 19.52944, 32.78393, 47.20572, 62.03824},
		1e-5)

	EqualCallTot(
		"deltas",
		t, spots,
		func(price option.Money) float64 { return option.Delta(pricing)(callOpt, price, 0) },
		[]float64{0.00000, 0.00001, 0.00267, 0.05724, 0.26903, 0.57924, 0.81520, 0.93511, 0.98004, 0.99455},
		1e-5)

	EqualCallTot(
		"gammas",
		t, spots,
		func(price option.Money) float64 { return option.Gamma(pricing)(callOpt, price, 0) },
		[]float64{0.00000, 0.00000, 0.00020, 0.00276, 0.00840, 0.01236, 0.00792, 0.00293, 0.00113, 0.00040},
		1e-5)

	EqualCallTot(
		"thetas",
		t, spots,
		func(price option.Money) float64 { return option.Theta(pricing)(callOpt, price, 0) },
		[]float64{-0.00000, -0.00011, -0.04534, -0.84290, -3.26983, -4.88992, -4.48796, -3.46209, -2.58404, -2.17636},
		1e-5)

	EqualCallTot(
		"rhos",
		t, spots,
		func(price option.Money) float64 { return option.Rho(pricingParameters)(callOpt, price, 0) },
		[]float64{0.00000, 0.00022, 0.13612, 3.64857, 20.40996, 49.01760, 74.27877, 88.73383, 94.92619, 97.10363},
		1e-5)
}

func TestPutCallParity(t *testing.T) {
	R := option.Rate(0.02)
	pricingParameters := option.PricingParameters{0.2, R}
	pricing := option.Price(pricingParameters)

	callOpt := &option.EuropeanCallOption{EuropeanOption: option.EuropeanOption{Strike: 100, T: 1}}
	putOpt := &option.EuropeanPutOption{EuropeanOption: option.EuropeanOption{Strike: 100, T: 1}}

	spot := option.Money(100)
	C := pricing(callOpt, spot, 0)
	P := pricing(putOpt, spot, 0)
	PVS := float64(spot) * math.Exp(-1*float64(R))

	if !cmp.Equal(PVS, float64(P-C+spot), absCmp(1e-5)) {
		t.Errorf("c %f\np %f\nP - C + spot = %f\nPV(x) = %f\n", C, P, P-C+spot, PVS)
	}
}
