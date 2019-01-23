package option

import (
	m "../measures"
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

func equalTol(tol float64) func(t *testing.T, got, want []float64, msg string) {
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

func equalCallTot(msg string, t *testing.T, xs []m.Money, f func(price m.Money) float64, want []float64, tol float64) {
	vals := make([]float64, len(xs))
	for i := 0; i < len(xs); i++ {
		vals[i] = f(xs[i])
	}
	equalTol(tol)(t, vals, want, msg)
}

func TestEuropeanCallATMStatsAreCorrect(t *testing.T) {
	pricingParameters := PricingParameters{0.2, 0.02}
	pricing := BinomialPricing(pricingParameters)

	var callOpt Option = &EuropeanCallOption{EuropeanOption{VanillaOption{100, 1}}}

	spots := make([]m.Money, 10)
	for i := 0; i < 10; i++ {
		spots[i] = m.Money(100.0 + float64(i-5)*15.0)
	}

	equalCallTot(
		"payoffs",
		t, spots,
		func(price m.Money) float64 { return float64(callOpt.Payoff(price)) },
		[]float64{0, 0, 0, 0, 0, 0, 15, 30, 45, 60},
		1e-5)

	equalCallTot(
		"prices",
		t, spots,
		func(price m.Money) float64 { return float64(pricing(callOpt, price, 0)) },
		[]float64{0.00000, 0.00001, 0.00816, 0.31394, 2.54571, 8.91424, 19.52944, 32.78393, 47.20572, 62.03824},
		1e-5)

	equalCallTot(
		"deltas",
		t, spots,
		func(price m.Money) float64 { return Delta(pricing)(callOpt, price, 0) },
		[]float64{0.00000, 0.00001, 0.00267, 0.05724, 0.26903, 0.57924, 0.81520, 0.93511, 0.98004, 0.99455},
		1e-5)

	equalCallTot(
		"gammas",
		t, spots,
		func(price m.Money) float64 { return Gamma(pricing)(callOpt, price, 0) },
		[]float64{0.00000, 0.00000, 0.00071, 0.00788, 0.01976, 0.02473, 0.01378, 0.00450, 0.00156, 0.00050},
		1e-5)

	equalCallTot(
		"thetas",
		t, spots,
		func(price m.Money) float64 { return Theta(pricing)(callOpt, price, 0) },
		[]float64{-0.00000, -0.00011, -0.04534, -0.84290, -3.26983, -4.88992, -4.48796, -3.46209, -2.58404, -2.17636},
		1e-5)

	equalCallTot(
		"rhos",
		t, spots,
		func(price m.Money) float64 { return Rho(BinomialPricing, pricingParameters)(callOpt, price, 0) },
		[]float64{0.00000, 0.00022, 0.13612, 3.64857, 20.40996, 49.01760, 74.27877, 88.73383, 94.92619, 97.10363},
		1e-5)
}

func TestPutCallParity(t *testing.T) {
	R := m.Rate(0.02)
	pricingParameters := PricingParameters{0.2, R}
	pricing := BinomialPricing(pricingParameters)

	callOpt := &EuropeanCallOption{EuropeanOption: EuropeanOption{VanillaOption{S: 100, T: 1}}}
	putOpt := &EuropeanPutOption{EuropeanOption: EuropeanOption{VanillaOption{S: 100, T: 1}}}

	spot := m.Money(100)
	C := pricing(callOpt, spot, 0)
	P := pricing(putOpt, spot, 0)
	PVS := float64(spot) * math.Exp(-1*float64(R))

	if !cmp.Equal(PVS, float64(P-C+spot), absCmp(1e-5)) {
		t.Errorf("c %f\np %f\nP - C + spot = %f\nPV(x) = %f\n", C, P, P-C+spot, PVS)
	}
}

func TestAmericanPutATMStatsAreCorrect(t *testing.T) {
	pricingParameters := PricingParameters{0.2, 0.02}
	pricing := BinomialPricing(pricingParameters)

	var putOpt Option = &AmericanPutOption{AmericanOption{VanillaOption{100, 1}}}

	spots := make([]m.Money, 10)
	for i := 0; i < 10; i++ {
		spots[i] = m.Money(100.0 + float64(i-5)*15.0)
	}

	equalCallTot(
		"payoffs",
		t, spots,
		func(price m.Money) float64 { return float64(putOpt.Payoff(price)) },
		[]float64{75.00000, 60.00000, 45.00000, 30.00000, 15.00000, 0.00000, 0.00000, 0.00000, 0.00000, 0.00000},
		1e-5)

	equalCallTot(
		"prices",
		t, spots,
		func(price m.Money) float64 { return float64(pricing(putOpt, price, 0)) },
		[]float64{75.00000, 60.00000, 45.00000, 30.00000, 16.16243, 7.10987, 2.59578, 0.81507, 0.22814, 0.05866},
		1e-5)

	equalCallTot(
		"deltas",
		t, spots,
		func(price m.Money) float64 { return Delta(pricing)(putOpt, price, 0) },
		[]float64{-1.00000, -1.00000, -1.00000, -1.00000, -0.77702, -0.43565, -0.18898, -0.06604, -0.02022, -0.00550},
		1e-5)

	equalCallTot(
		"gammas",
		t, spots,
		func(price m.Money) float64 { return Gamma(pricing)(putOpt, price, 0) },
		[]float64{0.00000, 0.00000, -0.00000, -0.00000, 0.02310, 0.02482, 0.01400, 0.00463, 0.00159, 0.00051},
		1e-5)

	equalCallTot(
		"thetas",
		t, spots,
		func(price m.Money) float64 { return Theta(pricing)(putOpt, price, 0) },
		[]float64{-0.00000, -0.00000, -0.00000, 0.00000, -1.67956, -3.13576, -2.62546, -1.53017, -0.63327, -0.21920},
		1e-5)

	equalCallTot(
		"rhos",
		t, spots,
		func(price m.Money) float64 { return Rho(BinomialPricing, pricingParameters)(putOpt, price, 0) },
		[]float64{-0.00000, 0.00000, -0.00000, 0.00000, -39.06065, -38.36523, -21.09362, -8.68014, -2.96423, -0.88939},
		1e-5)
}

func TestCompareGridVsBinomial(t *testing.T) {
	putOpt := &AmericanPutOption{AmericanOption{VanillaOption{100, 1}}}
	callOpt := &EuropeanCallOption{EuropeanOption{VanillaOption{100, 1}}}

	checkGridVsBinomial(t, putOpt, m.Money(100))
	checkGridVsBinomial(t, callOpt, m.Money(100))

	checkGridVsBinomial(t, putOpt, m.Money(500))
	checkGridVsBinomial(t, callOpt, m.Money(500))

	checkGridVsBinomial(t, putOpt, m.Money(20))
	checkGridVsBinomial(t, callOpt, m.Money(20))
}

func checkGridVsBinomial(t *testing.T, opt Option, spot m.Money) {
	parameters := PricingParameters{0.2, 0.02}
	binomial := BinomialPricing(parameters)
	grid := GridPricing(parameters)

	compareTwo("Bin vs Grid", t, "pricing", float64(binomial(opt, spot, 0)), float64(grid(opt, spot, 0)), absCmp(0.005))
	compareTwo("Bin vs Grid", t, "Delta", Delta(binomial)(opt, spot, 0), Delta(grid)(opt, spot, 0), absCmp(0.0001))
	compareTwo("Bin vs Grid", t, "Gamma", Gamma(binomial)(opt, spot, 0), Gamma(grid)(opt, spot, 0), absCmp(0.01))
	compareTwo("Bin vs Grid", t, "Theta", Theta(binomial)(opt, spot, 0), Theta(grid)(opt, spot, 0), absCmp(0.005))
	compareTwo("Bin vs Grid", t, "Delta", Rho(BinomialPricing, parameters)(opt, spot, 0), Rho(GridPricing, parameters)(opt, spot, 0), absCmp(0.1))
}

func compareTwo(desc string, t *testing.T, msg string, binv, gridv float64, comp cmp.Option) {
	if !cmp.Equal(gridv, binv, comp) {
		t.Errorf("%s %s is diff: %f = %f, diff = %f\n", desc, msg, binv, gridv, math.Abs(binv-gridv))
	}
}

func TestCompareGridVsEuropeanMC(t *testing.T) {
	opt := &EuropeanCallOption{EuropeanOption{VanillaOption{100, 1}}}
	spot := m.Money(100)

	parameters := PricingParameters{0.2, 0.02}

	grid := GridPricing(parameters)
	mc := EuropeanMCPricing(parameters)

	compareTwo("Grid vs MC", t, "pricing", float64(grid(opt, spot, 0)), float64(mc(opt, spot, 0)), absCmp(0.05))
	compareTwo("Grid vs MC", t, "Delta", Delta(grid)(opt, spot, 0), Delta(mc)(opt, spot, 0), absCmp(0.001))
	compareTwo("Grid vs MC", t, "Gamma", Gamma(grid)(opt, spot, 0), Gamma(mc)(opt, spot, 0), absCmp(0.01))
	compareTwo("Grid vs MC", t, "Theta", Theta(grid)(opt, spot, 0), Theta(mc)(opt, spot, 0), absCmp(0.05))
	compareTwo("Grid vs MC", t, "Delta", Rho(BinomialPricing, parameters)(opt, spot, 0), Rho(GridPricing, parameters)(opt, spot, 0), absCmp(0.1))
}

func TestImplyVol(t *testing.T) {
	opt := &EuropeanCallOption{EuropeanOption{VanillaOption{100, 1}}}
	spot := m.Money(100)
	R := m.Rate(0.02)

	prices := []m.Money{6, 5.6, 7, 8}
	for _, price := range prices {
		checkVolImplyAtPrice(R, opt, spot, price, t, 1e-4)
	}
	extremePrices := []m.Money{1, 20}
	for _, price := range extremePrices {
		checkVolImplyAtPrice(R, opt, spot, price, t, 1e-3)
	}
}

func checkVolImplyAtPrice(R m.Rate, opt Option, spot m.Money, price m.Money, t *testing.T, tol float64) {
	implGrid, err := ImplyVol(GridPricing, R)(opt, spot, 0)(price)
	if err != nil || math.IsNaN(implGrid) {
		t.Error("failed imply for grid", price, implGrid, err)
	}
	implBinomial, err := ImplyVol(BinomialPricing, R)(opt, spot, 0)(price)
	if err != nil || math.IsNaN(implBinomial) {
		t.Error("failed imply for binomial", price, implBinomial, err)
	}
	if !cmp.Equal(implGrid, implBinomial, absCmp(tol)) {
		t.Errorf("imply is different at price=%f: grid=%f != binomial=%f\n", price, implGrid, implBinomial)
	}
}

func TestExpiredPricing(t *testing.T) {
	opt := &EuropeanCallOption{EuropeanOption{VanillaOption{100, 1}}}
	spot := m.Money(100)

	binPrice := BinomialPricing(PricingParameters{0.2, 0.02})(opt, spot, 2)
	gridPrice := BinomialPricing(PricingParameters{0.2, 0.02})(opt, spot, 2)
	mcPrice := BinomialPricing(PricingParameters{0.2, 0.02})(opt, spot, 2)
	if !cmp.Equal(binPrice, gridPrice, absCmp(1e-10)) || !cmp.Equal(gridPrice, mcPrice, absCmp(1e-10)) {
		t.Errorf("Differences in pricing of expired option: %f != %f or %f != %f\n", binPrice, gridPrice, gridPrice, mcPrice)
	}
}
