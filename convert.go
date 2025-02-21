package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
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
		if m.Outgoing != nil && t.Amount.IsNegative() == *m.Outgoing {
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

func ConvertExpense(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
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
	liability.AmountNum = Amount{t.Amount.Neg()}

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

	expense.AmountNum = Amount{t.Amount.Neg()}
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

func ConvertInternalTransferExpense(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	override := FindOverride(t, config.Overrides.Transfer)
	if override != nil && override.AlwaysPair != nil && *override.AlwaysPair == true {
		log.Fatal("No pair for %s", t)
	}
	var expense DoubleEntryTransaction
	var liability DoubleEntryTransaction

	expense.Date = t.Date
	liability.Date = t.Date

	expense.Description = t.CreateDescription()
	liability.Description = t.CreateDescription()
	expense.Notes = t.CreateNote()
	liability.Notes = t.CreateNote()

	expense.AmountNum = t.Amount
	liability.AmountNum = Amount{t.Amount.Neg()}

	expense.AccountName = "TRANSFER OUT"
	if override != nil && override.SplitAccount != nil {
		expense.AccountName = *override.SplitAccount
	}
	liability.AccountName = t.Account
	if override != nil && override.Account != nil {
		liability.AccountName = *override.Account
	}

	return expense, liability
}

func ConvertInternalTransferIncome(t *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	override := FindOverride(t, config.Overrides.Transfer)
	if override != nil && override.AlwaysPair != nil && *override.AlwaysPair == true {
		log.Fatal("No pair for %s", t)
	}
	var expense DoubleEntryTransaction
	var liability DoubleEntryTransaction

	expense.Date = t.Date
	liability.Date = t.Date

	expense.Description = t.CreateDescription()
	liability.Description = t.CreateDescription()
	expense.Notes = t.CreateNote()
	liability.Notes = t.CreateNote()

	expense.AmountNum = Amount{t.Amount.Neg()}
	liability.AmountNum = t.Amount

	expense.AccountName = t.Account
	if override != nil && override.Account != nil {
		expense.AccountName = *override.Account
	}
	liability.AccountName = "TRANSFER INCOME"
	if override != nil && override.SplitAccount != nil {
		liability.AccountName = *override.SplitAccount
	}

	return expense, liability
}

func ConvertInternalTransferPair(p *CopilotTransaction, n *CopilotTransaction, config Config) (DoubleEntryTransaction, DoubleEntryTransaction) {
	// override_p := FindOverride(p, config.Overrides.Transfer)
	// override_n := FindOverride(n, config.Overrides.Transfer)

	name := p.Name + "/" + n.Name

	var expense DoubleEntryTransaction
	var liability DoubleEntryTransaction

	expense.Date = p.Date
	liability.Date = p.Date

	expense.Description = name
	liability.Description = name
	expense.Notes = p.CreateNote()
	liability.Notes = p.CreateNote()

	expense.AmountNum = p.Amount
	liability.AmountNum = Amount{p.Amount.Neg()}

	expense.AccountName = n.Account
	liability.AccountName = p.Account

	return expense, liability
}

func ProcessTransfers(transfers []*CopilotTransaction, config Config) []DoubleEntryTransaction {
	var negatives, positives []*CopilotTransaction

	for _, t := range transfers {
		if t.Amount.Sign() < 0 {
			negatives = append(negatives, t)
		} else {
			positives = append(positives, t)
		}
	}

	var used_negatives []*CopilotTransaction

	var converted []DoubleEntryTransaction

	for _, p := range positives {
		var paired bool

		for _, n := range negatives {
			if p.Amount.Add(n.Amount.Decimal).IsZero() {
				if math.Abs(p.Date.Time.Sub(n.Date.Time).Hours()) < 5*24 {
					used_negatives = append(used_negatives, n)
					paired = true
					one, two := ConvertInternalTransferPair(p, n, config)
					converted = append(converted, one, two)
				}
			}
		}

		if !paired {
			one, two := ConvertInternalTransferExpense(p, config)
			converted = append(converted, one, two)

		}
	}

	for _, n := range negatives {
		if slices.Contains(used_negatives, n) {
			continue
		}
		one, two := ConvertInternalTransferIncome(n, config)
		converted = append(converted, one, two)
	}
	return converted
}

func Convert(transactions []*CopilotTransaction, config Config, firstDay time.Time, lastDay time.Time) []DoubleEntryTransaction {
	var converted []DoubleEntryTransaction

	var transfers []*CopilotTransaction
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
			one, two := ConvertExpense(t, config)
			converted = append(converted, one, two)
		case INCOME:
			one, two := ConvertIncome(t, config)
			converted = append(converted, one, two)
		case INTERNAL_TRANSFER:
			transfers = append(transfers, t)

		}
	}

	converted = append(converted, ProcessTransfers(transfers, config)...)

	sort.SliceStable(converted, func(i, j int) bool {
		return converted[i].Date.Time.After(converted[j].Date.Time)
	})
	return converted
}
