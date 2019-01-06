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
	pricing := option.BinomialPricing(parameters)

	gridPricing := option.GridPricing(parameters)

	printGreekDiagATM(
		parameters, pricing,
		&option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		parameters, gridPricing,
		&option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		parameters, pricing,
		&option.EuropeanPutOption{option.EuropeanOption{option.VanillaOption{100, 1}}})

	printGreekDiagATM(
		parameters, pricing,
		&option.AmericanPutOption{option.AmericanOption{option.VanillaOption{100, 1}}})
}

func printGreekDiagATM(pricingParameters option.PricingParameters, pricing option.Pricing, opt option.Option) {
	fmt.Printf("================= %s ================\n", reflect.TypeOf(opt))
	spot := option.Money(100)
	fmt.Printf("S=%f V(T)=%f\n", float64(spot), opt.Payoff(spot))
	fmt.Printf("V(0)=%f\n", pricing(opt, spot, 0))
	fmt.Printf("Delta(0)=%f\n", option.Delta(pricing)(opt, spot, 0))
	fmt.Printf("Gamma(0)=%f\n", option.Gamma(pricing)(opt, spot, 0))
	fmt.Printf("Rho(0)=%f\n", option.Rho(option.BinomialPricing, pricingParameters)(opt, spot, 0))
	fmt.Printf("Theta(0)=%f\n", option.Theta(pricing)(opt, spot, 0.0))
}
