package monotone_convex

/*
	implementation based on Hagan, Patrick S., and Graeme West. "Methods for constructing a yield curve." Wilmott Magazine, May (2008): 70-81.
*/

import (
	"log"
	"math"
	"sort"
)

type mcInput struct {
	terms  []float64
	values []float64
}

func (inp *mcInput) N() int {
	return len(inp.terms)
}

func (inp *mcInput) Terms(i int) float64 {
	if i <= 0 {
		return 0
	} else {
		return inp.terms[i-1]
	}
}

func (inp *mcInput) Values(i int) float64 {
	if i <= 0 {
		return inp.values[0]
	} else {
		return inp.values[i-1]
	}
}

type initialFI struct {
	mcInput

	fD                 []float64
	interpolantAtNodeD []float64

	f []float64
}

func SpotRateInterpolator(terms []float64, rates []float64) func(Term float64) float64 {
	if len(terms) != len(rates) {
		log.Fatalf("must have corresponding length of terms and rates! %d != %d\n", len(terms), len(rates))
	}

	e := estimateInitialFI(mcInput{terms, rates})
	return func(Term float64) float64 { return spotRate(Term, e) }
}

func ForwardRateInterpolator(terms []float64, rates []float64) func(Term float64) float64 {
	if len(terms) != len(rates) {
		log.Fatalf("must have corresponding length of terms and rates! %d != %d\n", len(terms), len(rates))
	}
	e := estimateInitialFI(mcInput{terms, rates})
	return func(Term float64) float64 { return forwardRate(Term, e) }
}

func spotRate(Term float64, e initialFI) float64 {
	// 'numbering refers to Wilmott paper, functions are integrated.
	if Term <= 0 {
		return e.f[0]
	}
	if Term > e.Terms(e.N()) {
		return spotRate(e.Terms(e.N()), e)*e.Terms(e.N())/Term + forwardRate(e.Terms(e.N()), e)*(1-e.Terms(e.N())/Term)
	}

	i, x, g0, g1 := initialInterpolators(e, Term)
	G := adjustedGIntegrated(x, g0, g1)

	//'(12)
	return 1 / Term * (e.Terms(i)*e.interpolantAtNodeD[i] + (Term-e.Terms(i))*e.fD[i+1] + (e.Terms(i+1)-e.Terms(i))*G)
}

func initialInterpolators(e initialFI, Term float64) (i int, x float64, g0 float64, g1 float64) {
	i = e.lastTermIndexBefore(Term)
	// 'the x in (25)
	x = (Term - e.Terms(i)) / (e.Terms(i+1) - e.Terms(i))
	g0 = e.f[i] - e.fD[i+1]
	g1 = e.f[i+1] - e.fD[i+1]
	return
}

func adjustedGIntegrated(x float64, g0 float64, g1 float64) float64 {
	if x == 0 || x == 1 {
		return 0
	} else if (g0 < 0 && -0.5*g0 <= g1 && g1 <= -2*g0) || (g0 > 0 && -0.5*g0 >= g1 && g1 >= -2*g0) {
		// 'zone (i)
		return g0*(x-2*math.Pow(x, 2)+math.Pow(x, 3)) + g1*(-math.Pow(x, 2)+math.Pow(x, 3))
	} else if (g0 < 0 && g1 > -2*g0) || (g0 > 0 && g1 < -2*g0) {
		//'zone (ii)
		// '(29)
		eta := (g1 + 2*g0) / (g1 - g0)
		// '(28)
		if x <= eta {
			return g0 * x
		} else {
			return g0*x + (g1-g0)*math.Pow(x-eta, 3)/math.Pow(1-eta, 2)/3
		}
	} else if (g0 > 0 && 0 > g1 && g1 > -0.5*g0) || (g0 < 0 && 0 < g1 && g1 < -0.5*g0) {
		// 'zone (iii)
		// '(31)
		eta := 3 * g1 / (g1 - g0)
		//'(30)
		if x < eta {
			return g1*x - 1/3*(g0-g1)*(math.Pow(eta-x, 3)/math.Pow(eta, 2)-eta)
		} else {
			return (2/3*g1+1/3*g0)*eta + g1*(x-eta)
		}
	} else if g0 == 0 || g1 == 0 {
		return 0
	} else {
		// 'zone (iv)
		// '(33)
		eta := g1 / (g1 + g0)
		// '(34)
		A := -g0 * g1 / (g0 + g1)
		// '(32)
		if x <= eta {
			return A*x - 1/3*(g0-A)*(math.Pow(eta-x, 3)/math.Pow(eta, 2)-eta)
		} else {
			return (2/3*A+1/3*g0)*eta + A*(x-eta) + (g1-A)/3*math.Pow(x-eta, 3)/math.Pow(1-eta, 2)
		}
	}
}

