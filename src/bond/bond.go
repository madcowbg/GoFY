package bond

import (
	m "../measures"
	"fmt"
	"math"
)

type Expirable struct {
	Maturity m.Time
}

func (expirable Expirable) TimeToExpiration(t m.Time) float64 {
	return math.Max(0, float64(expirable.Maturity-t))
}

type ZeroCouponBond struct {
	Expirable
}

func (bond ZeroCouponBond) CurrentYield(t m.Time, yield m.Rate) m.Rate {
	return 0.0
}

func (bond ZeroCouponBond) Price(t m.Time, yield m.Rate) m.Money {
	return m.Money(math.Exp(-bond.TimeToExpiration(t) * float64(yield)))
}

func (bond ZeroCouponBond) YieldToMaturity(t m.Time, price m.Money) m.Rate {
	return m.Rate(-math.Log(float64(price)) / bond.TimeToExpiration(t))
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
	fmt.Printf("start: %f, to: %f, cnt: %d\n", start, to, cnt)
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
	fmt.Printf("next coupon: %f\n", nextCoupon)
	couponCashflows := bond.Coupon.Cashflows(nextCoupon, bond.Maturity)

	cashflows := make([]Cashflow, len(couponCashflows)+1)
	copy(cashflows, couponCashflows)
	cashflows[len(cashflows)-1] = Cashflow{bond.Maturity, 1}

	return cashflows
}
