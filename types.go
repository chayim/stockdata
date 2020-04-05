package stockdata

import "time"

type Symbol struct {
	Symbol   string
	Name     string
	Sector   string
	Industry string
	ETF      bool
	DRIP     bool
}

type ExDividend struct {
	Symbol         string
	ExDividendDate time.Time
	PayableDate    time.Time
	Dividend       float64
	Yield          float64
	Currency       string
}
