package monotone_convex

import (
	m "../../measures"
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
)

func absCmp(tol float64) cmp.Option {
	return cmp.Comparer(func(x, y float64) bool {
		return (math.IsNaN(x) && math.IsNaN(y)) || math.Abs(x-y) < tol
	})
}

func TestMonotoneConvexGeneralCaseResults(t *testing.T) {
	Terms := []m.Time{1, 2, 3, 4, 5}
	Values := []m.Rate{0.03, 0.04, 0.047, 0.06, 0.06}

	spotInterpolator := SpotRateInterpolator(0)(Terms, Values)
	forwardInterpolator := ForwardRateInterpolator(0)(Terms, Values)

	tenors := []m.Time{0.1, 0.2, 0.3, 0.4, 0.5, 0.99, 1, 1.01, 2, 2.01, 3, 4, 4.9, 5, 5.1, 10}
	interpolated := interpolateOnArray(spotInterpolator, tenors)
	forward := interpolateOnArray(forwardInterpolator, tenors)

	expectedInterpolated := []m.Rate{
		0.02505, 0.025199999999999997, 0.02545, 0.025799999999999997, 0.02625, 0.0299005, 0.03, 0.03010044108910891,
		0.04, 0.040077114427860695, 0.047, 0.06, 0.060196989795918365, 0.06, 0.05980882352941176, 0.055125}

	if !cmp.Equal(interpolated, expectedInterpolated, absCmp(1e-14)) {
		t.Errorf("wrong spot rates:\n got %v\n expected %v\n", interpolated, expectedInterpolated)
	}

	expectedForward := []m.Rate{0.02515, 0.025599999999999998, 0.02635, 0.0274, 0.028749999999999998, 0.0397015,
		0.04, 0.04028865, 0.05550000000000001, 0.05550000000000001, 0.07999999999999999, 0.07949999999999999,
		0.050542500000000004, 0.05025, 0.05025, 0.05025}
	if !cmp.Equal(forward, expectedForward, absCmp(1e-14)) {
		t.Errorf("wrong forward rates:\n got %v\n expected %v\n", forward, expectedForward)
	}
}

func TestMonotoneConvexAmelioratedGeneralCaseResults(t *testing.T) {
	Terms := []m.Time{1, 2, 3, 4, 5}
	Values := []m.Rate{0.03, 0.04, 0.047, 0.06, 0.06}

	lambda := 1.0
	spotInterpolator := SpotRateInterpolator(lambda)(Terms, Values)
	forwardInterpolator := ForwardRateInterpolator(lambda)(Terms, Values)

	tenors := []m.Time{0.1, 0.2, 0.3, 0.4, 0.5, 0.99, 1, 1.01, 2, 2.01, 3, 4, 4.9, 5, 5.1, 10}
	interpolated := interpolateOnArray(spotInterpolator, tenors)
	forward := interpolateOnArray(forwardInterpolator, tenors)

	expectedInterpolated := []m.Rate{
		0.02505, 0.025199999999999997, 0.02545, 0.025799999999999997, 0.02625, 0.0299005, 0.03, 0.0301004900990099,
		0.04, 0.04007462686567164, 0.047, 0.06, 0.060196989795918365, 0.06, 0.05980882352941176, 0.055125}

	if !cmp.Equal(interpolated, expectedInterpolated, absCmp(1e-14)) {
		t.Errorf("wrong spot rates:\n got %v\n expected %v\n", interpolated, expectedInterpolated)
	}

	expectedForward := []m.Rate{
		0.02515, 0.025599999999999998, 0.02635, 0.0274, 0.028749999999999998, 0.0397015, 0.04, 0.0402985,
		0.05500000000000001, 0.05500000000000001, 0.07999999999999999, 0.07949999999999999, 0.050542500000000004,
		0.05025, 0.05025, 0.05025}
	if !cmp.Equal(forward, expectedForward, absCmp(1e-14)) {
		t.Errorf("wrong forward rates:\n got %v\n expected %v\n", forward, expectedForward)
	}
}

func interpolateOnArray(interpolator func(m.Time) m.Rate, tenors []m.Time) []m.Rate {
	result := make([]m.Rate, len(tenors))
	for i, t := range tenors {
		result[i] = interpolator(t)
	}
	return result
}

func TestMonotoneConvexWithNillsPanic(t *testing.T) {
	inp := mcInput{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	estimateInitialFI(inp)
}
