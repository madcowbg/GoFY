package main

import (
	"./option"
	"fmt"
	"reflect"
)

func main() {
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
	spot := option.Money(100)
	fmt.Printf("S=%f V(T)=%f\n", float64(spot), opt.Payoff(spot))
	fmt.Printf("V(0)=%f\n", pricing(opt, spot, 0))
	fmt.Printf("Delta(0)=%f\n", option.Delta(pricing)(opt, spot, 0))
	fmt.Printf("Gamma(0)=%f\n", option.Gamma(pricing)(opt, spot, 0))
	fmt.Printf("Rho(0)=%f\n", option.Rho(pricingFun, pricingParameters)(opt, spot, 0))
	fmt.Printf("Theta(0)=%f\n", option.Theta(pricing)(opt, spot, 0.0))
}
