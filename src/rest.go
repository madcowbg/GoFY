package main

import (
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
	log.Fatal(http.ListenAndServe(":10001", nil))
}

func main() {
	handleRequests()

	//book := &generated.AddressBook{}
	//proto.Unmarshal(in, book)
}
