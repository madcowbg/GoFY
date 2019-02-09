package main

import (
	"./bond"
	"./bond/monotone_convex"
	m "./measures"
	"./option"
	"fmt"
	"reflect"
)

//func main() {
//	run_bond()
//	run_option()
//	run_mc()
//}

func run_mc() {
	Terms := []m.Time{1, 2, 3, 4, 5}
	Rates := []m.Rate{0.03, 0.04, 0.047, 0.06, 0.06}

	fmt.Printf("term\trate\tfwd\n")
	for i := 0; i < 520; i++ {
		t := m.Time(i) / 100
		fmt.Printf(
			"%f\t%f\t%f\n",
			t,
			monotone_convex.SpotRateInterpolator(Terms, Rates)(t),
			monotone_convex.ForwardRateInterpolator(Terms, Rates)(t))
	}
}

func run_bond() {
	zcBond := bond.ZeroCouponBond{bond.Expirable{1}}
	rate := m.Rate(0.1)

	fmt.Printf("T=%f, r=%f, Price=%f\n", zcBond.Maturity, rate, zcBond.Price(0, rate))
}

func run_option() {
	// dec := json.NewDecoder(os.Stdin)
	// enc := json.NewEncoder(os.Stdout)
	parameters := option.PricingParameters{0.2, 0.02}

	printGreekDiagATM(
		option.BinomialPricing,
		parameters,
		&option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		option.GridPricing,
		parameters,
		&option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		option.EuropeanMCPricing,
		parameters,
		&option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		option.BinomialPricing,
		parameters,
		&option.EuropeanPutOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		option.BinomialPricing,
		parameters,
		&option.AmericanPutOption{option.AmericanOption{option.VanillaOption{100, 1}}})
}

func printGreekDiagATM(pricingFun func(parameters option.PricingParameters) option.Pricing, pricingParameters option.PricingParameters, opt option.Option) {
	pricing := pricingFun(pricingParameters)

	fmt.Printf("================= %s ================\n", reflect.TypeOf(opt))
	spot := m.Money(100)
	fmt.Printf("S=%f V(T)=%f\n", float64(spot), opt.Payoff(spot))
	fmt.Printf("V(0)=%f\n", pricing(opt, spot, 0))
	fmt.Printf("Delta(0)=%f\n", option.Delta(pricing)(opt, spot, 0))
	fmt.Printf("Gamma(0)=%f\n", option.Gamma(pricing)(opt, spot, 0))
	fmt.Printf("Rho(0)=%f\n", option.Rho(pricingFun, pricingParameters)(opt, spot, 0))
	fmt.Printf("Theta(0)=%f\n", option.Theta(pricing)(opt, spot, 0.0))
}
