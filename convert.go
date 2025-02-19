package main

import (
	"fmt"
	"strings"
)

func (t *CopilotTransaction) CreateDescription() string {
	var parts []string
	if len(t.Name) > 0 {
		parts = append(parts, t.Name)
	}
	if len(t.Recurring) > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", t.Recurring))
	}

	return strings.Join(parts, " ")
}

func SplitTags(t string) []string {
	var tags []string
	for _, tag := range strings.Split(t, ",") {
		trimmed := strings.TrimSpace(tag)
		underscored := strings.ReplaceAll(trimmed, " ", "_")
		hashed := "#" + underscored
		tags = append(tags, hashed)
	}
	return tags
}

func (t *CopilotTransaction) CreateNote() string {
	var parts []string
	if len(t.Note) > 0 {
		parts = append(parts, t.Note)
	}
	if len(t.Tags) > 0 {
		tags := SplitTags(t.Tags)
		joined := strings.Join(tags, " ")
		parts = append(parts, "["+joined+"]")
	}

	return strings.Join(parts, " ")
}

func ConvertRegular(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	// TODO Overrides
	var expense DoubleEntryTransaction
	var liability DoubleEntryTransaction

	expense.Date = t.Date
	liability.Date = t.Date

	expense.Description = t.CreateDescription()
	liability.Description = t.CreateDescription()
	expense.Notes = t.CreateNote()
	liability.Notes = t.CreateNote()

	expense.AmountNum = t.Amount
	liability.AmountNum = -t.Amount

	expense.AccountName = t.Category
	// Override
	liability.AccountName = t.Account

	return expense, liability
}

func Convert(transactions []*CopilotTransaction, config Config) []DoubleEntryTransaction {
	var converted []DoubleEntryTransaction

	for _, t := range transactions {
		if t.Status == PENDING {
			continue
		}
		// TODO: Add Date Filter

		if t.Type == REGULAR {
			one, two := ConvertRegular(t, config)
			fmt.Println(one)
			fmt.Println(two)
		}
		fmt.Println(t)
	}

	return converted
}
