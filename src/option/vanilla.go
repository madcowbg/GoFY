package option

import m "../measures"

type VanillaOption struct {
	S m.Money
	T m.Time
}

func (option *VanillaOption) Expiration() m.Time {
	return option.T
}

func (option *VanillaOption) Strike() m.Money {
	return option.S
}
