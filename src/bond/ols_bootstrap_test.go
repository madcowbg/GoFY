package bond

import (
	m "../measures"
	"./monotone_convex"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestMCOLSBootstrapNotes(t *testing.T) {
	dateLayout := "1/2/2006"
	ed, _ := time.Parse(dateLayout, "1/11/2019")
	yearStart, _ := time.Parse(dateLayout, "12/31/2018")
	t0 := m.Time(daysBetween(yearStart, ed) / 365.0)

	quotes := demoNotesQuotes(t.Error)
	quotes = quotes[0:]

	bonds, yields := bondsFromQuotes(yearStart, quotes)

	tenors := []m.Time{1.0 / 365, 1.0 / 12, 2.0 / 12, 3.0 / 12, 6.0 / 12, 1, 2, 3, 4, 5, 10, 20, 30}
	spotCurve := OLSBootstrapFromFixedCoupon(monotone_convex.SpotRateInterpolator, yields, bonds, t0, tenors)

	expectedRates := []m.Rate{0.011219999790191169, 0.025820000171661396, 0.023740000724792664, 0.024489998817443837, 0.02512071150322988, 0.025823514641674032, 0.02551790456510539, 0.02506720947283312, 0.02513419625934795, 0.025304471150183236, 0.027458742246813946, 0.029973274900330017, 0.031206454632832045}

	if !cmp.Equal(spotCurve.Tenors, tenors, absCmp(0.00001)) {
		t.Errorf(
			"bootstrapped tenors wrong:\n by yields %v\n expected %v\n%s\n",
			spotCurve.Tenors, tenors, cmp.Diff(spotCurve.Tenors, tenors, absCmp(0.00001)))
	}

	if !cmp.Equal(spotCurve.Rates, expectedRates, absCmp(0.00001)) {
		t.Errorf(
			"bootstrapped rates wrong:\n by yields %v\n expected %v\n%s\n",
			spotCurve.Rates, expectedRates, cmp.Diff(spotCurve.Rates, expectedRates, absCmp(0.00001)))
	}
}
