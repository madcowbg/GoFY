package measures

type Time float64
type Money float64
type Return float64
type Rate float64

type DiscountFactor func(t Time) Money
type SpotRate func(t Time) Rate
