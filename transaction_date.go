package main

import (
	"time"
)

const (
	TRANSACTION_DATE_FORMAT = "2006-01-02"
)

type TransactionDate struct {
	time.Time
}

func (date TransactionDate) String() string {
	return date.Time.Format(TRANSACTION_DATE_FORMAT)
}

func (date *TransactionDate) MarshalCSV() (string, error) {
	return date.Time.Format(TRANSACTION_DATE_FORMAT), nil
}

func (date *TransactionDate) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(TRANSACTION_DATE_FORMAT, csv)
	return err
}
