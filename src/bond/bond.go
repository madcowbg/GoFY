package bond

import (
	m "../measures"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type Bond interface {
	TimeToExpiration(t m.Time) m.Time
	CurrentYield(t m.Time, yield m.Rate) m.Rate
	Price(t m.Time, yield m.Rate) m.Money
	YieldToMaturity(t m.Time, price m.Money) m.Rate
}

type Expirable struct {
	Maturity m.Time
}

func (expirable Expirable) TimeToExpiration(t m.Time) m.Time {
	return m.Time(math.Max(0, float64(expirable.Maturity-t)))
}

type ZeroCouponBond struct {
	Expirable
}

func (bond *ZeroCouponBond) CurrentYield(t m.Time, yield m.Rate) m.Rate {
	return 0.0
}

func (bond *ZeroCouponBond) Price(t m.Time, yield m.Rate) m.Money {
	return m.Money(math.Exp(-float64(bond.TimeToExpiration(t)) * float64(yield)))
}

func (bond *ZeroCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	return m.Rate(-math.Log(float64(price)) / float64(bond.TimeToExpiration(t)))
}

type FixedCouponTerm struct {
	Frequency float64
	PerAnnum  m.Money
}

type Cashflow struct {
	Time   m.Time
	Amount m.Money
}

func (coupon FixedCouponTerm) Cashflows(start m.Time, to m.Time) []Cashflow {
	cnt := int(math.Ceil(float64(to-start) * coupon.Frequency))
	result := make([]Cashflow, cnt)
	for i := 0; i < cnt; i++ {
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

func (bond FixedCouponBond) RemainingCashflows(t m.Time) []Cashflow {
	if t > bond.Maturity {
		return []Cashflow{}
	}

	if t < bond.IssueTime {
		t = bond.IssueTime
	}

	nextCoupon := bond.Coupon.NextCoupon(bond.IssueTime, t)
	couponCashflows := bond.Coupon.Cashflows(nextCoupon, bond.Maturity)

	cashflows := make([]Cashflow, len(couponCashflows)+1)
	copy(cashflows, couponCashflows)
	cashflows[len(cashflows)-1] = Cashflow{bond.Maturity, 1}

	return cashflows
}

func (bond *FixedCouponBond) CurrentYield(t m.Time, yield m.Rate) m.Rate {
	if t < bond.IssueTime || t > bond.Maturity {
		return 0.0
	}
	return m.Rate(float64(bond.Coupon.PerAnnum) / float64(bond.Price(t, yield)))
}

func (bond *FixedCouponBond) Price(t m.Time, yield m.Rate) m.Money {
	price := m.Money(0.0)
	for _, cashflow := range bond.RemainingCashflows(t) {
		price += cashflow.Amount * (&ZeroCouponBond{Expirable{cashflow.Time}}).Price(t, yield)
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