func forwardRate(Term float64, e initialFI) float64 {
	// 'numbering refers to Wilmott paper
	if Term <= 0 {
		return e.f[0]
	}
	if Term > e.Terms(e.N()) {
		return forwardRate(e.Terms(e.N()), e)
	}

	i, x, g0, g1 := initialInterpolators(e, Term)
	G := adjustedG(x, g0, g1)

	// '(26)
	return G + e.fD[i+1]
}

func adjustedG(x float64, g0 float64, g1 float64) float64 {
	if x == 0 {
		return g0
	} else if x == 1 {
		return g1
	} else if (g0 < 0 && -0.5*g0 <= g1 && g1 <= -2*g0) || (g0 > 0 && -0.5*g0 >= g1 && g1 >= -2*g0) {
		// 'zone (i)
		return g0*(1-4*x+3*math.Pow(x, 2)) + g1*(-2*x+3*math.Pow(x, 2))
	} else if (g0 < 0 && g1 > -2*g0) || (g0 > 0 && g1 < -2*g0) {
		// 'zone (ii)
		// '(29)
		eta := (g1 + 2*g0) / (g1 - g0)
		// '(28)
		if x <= eta {
			return g0
		} else {
			return g0 + (g1-g0)*math.Pow((x-eta)/(1-eta), 2)
		}
	} else if (g0 > 0 && 0 > g1 && g1 > -0.5*g0) || (g0 < 0 && 0 < g1 && g1 < -0.5*g0) {
		// 'zone (iii)
		// '(31)
		eta := 3 * g1 / (g1 - g0)
		// '(30)
		if x < eta {
			return g1 + (g0-g1)*math.Pow((eta-x)/eta, 2)
		} else {
			return g1
		}
	} else if g0 == 0 || g1 == 0 {
		return 0
	} else {
		// 'zone (iv)
		// '(33)
		eta := g1 / (g1 + g0)
		// '(34)
		A := -g0 * g1 / (g0 + g1)
		// '(32)
		if x <= eta {
			return A + (g0-A)*math.Pow((eta-x)/eta, 2)
		} else {
			return A + (g1-A)*math.Pow((eta-x)/(1-eta), 2)
		}
	}
}

func bound(Minimum float64, Variable float64, Maximum float64) float64 {
	return math.Max(Minimum, math.Min(Variable, Maximum))
}

func (e *initialFI) lastTermIndexBefore(Term float64) int {
	i := sort.SearchFloat64s(e.terms, Term)

	if i >= 1 && Term == e.terms[i-1] {
		return i - 1
	}

	return i
}

func estimateInitialFI(inp mcInput) initialFI {
	fD := make([]float64, inp.N()+1)
	interpolantAtNodeD := make([]float64, inp.N()+1)
	f := make([]float64, inp.N()+1)

	// 'step 1
	for j := 1; j < inp.N()+1; j++ {
		fD[j] = (inp.Terms(j)*inp.Values(j) - inp.Terms(j-1)*inp.Values(j-1)) / (inp.Terms(j) - inp.Terms(j-1))
		interpolantAtNodeD[j] = inp.Values(j)
	}

	// 'f_i estimation under the unameliorated method
	// 'numbering refers to Wilmott paper
	// 'step 2
	// '(22)
	for j := 1; j < inp.N(); j++ {
		f[j] = (inp.Terms(j)-inp.Terms(j-1))/(inp.Terms(j+1)-inp.Terms(j-1))*fD[j+1] +
			(inp.Terms(j+1)-inp.Terms(j))/(inp.Terms(j+1)-inp.Terms(j-1))*fD[j]
	}
	// '(23)
	f[0] = fD[1] - 0.5*(f[1]-fD[1])
	// '(24)
	f[inp.N()] = fD[inp.N()] - 0.5*(f[inp.N()-1]-fD[inp.N()])

	// 'step 3
	f[0] = bound(0, f[0], 2*fD[1])
	for j := 1; j < inp.N(); j++ {
		f[j] = bound(0, f[j], 2*math.Min(fD[j], fD[j+1]))
	}

	f[inp.N()] = bound(0, f[inp.N()], 2*fD[inp.N()])

	return initialFI{mcInput: inp, fD: fD, interpolantAtNodeD: interpolantAtNodeD, f: f}
}