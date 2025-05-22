package main

import (
	"os"
	"text/template"
)

type OverrideData struct {
	Description string
	Account     string
	Negative    bool
}

const (
	kExpenseOverrideTemplate = `{
  "overrides": {
    "expense": [
      {
        "match": {
          "description": "{{ .Description }}",
          "account": "{{ .Account }}",
          "outgoing":  {{ if .Negative }}false{{else}}true{{end}}
        },
        "split_account": <expense account>
      }
    ]
  }
}`
	kIncomeOverrideTemplate = `{
  "overrides": {
    "income": [
      {
        "match": {
          "description": "{{ .Description }}",
          "account": "{{ .Account }}",
          "outgoing":  {{ if .Negative }}false{{else}}true{{end}}
        },
        "split_account": <income source account>
      }
    ]
  }
}`
	kTransferOutOverrideTemplate = `{
  "overrides": {
    "transfer": [
      {
        "match": {
          "description": "{{ .Description }}",
          "account": "{{ .Account }}",
          "outgoing":  true
        },
        "split_account": <transfer destination account>
      }
    ]
  }
}`
	kTransferInOverrideTemplate = `{
  "overrides": {
    "transfer": [
      {
        "match": {
          "description": "{{ .Description }}",
          "account": "{{ .Account }}",
          "outgoing":  false
        },
        "split_account": <transfer source account>
      }
    ]
  }
}`
)

func ErrorWithConfigTemplate(tmpl string, transaction *CopilotTransaction) {
	var o OverrideData
	o.Description = transaction.Name
	o.Account = transaction.Account
	o.Negative = transaction.Amount.IsNegative()

	t := template.Must(template.New("config").Parse(tmpl))
	os.Stderr.WriteString("config needs to be updated to resolve this transaction:\n")
	t.Execute(os.Stderr, o)
	os.Stderr.WriteString("\n")
	os.Exit(1)
}
