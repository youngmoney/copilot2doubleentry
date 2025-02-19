package main

import (
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
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

type CopilotTransactionStatus int64

const (
	PENDING CopilotTransactionStatus = 0
	POSTED  CopilotTransactionStatus = 1
)

func (status CopilotTransactionStatus) String() string {
	switch status {
	case PENDING:
		return "pending"
	case POSTED:
		return "posted"
	default:
		return fmt.Sprintf("%d", int64(status))
	}
}

func (status *CopilotTransactionStatus) MarshalCSV() (string, error) {
	switch *status {
	case PENDING:
		return "pending", nil
	case POSTED:
		return "posted", nil
	default:
		return "", errors.New("invalid CopilotTransactionStatus")
	}
}

func (status *CopilotTransactionStatus) UnmarshalCSV(csv string) (err error) {
	switch csv {
	case "pending":
		*status = PENDING
	case "posted":
		*status = POSTED
	default:
		return errors.New(fmt.Sprintf("invalid CopilotTransactionStatus string %s", csv))
	}
	return nil

}

type CopilotTransactionType int64

const (
	INCOME            CopilotTransactionType = 0
	REGULAR           CopilotTransactionType = 1
	INTERNAL_TRANSFER CopilotTransactionType = 2
)

func (t CopilotTransactionType) String() string {
	switch t {
	case INCOME:
		return "income"
	case REGULAR:
		return "regular"
	case INTERNAL_TRANSFER:
		return "internal transfer"
	default:
		return fmt.Sprintf("%d", int64(t))
	}
}

func (t *CopilotTransactionType) MarshalCSV() (string, error) {
	switch *t {
	case INCOME:
		return "income", nil
	case REGULAR:
		return "regular", nil
	case INTERNAL_TRANSFER:
		return "internal transfer", nil
	default:
		return "", errors.New("invalid CopilotTransactionType")
	}
}

func (t *CopilotTransactionType) UnmarshalCSV(csv string) (err error) {
	switch csv {
	case "income":
		*t = INCOME
	case "regular":
		*t = REGULAR
	case "internal transfer":
		*t = INTERNAL_TRANSFER
	default:
		return errors.New(fmt.Sprintf("invalid CopilotTransactionType string %s", csv))
	}
	return nil

}

type CopilotTransaction struct {
	Date TransactionDate `csv:"date"`
	Name string          `csv:"name"`
	// TODO: Make this fixed precision
	Amount         float64                  `csv:"amount"`
	Status         CopilotTransactionStatus `csv:"status"`
	Category       string                   `csv:"category"`
	ParentCategory string                   `csv:"parent category"`
	Excluded       bool                     `csv:"excluded"`
	Tags           string                   `csv:"tags"`
	Type           CopilotTransactionType   `csv:"type"`
	Account        string                   `csv:"account"`
	AccountMask    string                   `csv:"account mask"`
	Note           string                   `csv:"note"`
	Recurring      string                   `csv:"recurring"`
}

func ParseDateFlag(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func ReadCopilot(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	transactions := []*CopilotTransaction{}

	if err := gocsv.UnmarshalFile(file, &transactions); err != nil { // Load clients from file
		panic(err)
	}
	for _, t := range transactions {
		fmt.Println(t)
	}
}
