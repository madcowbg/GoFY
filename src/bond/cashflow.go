package bond

import (
	m "../measures"
)

type Cashflow struct {
	Time   m.Time
	Amount m.Money
}

func (cashflow Cashflow) Price(t m.Time, rate m.Rate) m.Money {
	return cashflow.Amount * (&ZeroCouponBond{Expirable{cashflow.Time}}).Price(t, rate)
}
