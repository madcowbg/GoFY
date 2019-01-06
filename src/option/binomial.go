package option

import "math"

func binomialModel(option Option, spot Money, t Time, sigmaX Return, rf Rate) Money {
	if t >= option.Expiration() {
		return option.Payoff(spot)
	}

	sigma := float64(sigmaX)
	r := float64(rf)

	nsteps := 1000
	step := (float64(option.Expiration()) - float64(t)) / float64(nsteps)

	discountFactor := math.Exp(-r * step)

	temp2 := 0.5 * (discountFactor + math.Exp((r+sigma*sigma)*step))

	u := temp2 + math.Sqrt(math.Pow(temp2, 2)-1)
	d := 1 / u
	p := (math.Exp(r*step) - d) / (u - d)

	S := make([]Money, nsteps+1)
	S[0] = spot

	for n := 1; n <= nsteps; n++ {
		for j := n; j >= 1; j-- {
			S[j] = Money(u * float64(S[j-1]))
		}
		S[0] = Money(d * float64(S[0]))
	}

	V := make([]Money, nsteps+1)
	for j := 0; j <= nsteps; j++ {
		V[j] = option.Payoff(S[j])
	}

	for n := nsteps; n >= 1; n-- {
		for j := 0; j < n; j++ {
			S[j] = Money(float64(S[j]) / d)
		}
		S[n] = Money(math.NaN())

		for j := 0; j < n; j++ {
			V[j] = option.EarlyExcercise(
				option.Payoff(S[j]),
				Money((p*float64(V[j+1])+(1-p)*float64(V[j]))*discountFactor))
		}
	}
	return V[0]
}
