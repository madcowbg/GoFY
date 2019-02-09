package bond

import (
	m "../measures"
	"encoding/csv"
	"github.com/google/go-cmp/cmp"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNaiveBootstrapBills(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")

	quotes := demoQuotes(t.Error)
	ttms, yields := yieldsFromQuotes(ed, m.Time(0), quotes)

	fwd := NaiveBootstrapFromZCYields(yields, ttms)
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

	dfs := []m.Money{
		DFByConstantRateInterpolation(fwd)(0.4),
		DFByConstantRateInterpolation(fwd)(ttms[0]),
		DFByConstantRateInterpolation(fwd)(ttms[5]),
		DFByConstantRateInterpolation(fwd)(ttms[10]),
	}
	expectedDfs := []m.Money{0.9902201368982367, 0.9997428276077888, 0.9986885319701507, 0.9974410331746054}
	if !cmp.Equal(dfs, expectedDfs, absCmp(1e-8)) {
		t.Errorf("bootstrapped discount factors wrong:\n got %v\n expected %v\n", dfs, expectedDfs)
	}

	bootsrappedYields := make([]float64, len(yields))
	expectedYields := make([]float64, len(yields))
	for i := range ttms {
		bootsrappedYields[i] = float64(AsRate(DFByConstantRateInterpolation(fwd))(ttms[i]))
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
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			errFun(err)
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

func TestNaiveBootstrapNotes(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")
	yearStart, _ := time.Parse(dateLayout, "12/31/2018")
	t0 := m.Time(daysBetween(yearStart, ed) / 365.0)

	quotes := demoNotesQuotes(t.Error)
	bonds, yields := bondsFromQuotes(yearStart, quotes)

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

	yieldByPrice := make([]float64, len(bonds))
	quotedYield := make([]float64, len(bonds))
	for i := range bonds {
		yieldByPrice[i] = float64(bonds[i].YieldToMaturity(t0, m.Money(dirtyPrices[i])))
		quotedYield[i] = float64(quotes[i].yield)
	}
	if !cmp.Equal(yieldByPrice, quotedYield, absCmp(0.006)) {
		t.Errorf(
			"bootstrapped yields wrong:\n by yields %v\n quoted dirty %v\n%s\n",
			yieldByPrice, quotedYield, cmp.Diff(yieldByPrice, quotedYield, absCmp(0.006)))
	}

	curve := NaiveBootstrapFromFixedCoupon(yields, bonds, t0)
	zeroCouponCurve := SpotCurveByConstantRateInterpolation(curve)

	zeroRatesAsFloat := make([]float64, len(zeroCouponCurve.Rates))
	for i, rate := range zeroCouponCurve.Rates {
		zeroRatesAsFloat[i] = float64(rate)
	}

	expectedRates := []float64{
		0.011220, 0.021730, 0.025820, 0.023800, 0.023740, 0.024500, 0.024490, 0.024830, 0.025280, 0.024420, 0.024830,
		0.024470, 0.025124, 0.025235, 0.025983, 0.026015, 0.026013, 0.025884, 0.026083, 0.026214, 0.026107, 0.026117,
		0.025974, 0.025997, 0.025827, 0.025866, 0.025981, 0.025973, 0.025892, 0.025923, 0.025961, 0.025591, 0.025924,
		0.026013, 0.025824, 0.025835, 0.025835, 0.025981, 0.025800, 0.025599, 0.025699, 0.025749, 0.025686, 0.025710,
		0.025474, 0.025465, 0.025446, 0.025520, 0.025426, 0.025364, 0.025424, 0.025221, 0.025261, 0.025351, 0.025197,
		0.025538, 0.025232, 0.025126, 0.025059, 0.025243, 0.025099, 0.025359, 0.025190, 0.025036, 0.025078, 0.025113,
		0.025078, 0.025289, 0.024968, 0.025087, 0.025009, 0.025069, 0.025061, 0.024983, 0.025007, 0.025009, 0.025080,
		0.025099, 0.025052, 0.025011, 0.025042, 0.025238, 0.025092, 0.025083, 0.025112, 0.025261, 0.025094, 0.025168,
		0.025134, 0.025401, 0.025110, 0.025103, 0.025152, 0.025224, 0.025082, 0.025147, 0.025209, 0.025434, 0.025069,
		0.025219, 0.025174, 0.025212, 0.025093, 0.025189, 0.025304, 0.025345, 0.025377, 0.025397, 0.025415, 0.025432,
		0.025453, 0.025483, 0.025507, 0.025397, 0.025557, 0.025641, 0.025630, 0.025665, 0.025698, 0.025698, 0.025902,
		0.025807, 0.025819, 0.025880, 0.025867, 0.025854, 0.025930, 0.026018, 0.025927, 0.025886, 0.025895, 0.026096,
		0.025897, 0.025947, 0.026072, 0.026187, 0.026036, 0.026348, 0.026370, 0.026548, 0.026773, 0.027001, 0.026935,
		0.026992, 0.027310, 0.027405, 0.027460, 0.027487, 0.027681, 0.027741, 0.028230, 0.028580, 0.028608, 0.028980,
		0.029077, 0.029475, 0.029817, 0.029994, 0.030162, 0.030233, 0.030357, 0.030348, 0.030559, 0.030580, 0.030688,
		0.030689, 0.030710, 0.030735, 0.030864, 0.030934, 0.030980, 0.031117, 0.031088, 0.031212, 0.031216, 0.031237,
		0.031135, 0.031073, 0.031043, 0.031101, 0.031057, 0.031123, 0.031058, 0.031194, 0.031159, 0.031120, 0.031150,
		0.031082, 0.031171, 0.031220, 0.031229, 0.031280, 0.031292, 0.031216, 0.031209}

	log.Printf("exp %v\n", expectedRates)
	if !cmp.Equal(zeroRatesAsFloat, expectedRates, absCmp(0.00001)) {
		t.Errorf(
			"bootstrapped yields wrong:\n by yields %v\n expected %v\n%s\n",
			zeroRatesAsFloat, expectedRates, cmp.Diff(zeroRatesAsFloat, expectedRates, absCmp(0.00001)))
	}
}

func bondsFromQuotes(yearStart time.Time, quotes []WSJUSNoteQuote) ([]*FixedCouponBond, []m.Rate) {
	bonds := make([]*FixedCouponBond, len(quotes))
	yields := make([]m.Rate, len(quotes))
	for i := range bonds {
		ttm := m.Time(daysBetween(yearStart, quotes[i].maturity) / 365.0)
		bonds[i] = &FixedCouponBond{
			Expirable: Expirable{ttm},
			IssueTime: ttm - m.Time(math.Ceil(float64(ttm))),
			Coupon: FixedCouponTerm{
				Frequency: 2.0, // semiannual
				PerAnnum:  quotes[i].coupon,
			}}
		yields[i] = quotes[i].yield
	}
	return bonds, yields
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
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			errFun(err)
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
