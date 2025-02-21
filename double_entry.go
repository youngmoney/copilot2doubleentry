package main

type DoubleEntryTransaction struct {
	Date        TransactionDate `csv:"Date"`
	Description string          `csv:"Description"`
	Notes       string          `csv:"Notes"`
	Memo        string          `csv:"Memo"`
	AccountName string          `csv:"Account Name"`
	AmountNum   Amount          `csv:"Amount Num."`
}
