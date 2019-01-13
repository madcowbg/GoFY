package bond

import (
	m "../measures"
	"encoding/csv"
	"github.com/google/go-cmp/cmp"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBootstrap(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")

	quotes := demoQuotes(t.Error)
	ttms, yields := zeroCouponBondsFromQuotes(ed, m.Time(0), quotes)

	fwd := BootstrapForwardRates(yields, ttms)
	expectedFwd := []m.Rate{
		0.02346999883651903, 0.022810003757470675, 0.02465799808501816, 0.023889999389650945, 0.023638000488285734,
		0.025120000839212616, 0.023499996662143332, 0.020080015659336724, 0.026779997348779297, 0.0196599960327167,
		0.025068002223963708, 0.025210008621216123, 0.024131993770601837, 0.022850029468532897, 0.02294000196456491,
		0.02774998426436658, 0.024305716582712206, 0.022985699517382125, 0.024068588529318022, 0.025274276733399444,
		0.024851433890203313, 0.026495709078655748, 0.02208141292844351, 0.026380027021680456, 0.02851426226751792,
		0.02419288226536361, 0.023398564202444343, 0.028242868014746035, 0.022782844815931923, 0.025444274629869784,
		0.025295718738008243, 0.023208561284203472, 0.03048623859882728, 0.03148668805758071, 0.026681439535961846,
		0.026287151064191588, 0.023711418764930856, 0.02577999694006839, 0.029670005525862027, 0.027112858636038743,
		0.027499992506844165}

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
	expectedDfs := []float64{0.9902103054069268, 0.9997392562347156, 0.9986703292516028, 0.9974055381135003}
	if !cmp.Equal(dfs, expectedDfs, absCmp(1e-8)) {
		t.Errorf("bootstrapped discount factors wrong:\n got %v\n expected %v\n", dfs, expectedDfs)
	}

	bootsrappedYields := make([]float64, len(yields))
	expectedYields := make([]float64, len(yields))
	for i := range ttms {
		bootsrappedYields[i] = float64(fwd.Yield(ttms[i]))
		expectedYields[i] = float64(yields[i])
	}
	if !cmp.Equal(bootsrappedYields, expectedYields, absCmp(1e-8)) {
		t.Errorf("bootstrapped yields wrong:\n got %v\n expected %v\n%s\n", bootsrappedYields, expectedYields, cmp.Diff(bootsrappedYields, yields, absCmp(1e-8)))
	}
}

func zeroCouponBondsFromQuotes(ed time.Time, t0 m.Time, quotes []WSJUSNoteQuote) ([]m.Time, []m.Rate) {
	ttms := make([]m.Time, len(quotes))
	for i, quote := range quotes {
		ttms[i] = t0 + m.Time(daysBetween(ed, quote.date)/360.0)
	}
	yields := make([]m.Rate, len(quotes))
	for i := range quotes {
		yields[i] = quotes[i].yield
	}
	return ttms, yields
}

func demoQuotes(errFun func(args ...interface{})) []WSJUSNoteQuote {
	dateLayout := "1/2/2006"

	// headers := `Maturity,Bid,Asked,Chg,Asked yield`
	data := `1/15/2019,2.325,2.315,0.013,2.347
1/17/2019,2.303,2.293,0.005,2.325
1/22/2019,2.365,2.355,0.008,2.389
1/24/2019,2.365,2.355,0.005,2.389
1/29/2019,2.358,2.348,-0.002,2.382
1/31/2019,2.37,2.36,0.002,2.395
2/5/2019,2.36,2.35,0.002,2.386
2/7/2019,2.333,2.323,0.005,2.358
2/12/2019,2.38,2.37,0.012,2.408
2/14/2019,2.355,2.345,0.002,2.382
2/19/2019,2.37,2.36,0.002,2.398
2/21/2019,2.375,2.365,0.005,2.404
2/26/2019,2.375,2.365,0.005,2.405
2/28/2019,2.37,2.36,0.01,2.4
3/5/2019,2.36,2.35,-0.008,2.39
3/7/2019,2.373,2.363,0.015,2.404
3/14/2019,2.375,2.365,0.013,2.407
3/21/2019,2.363,2.353,0.002,2.396
3/28/2019,2.363,2.353,0.01,2.397
4/4/2019,2.373,2.363,0.01,2.408
4/11/2019,2.378,2.368,-0.007,2.414
4/18/2019,2.393,2.383,-0.005,2.431
4/25/2019,2.378,2.368,0.002,2.416
5/2/2019,2.39,2.38,-0.002,2.43
5/9/2019,2.413,2.403,-0.002,2.455
5/16/2019,2.41,2.4,unch.,2.453
5/23/2019,2.403,2.393,unch.,2.447
5/30/2019,2.42,2.41,0.002,2.466
6/6/2019,2.41,2.4,-0.01,2.457
6/13/2019,2.413,2.403,-0.002,2.461
6/20/2019,2.415,2.405,-0.005,2.464
6/27/2019,2.408,2.398,unch.,2.458
7/5/2019,2.433,2.423,unch.,2.485
7/11/2019,2.453,2.443,0.002,2.507
7/18/2019,2.458,2.448,0.007,2.513
8/15/2019,2.468,2.458,0.007,2.528
9/12/2019,2.445,2.435,unch.,2.51
10/10/2019,2.448,2.438,0.012,2.517
11/7/2019,2.483,2.473,-0.01,2.559
12/5/2019,2.49,2.48,-0.012,2.572
1/2/2020,2.498,2.488,-0.015,2.586`

	reader := csv.NewReader(strings.NewReader(data))
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

		yield100, err := strconv.ParseFloat(line[4], 32)
		if err != nil {
			errFun(err)
		}
		quotes = append(quotes, WSJUSNoteQuote{date: v, yield: m.Rate(yield100 / 100)})
	}
	return quotes
}

type WSJUSNoteQuote struct {
	date  time.Time
	yield m.Rate
}

func daysBetween(ed, v time.Time) float64 {
	return float64(v.Sub(ed) / (24 * time.Hour))
}
