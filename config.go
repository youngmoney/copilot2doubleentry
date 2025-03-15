package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type TransactionMatch struct {
	DescriptionRegex *regexp.Regexp `json:"description"`
	AccountRegex     *regexp.Regexp `json:"account"`
	// if positive in copilot
	Outgoing *bool `json:"outgoing"`
}

type Override struct {
	Match        TransactionMatch `json:"match"`
	Account      *string          `json:"account"`
	SplitAccount *string          `json:"split_account"`
	AlwaysPair   *bool            `json:"always_pair"`
}

type Overrides struct {
	Income   []Override `json:"income"`
	Expense  []Override `json:"expense"`
	Transfer []Override `json:"transfer"`
}

type Config struct {
	Overrides Overrides `json:"overrides"`
}

func ReadConfig(filename string) Config {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config := Config{}
	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return config
}
