package main

import "C"
import (
	"./bond"
	"./bond/monotone_convex"
	"./measures"
	"./option"
	"./proto/generated"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"net/http/httputil"
)

type Calculator func([]byte) ([]byte, error)
type Handler func(w http.ResponseWriter, r *http.Request)

func endpointHandler(endpoint string, handler Calculator) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Welcome to API!: "+endpoint)
			fmt.Printf("call received without actual request...\n")
		case "POST":
			//dumpContent(r)

			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r.Body)

			result, err := handler(buffer.Bytes())

			if err != nil {
				fmt.Println("Error calculating result:", err)
				fmt.Println(w, "Error calculating result:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/protobuf")
			w.Write(result)
		}
	}
}

func handleCalculateOptionAnalytics(input []byte) ([]byte, error) {
	request := generated.RequestCalculateOptionAnalytics{}
	err := proto.Unmarshal(input, &request)
	if err != nil {
		fmt.Println("Error unmarshalling request:", err)
		return nil, err
	}

	opt := readOptionTnC(request.TermsAndConditions)
	analytics := calculateOptionAnalyticsResponse(opt, asPricingParameters(request.StateOfWorld.Parameters), measures.Money(*request.StateOfWorld.Spot), measures.Time(*request.StateOfWorld.Time))

	return proto.Marshal(analytics)
}

func calculateOptionAnalyticsResponse(opt option.Option, parameters option.PricingParameters, spot measures.Money, t measures.Time) *generated.ResponseCalculateOptionAnalytics {
	pricing := option.BinomialPricing(parameters)
	return &generated.ResponseCalculateOptionAnalytics{
		Price:     proto.Float32(float32(pricing(opt, spot, t))),
		Delta:     proto.Float32(float32(option.Delta(pricing)(opt, spot, t))),
		Gamma:     proto.Float32(float32(option.Gamma(pricing)(opt, spot, t))),
		Theta:     proto.Float32(float32(option.Theta(pricing)(opt, spot, t))),
		Rho:       proto.Float32(float32(option.Rho(option.BinomialPricing, parameters)(opt, spot, t))),
		Intrinsic: proto.Float32(float32(opt.Payoff(spot))),
	}
}

func asPricingParameters(parameters *generated.PricingParameters) option.PricingParameters {
	return option.PricingParameters{R: measures.Rate(*parameters.R), Sigma: measures.Return(*parameters.Sigma)}
}

func readOptionTnC(tnc *generated.OptionTermsAndConditions) option.Option {
	maturity := measures.Time(*tnc.T)
	strike := measures.Money(*tnc.S)

	switch *tnc.Parity {
	case generated.OptionParity_Call:
		switch *tnc.Type {
		case generated.OptionType_American:
			return &option.AmericanCallOption{option.AmericanOption{option.VanillaOption{strike, maturity}}}
		case generated.OptionType_European:
			return &option.EuropeanCallOption{option.EuropeanOption{option.VanillaOption{strike, maturity}}}
		default:
			log.Fatalf("wrong option type: %s", *tnc.Type)
		}
	case generated.OptionParity_Put:
		switch *tnc.Type {
		case generated.OptionType_American:
			return &option.AmericanPutOption{option.AmericanOption{option.VanillaOption{strike, maturity}}}
		case generated.OptionType_European:
			return &option.EuropeanPutOption{option.EuropeanOption{option.VanillaOption{strike, maturity}}}
		default:
			log.Fatalf("wrong option type: %s", *tnc.Type)
		}

	default:
		log.Fatalf("wrong option parity: %s", *tnc.Parity)
	}
	return nil
}

