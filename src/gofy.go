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
	Price     m.Money
	Delta     float64
	Gamma     float64
	Theta     float64
	Rho       float64
	Intrinsic m.Money
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
	pricing := o.BinomialPricing(parameters)
	return SimpleAnalytics{
		Price:     pricing(opt, spot, t),
		Delta:     o.Delta(pricing)(opt, spot, t),
		Gamma:     o.Gamma(pricing)(opt, spot, t),
		Theta:     o.Theta(pricing)(opt, spot, t),
		Rho:       o.Rho(o.BinomialPricing, parameters)(opt, spot, t),
		Intrinsic: opt.Payoff(spot),
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

type BootstrapOutput struct {
	SpotCurve                b.FixedSpotCurve
	InterpolatedSpotCurve    b.FixedSpotCurve
	InterpolatedForwardCurve b.FixedForwardRateCurve
}

//export bootstrapCurve
func bootstrapCurve(method BootstrapMethod, lambda float64, t0 float64, BootstrapData string, TenorData string, OutputTenors string) *C.char {
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

	var output BootstrapOutput
	switch method {
	case MonotoneConvex:
		var tenorsDefs TenorDefs
		json.Unmarshal([]byte(TenorData), &tenorsDefs)
		if len(tenorsDefs.Tenors) == 0 {
			log.Fatalf("Invalid tenor request, needed for monotone convex! : %v\n", tenorsDefs.Tenors)
		}

		tenors := tenorsDefs.Tenors

		spotCurve := b.OLSBootstrapFromFixedCoupon(mc.SpotRateInterpolator(lambda), yields, bonds, m.Time(t0), tenors)

		var outputTenorDefs TenorDefs
		json.Unmarshal([]byte(OutputTenors), &outputTenorDefs)

		var interpolatedSpot b.FixedSpotCurve
		var interpolatedForward b.FixedForwardRateCurve
		if len(outputTenorDefs.Tenors) > 0 {
			interpolatedSpot = b.FixedSpotCurve{
				Tenors: outputTenorDefs.Tenors,
				Rates: b.InterpolateOnArray(
					mc.SpotRateInterpolator(lambda)(spotCurve.Tenors, spotCurve.Rates),
					outputTenorDefs.Tenors)}
			interpolatedForward = b.FixedForwardRateCurve{
				Tenors: outputTenorDefs.Tenors,
				Rates: b.InterpolateOnArray(
					mc.ForwardRateInterpolator(lambda)(spotCurve.Tenors, spotCurve.Rates),
					outputTenorDefs.Tenors)}
		}

		output = BootstrapOutput{
			SpotCurve:                *spotCurve,
			InterpolatedSpotCurve:    interpolatedSpot,
			InterpolatedForwardCurve: interpolatedForward,
		}
		break
	case Naive:
		forwardCurve := b.NaiveBootstrapFromFixedCoupon(yields, bonds, m.Time(t0))
		spotCurve := b.SpotCurveByConstantRateInterpolation(forwardCurve)

		output = BootstrapOutput{SpotCurve: *spotCurve}
		break
	default:
		return C.CString(fmt.Sprintf("Invalid bootstrap method: %s\n", method))
	}

	res, error := json.Marshal(output)
	if error != nil {
		return C.CString(fmt.Sprint(error))
	}
	return C.CString(string(res))
}

/*func main() {
	// just needed otherwise won't compile dll
}*/
