package option

import (
	m "../measures"
	"gonum.org/v1/gonum/stat"
	"math"
	"math/rand"
)

func NoEarlyExcerciseMonteCarloModel(npaths int, nsteps int, seed int64) func(option Option, spot m.Money, t m.Time, sigmaX m.Return, rf m.Rate) m.Money {
	return func(option Option, spot m.Money, t m.Time, sigmaX m.Return, rf m.Rate) m.Money {
		if t >= option.Expiration() {
			return option.Payoff(spot)
		}

		sigma := float64(sigmaX)
		r := float64(rf)

		dt := float64(option.Expiration()-t) / float64(nsteps)

		S := generateExpiryS(npaths, nsteps, spot, r, sigma, dt, rand.New(rand.NewSource(seed)))

		V := make([]float64, npaths)
		for i := 0; i < npaths; i++ {
			V[i] = math.Exp(-r*float64(option.Expiration()-t)) * float64(option.Payoff(m.Money(S[i])))
		}

		return m.Money(stat.Mean(V, nil))
	}
}

func generateExpiryS(npaths int, nsteps int, spot m.Money, r float64, sigma float64, dt float64, randSrc *rand.Rand) []float64 {
	S := make([]float64, npaths)
	for p := 0; p < npaths; p++ {
		S[p] = float64(spot)

		for k := 0; k < nsteps; k++ {
			// dS = r dt + sigma * S * dX
			S[p] *= math.Exp((r-sigma*sigma/2)*dt + sigma*math.Sqrt(dt)*randSrc.NormFloat64())
		}
	}
	return S
}
