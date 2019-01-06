package option

type VanillaOption struct {
	Strike Money
	T      Time
}

func (option *VanillaOption) Expiration() Time {
	return option.T
}
