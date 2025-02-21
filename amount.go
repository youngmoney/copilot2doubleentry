package main

import (
	"github.com/shopspring/decimal"
)

type Amount struct{ decimal.Decimal }

func (d *Amount) MarshalCSV() (string, error) {
	if d.IsInteger() {
		return d.StringFixed(0), nil
	}
	return d.StringFixed(2), nil
}
