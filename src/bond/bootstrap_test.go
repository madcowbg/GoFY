package bond

import (
	m "../measures"
	"encoding/csv"
	"github.com/google/go-cmp/cmp"
	"io"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBootstrapBills(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")

	quotes := demoQuotes(t.Error)
	ttms, yields := yieldsFromQuotes(ed, m.Time(0), quotes)

	fwd := BootstrapForwardRates(yields, ttms)
	expectedFwd := []m.Rate{
		0.023469998836519075, 0.02281000375745703, 0.024657998085030178, 0.023889999389644142, 0.023638000488280967,
		0.0251200008392252, 0.0234999966621385, 0.020080015659347413, 0.026779997348784976, 0.019659996032726456,
		0.025068002223964946, 0.02521000862121547, 0.02413199377060504, 0.022850029468534237, 0.02294000196456842,
		0.027749984264386275, 0.02430571658270802, 0.022985699517381823, 0.024068588529312856, 0.02527427673340311,
		0.02485143389020457, 0.026495709078651848, 0.022081412928442193, 0.026380027021682875, 0.028514262267522822,
		0.024192882265362606, 0.02339856420245001, 0.028242868014741847, 0.022782844815934254, 0.025444274629866332,
		0.025295718738010398, 0.023208561284199646, 0.030486238598823857, 0.03148668805758675, 0.026681439535955982,
		0.026287151064192247, 0.023711418764933094, 0.02577999694006619, 0.029670005525861926, 0.0271128586360389,
		0.027499992506844508}

	if !cmp.Equal(fwd.Rates, expectedFwd, absCmp(1e-8)) {
		t.Errorf(
			"forward rates bootstrapped wrong:\n got %v\n expected %v\n%s\n",
			fwd.Rates, expectedFwd, cmp.Diff(fwd.Rates, expectedFwd, absCmp(1e-8)))
	}

	dfs := []float64{
		fwd.DiscountFactor(0.4),
		fwd.DiscountFactor(ttms[0]),
		fwd.DiscountFactor(ttms[5]),
		fwd.DiscountFactor(ttms[10]),
	}
	expectedDfs := []float64{0.9902201368982367, 0.9997428276077888, 0.9986885319701507, 0.9974410331746054}
	if !cmp.Equal(dfs, expectedDfs, absCmp(1e-8)) {
		t.Errorf("bootstrapped discount factors wrong:\n got %v\n expected %v\n", dfs, expectedDfs)
	}

	bootsrappedYields := make([]float64, len(yields))
	expectedYields := make([]float64, len(yields))
	for i := range ttms {
		bootsrappedYields[i] = float64(Yield(fwd.DiscountFactor, ttms[i]))
		expectedYields[i] = float64(yields[i])
	}
	if !cmp.Equal(bootsrappedYields, expectedYields, absCmp(1e-8)) {
		t.Errorf("bootstrapped yields wrong:\n got %v\n expected %v\n%s\n", bootsrappedYields, expectedYields, cmp.Diff(bootsrappedYields, yields, absCmp(1e-8)))
	}
}

func yieldsFromQuotes(ed time.Time, t0 m.Time, quotes []WSJUSBillQuote) ([]m.Time, []m.Rate) {
	ttms := make([]m.Time, len(quotes))
	for i, quote := range quotes {
		ttms[i] = t0 + m.Time(daysBetween(ed, quote.maturity)/365.0)
	}
	yields := make([]m.Rate, len(quotes))
	for i := range quotes {
		yields[i] = quotes[i].yield
	}
	return ttms, yields
}

func demoQuotes(errFun func(args ...interface{})) []WSJUSBillQuote {
	dateLayout := "1/2/2006"

	reader := csv.NewReader(strings.NewReader(billsData))
	var quotes []WSJUSBillQuote
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			errFun(error)
		}

		v, err := time.Parse(dateLayout, line[0])
		if err != nil {
			errFun(err)
		}

		yield100, err := strconv.ParseFloat(line[4], 32)
		if err != nil {
			errFun(err)
		}
		quotes = append(quotes, WSJUSBillQuote{maturity: v, yield: m.Rate(yield100 / 100)})
	}
	return quotes
}

