package main

import "C"
import (
	"encoding/json"
	"log"
)
import (
	m "./measures"
	o "./option"
)

type InstrumentType string

const (
	Option InstrumentType = "Option"
)

type OptionParity string

const (
	Call OptionParity = "Call"
	Put  OptionParity = "Put"
)

type ExcerciseType string

const (
	European ExcerciseType = "European"
	American ExcerciseType = "American"
)

type OptionTnC struct {
	S      m.Money
	T      m.Time
	Parity OptionParity
	Type   ExcerciseType
}

func readOptionType(optionTnC OptionTnC) o.Option {
	switch optionTnC.Parity {
	case Call:
		{
			switch optionTnC.Type {
			case American:
				return &o.AmericanCallOption{o.AmericanOption{o.VanillaOption{optionTnC.S, optionTnC.T}}}
			case European:
				return &o.EuropeanCallOption{o.EuropeanOption{o.VanillaOption{optionTnC.S, optionTnC.T}}}
			default:
				log.Fatalf("wrong option type: %s", optionTnC.Type)
			}
		}
	case Put:
		{
			switch optionTnC.Type {
			case American:
				return &o.AmericanPutOption{o.AmericanOption{o.VanillaOption{optionTnC.S, optionTnC.T}}}
			case European:
				return &o.EuropeanPutOption{o.EuropeanOption{o.VanillaOption{optionTnC.S, optionTnC.T}}}
			default:
				log.Fatalf("wrong option type: %s", optionTnC.Type)
			}
		}
	default:
		log.Fatalf("wrong option parity: %s", optionTnC.Parity)
	}
	return nil
}

type SimpleAnalytics struct {
	Price m.Money
	Delta float64
	Gamma float64
	Theta float64
	Rho   float64
}

type StateOfWorld struct {
	Parameters o.PricingParameters
	Spot       m.Money
	Time       m.Time
}

//export calculateAnalytics
func calculateAnalytics(instrumentType string, termsAndConditions string, StateOfWorldJSON string) *C.char {
	var world StateOfWorld
	json.Unmarshal([]byte(StateOfWorldJSON), &world)

	switch InstrumentType(instrumentType) {
	case Option:
		{
			var optionTnC OptionTnC

			err := json.Unmarshal([]byte(termsAndConditions), &optionTnC)
			if err != nil {
				log.Fatalf("bad T&C: %s", termsAndConditions)
			}
			opt := readOptionType(optionTnC)
			analytics := calculateOptionAnalytics(opt, world.Parameters, world.Spot, world.Time)

			res, error := json.Marshal(&analytics)
			if error != nil {
				log.Fatal(error)
			}
			return C.CString(string(res))
		}
	default:
		{
			log.Fatalf("invalid instrument type: %s", instrumentType)
		}
	}
	return C.CString("BAD CALL!")
}

func calculateOptionAnalytics(opt o.Option, parameters o.PricingParameters, spot m.Money, t m.Time) SimpleAnalytics {
	pricing := o.GridPricing(parameters)
	return SimpleAnalytics{
		Price: pricing(opt, spot, t),
		Delta: o.Delta(pricing)(opt, spot, t),
		Gamma: o.Gamma(pricing)(opt, spot, t),
		Theta: o.Theta(pricing)(opt, spot, t),
		Rho:   o.Rho(o.GridPricing, parameters)(opt, spot, t),
	}
}

func main() {
	// just needed otherwise won't compile dll
}
