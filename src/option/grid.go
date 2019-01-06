package option

func FiniteDifferenceGrid(NAS int, SInf float64) func(option Option, t Time, sigma Return, rf Rate) (S []float64, V [][]float64) {
	return func(option Option, t Time, sigma Return, rf Rate) (S []float64, V [][]float64) {
		Vol := float64(sigma)
		RF := float64(rf)

		dS, dt, NTS := gridParameters(SInf, NAS, Vol, float64(option.Expiration()-t))
		S, V = grids(NAS, NTS)

		initializeBoundaryCondition(V, S, dS, option)

		for k := 1; k <= NTS; k++ {
			updateGridInteriorAtTime(k, V, S, dS, Vol, RF, dt)
			updateGridBoundaryAtTime(k, V, RF, dt)
			updateForEarlyExcerciseAtTime(k, option, V, S)
		}

		return S, V
	}
}

func initializeBoundaryCondition(V [][]float64, S []float64, dS float64, option Option) {
	NAS := len(V) - 1
	for i := 0; i <= NAS; i++ {
		S[i] = float64(i) * dS
		V[i][0] = float64(option.Payoff(Money(S[i])))
	}
}

func updateGridBoundaryAtTime(k int, V [][]float64, RF float64, dt float64) {
	NAS := len(V) - 1
	// Boundary condition at S=0
	V[0][k] = V[0][k-1] * (1 - RF*dt)
	// Boundary condition at S=infinity
	V[NAS][k] = 2*V[NAS-1][k] - V[NAS-2][k]
}

func updateGridInteriorAtTime(k int, V [][]float64, S []float64, dS float64, Vol float64, RF float64, dt float64) {
	for i := 1; i < len(V)-1; i++ {
		Delta := (V[i+1][k-1] - V[i-1][k-1]) / (2.0 * dS)
		Gamma := (V[i+1][k-1] - 2*V[i][k-1] + V[i-1][k-1]) / (dS * dS)
		// Black-Scholes to derive Theta
		Theta := (-0.5 * Vol * Vol * S[i] * S[i] * Gamma) - (RF * S[i] * Delta) + (RF * V[i][k-1])
		V[i][k] = V[i][k-1] - dt*Theta
	}
}

func updateForEarlyExcerciseAtTime(k int, option Option, V [][]float64, S []float64) {
	for i := 0; i < len(V); i++ {
		V[i][k] = float64(option.EarlyExcercise(option.Payoff(Money(S[i])), Money(V[i][k])))
	}
}

func grids(NAS int, NTS int) ([]float64, [][]float64) {
	S := make([]float64, NAS+1)
	V := make([][]float64, NAS+1)
	for n := 0; n <= NAS; n++ {
		V[n] = make([]float64, NTS+1)
	}
	return S, V
}

func gridParameters(SInf float64, NAS int, Vol float64, TimeToExpiration float64) (float64, float64, int) {
	dS := SInf / float64(NAS)

	wantedDt := 0.9 / (Vol * Vol * float64(NAS) * float64(NAS))
	NTS := int(TimeToExpiration/wantedDt + 1)
	dt := TimeToExpiration / float64(NTS)
	return dS, dt, NTS
}
