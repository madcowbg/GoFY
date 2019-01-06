package option

type VanillaOption struct {
	S Money
	T Time
}

func (option *VanillaOption) Expiration() Time {
	return option.T
}

func (option *VanillaOption) Strike() Money {
	return option.S
}
