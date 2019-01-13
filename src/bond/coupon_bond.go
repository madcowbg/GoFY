package bond

import (
	m "../measures"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type FixedCouponTerm struct {
	Frequency float64
	PerAnnum  m.Money
}

func (coupon FixedCouponTerm) Cashflows(start m.Time, to m.Time) []Cashflow {
	cnt := int(math.Ceil(float64(to-start) * coupon.Frequency))
	result := make([]Cashflow, cnt+1)
	for i := 0; i <= cnt; i++ {
		result[i] = Cashflow{
			Time:   start + m.Time(float64(i)/coupon.Frequency),
			Amount: coupon.PerAnnum / m.Money(coupon.Frequency)}
	}
	return result
}

func (coupon FixedCouponTerm) NextCoupon(startTime m.Time, t m.Time) m.Time {
	if t <= startTime {
		return startTime
	}

	cnt := math.Ceil(float64(t-startTime) * coupon.Frequency)
	return startTime + m.Time(cnt/coupon.Frequency)
}

type FixedCouponBond struct {
	Expirable
	IssueTime m.Time
	Coupon    FixedCouponTerm
}

const EPS = 1e-10

func (bond FixedCouponBond) RemainingCashflows(t m.Time) []Cashflow {
	if t > bond.Maturity {
		return []Cashflow{}
	}

	if t <= bond.IssueTime {
		t = bond.IssueTime + EPS
	}

	nextCoupon := bond.Coupon.NextCoupon(bond.IssueTime, t)
	couponCashflows := bond.Coupon.Cashflows(nextCoupon, bond.Maturity)

	cashflows := make([]Cashflow, len(couponCashflows)+1)
	copy(cashflows, couponCashflows)
	cashflows[len(cashflows)-1] = Cashflow{bond.Maturity, 1}

	return cashflows
}

func (bond *FixedCouponBond) CurrentYield(t m.Time, rate m.Rate) m.Rate {
	if t < bond.IssueTime || t > bond.Maturity {
		return 0.0
	}
	return m.Rate(float64(bond.Coupon.PerAnnum) / float64(bond.Price(t, rate)))
}

func (bond *FixedCouponBond) Price(t m.Time, rate m.Rate) m.Money {
	price := m.Money(0.0)
	for _, cashflow := range bond.RemainingCashflows(t) {
		price += cashflow.Price(t, rate)
	}
	return price
}

func (bond *FixedCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	if t > bond.Maturity {
		return m.Rate(math.NaN())
	}

	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			return math.Abs(float64(bond.Price(t, m.Rate(x[0])) - price))
		},
	}
	result, err := optimize.Minimize(problem, []float64{0.0}, nil, &optimize.NelderMead{})
	if err != nil {
		return m.Rate(math.NaN())
	}
	return m.Rate(result.X[0])
}
