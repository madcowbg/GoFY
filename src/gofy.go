package main

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
)
import (
	b "./bond"
	mc "./bond/monotone_convex"
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

type ImpliedVolAnalytics struct {
	Analytics  SimpleAnalytics
	ImpliedVol m.Return
}

//export implyVol
func implyVol(instrumentType string, termsAndConditions string, StateOfWorldJSON string, optionPrice float64) *C.char {
	var world StateOfWorld
	json.Unmarshal([]byte(StateOfWorldJSON), &world)

	switch InstrumentType(instrumentType) {
	case Option:
		{
			var optionTnC OptionTnC

			err := json.Unmarshal([]byte(termsAndConditions), &optionTnC)
			if err != nil {
				return C.CString(fmt.Sprintf("bad T&C: %s", termsAndConditions))
			}
			opt := readOptionType(optionTnC)

			analytics, error := calculateOptionImplyVol(opt, world.Parameters, world.Spot, world.Time, m.Money(optionPrice))
			if error != nil {
				return C.CString(fmt.Sprint(error))
			}

			res, error := json.Marshal(*analytics)
			if error != nil {
				return C.CString(fmt.Sprint(error))
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

func calculateOptionImplyVol(option o.Option, parameters o.PricingParameters, spot m.Money, t m.Time, optionPrice m.Money) (*ImpliedVolAnalytics, error) {
	impliedVol, err := o.ImplyVol(o.BinomialPricing, parameters.R)(option, spot, t)(optionPrice)
	if err != nil {
		return nil, err
	}

	newParameters := o.PricingParameters{R: parameters.R, Sigma: m.Return(impliedVol)}

	return &ImpliedVolAnalytics{
		Analytics:  calculateOptionAnalytics(option, newParameters, spot, t),
		ImpliedVol: m.Return(impliedVol),
	}, nil
}

type CurveBootstrapData struct {
	BondDefinitions []CouponBondDef
	Yields          []m.Rate
}

type CouponBondDef struct {
	IssueTime       m.Time
	Maturity        m.Time
	CouponFrequency float64
	Coupon          m.Money
}

type TenorDefs struct {
	Tenors []m.Time
}

type BootstrapMethod string

const (
	Naive          BootstrapMethod = "Naive"
	MonotoneConvex BootstrapMethod = "MonotoneConvex"
)

//export bootstrapCurve
func bootstrapCurve(method BootstrapMethod, t0 float64, BootstrapData string, TenorData string) *C.char {
	var bootstrapData CurveBootstrapData
	json.Unmarshal([]byte(BootstrapData), &bootstrapData)
	if len(bootstrapData.BondDefinitions) != len(bootstrapData.Yields) {
		log.Fatalf(
			"bad format - different length of bonds and qoutes: %d != %d\n",
			len(bootstrapData.BondDefinitions),
			len(bootstrapData.Yields))
	}

	bonds := make([]*b.FixedCouponBond, len(bootstrapData.BondDefinitions))
	for i, bondDef := range bootstrapData.BondDefinitions {
		bonds[i] = &b.FixedCouponBond{
			Expirable: b.Expirable{Maturity: bondDef.Maturity},
			IssueTime: bondDef.IssueTime,
			Coupon: b.FixedCouponTerm{
				Frequency: bondDef.CouponFrequency,
				PerAnnum:  bondDef.Coupon,
			},
		}
	}

	yields := bootstrapData.Yields

	var spotCurve *b.FixedSpotCurve
	switch method {
	case MonotoneConvex:
		var tenorsDefs TenorDefs
		json.Unmarshal([]byte(TenorData), &tenorsDefs)
		if len(tenorsDefs.Tenors) == 0 {
			log.Fatalf("Invalid tenor request, needed for monotone convex! : %v\n", tenorsDefs.Tenors)
		}

		tenors := tenorsDefs.Tenors

		spotCurve = b.OLSBootstrapFromFixedCoupon(mc.SpotRateInterpolator, yields, bonds, m.Time(t0), tenors)
		break
	case Naive:
		forwardCurve := b.NaiveBootstrapFromFixedCoupon(yields, bonds, m.Time(t0))
		spotCurve = b.SpotCurveByConstantRateInterpolation(forwardCurve)
		break
	default:
		return C.CString(fmt.Sprintf("Invalid bootstrap method: %s\n", method))
	}

	res, error := json.Marshal(spotCurve)
	if error != nil {
		return C.CString(fmt.Sprint(error))
	}
	return C.CString(string(res))
}

func main() {
	// just needed otherwise won't compile dll
}
