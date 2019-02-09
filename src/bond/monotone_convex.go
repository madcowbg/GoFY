package bond

import (
	"math"
	"sort"
)

type MCInput struct {

	// 'extend the curve to time 0, for the purpose of calculating forward at time 1
	Terms  []float64 //Note: 0, t1, t2, ...
	Values []float64 //Note: v1, v1, v2, ...

	N int /*last*/
}

type f_i_estimates struct {
	MCInput

	fdiscrete          []float64
	dInterpolantatNode []float64

	f []float64
}

func Interpolant(Term float64, e f_i_estimates) float64 {
	// 'numbering refers to Wilmott paper, functions are integrated.
	if Term <= 0 {
		return e.f[0]
	}
	if Term > e.Terms[e.N] {
		return Interpolant(e.Terms[e.N], e)*e.Terms[e.N]/Term + Forward(e.Terms[e.N], e)*(1-e.Terms[e.N]/Term)
	}

	i := e.lastTenorBefore(Term)
	// 'the x in (25)
	x := (Term - e.Terms[i]) / (e.Terms[i+1] - e.Terms[i])
	g0 := e.f[i] - e.fdiscrete[i+1]
	g1 := e.f[i+1] - e.fdiscrete[i+1]

	var G float64 // todo make function
	if x == 0 || x == 1 {
		G = 0
	} else if (g0 < 0 && -0.5*g0 <= g1 && g1 <= -2*g0) || (g0 > 0 && -0.5*g0 >= g1 && g1 >= -2*g0) {
		// 'zone (i)
		G = g0*(x-2*math.Pow(x, 2)+math.Pow(x, 3)) + g1*(-math.Pow(x, 2)+math.Pow(x, 3))
	} else if (g0 < 0 && g1 > -2*g0) || (g0 > 0 && g1 < -2*g0) {
		//'zone (ii)
		// '(29)
		eta := (g1 + 2*g0) / (g1 - g0)
		// '(28)
		if x <= eta {
			G = g0 * x
		} else {
			G = g0*x + (g1-g0)*math.Pow(x-eta, 3)/math.Pow(1-eta, 2)/3
		}
	} else if (g0 > 0 && 0 > g1 && g1 > -0.5*g0) || (g0 < 0 && 0 < g1 && g1 < -0.5*g0) {
		// 'zone (iii)
		// '(31)
		eta := 3 * g1 / (g1 - g0)
		//'(30)
		if x < eta {
			G = g1*x - 1/3*(g0-g1)*(math.Pow(eta-x, 3)/math.Pow(eta, 2)-eta)
		} else {
			G = (2/3*g1+1/3*g0)*eta + g1*(x-eta)
		}
	} else if g0 == 0 || g1 == 0 {
		G = 0
	} else {
		// 'zone (iv)
		// '(33)
		eta := g1 / (g1 + g0)
		// '(34)
		A := -g0 * g1 / (g0 + g1)
		// '(32)
		if x <= eta {
			G = A*x - 1/3*(g0-A)*(math.Pow(eta-x, 3)/math.Pow(eta, 2)-eta)
		} else {
			G = (2/3*A+1/3*g0)*eta + A*(x-eta) + (g1-A)/3*math.Pow(x-eta, 3)/math.Pow(1-eta, 2)
		}
	}

	//'(12)
	return 1 / Term * (e.Terms[i]*e.dInterpolantatNode[i] + (Term-e.Terms[i])*e.fdiscrete[i+1] + (e.Terms[i+1]-e.Terms[i])*G)
}