type WSJUSBillQuote struct {
	maturity time.Time
	yield    m.Rate
}

func daysBetween(ed, v time.Time) float64 {
	return float64(v.Sub(ed) / (24 * time.Hour))
}

func TestBootstrapNotes(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")
	yearStart, _ := time.Parse(dateLayout, "12/31/2018")
	t0 := m.Time(daysBetween(yearStart, ed) / 365.0)

	// fmt.Println(t0)
	quotes := demoNotesQuotes(t.Error)
	bonds := bondsFromQuotes(yearStart, quotes)

	dirtyPrices := make([]float64, len(bonds))
	for i := range bonds {
		dirtyPrices[i] = float64(quotes[i].price + bonds[i].AccruedInterest(t0))
	}

	priceByYield := make([]float64, len(bonds))
	for i := range bonds {
		priceByYield[i] = float64(bonds[i].Price(t0, quotes[i].yield))
	}

	if !cmp.Equal(priceByYield, dirtyPrices, absCmp(1e-2)) {
		t.Errorf(
			"bootstrapped yields wrong:\n by yields %v\n quoted dirty %v\n%s\n",
			priceByYield, dirtyPrices, cmp.Diff(priceByYield, dirtyPrices, absCmp(1e-2)))
	}

	//for i := range bonds {
	//	fmt.Printf(
	//		"%f\t%f\t%f\t%f\t%f\t%f\t%v\n",
	//		quotes[i].price,
	//		dirtyPrices[i],
	//		bonds[i].Price(t0, quotes[i].yield),
	//		quotes[i].yield,
	//		bonds[i].YieldToMaturity(t0, m.Money(dirtyPrices[i])),
	//		bonds[i].YieldToMaturity(t0, quotes[i].price),
	//		bonds[i].RemainingCashflows(t0),
	//	)
	//}
}

func bondsFromQuotes(yearStart time.Time, quotes []WSJUSNoteQuote) []*FixedCouponBond {
	result := make([]*FixedCouponBond, len(quotes))
	for i := range result {
		ttm := m.Time(daysBetween(yearStart, quotes[i].maturity) / 365.0)
		result[i] = &FixedCouponBond{
			Expirable: Expirable{ttm},
			IssueTime: ttm - m.Time(math.Ceil(float64(ttm))),
			Coupon: FixedCouponTerm{
				Frequency: 2.0, // semiannual
				PerAnnum:  quotes[i].coupon,
			}}
	}
	return result
}

type WSJUSNoteQuote struct {
	maturity time.Time
	coupon   m.Money
	yield    m.Rate
	price    m.Money
}

func demoNotesQuotes(errFun func(args ...interface{})) []WSJUSNoteQuote {
	dateLayout := "1/2/2006"

	reader := csv.NewReader(strings.NewReader(notesData))
	var quotes []WSJUSNoteQuote
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			errFun(error)
		}

		v, err := time.Parse(dateLayout, line[0])
		if err != nil {
			errFun(err)
		}

		yield100, err := strconv.ParseFloat(line[5], 32)
		if err != nil {
			errFun(err)
		}

		coupon100, err := strconv.ParseFloat(line[1], 32)
		if err != nil {
			errFun(err)
		}

		asked100, err := strconv.ParseFloat(line[3], 32)
		if err != nil {
			errFun(err)
		}

		quotes = append(quotes, WSJUSNoteQuote{
			maturity: v,
			coupon:   m.Money(coupon100 / 100.0),
			yield:    m.Rate(yield100 / 100),
			price:    m.Money(asked100 / 100)})
	}

	return quotes
}
