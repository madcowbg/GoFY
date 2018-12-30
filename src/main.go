package main

import (
	"./option"
	"fmt"
)

func main() {
	// dec := json.NewDecoder(os.Stdin)
	// enc := json.NewEncoder(os.Stdout)
	opt := option.Call(0.2, 0.02, 100, 1)
	for i := 1; i < 20; i++ {
		spot := option.Money(i * 10.0)
		fmt.Printf("S=%f V(T)=%f\n", float64(spot), opt.Payoff(spot))
		fmt.Printf("V(0)=%f\n", opt.Price(spot, 0))
		fmt.Printf("Delta(0)=%f\n", option.Delta(opt, spot, 0))
		fmt.Printf("Gamma(0)=%f\n", option.Gamma(opt, spot, 0))
		fmt.Printf("Rho(0)=%f\n", opt.Rho(spot, 0))
		fmt.Printf("Theta(0)=%f\n", option.Theta(opt, spot, 0.0))
	}
}