func Forward(Term float64, e f_i_estimates) float64 {
	// 'numbering refers to Wilmott paper
	if Term <= 0 {
		return e.f[0]
	}
	if Term > e.Terms[e.N] {
		return Forward(e.Terms[e.N], e)
	}

	i := e.lastTenorBefore(Term)

	// 'the x in (25)
	x := (Term - e.Terms[i]) / (e.Terms[i+1] - e.Terms[i])
	g0 := e.f[i] - e.fdiscrete[i+1]
	g1 := e.f[i+1] - e.fdiscrete[i+1]

	var G float64 // todo make function
	if x == 0 {
		G = g0
	} else if x == 1 {
		G = g1
	} else if (g0 < 0 && -0.5*g0 <= g1 && g1 <= -2*g0) || (g0 > 0 && -0.5*g0 >= g1 && g1 >= -2*g0) {
		// 'zone (i)
		G = g0*(1-4*x+3*math.Pow(x, 2)) + g1*(-2*x+3*math.Pow(x, 2))
	} else if (g0 < 0 && g1 > -2*g0) || (g0 > 0 && g1 < -2*g0) {
		// 'zone (ii)
		// '(29)
		eta := (g1 + 2*g0) / (g1 - g0)
		// '(28)
		if x <= eta {
			G = g0
		} else {
			G = g0 + (g1-g0)*math.Pow((x-eta)/(1-eta), 2)
		}
	} else if (g0 > 0 && 0 > g1 && g1 > -0.5*g0) || (g0 < 0 && 0 < g1 && g1 < -0.5*g0) {
		// 'zone (iii)
		// '(31)
		eta := 3 * g1 / (g1 - g0)
		// '(30)
		if x < eta {
			G = g1 + (g0-g1)*math.Pow((eta-x)/eta, 2)
		} else {
			G = g1
		}
	} else if g0 == 0 || g1 == 0 {
		G = 0
	} else {
		// 'zone (iv)
		// '(33)
		eta := g1 / (g1 + g0)
		// '(34)
		A := -g0 * g1 / (g0 + g1)
		// '(32)
		if x <= eta {
			G = A + (g0-A)*math.Pow((eta-x)/eta, 2)
		} else {
			G = A + (g1-A)*math.Pow((eta-x)/(1-eta), 2)
		}
	}
	// '(26)
	return G + e.fdiscrete[i+1]
}

func bound(Minimum float64, Variable float64, Maximum float64) float64 {
	return math.Max(Minimum, math.Min(Variable, Maximum))
}

func (e *f_i_estimates) lastTenorBefore(Term float64) int {
	i := sort.SearchFloat64s(e.Terms, Term)
	if i == 0 {
		return 0
	}

	if i >= 1 && Term == e.Terms[i-1] {
		return i - 2
	}

	return i - 1
}

func InitialFIEstimates(inp MCInput) f_i_estimates {
	fdiscrete := make([]float64, len(inp.Terms))
	dInterpolantatNode := make([]float64, len(inp.Terms))
	f := make([]float64, len(inp.Terms))

	// 'step 1
	for j := 1; j < len(inp.Terms); j++ {
		fdiscrete[j] = (inp.Terms[j]*inp.Values[j] - inp.Terms[j-1]*inp.Values[j-1]) / (inp.Terms[j] - inp.Terms[j-1])
		dInterpolantatNode[j] = inp.Values[j]
	}

	// 'f_i estimation under the unameliorated method
	// 'numbering refers to Wilmott paper
	// 'step 2
	// '(22)
	for j := 1; j < len(f)-1; j++ {
		f[j] = (inp.Terms[j]-inp.Terms[j-1])/(inp.Terms[j+1]-inp.Terms[j-1])*fdiscrete[j+1] + (inp.Terms[j+1]-inp.Terms[j])/(inp.Terms[j+1]-inp.Terms[j-1])*fdiscrete[j]
	}
	// '(23)
	f[0] = fdiscrete[1] - 0.5*(f[1]-fdiscrete[1])
	// '(24)
	f[len(f)-1] = fdiscrete[len(f)-1] - 0.5*(f[len(f)-2]-fdiscrete[len(f)-1])

	// 'step 3
	f[0] = bound(0, f[0], 2*fdiscrete[1])
	for j := 1; j < inp.N; j++ {
		f[j] = bound(0, f[j], 2*math.Min(fdiscrete[j], fdiscrete[j+1]))
	}

	f[inp.N] = bound(0, f[inp.N], 2*fdiscrete[inp.N])

	return f_i_estimates{MCInput: inp, fdiscrete: fdiscrete, dInterpolantatNode: dInterpolantatNode, f: f}
}
