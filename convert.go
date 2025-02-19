package main

import (
	"fmt"
	"strings"
	"time"
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

func FindOverride(t *CopilotTransaction, overrides []Override) *Override {
	for _, o := range overrides {
		m := o.Match
		if m.Outgoing != nil && t.Amount < 0 == *m.Outgoing {
			continue
		}
		if m.DescriptionRegex != nil && !m.DescriptionRegex.MatchString(t.Name) {
			continue
		}
		if m.AccountRegex != nil && !m.AccountRegex.MatchString(t.Account) {
			continue
		}
		return &o
	}

	return nil
}

func ConvertRegular(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	override := FindOverride(t, config.Overrides.Expense)
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
	if override != nil && override.Account != nil {
		expense.AccountName = *override.Account
	}
	liability.AccountName = t.Account
	if override != nil && override.SplitAccount != nil {
		liability.AccountName = *override.SplitAccount
	}

	return expense, liability
}

func ConvertIncome(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	override := FindOverride(t, config.Overrides.Income)
	var expense DoubleEntryTransaction
	var liability DoubleEntryTransaction

	expense.Date = t.Date
	liability.Date = t.Date

	expense.Description = t.CreateDescription()
	liability.Description = t.CreateDescription()
	expense.Notes = t.CreateNote()
	liability.Notes = t.CreateNote()

	expense.AmountNum = -t.Amount
	liability.AmountNum = t.Amount

	expense.AccountName = t.Account
	if override != nil && override.Account != nil {
		expense.AccountName = *override.Account
	}
	liability.AccountName = "INCOME INCOME"
	if override != nil && override.SplitAccount != nil {
		liability.AccountName = *override.SplitAccount
	}

	return expense, liability
}

func Convert(transactions []*CopilotTransaction, config Config, firstDay time.Time, lastDay time.Time) []DoubleEntryTransaction {
	var converted []DoubleEntryTransaction

	for _, t := range transactions {
		if t.Status == PENDING {
			continue
		}

		if !firstDay.IsZero() && t.Date.Before(firstDay) {
			continue
		}
		if !lastDay.IsZero() && t.Date.After(lastDay) {
			continue
		}

		switch t.Type {
		case REGULAR:
			one, two := ConvertRegular(t, config)
			converted = append(converted, one, two)
		case INCOME:
			one, two := ConvertIncome(t, config)
			converted = append(converted, one, two)

		}
	}

	return converted
}
