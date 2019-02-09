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

	expectedRates := []m.Rate{0.011348562987947706, 0.025948563369417896, 0.023868563922549185, 0.02461856201520038, 0.025249274700986413, 0.025952077839430546, 0.025646467762861905, 0.02498989752616197, 0.024602996143495006, 0.025254609222769057, 0.027331484259186935, 0.02960183809808654, 0.031335017830588556}

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
