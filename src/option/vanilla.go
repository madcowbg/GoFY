package option

type VanillaOption struct {
	Strike Money
	T      Time
}

func (option *VanillaOption) Maturity() Time {
	return option.T
}