func handleBootstrapCurve(input []byte) ([]byte, error) {
	request := generated.RequestBootstrapCurve{}
	err := proto.Unmarshal(input, &request)
	if err != nil {
		return nil, err
	}

	bootstrapData := request.BootstrapData
	if len(bootstrapData.BondDefinitions) != len(bootstrapData.Yields) {
		return nil, fmt.Errorf(
			"bad format - different length of bonds and qoutes: %d != %d\n",
			len(bootstrapData.BondDefinitions),
			len(bootstrapData.Yields))
	}

	bonds := make([]*bond.FixedCouponBond, len(bootstrapData.BondDefinitions))
	for i, bondDef := range bootstrapData.BondDefinitions {
		bonds[i] = &bond.FixedCouponBond{
			Expirable: bond.Expirable{Maturity: measures.Time(*bondDef.Maturity)},
			IssueTime: measures.Time(*bondDef.IssueTime),
			Coupon: bond.FixedCouponTerm{
				Frequency: float64(*bondDef.CouponFrequency),
				PerAnnum:  measures.Money(*bondDef.Coupon),
			},
		}
	}

	yields := make([]measures.Rate, len(bootstrapData.Yields))
	for i, y := range bootstrapData.Yields {
		yields[i] = measures.Rate(y)
	}

	switch *request.Method {
	case generated.BootstrapMethod_MonotoneConvex:
		tenorsDefs := *request.TenorData
		if len(tenorsDefs.Tenors) == 0 {
			log.Fatalf("Invalid tenor request, needed for monotone convex! : %v\n", tenorsDefs.Tenors)
		}

		tenors := asTime(tenorsDefs.Tenors)

		spotCurve := bond.OLSBootstrapFromFixedCoupon(monotone_convex.SpotRateInterpolator(*request.Lambda), yields, bonds, measures.Time(*request.T0), tenors)

		outputTenorDefs := request.OutputTenors

		var interpolatedSpot bond.FixedSpotCurve
		var interpolatedForward bond.FixedForwardRateCurve
		if len(outputTenorDefs.Tenors) > 0 {
			interpolatedSpot = bond.FixedSpotCurve{
				Tenors: asTime(outputTenorDefs.Tenors),
				Rates: bond.InterpolateOnArray(
					monotone_convex.SpotRateInterpolator(*request.Lambda)(spotCurve.Tenors, spotCurve.Rates),
					asTime(outputTenorDefs.Tenors))}
			interpolatedForward = bond.FixedForwardRateCurve{
				Tenors: asTime(outputTenorDefs.Tenors),
				Rates: bond.InterpolateOnArray(
					monotone_convex.ForwardRateInterpolator(*request.Lambda)(spotCurve.Tenors, spotCurve.Rates),
					asTime(outputTenorDefs.Tenors))}
		}

		return proto.Marshal(&generated.ResponseBootstrapCurve{
			SpotCurve:                asProtoFSC(spotCurve),
			InterpolatedSpotCurve:    asProtoFSC(&interpolatedSpot),
			InterpolatedForwardCurve: asProtoFFRC(&interpolatedForward),
		})
	case generated.BootstrapMethod_Naive:
		forwardCurve := bond.NaiveBootstrapFromFixedCoupon(yields, bonds, measures.Time(*request.T0))
		spotCurve := bond.SpotCurveByConstantRateInterpolation(forwardCurve)

		return proto.Marshal(&generated.ResponseBootstrapCurve{
			SpotCurve: asProtoFSC(spotCurve)})
	default:
		return nil, fmt.Errorf("Invalid bootstrap method: %s\n", *request.Method)
	}
}

func asProtoFSC(curve *bond.FixedSpotCurve) *generated.Curve {
	return &generated.Curve{
		Rates:  asFloatRate(curve.Rates),
		Tenors: asFloatTime(curve.Tenors),
	}
}

func asFloatRate(a []measures.Rate) []float32 {
	res := make([]float32, len(a))
	for i, t := range a {
		res[i] = float32(t)
	}
	return res
}

func asFloatTime(a []measures.Time) []float32 {
	res := make([]float32, len(a))
	for i, t := range a {
		res[i] = float32(t)
	}
	return res
}

func asProtoFFRC(curve *bond.FixedForwardRateCurve) *generated.Curve {
	return &generated.Curve{
		Rates:  asFloatRate(curve.Rates),
		Tenors: asFloatTime(curve.Tenors),
	}
}

func asTime(a []float32) []measures.Time {
	res := make([]measures.Time, len(a))
	for i, t := range a {
		res[i] = measures.Time(t)
	}
	return res
}

func dumpContent(request *http.Request) {
	output, err := httputil.DumpRequest(request, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println("-----------------------")
	fmt.Println(string(output))
	fmt.Println("-----------------------")

}

func handleRequests() {
	http.HandleFunc("/calculateOptionAnalytics", endpointHandler("/calculateOptionAnalytics", handleCalculateOptionAnalytics))
	http.HandleFunc("/calculateBootstrapCurve", endpointHandler("/calculateBootstrapCurve", handleBootstrapCurve))
	log.Fatal(http.ListenAndServe(":10001", nil))
}

func main() {
	handleRequests()

	//book := &generated.AddressBook{}
	//proto.Unmarshal(in, book)
}
